package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

type volumeRequest struct {
	Name                      string                 `structs:"name"`
	WorkingEnvironmentID      string                 `structs:"workingEnvironmentId"`
	SvmName                   string                 `structs:"svmName"`
	AggregateName             string                 `structs:"aggregateName"`
	Size                      size                   `structs:"size"`
	SnapshotPolicyName        string                 `structs:"snapshotPolicyName,omitempty"`
	EnableThinProvisioning    bool                   `structs:"enableThinProvisioning"`
	EnableCompression         bool                   `structs:"enableCompression"`
	EnableDeduplication       bool                   `structs:"enableDeduplication"`
	ExportPolicyInfo          exportPolicyInfo       `structs:"exportPolicyInfo,omitempty"`
	ID                        string                 `structs:"uuid"`
	NewAggregate              bool                   `structs:"newAggregate"`
	CapacityTier              string                 `structs:"capacityTier,omitempty"`
	ProviderVolumeType        string                 `structs:"providerVolumeType"`
	TieringPolicy             string                 `structs:"tieringPolicy,omitempty"`
	NumOfDisks                float64                `structs:"maxNumOfDisksApprovedToAdd,omitempty"`
	AutoVsaCapacityManagement bool                   `structs:"autoVsaCapacityManagement"`
	DiskSize                  size                   `structs:"diskSize,omitempty"`
	Iops                      int                    `structs:"iops,omitempty"`
	WorkingEnvironmentType    string                 `structs:"workingEnvironmentType,omitempty"`
	ShareInfo                 shareInfoRequest       `structs:"shareInfo,omitempty"`
	ShareInfoUpdate           shareInfoUpdateRequest `structs:"shareInfo,omitempty"`
	IscsiInfo                 iscsiInfo              `structs:"iscsiInfo,omitempty"`
}

type volumeResponse struct {
	Name                   string              `json:"name"`
	SvmName                string              `json:"svmName"`
	AggregateName          string              `json:"aggregateName"`
	Size                   size                `json:"size"`
	SnapshotPolicyName     string              `json:"snapshotPolicy"`
	EnableThinProvisioning bool                `json:"thinProvisioning"`
	EnableCompression      bool                `json:"compression"`
	EnableDeduplication    bool                `json:"deduplication"`
	ExportPolicyInfo       exportPolicyInfo    `json:"exportPolicyInfo"`
	ID                     string              `json:"uuid"`
	CapacityTier           string              `json:"capacityTier,omitempty"`
	TieringPolicy          string              `json:"tieringPolicy,omitempty"`
	ProviderVolumeType     string              `json:"providerVolumeType"`
	Iops                   int                 `json:"iops"`
	ShareInfo              []shareInfoResponse `json:"shareInfo"`
}

type exportPolicyInfo struct {
	Name       string   `structs:"name,omitempty"`
	PolicyType string   `structs:"policyType,omitempty"`
	Ips        []string `structs:"ips,omitempty"`
	NfsVersion []string `structs:"nfsVersion,omitempty"`
}

type size struct {
	Size float64 `structs:"size"`
	Unit string  `structs:"unit"`
}

type quoteRequest struct {
	Size                   size             `structs:"size"`
	WorkingEnvironmentID   string           `structs:"workingEnvironmentId"`
	SvmName                string           `structs:"svmName"`
	AggregateName          string           `structs:"aggregateName,omitempty"`
	EnableThinProvisioning bool             `structs:"enableThinProvisioning"`
	EnableCompression      bool             `structs:"enableCompression"`
	EnableDeduplication    bool             `structs:"enableDeduplication"`
	ExportPolicyInfo       exportPolicyInfo `structs:"exportPolicyInfo,omitempty"`
	SnapshotPolicyName     string           `structs:"snapshotPolicyName"`
	Name                   string           `structs:"name"`
	CapacityTier           string           `structs:"capacityTier,omitempty"`
	ProviderVolumeType     string           `structs:"providerVolumeType"`
	TieringPolicy          string           `structs:"tieringPolicy,omitempty"`
	VerifyNameUniqueness   bool             `structs:"verifyNameUniqueness"`
	Iops                   int              `structs:"iops,omitempty"`
	WorkingEnvironmentType string           `structs:"workingEnvironmentType"`
}

type shareInfoRequest struct {
	ShareName     string        `structs:"shareName,omitempty"`
	AccessControl accessControl `structs:"accessControl,omitempty"`
}

type accessControl struct {
	Permission string   `structs:"permission,omitempty"`
	Users      []string `structs:"users,omitempty"`
}

// cifs volume response has different strcuture comparing to cifs create.
type shareInfoResponse struct {
	ShareName         string                      `json:"shareName"`
	AccessControlList []accessControlListResponse `json:"accessControlList"`
}

type accessControlListResponse struct {
	Permission string   `json:"permission"`
	Users      []string `json:"users"`
}

type shareInfoUpdateRequest struct {
	ShareName         string              `structs:"shareName,omitempty"`
	AccessControlList []accessControlList `structs:"accessControlList,omitempty"`
}

type accessControlList struct {
	Permission string   `structs:"permission,omitempty"`
	Users      []string `structs:"users,omitempty"`
}

type iscsiInfo struct {
	OsName                string `structs:"osName,omitempty"`
	IgroupCreationRequest struct {
		Initiators []string `structs:"initiators,omitempty"`
		IgroupName string   `structs:"igroupName,omitempty"`
	} `structs:"igroupCreationRequest,omitempty"`
	Igroups []string `structs:"igroups,omitempty"`
}

type initiator struct {
	AliasName              string `structs:"aliasName,omitempty"`
	Iqn                    string `structs:"iqn,omitempty"`
	WorkingEnvironmentId   string `structs:"workingEnvironmentId,omitempty"`
	SvmName                string `structs:"svmName,omitempty"`
	WorkingEnvironmentType string `structs:"workingEnvironmentType,omitempty"`
}

