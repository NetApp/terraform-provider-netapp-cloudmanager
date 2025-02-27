package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

type cifsRequest struct {
	Domain               string   `structs:"activeDirectoryDomain"`
	Username             string   `structs:"activeDirectoryUsername"`
	Password             string   `structs:"activeDirectoryPassword"`
	DNSDomain            string   `structs:"dnsDomain"`
	IPAddresses          []string `structs:"ipAddresses"`
	NetBIOS              string   `structs:"netBIOS"`
	OrganizationalUnit   string   `structs:"organizationalUnit"`
	WorkingEnvironmentID string   `structs:"workingEnvironmentId"`
	SvmName              string   `structs:"svmName,omitempty"`
}

type cifsResponse struct {
	Domain             string   `json:"activeDirectoryDomain"`
	DNSDomain          string   `json:"dnsDomain"`
	IPAddresses        []string `json:"ipAddresses"`
	NetBIOS            string   `json:"netBIOS"`
	OrganizationalUnit string   `json:"organizationalUnit"`
	AuthenticationType string   `json:"authenticationType"`
}

type cifsDeleteRequest struct {
	Username string `structs:"activeDirectoryUsername,omitempty"`
	Password string `structs:"activeDirectoryPassword,omitempty"`
	SvmName  string `structs:"svmName,omitempty"`
}

func (c *Client) createCIFS(cifs cifsRequest, clientID string) error {
	baseURL, _, err := c.getAPIRoot(cifs.WorkingEnvironmentID, clientID, true, "")
	hostType := "CloudManagerHost"
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/working-environments/%s/cifs", baseURL, cifs.WorkingEnvironmentID)
	param := structs.Map(cifs)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createCIFS request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createCifs")
	if responseError != nil {
		return responseError
	}
	err = c.waitOnCompletion(onCloudRequestID, "cifs", "create", 10, 10, clientID)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getCIFS(cifs cifsRequest, clientID string) ([]cifsResponse, error) {
	var result []cifsResponse
	baseURL, _, err := c.getAPIRoot(cifs.WorkingEnvironmentID, clientID, true, "")
	if err != nil {
		return result, err
	}

	baseURL = fmt.Sprintf("%s/working-environments/%s/cifs?svm=%s", baseURL, cifs.WorkingEnvironmentID, cifs.SvmName)
	hostType := "CloudManagerHost"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getCIFS request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getCIFS")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getCIFS ", err)
		return result, err
	}
	return result, nil
}

func (c *Client) deleteCIFS(cifs cifsDeleteRequest, workingEnvironmentID, clientID string) error {
	baseURL, _, err := c.getAPIRoot(workingEnvironmentID, clientID, true, "")
	hostType := "CloudManagerHost"
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/working-environments/%s/delete-cifs", baseURL, workingEnvironmentID)
	param := structs.Map(cifs)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType, clientID)
	if err != nil {
		log.Print("deleteCIFS request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "deleteCIFS")
	if responseError != nil {
		return responseError
	}

	return nil
}
