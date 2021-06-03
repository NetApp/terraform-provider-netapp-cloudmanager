package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/fatih/structs"
)

// createAggregateRequest the users input for creating an Aggregate
type createAggregateRequest struct {
	Name                 string   `structs:"name"`
	WorkingEnvironmentID string   `structs:"workingEnvironmentId"`
	NumberOfDisks        int      `structs:"numberOfDisks"`
	DiskSize             diskSize `structs:"diskSize"`
	HomeNode             string   `structs:"homeNode,omitempty"`
	ProviderVolumeType   string   `structs:"providerVolumeType,omitempty"`
	CapacityTier         string   `structs:"capacityTier,omitempty"`
	Iops                 int      `structs:"iops,omitempty"`
	Throughput           int      `structs:"throughput,omitempty"`
}

// diskSize struct
type diskSize struct {
	Size int    `structs:"size"`
	Unit string `structs:"unit"`
}

// aggregateResult from aggregate request
type aggregateResult struct {
	Name              string           `json:"name"`
	AvailableCapacity capacity         `json:"availableCapacity"`
	TotalCapacity     capacity         `json:"totalCapacity"`
	UsedCapacity      capacity         `json:"usedCapacity"`
	Volumes           []volume         `json:"volumes"`
	ProviderVolumes   []providerVolume `json:"providerVolumes"`
	Disks             []disk           `json:"disks"`
	State             string           `json:"state"`
	EncryptionType    string           `json:"encryptionType"`
	EncryptionKeyID   string           `json:"encryptionKeyId"`
	IsRoot            bool             `json:"isRoot"`
	HomeNode          string           `json:"homeNode"`
	OwnerNode         string           `json:"ownerNode"`
	CapacityTier      string           `json:"capacityTier"`
	CapacityTierUsed  capacity         `json:"capacityTierUsed"`
	SidlEnabled       bool             `json:"sidlEnabled"`
	SnaplockType      string           `json:"snaplockType"`
}

type capacity struct {
	Size float64 `json:"size"`
	Unit string  `json:"unit"`
}

type volume struct {
	Name            string   `json:"name"`
	TotalSize       capacity `json:"totalSize"`
	UsedSize        capacity `json:"usedSize"`
	ThinProvisioned bool     `json:"thinProvisioned"`
	IsClone         bool     `json:"isClone"`
	RootVolume      bool     `json:"rootVolume"`
}

type providerVolume struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Size       capacity `json:"size"`
	State      string   `json:"state"`
	Device     string   `json:"device"`
	InstanceID string   `json:"instanceId"`
	DiskType   string   `json:"diskType"`
	Encrypted  bool     `json:"encrypted"`
	Iops       int      `json:"iops"`
	Throughput int      `json:"throughput"`
}

type disk struct {
	Name             string           `json:"name"`
	Position         string           `json:"position"`
	OwnerNode        string           `json:"ownerNode"`
	Device           string           `json:"device"`
	VMDiskProperties vmDiskProperties `json:"vmDiskProperties"`
}

type vmDiskProperties struct {
	ObjectName         string `json:"objectName"`
	StorageAccountName string `json:"storageAccountName"`
	ContainerName      string `json:"containerName"`
}

type aggregateRequest struct {
	WorkingEnvironmentID string `structs:"workingEnvironmentId"`
}

type deleteAggregateRequest struct {
	WorkingEnvironmentID string `structs:"workingEnvironmentId"`
	Name                 string `structs:"name"`
}

type updateAggregateRequest struct {
	WorkingEnvironmentID string `structs:"workingEnvironmentId"`
	Name                 string `structs:"name"`
	NumberOfDisks        int    `structs:"numberOfDisks"`
}

// get aggregate by workingEnvironmentId+aggregate name
func (c *Client) getAggregate(request aggregateRequest, name string) (aggregateResult, error) {
	log.Printf("getAggregate %s...", name)
	hostType := "CloudManagerHost"

	var baseURL string
	rootURL, cloudProviderName, err := c.getAPIRoot(request.WorkingEnvironmentID)

	if err != nil {
		log.Print("getAggregate: Cannot get API root.")
		return aggregateResult{}, err
	}

	if cloudProviderName != "Amazon" {
		baseURL = fmt.Sprintf("%s/aggregates/%s", rootURL, request.WorkingEnvironmentID)
	} else {
		baseURL = fmt.Sprintf("%s/aggregates?workingEnvironmentId=%s", rootURL, request.WorkingEnvironmentID)
	}

	var aggregates []aggregateResult

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("getAggregate request failed. Response %v, err %v", response, err)
		return aggregateResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "getAggregate")
	if responseError != nil {
		return aggregateResult{}, responseError
	}

	if err := json.Unmarshal(response, &aggregates); err != nil {
		log.Print("Failed to unmarshall response from getAggregates")
		return aggregateResult{}, err
	}

	log.Printf("getAggregate: get list of aggregates. %v", aggregates)
	log.Printf("Find the match one. %v", name)

	for i := range aggregates {
		if aggregates[i].Name == name {
			log.Printf("Found aggregate: %v", aggregates[i])
			return aggregates[i], nil
		}
	}
	log.Print("Cannot find the aggregate")

	return aggregateResult{}, nil
}

