package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
)

// createCVOOnPremDetails the users input for creating a CVO
type createCVOOnPremDetails struct {
	Name            string `structs:"name"`
	WorkspaceID     string `structs:"tenantId"`
	ClusterAddress  string `structs:"clusterAddress"`
	ClusterUserName string `structs:"clusterUserName"`
	ClusterPassword string `structs:"clusterPassword"`
	Location        string `structs:"location"`
}

type cvoOnPremList struct {
	CVO []cvoOnPremResult `json:"onPremWorkingEnvironments"`
}

type cvoOnPremResult struct {
	PublicID string `json:"publicId"`
}

func (c *Client) createCVOOnPrem(cvoDetails createCVOOnPremDetails, clientID string) (cvoResult, error) {

	log.Print("createCVO")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createCVO request, failed to get AccessToken: ", err)
		return cvoResult{}, err
	}
	c.Token = accessTokenResult.Token

	if cvoDetails.WorkspaceID == "" {
		tenantID, err := c.getTenant(clientID)
		if err != nil {
			log.Print("getTenant request failed ", err)
			return cvoResult{}, err
		}
		log.Print("tenant result ", tenantID)
		cvoDetails.WorkspaceID = tenantID
	}

	baseURL := "/occm/api/onprem/working-environments"
	creationWaitTime := 60
	retries := 60
	hostType := "CloudManagerHost"
	params := structs.Map(cvoDetails)

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("createCVO request failed: %v, %v", statusCode, err)
		return cvoResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createCVO")
	if responseError != nil {
		return cvoResult{}, responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "create", retries, creationWaitTime, clientID)
	if err != nil {
		return cvoResult{}, err
	}

	var result cvoResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createCVO ", err)
		return cvoResult{}, err
	}

	return result, nil
}

func (c *Client) getCVOOnPremByID(id string, clientID string) (map[string]interface{}, error) {

	log.Print("getCVOOnPremByID")

	baseURL := fmt.Sprintf("/occm/api/onprem/working-environments/%s?fields=*", id)

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCVOOnPremByID request, failed to get AccessToken: ", err)
		return nil, err
	}
	c.Token = accessTokenResult.Token

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("getCVOOnPremByID request failed: %v, %v", statusCode, err)
		return nil, err
	}

	responseError := apiResponseChecker(statusCode, response, "getCVOOnPremByID")
	if responseError != nil {
		return nil, responseError
	}
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getCVOOnPremByID ", err)
		return nil, err
	}

	return result, nil
}

func (c *Client) getCVOOnPrem(id string, clientID string) (string, error) {

	log.Print("getCVOOnPrem")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCVOOnPrem request, failed to get AccessToken: ", err)
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := "/occm/api/working-environments"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("getCVOOnPrem request failed: %v, %v", statusCode, err)
		return "", err
	}

	log.Print(string(response))
	responseError := apiResponseChecker(statusCode, response, "getCVOOnPrem")
	if responseError != nil {
		return "", responseError
	}

	var result cvoOnPremList
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getCVOOnPrem ", err)
		return "", err
	}

	for _, cvoID := range result.CVO {
		if cvoID.PublicID == id {
			return cvoID.PublicID, nil
		}
	}

	return "", nil
}

func (c *Client) deleteCVOOnPrem(id string, clientID string) error {

	log.Print("deleteCVOOnPrem")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteCVOOnPrem request, failed to get AccessToken: ", err)
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/occm/api/onprem/working-environments/%s", id)

	hostType := "CloudManagerHost"
	creationWaitTime := 60
	retries := 40

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("deleteCVOOnPrem request failed: %v, %v", statusCode, err)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteCVOOnPrem")
	if responseError != nil {
		return responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "delete", retries, creationWaitTime, clientID)
	if err != nil {
		return err
	}

	return nil
}
