package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

type snapMirrorRequest struct {
	ReplicationRequest replicationRequest `structs:"replicationRequest"`
	ReplicationVolume  replicationVolume  `structs:"replicationVolume"`
}

type replicationRequest struct {
	SourceWorkingEnvironmentID      string   `structs:"sourceWorkingEnvironmentId"`
	DestinationWorkingEnvironmentID string   `structs:"destinationWorkingEnvironmentId"`
	SourceInterclusterLifIps        []string `structs:"sourceInterclusterLifIps"`
	DestinationInterclusterLifIps   []string `structs:"destinationInterclusterLifIps"`
	PolicyName                      string   `structs:"policyName"`
	ScheduleName                    string   `structs:"scheduleName,omitempty"`
	MaxTransferRate                 int      `structs:"maxTransferRate,omitempty"`
}

type replicationVolume struct {
	SourceSvmName                 string  `structs:"sourceSvmName"`
	SourceVolumeName              string  `structs:"sourceVolumeName"`
	DestinationVolumeName         string  `structs:"destinationVolumeName"`
	DestinationAggregateName      string  `structs:"destinationAggregateName"`
	NumOfDisksApprovedToAdd       float64 `structs:"numOfDisksApprovedToAdd"`
	AdvancedMode                  bool    `structs:"advancedMode"`
	DestinationSvmName            string  `structs:"destinationSvmName,omitempty"`
	DestinationProviderVolumeType string  `structs:"destinationProviderVolumeType,omitempty"`
	DestinationCapacityTier       string  `structs:"destinationCapacityTier,omitempty"`
	Iops                          int     `structs:"iops,omitempty"`
	Throughput                    int     `structs:"throughput,omitempty"`
}

type interclusterlif struct {
	Interclusterlif     []interClusterLifsAddress `json:"interClusterLifs"`
	PeerInterclusterlif []interClusterLifsAddress `json:"peerInterClusterLifs"`
}

type interClusterLifsAddress struct {
	Address string `json:"address"`
}

type snapMirrorStatusResponse struct {
	Destination destination `structs:"destination"`
}

type destination struct {
	VolumeName string `structs:"volumeName"`
}

func (c *Client) getInterclusterlifs(snapMirror snapMirrorRequest) (interclusterlif, error) {
	baseURL := fmt.Sprintf("/occm/api/replication/intercluster-lifs?peerWorkingEnvironmentId=%s&workingEnvironmentId=%s", snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID, snapMirror.ReplicationRequest.SourceWorkingEnvironmentID)
	hostType := "CloudManagerHost"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("intercluster-lifs reading failed ", statusCode)
		return interclusterlif{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "intercluster-lifs")
	if responseError != nil {
		return interclusterlif{}, responseError
	}
	var interclusterlifsResponse interclusterlif

	if err := json.Unmarshal(response, &interclusterlifsResponse); err != nil {
		log.Print("Failed to unmarshall response from interclusterlif ", err)
		return interclusterlif{}, err
	}

	return interclusterlifsResponse, nil
}

func (c *Client) buildSnapMirrorCreate(snapMirror snapMirrorRequest, sourceWorkingEnvironmentType string, destWorkingEnvironmentType string) (snapMirrorRequest, error) {

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createSnapMirror request, failed to get AccessToken")
		return snapMirrorRequest{}, err
	}
	c.Token = accessTokenResult.Token

	interclusterlifsResponse, err := c.getInterclusterlifs(snapMirror)
	if err != nil {
		log.Print("intercluster-lifs reading failed")
		return snapMirrorRequest{}, err
	}

	var volumeSource []volumeResponse
	volumeS := volumeRequest{}
	volumeS.WorkingEnvironmentID = snapMirror.ReplicationRequest.SourceWorkingEnvironmentID
	volumeS.Name = snapMirror.ReplicationVolume.SourceVolumeName

	if sourceWorkingEnvironmentType != "ON_PREM" {
		volumeSource, err = c.getVolume(volumeS)
		if err != nil {
			log.Print("Error reading source volume")
			return snapMirrorRequest{}, err
		}
	} else {
		volumeSource, err = c.getVolumeForOnPrem(volumeS)
		if err != nil {
			log.Print("Error reading source volume")
			return snapMirrorRequest{}, err
		}
	}

	if len(volumeSource) == 0 {
		log.Print("source volume not found")
		return snapMirrorRequest{}, fmt.Errorf("source volume not found")
	}
	volFound := false
	var volDestQuote volumeResponse
	var sourceVolume volumeResponse
	for _, vol := range volumeSource {
		if vol.Name == volumeS.Name {
			volFound = true
			volDestQuote = vol
			sourceVolume = vol
			if snapMirror.ReplicationVolume.SourceSvmName != "" && vol.SvmName != snapMirror.ReplicationVolume.SourceSvmName {
				volFound = false
			}
		}
	}

	if volFound == false {
		log.Print("source volume not found")
		return snapMirrorRequest{}, fmt.Errorf("source volume not found")
	}

	if destWorkingEnvironmentType != "ON_PREM" {
		quote := c.buildQuoteRequest(snapMirror, volDestQuote, snapMirror.ReplicationRequest.SourceWorkingEnvironmentID, sourceWorkingEnvironmentType, snapMirror.ReplicationVolume.DestinationVolumeName, snapMirror.ReplicationVolume.DestinationSvmName, snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID)
		quoteResponse, err := c.quoteVolume(quote)
		if err != nil {
			log.Printf("Error quoting destination volume")
			return snapMirrorRequest{}, err
		}
		snapMirror.ReplicationVolume.NumOfDisksApprovedToAdd = quoteResponse["numOfDisks"].(float64)
		if snapMirror.ReplicationVolume.DestinationAggregateName != "" {
			snapMirror.ReplicationVolume.AdvancedMode = true
		} else {
			snapMirror.ReplicationVolume.AdvancedMode = false
			snapMirror.ReplicationVolume.DestinationAggregateName = quoteResponse["aggregateName"].(string)
		}
		if quote.Iops != 0 {
			snapMirror.ReplicationVolume.Iops = quote.Iops
		}
		if quote.Throughput != 0 {
			snapMirror.ReplicationVolume.Throughput = quote.Throughput
		}
	}

	var sourceInterclusterLifIps []string
	var destinationInterclusterLifIps []string
	sourceInterclusterLifIps = append(sourceInterclusterLifIps, interclusterlifsResponse.Interclusterlif[0].Address)
	snapMirror.ReplicationRequest.SourceInterclusterLifIps = sourceInterclusterLifIps
	destinationInterclusterLifIps = append(destinationInterclusterLifIps, interclusterlifsResponse.PeerInterclusterlif[0].Address)
	snapMirror.ReplicationRequest.DestinationInterclusterLifIps = destinationInterclusterLifIps
	snapMirror.ReplicationVolume.SourceSvmName = sourceVolume.SvmName
	snapMirror.ReplicationVolume.SourceVolumeName = sourceVolume.Name

	if snapMirror.ReplicationVolume.DestinationProviderVolumeType == "" {
		snapMirror.ReplicationVolume.DestinationProviderVolumeType = sourceVolume.ProviderVolumeType
	}

	err = c.createSnapMirror(snapMirror, destWorkingEnvironmentType)
	if err != nil {
		log.Printf("Error creating snapmirror")
		return snapMirrorRequest{}, err
	}

	return snapMirror, nil
}

