package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

type apiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// only list what is needed
type workingEnvironmentInfo struct {
	Name                   string `json:"name"`
	PublicID               string `json:"publicId"`
	CloudProviderName      string `json:"cloudProviderName"`
	IsHA                   bool   `json:"isHA"`
	WorkingEnvironmentType string `json:"workingEnvironmentType"`
	SvmName                string `json:"svmName"`
}

type workingEnvironmentResult struct {
	VsaWorkingEnvironment       []workingEnvironmentInfo `json:"vsaWorkingEnvironments"`
	OnPremWorkingEnvironments   []workingEnvironmentInfo `json:"onPremWorkingEnvironments"`
	AzureVsaWorkingEnvironments []workingEnvironmentInfo `json:"azureVsaWorkingEnvironments"`
	GcpVsaWorkingEnvironments   []workingEnvironmentInfo `json:"gcpVsaWorkingEnvironments"`
}

// Check HTTP response code, return error if HTTP request is not successed.
func apiResponseChecker(statusCode int, response []byte, funcName string) error {

	if statusCode >= 300 || statusCode < 200 {
		log.Printf("%s request failed", funcName)
		return fmt.Errorf("code: %d, message: %s", statusCode, string(response))
	}

	return nil

}

func (c *Client) checkTaskStatus(id string) (int, error) {

	log.Printf("checkTaskStatus: %s", id)

	baseURL := fmt.Sprintf("/occm/api/audit/activeTask/%s", id)

	hostType := "CloudManagerHost"

	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(1 * time.Second)
				networkRetries--
			} else {
				log.Print("checkTaskStatus request failed ", code)
				return 0, err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}

	responseError := apiResponseChecker(statusCode, response, "checkTaskStatus")
	if responseError != nil {
		return 0, responseError
	}

	var result cvoStatusResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkTaskStatus ", err)
		return 0, err
	}

	return result.Status, nil
}

func (c *Client) waitOnCompletion(id string, actionName string, task string, retries int, waitInterval int) error {
	for {
		cvoStatus, err := c.checkTaskStatus(id)
		if err != nil {
			return err
		}
		if cvoStatus == 1 {
			return nil
		} else if cvoStatus == -1 {
			return fmt.Errorf("Failed to %s %s", task, actionName)
		} else if cvoStatus == 0 {
			if retries == 0 {
				log.Print("Taking too long to ", task, actionName)
				return fmt.Errorf("Taking too long for %s to %s or not properly setup", actionName, task)
			}
			log.Printf("Sleep for %d seconds", waitInterval)
			time.Sleep(time.Duration(waitInterval) * time.Second)
			retries--
		}

	}
}

// get working environment information by working environment id
func (c *Client) getWorkingEnvironmentInfo(id string) (workingEnvironmentInfo, error) {
	baseURL := fmt.Sprintf("/occm/api/working-environments/%s", id)
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("getWorkingEnvironmentInfo %s request failed (%d)", id, statusCode)
		return workingEnvironmentInfo{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getWorkingEnvironmentInfo")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var result workingEnvironmentInfo
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}

	return result, nil
}

func findWE(name string, weList []workingEnvironmentInfo) (workingEnvironmentInfo, error) {

	for i := range weList {
		if weList[i].Name == name {
			log.Printf("Found working environment: %v", weList[i])
			return weList[i], nil
		}
	}
	return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment %s in the list", name)
}

func (c *Client) findWorkingEnvironmentByName(name string) (workingEnvironmentInfo, error) {

	// check working environment exists or not
	baseURL := fmt.Sprintf("/occm/api/working-environments/exists/%s", name)
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("findWorkingEnvironmentByName request failed. (check exists) ", statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "findWorkingEnvironmentByName")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	// get working environment information
	baseURL = "/occm/api/working-environments"
	statusCode, response, _, err = c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("findWorkingEnvironmentByName %s request failed (%d)", name, statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError = apiResponseChecker(statusCode, response, "findWorkingEnvironmentByName")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var workingEnvironments workingEnvironmentResult
	if err := json.Unmarshal(response, &workingEnvironments); err != nil {
		log.Print("Failed to unmarshall response from findWorkingEnvironmentByName")
		return workingEnvironmentInfo{}, err
	}

	var workingEnvironment workingEnvironmentInfo
	workingEnvironment, err = findWE(name, workingEnvironments.VsaWorkingEnvironment)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.OnPremWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.AzureVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.GcpVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}

	log.Printf("Cannot find the working environment %s", name)

	return workingEnvironmentInfo{}, err
}

func (c *Client) getAPIRoot(workingEnvironmentID string) (string, string, error) {

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", "", err
		}
		c.Token = accesTokenResult.Token
	}
	workingEnvDetail, err := c.getWorkingEnvironmentInfo(workingEnvironmentID)
	if err != nil {
		log.Print("Cannot get working environment information.")
		return "", "", err
	}
	log.Printf("Working environment %v", workingEnvDetail)

	var baseURL string
	if workingEnvDetail.CloudProviderName != "Amazon" {
		if workingEnvDetail.IsHA {
			baseURL = fmt.Sprintf("/occm/api/%s/ha", strings.ToLower(workingEnvDetail.CloudProviderName))
		} else {
			baseURL = fmt.Sprintf("/occm/api/%s/vsa", strings.ToLower(workingEnvDetail.CloudProviderName))
		}
	} else {
		if workingEnvDetail.IsHA {
			baseURL = "/occm/api/aws/ha"
		} else {
			baseURL = "/occm/api/vsa"
		}
	}
	log.Printf("API root = %s", baseURL)
	return baseURL, workingEnvDetail.CloudProviderName, nil
}