func (c *Client) createVolume(vol volumeRequest) error {
	baseURL, _, err := c.getAPIRoot(vol.WorkingEnvironmentID)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/volumes?createAggregateIfNotFound=true", baseURL)
	hostType := "CloudManagerHost"
	param := structs.Map(vol)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType)
	if err != nil {
		log.Print("createVolume request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createVolume")
	if responseError != nil {
		return responseError
	}
	err = c.waitOnCompletion(onCloudRequestID, "volume", "create", 40, 10)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteVolume(vol volumeRequest) error {
	baseURL, _, err := c.getAPIRoot(vol.WorkingEnvironmentID)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/volumes/%s/%s/%s", baseURL, vol.WorkingEnvironmentID, vol.SvmName, vol.Name)
	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType)
	responseError := apiResponseChecker(statusCode, response, "deleteVolume")
	if responseError != nil {
		return responseError
	}
	if err != nil {
		log.Print("deleteVolume request failed ", statusCode)
		return err
	}
	return nil
}

func (c *Client) getVolume(vol volumeRequest) ([]volumeResponse, error) {
	var result []volumeResponse
	hostType := "CloudManagerHost"
	baseURL, _, err := c.getAPIRoot(vol.WorkingEnvironmentID)
	if err != nil {
		return result, err
	}
	baseURL = fmt.Sprintf("%s/volumes?workingEnvironmentId=%s", baseURL, vol.WorkingEnvironmentID)

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getVolume request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getVolume")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolume ", err)
		return result, err
	}

	return result, nil
}

func (c *Client) getVolumeByID(request volumeRequest) (volumeResponse, error) {
	res, err := c.getVolume(request)
	if err != nil {
		return volumeResponse{}, err
	}
	for _, vol := range res {
		if vol.ID == request.ID {
			return vol, nil
		}
	}
	return volumeResponse{}, fmt.Errorf("Error fetching volume: volume doesn't exist")
}

func (c *Client) updateVolume(request volumeRequest) error {
	hostType := "CloudManagerHost"
	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/volumes/%s/%s/%s", baseURL, request.WorkingEnvironmentID, request.SvmName, request.Name)
	params := structs.Map(request)
	statusCode, response, _, err := c.CallAPIMethod("PUT", baseURL, params, c.Token, hostType)

	responseError := apiResponseChecker(statusCode, response, "updateVolume")
	if responseError != nil {
		return responseError
	}

	if err != nil {
		log.Print("updateVolume request failed ", statusCode)
		return err
	}
	return nil
}

func (c *Client) quoteVolume(request quoteRequest) (map[string]interface{}, error) {
	hostType := "CloudManagerHost"
	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID)
	if err != nil {
		return nil, err
	}
	baseURL = fmt.Sprintf("%s/volumes/quote", baseURL)
	params := structs.Map(request)

	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("quoteVolume request failed ", statusCode)
		return nil, err
	}
	responseError := apiResponseChecker(statusCode, response, "quoteVolume")
	if responseError != nil {
		return nil, responseError
	}
	var result map[string]interface{}
	json.Unmarshal(response, &result)
	return result, nil

}

func (c *Client) createInitiator(request initiator) error {
	hostType := "CloudManagerHost"
	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentId)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/volumes/initiator", baseURL)
	params := structs.Map(request)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("createInitiator request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createInitiator")
	if responseError != nil {
		return responseError
	}
	return nil
}

func (c *Client) getInitiator(request initiator) ([]initiator, error) {
	hostType := "CloudManagerHost"
	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentId)
	var result []initiator
	if err != nil {
		return result, err
	}
	baseURL = fmt.Sprintf("%s/volumes/initiator", baseURL)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("createInitiator request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "createInitiator")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getInitiator ", err)
		return result, err
	}
	return result, nil
}

type igroup struct {
	IgroupName             string   `json:"igroupName"`
	OsType                 string   `json:"osType"`
	PortsetName            string   `json:"portsetName"`
	IgroupType             string   `json:"igroupType"`
	Initiators             []string `json:"initiators"`
	WorkingEnvironmentId   string   `structs:"workingEnvironmentId"`
	SvmName                string   `structs:"svmName"`
	WorkingEnvironmentType string   `structs:"workingEnvironmentType,omitempty"`
}

func (c *Client) getIgroups(request igroup) ([]igroup, error) {
	hostType := "CloudManagerHost"
	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentId)
	var result []igroup
	if err != nil {
		return result, err
	}
	baseURL = fmt.Sprintf("%s/volumes/igroups/%s/%s", baseURL, request.WorkingEnvironmentId, request.SvmName)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getIgroups request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getIgroups")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getIgroups ", err)
		return result, err
	}
	return result, nil
}

func (c *Client) checkCifsExists(workingEnvironmentID string, svm string) (bool, error) {
	hostType := "CloudManagerHost"
	baseURL, _, err := c.getAPIRoot(workingEnvironmentID)
	var result []map[string]interface{}
	if err != nil {
		return false, err
	}
	baseURL = fmt.Sprintf("%s/working-environments/%s/cifs?svm=%s", baseURL, workingEnvironmentID, svm)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("chkeckCifsExists request failed ", statusCode)
		return false, err
	}
	responseError := apiResponseChecker(statusCode, response, "getIgroups")
	if responseError != nil {
		return false, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from chkeckCifsExists ", err)
		return false, err
	}
	if len(result) > 0 {
		return true, nil
	}
	return false, nil
}

func convertSizeUnit(size float64, from string, to string) float64 {
	if from == "GB" && to == "TB" {
		size = size / 1024
	}
	return size
}