func (c *Client) buildQuoteRequest(snapMirror snapMirrorRequest, vol volumeResponse, sourceWorkingEnvironmentID string, sourceWorkingEnvironmentType string, name string, svm string, workingEnvironmentID string) quoteRequest {
	var quote quoteRequest

	quote.Name = name
	quote.Size.Size = vol.Size.Size
	quote.Size.Unit = vol.Size.Unit
	quote.SnapshotPolicyName = vol.SnapshotPolicyName
	quote.EnableDeduplication = vol.EnableDeduplication
	quote.EnableThinProvisioning = vol.EnableThinProvisioning
	quote.EnableCompression = vol.EnableCompression
	quote.VerifyNameUniqueness = true
	quote.ReplicationFlow = true
	quote.WorkingEnvironmentID = workingEnvironmentID
	quote.SvmName = svm

	aggregate, err := c.getAggregate(aggregateRequest{WorkingEnvironmentID: sourceWorkingEnvironmentID}, vol.AggregateName, sourceWorkingEnvironmentType)
	if err != nil {
		log.Printf("Error getting aggregate. aggregate name = %v", vol.AggregateName)
	}
	if len(aggregate.ProviderVolumes) != 0 {
		// Iops and Throughput values are the same if the volumes under the same aggregate
		if aggregate.ProviderVolumes[0].DiskType == "gp3" || aggregate.ProviderVolumes[0].DiskType == "io1" || aggregate.ProviderVolumes[0].DiskType == "io2" {
			quote.Iops = aggregate.ProviderVolumes[0].Iops
		}
		if aggregate.ProviderVolumes[0].DiskType == "gp3" {
			quote.Throughput = aggregate.ProviderVolumes[0].Throughput
		}
	}

	if snapMirror.ReplicationVolume.DestinationCapacityTier != "" {
		quote.CapacityTier = snapMirror.ReplicationVolume.DestinationCapacityTier
	}
	if snapMirror.ReplicationVolume.DestinationProviderVolumeType == "" {
		quote.ProviderVolumeType = vol.ProviderVolumeType
	} else {
		quote.ProviderVolumeType = snapMirror.ReplicationVolume.DestinationProviderVolumeType
	}

	return quote
}

func (c *Client) createSnapMirror(sm snapMirrorRequest, destWorkingEnvironmentType string) error {
	var baseURL string
	if destWorkingEnvironmentType != "ON_PREM" {
		baseURL = "/occm/api/replication/vsa"
	} else {
		baseURL = "/occm/api/replication/onprem"
	}
	hostType := "CloudManagerHost"

	params := structs.Map(sm)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("createSnapMirror request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createSnapMirror")
	if responseError != nil {
		return responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "snapmirror", "create", 10, 10)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteSnapMirror(snapMirror snapMirrorRequest) error {

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteSnapMirror request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token
	baseURL := fmt.Sprintf("/occm/api/replication/%s/%s/%s", snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID, snapMirror.ReplicationVolume.DestinationSvmName, snapMirror.ReplicationVolume.DestinationVolumeName)
	hostType := "CloudManagerHost"

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("deleteSnapMirror request failed with statusCode:%v, Error:%v", statusCode, err)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteSnapMirror")
	if responseError != nil {
		return responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "snapmirror", "delete", 10, 10)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getSnapMirror(snapMirror snapMirrorRequest, vol string) (string, error) {

	var result []snapMirrorStatusResponse

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createSnapMirror request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/occm/api/replication/status/%s", snapMirror.ReplicationRequest.SourceWorkingEnvironmentID)

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getSnapMirror request failed ", statusCode)
		return "", err
	}
	responseError := apiResponseChecker(statusCode, response, "getSnapMirror")
	if responseError != nil {
		return "", responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getSnapMirror ", err)
		return "", err
	}
	for _, sm := range result {
		if sm.Destination.VolumeName == vol {
			return vol, nil
		}
	}

	return "", nil
}