// create aggregate
func (c *Client) createAggregate(request *createAggregateRequest) (aggregateResult, error) {
	log.Print("createAggregate... ")
	params := structs.Map(request)
	hostType := "CloudManagerHost"

	var baseURL string
	rootURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID)

	if err != nil {
		log.Print("createAggregate: Cannot get API root.")
		return aggregateResult{}, err
	}
	baseURL = fmt.Sprintf("%s/aggregates", rootURL)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("createAggregate request failed")
		return aggregateResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createAggregate")
	if responseError != nil {
		return aggregateResult{}, responseError
	}

	// wait for creation
	log.Print("Wait for aggregate creation.")
	err = c.waitOnCompletion(onCloudRequestID, "Aggregate", "create", 10, 60)
	log.Print("Finish waiting.")
	if err != nil {
		return aggregateResult{}, err
	}

	var aggregate aggregateResult
	aggregate, err = c.getAggregate(aggregateRequest{WorkingEnvironmentID: request.WorkingEnvironmentID}, request.Name)
	if err != nil {
		return aggregateResult{}, err
	}
	return aggregate, nil
}

// delete aggregate
func (c *Client) deleteAggregate(request deleteAggregateRequest) error {
	log.Print("On deleteAggregate... ")
	hostType := "CloudManagerHost"

	var baseURL string
	rootURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID)

	if err != nil {
		log.Print("deleteAggregate: Cannot get API root.")
		return err
	}

	baseURL = fmt.Sprintf("%s/aggregates/%s/%s", rootURL, request.WorkingEnvironmentID, request.Name)

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("deleteAggregate request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteAggregate")
	if responseError != nil {
		return responseError
	}

	log.Print("Wait for aggregate deletion.")
	err = c.waitOnCompletion(onCloudRequestID, "Aggregate", "delete", 10, 60)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) updateAggregate(request updateAggregateRequest) error {
	log.Print("updateAggregate... ")
	params := structs.Map(request)
	hostType := "CloudManagerHost"

	var baseURL string
	rootURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID)

	if err != nil {
		log.Print("updateAggregate: Cannot get API root.")
		return err
	}
	baseURL = fmt.Sprintf("%s/aggregates/%s/%s/disks", rootURL, request.WorkingEnvironmentID, request.Name)

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("updateAggregate request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "updateAggregate")
	if responseError != nil {
		return responseError
	}

	log.Print("Wait for aggregate update.")
	err = c.waitOnCompletion(onCloudRequestID, "Aggregate", "update", 10, 60)
	if err != nil {
		return err
	}

	return nil
}

// flattenCapacity: convert struct size + unit
func flattenCapacity(c capacity) interface{} {
	flattened := make(map[string]interface{})
	flattened["size"] = strconv.FormatFloat(c.Size, 'f', -1, 64)
	flattened["unit"] = c.Unit
	return flattened
}

func flattenDisks(d []disk) interface{} {
	dts := make([]map[string]interface{}, 0, len(d))
	for _, diskelement := range d {
		dt := make(map[string]interface{})
		dt["name"] = diskelement.Name
		dt["position"] = diskelement.Position
		dt["device"] = diskelement.Device
		dt["owner_node"] = diskelement.OwnerNode
		vdp := make(map[string]interface{})
		vdp["object_name"] = diskelement.VMDiskProperties.ObjectName
		vdp["storage_account_name"] = diskelement.VMDiskProperties.StorageAccountName
		vdp["container_name"] = diskelement.VMDiskProperties.ContainerName
		dt["vm_disk_properties"] = vdp
		dts = append(dts, dt)
	}
	return dts
}

func flattenVolumes(v []volume) interface{} {
	volumes := make([]map[string]interface{}, 0, len(v))
	for _, volume := range v {
		vol := make(map[string]interface{})
		vol["name"] = volume.Name
		vol["thin_provisioned"] = volume.ThinProvisioned
		vol["root_volume"] = volume.RootVolume
		vol["is_clone"] = volume.IsClone
		vol["total_size"] = flattenCapacity(volume.TotalSize)
		vol["used_size"] = flattenCapacity(volume.UsedSize)

		volumes = append(volumes, vol)
	}
	return volumes
}

func flattenProviderVolumes(v []providerVolume) interface{} {
	volumes := make([]map[string]interface{}, 0, len(v))
	for _, volume := range v {
		vol := make(map[string]interface{})
		vol["id"] = volume.ID
		vol["name"] = volume.Name
		vol["state"] = volume.State
		vol["device"] = volume.Device
		vol["instance_id"] = volume.InstanceID
		vol["disk_type"] = volume.DiskType
		vol["encrypted"] = volume.Encrypted
		vol["iops"] = volume.Iops
		vol["throughput"] = volume.Throughput
		vol["size"] = flattenCapacity(volume.Size)

		volumes = append(volumes, vol)
	}
	return volumes
}
