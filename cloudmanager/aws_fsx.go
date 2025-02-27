package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

// createAWSFSXDetails the users input for creating a FSX
type createAWSFSXDetails struct {
	Name                   string          `structs:"name"`
	AWSCredentials         string          `structs:"credentialsId"`
	WorkspaceID            string          `structs:"workspaceId"`
	Region                 string          `structs:"region"`
	PrimarySubnetID        string          `structs:"primarySubnetId"`
	SecondarySubnetID      string          `structs:"secondarySubnetId"`
	FSXAdminPassword       string          `structs:"fsxAdminPassword"`
	KmsKeyID               string          `structs:"kmsKeyId,omitempty"`
	MinimumSsdIops         int             `structs:"minimumSsdIops,omitempty"`
	EndpointIPAddressRange string          `structs:"endpointIpAddressRange,omitempty"`
	StorageCapacity        storageCapacity `structs:"storageCapacity"`
	RouteTableIds          []string        `structs:"routeTableIds,omitempty"`
	ThroughputCapacity     int             `structs:"throughputCapacity,omitempty"`
	SecurityGroupIds       []string        `structs:"securityGroupIds,omitempty"`
	AwsFSXTags             []fsxTags       `structs:"tags,omitempty"`
	TenantID               string
}

// fsxTags the input for requesting a FSX AWS
type fsxTags struct {
	TagKey   string `structs:"key"`
	TagValue string `structs:"value,omitempty"`
}

// storageCapacity the input for requesting a FSX AWS
type storageCapacity struct {
	Size int    `structs:"size"`
	Unit string `structs:"unit"`
}

// deleteAWSFSXDetails the users input for deleting a FSX AWS
type deleteAWSFSXDetails struct {
	InstanceID string
	Region     string
}

// fsxResult the users input for creating a FSX AWS
type fsxResult struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Region          string          `json:"region"`
	ProviderDetails providerDetails `json:"providerDetails"`
}

// fsxSVMResult the result for FSX SVM
type fsxSVMResult struct {
	Name string `json:"name"`
}

// fsxStatusResult for creating a fsx
type fsxStatusResult struct {
	ProviderDetails providerDetails `json:"providerDetails"`
	Error           string          `json:"error"`
}

// providerDetails for creating a fsx
type providerDetails struct {
	Status status `json:"status"`
}

// status for creating a fsx
type status struct {
	Status    string `json:"status"`
	Lifecycle string `json:"lifecycle"`
}

// recoverAWSFSXDetails the users input for importing a FSX
type recoverAWSFSXDetails struct {
	Name           string `structs:"name"`
	AWSCredentials string `structs:"credentialsId"`
	WorkspaceID    string `structs:"workspaceId"`
	Region         string `structs:"region"`
	FileSystemID   string `structs:"fileSystemId"`
}

// check if name tag exists
func hasNameTag(tags []fsxTags) bool {
	for _, v := range tags {
		if v.TagKey == "name" {
			return true
		}
	}
	return false
}

// expandUserTags converts set to fsxTag struct
func expandfsxTags(set *schema.Set) []fsxTags {
	tags := []fsxTags{}

	for _, v := range set.List() {
		tag := v.(map[string]interface{})
		fsxTag := fsxTags{}
		fsxTag.TagKey = tag["tag_key"].(string)
		fsxTag.TagValue = tag["tag_value"].(string)
		tags = append(tags, fsxTag)
	}
	return tags
}

func (c *Client) getAWSCredentialsID(name string, tenantID string) (string, error) {

	log.Print("getAWSCredentialsID ", tenantID)

	baseURL := fmt.Sprintf("/fsx-ontap/aws-credentials/%s", tenantID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, "")
	if err != nil {
		log.Print("getAWSCredentialsID request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getAWSCredentialsID")
	if responseError != nil {
		return "", responseError
	}

	var result []fsxResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getAWSCredentialsID ", err)
		return "", err
	}

	for _, fsxID := range result {
		if fsxID.Name == name {
			return fsxID.ID, nil
		}
	}

	return "", fmt.Errorf("aws_credentials_name not found")
}

func (c *Client) getAWSFSX(id string, tenantID string, isSaas bool, connectorIP string) (string, error) {

	log.Print("getAWSFSX")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getAWSFSX request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s", tenantID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, "")
	if err != nil {
		log.Print("getAWSFSX request failed ", statusCode, err)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getAWSFSX")
	if responseError != nil {
		return "", responseError
	}

	var result []fsxResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getAWSFSX ", err)
		return "", err
	}

	for _, fsxID := range result {
		if fsxID.ID == id {
			return fsxID.ID, nil
		}
	}

	return "", nil
}

func (c *Client) getAWSFSXByID(id string, tenantID string) (fsxResult, error) {

	log.Print("getAWSFSXByID")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getAWSFSXByID request, failed to get AccessToken")
		return fsxResult{}, err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s/%s?provider-details=true", tenantID, id)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, "")
	if err != nil {
		log.Print("getAWSFSXByID request failed ", statusCode, err)
		return fsxResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "getAWSFSXByID")
	if responseError != nil {
		return fsxResult{}, responseError
	}

	var result fsxResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getAWSFSXByID ", err)
		return fsxResult{}, err
	}

	return result, nil
}

