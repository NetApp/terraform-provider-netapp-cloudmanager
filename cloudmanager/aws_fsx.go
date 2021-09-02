package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
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
	AwsFSXTags             []userTags      `structs:"tags,omitempty"`
	TenantID               string
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
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) getAWSCredentialsID(name string, tenantID string) (string, error) {

	log.Print("getAWSCredentialsID ", tenantID)

	baseURL := fmt.Sprintf("/fsx-ontap/aws-credentials/%s", tenantID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
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

func (c *Client) getAWSFSX(id string, tenantID string) (string, error) {

	log.Print("getAWSFSX")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getAWSFSX request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	if tenantID == "" {
		id, err := c.getTenant()
		if err != nil {
			log.Print("getTenant request failed ", err)
			return "", err
		}
		log.Print("tenant result ", id)
		tenantID = id
	}

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s", tenantID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
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

func (c *Client) createAWSFSX(fsxDetails createAWSFSXDetails) (fsxResult, error) {

	log.Print("createFSX")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createFSX request, failed to get AccessToken")
		return fsxResult{}, err
	}
	c.Token = accessTokenResult.Token

	if fsxDetails.TenantID == "" {
		tenantID, err := c.getTenant()
		if err != nil {
			log.Print("getTenant request failed ", err)
			return fsxResult{}, err
		}
		log.Print("tenant result ", tenantID)
		fsxDetails.TenantID = tenantID
	}

	fsxDetails.AWSCredentials, err = c.getAWSCredentialsID(fsxDetails.AWSCredentials, fsxDetails.TenantID)
	if err != nil {
		log.Print("createFSX request failed ", err)
		return fsxResult{}, err
	}

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s", fsxDetails.TenantID)

	hostType := "CloudManagerHost"
	params := structs.Map(fsxDetails)

	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
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

	return result, nil
}

func (c *Client) deleteAWSFSX(id string, tenantID string) error {

	log.Print("deleteAWSFSX")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteAWSFSX request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	if tenantID == "" {
		id, err := c.getTenant()
		if err != nil {
			log.Print("getTenant request failed ", err)
			return err
		}
		log.Print("tenant result ", id)
		tenantID = id
	}

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s/%s", tenantID, id)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType)
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
