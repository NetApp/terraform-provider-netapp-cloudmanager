package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
)

func (c *Client) getCustomData(registerAgentTOService registerAgentTOServiceRequest) (string, error) {
	accesTokenResult, err := c.getAccessToken()
	if err != nil {
		return "", err
	}
	log.Print("getAccessToken  ", accesTokenResult.Token)
	c.Token = accesTokenResult.Token

	if c.AccountID == "" {
		accountID, err := c.getAccount()
		if err != nil {
			return "", err
		}
		log.Print("getAccount ", accountID)
		registerAgentTOService.AccountID = accountID
	} else {
		registerAgentTOService.AccountID = c.AccountID
	}

	userDataRespone, err := c.registerAgentTOServiceForAzure(registerAgentTOService)
	if err != nil {
		return "", err
	}

	c.ClientID = userDataRespone.ClientID
	c.AccountID = userDataRespone.AccountID

	userData := "{\n\t\"instanceName\": \"" + userDataRespone.Name + "\",\n\t\"company\": \"" + userDataRespone.Company + "\",\n\t\"clientId\": \"" + userDataRespone.ClientID + "\",\n\t\"clientSecret\": \"" + userDataRespone.ClientSecret + "\",\n\t\"systemId\": \"" + userDataRespone.UUID + "\",\n\t\"tenancyAccountId\": \"" + userDataRespone.AccountID + "\",\n\t\"proxySettings\": {\n\t\"proxyPassword\": \"" + userDataRespone.ProxySettings.ProxyPassword + "\",\n\t\"proxyUserName\": \"" + userDataRespone.ProxySettings.ProxyUserName + "\",\n\t\"proxyUrl\": \"" + userDataRespone.ProxySettings.ProxyURL + "\"\n}\n}"
	log.Print("userData ", userData)

	return userData, nil
}

func (c *Client) registerAgentTOServiceForAzure(registerAgentTOServiceRequest registerAgentTOServiceRequest) (createUserData, error) {

	baseURL := "/agents-mgmt/connector-setup"
	hostType := "CloudManagerHost"

	registerAgentTOServiceRequest.Placement.Provider = "AZURE"

	params := structs.Map(registerAgentTOServiceRequest)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("registerAgentTOService request failed ", statusCode)
		return createUserData{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "registerAgentTOService")
	if responseError != nil {
		return createUserData{}, responseError
	}

	var result createUserData
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from registerAgentTOService ", err)
		return createUserData{}, err
	}

	return result, nil
}

func (c *Client) deployAzureVM(occmDetails createOCCMDetails) (string, error) {

	var registerAgentTOService registerAgentTOServiceRequest
	registerAgentTOService.Name = occmDetails.Name
	registerAgentTOService.Placement.Region = occmDetails.Location
	registerAgentTOService.Company = occmDetails.Company
	if occmDetails.ProxyURL != "" {
		registerAgentTOService.Extra.Proxy.ProxyURL = occmDetails.ProxyURL
	}

	if occmDetails.ProxyUserName != "" {
		registerAgentTOService.Extra.Proxy.ProxyUserName = occmDetails.ProxyUserName
	}

	if occmDetails.ProxyPassword != "" {
		registerAgentTOService.Extra.Proxy.ProxyPassword = occmDetails.ProxyPassword
	}

	if occmDetails.VnetResourceGroup != "" {
		registerAgentTOService.Placement.Network = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s", occmDetails.SubscriptionID, occmDetails.VnetResourceGroup, occmDetails.VnetID)
	} else {
		registerAgentTOService.Placement.Network = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s", occmDetails.SubscriptionID, occmDetails.ResourceGroup, occmDetails.VnetID)
	}

	registerAgentTOService.Placement.Subnet = fmt.Sprintf("%s/subnets/%s", registerAgentTOService.Placement.Network, occmDetails.SubnetID)

	userData, err := c.getCustomData(registerAgentTOService)
	if err != nil {
		return "", err
	}

	log.Print(userData)

	c.UserData = userData

	instanceID, err := c.CallDeployAzureVM(occmDetails)
	if err != nil {
		return "", err
	}

	log.Print("Sleep for 2 minutes")
	time.Sleep(time.Duration(120) * time.Second)

	retries := 16
	for {
		occmResp, err := c.checkOCCMStatus()
		if err != nil {
			return "", err
		}
		if occmResp.Status == "active" {
			break
		} else {
			if retries == 0 {
				log.Print("Taking too long for status to be active")
				return "", fmt.Errorf("Taking too long for OCCM agent to be active or not properly setup")
			}
			time.Sleep(time.Duration(30) * time.Second)
			retries--
		}
	}

	return instanceID, nil
}

func (c *Client) getdeployAzureVM(occmDetails createOCCMDetails, id string) (string, error) {

	log.Print("getdeployAzureVM")

	res, err := c.CallGetAzureVM(occmDetails)
	if err != nil {
		return "", err
	}

	if res == id {
		return res, nil
	}

	return "", nil
}

func (c *Client) createOCCMAzure(occmDetails createOCCMDetails) (OCCMMResult, error) {

	log.Print("createOCCMAzure")
	_, err := c.deployAzureVM(occmDetails)
	if err != nil {
		return OCCMMResult{}, err
	}

	var result OCCMMResult
	result.ClientID = c.ClientID
	result.AccountID = c.AccountID

	return result, nil
}

func (c *Client) deleteOCCMAzure(request deleteOCCMDetails) error {

	err := c.CallDeleteAzureVM(request)
	if err != nil {
		return err
	}

	log.Print("Sleep for 30 seconds")
	time.Sleep(time.Duration(30) * time.Second)

	accesTokenResult, err := c.getAccessToken()
	c.Token = accesTokenResult.Token
	if err != nil {
		return err
	}

	retries := 30
	for {
		occmResp, err := c.checkOCCMStatus()
		if err != nil {
			return err
		}
		if occmResp.Status != "active" {
			break
		} else {
			if retries == 0 {
				log.Print("Taking too long for instance to finish terminating")
				return fmt.Errorf("Taking too long for instance to finish terminating")
			}
			time.Sleep(time.Duration(10) * time.Second)
			retries--
		}
	}

	if err := c.callOCCMDelete(); err != nil {
		return err
	}

	return nil
}