func (c *Client) importAWSFSX(fsxDetails createAWSFSXDetails, fileSystemID string) (string, error) {

	log.Print("importAWSFSX")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in importAWSFSX request, failed to get AccessToken")
		return "", fmt.Errorf("in importAWSFSX request, failed to get AccessToken: %s", err)
	}
	c.Token = accessTokenResult.Token

	fsxDetails.AWSCredentials, err = c.getAWSCredentialsID(fsxDetails.AWSCredentials, fsxDetails.TenantID)
	if err != nil {
		log.Print("importAWSFSX request failed ", err)
		return "", fmt.Errorf("importAWSFSX request failed: %s", err)
	}

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s/discover?credentials-id=%s&workspace-id=%s&region=%s",
		fsxDetails.TenantID, fsxDetails.AWSCredentials, fsxDetails.WorkspaceID, fsxDetails.Region)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, "")
	if err != nil {
		log.Print("importAWSFSX request failed ", statusCode)
		return "", fmt.Errorf("importAWSFSX request failed: %s", err)
	}

	responseError := apiResponseChecker(statusCode, response, "importAWSFSX")
	if responseError != nil {
		return "", responseError
	}

	var result []fsxResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from importAWSFSX ", err)
		return "", err
	}

	idFound := false
	for _, fsxID := range result {
		log.Print(string(fsxID.ID))
		if fsxID.ID == fileSystemID {
			idFound = true
			break
		}
	}

	if idFound == false {
		return "", fmt.Errorf("file_system_id provided could not be found")
	}

	recoverAWSFSXDetails := recoverAWSFSXDetails{}

	recoverAWSFSXDetails.AWSCredentials = fsxDetails.AWSCredentials
	recoverAWSFSXDetails.Name = fsxDetails.Name
	recoverAWSFSXDetails.Region = fsxDetails.Region
	recoverAWSFSXDetails.WorkspaceID = fsxDetails.WorkspaceID
	recoverAWSFSXDetails.FileSystemID = fileSystemID

	baseURL = fmt.Sprintf("/fsx-ontap/working-environments/%s/recover", fsxDetails.TenantID)

	params := structs.Map(recoverAWSFSXDetails)

	statusCode, response, _, err = c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, "")
	if err != nil {
		log.Print("importAWSFSX request failed ", statusCode)
		return "", fmt.Errorf("importAWSFSX request failed: %s", err)
	}

	responseError = apiResponseChecker(statusCode, response, "importAWSFSX")
	if responseError != nil {
		return "", responseError
	}

	return fileSystemID, nil
}

func (c *Client) createAWSFSX(fsxDetails createAWSFSXDetails) (fsxResult, error) {

	log.Print("createFSX")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createFSX request, failed to get AccessToken")
		return fsxResult{}, err
	}
	c.Token = accessTokenResult.Token

	fsxDetails.AWSCredentials, err = c.getAWSCredentialsID(fsxDetails.AWSCredentials, fsxDetails.TenantID)
	if err != nil {
		log.Print("createFSX request failed ", err)
		return fsxResult{}, err
	}

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s", fsxDetails.TenantID)

	creationWaitTime := 60
	creationRetryCount := 60
	hostType := "CloudManagerHost"
	params := structs.Map(fsxDetails)

	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, "")
	if err != nil {
		log.Print("createFSX request failed ", statusCode)
		return fsxResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createFSX")
	if responseError != nil {
		return fsxResult{}, responseError
	}

	var result fsxResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createFSX ", err)
		return fsxResult{}, err
	}

	err = c.waitOnCompletionFSX(result.ID, fsxDetails.TenantID, "FSX", "create", creationRetryCount, creationWaitTime)
	if err != nil {
		return fsxResult{}, err
	}

	return result, nil
}

func (c *Client) checkTaskStatusFSX(id string, tenantID string) (providerDetails, string, error) {

	log.Printf("checkTaskStatusFSX: %s", tenantID)

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s/%s?provider-details=true", tenantID, id)

	hostType := "CloudManagerHost"

	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, "")
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(2 * time.Second)
				networkRetries--
			} else {
				log.Print("checkTaskStatusFSX request failed ", code)
				return providerDetails{}, "", err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}

	responseError := apiResponseChecker(statusCode, response, "checkTaskStatusFSX")
	if responseError != nil {
		return providerDetails{}, "", responseError
	}

	var result fsxStatusResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkTaskStatusFSX ", err)
		return providerDetails{}, "", err
	}

	return result.ProviderDetails, result.Error, nil
}

func (c *Client) waitOnCompletionFSX(id string, tenantID string, actionName string, task string, retries int, waitInterval int) error {
	// below timeout is for handling for a situation in which the status does not exist yet
	time.Sleep(5 * time.Second)
	for {
		fsxStatus, failureErrorMessage, err := c.checkTaskStatusFSX(id, tenantID)
		if err != nil {
			return err
		}
		if fsxStatus.Status.Status == "ON" && fsxStatus.Status.Lifecycle != "CREATING" {
			return nil
		} else if fsxStatus.Status.Status == "FAILED" {
			return fmt.Errorf("Failed to %s %s, error: %s", task, actionName, failureErrorMessage)
		} else if retries == 0 {
			log.Print("Taking too long to ", task, actionName)
			return fmt.Errorf("Taking too long for %s to %s or not properly setup", actionName, task)
		}
		log.Printf("Sleep for %d seconds", waitInterval)
		time.Sleep(time.Duration(waitInterval) * time.Second)
		retries--
	}
}

func (c *Client) deleteAWSFSX(id string, tenantID string) error {

	log.Print("deleteAWSFSX")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteAWSFSX request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s/%s", tenantID, id)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, "")
	if err != nil {
		log.Print("deleteAWSFSX request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteAWSFSX")
	if responseError != nil {
		return responseError
	}

	return nil
}
