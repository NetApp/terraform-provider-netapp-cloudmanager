package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
)

func (c *Client) getCustomData(registerAgentTOService registerAgentTOServiceRequest, proxyCertificates []string, clientID string) (string, string, error) {
	accesTokenResult, err := c.getAccessToken()
	if err != nil {
		return "", "", err
	}
	log.Print("getAccessToken  ", accesTokenResult.Token)
	c.Token = accesTokenResult.Token

	if c.AccountID == "" {
		accountID, err := c.getAccount(clientID)
		if err != nil {
			return "", "", err
		}
		log.Print("getAccount ", accountID)
		registerAgentTOService.AccountID = accountID
	} else {
		registerAgentTOService.AccountID = c.AccountID
	}

	userDataRespone, err := c.registerAgentTOServiceForAzure(registerAgentTOService, clientID)
	if err != nil {
		return "", "", err
	}

	newClientID := userDataRespone.ClientID
	c.AccountID = userDataRespone.AccountID

	userDataRespone.ProxySettings.ProxyCertificates = proxyCertificates
	rawUserData, _ := json.MarshalIndent(userDataRespone, "", "\t")
	userData := string(rawUserData)
	log.Printf("getCustomData: userData %#v %s", userData, newClientID)

	return userData, newClientID, nil
}

func (c *Client) registerAgentTOServiceForAzure(registerAgentTOServiceRequest registerAgentTOServiceRequest, clientID string) (createUserData, error) {

	baseURL := "/agents-mgmt/connector-setup"
	hostType := "CloudManagerHost"

	registerAgentTOServiceRequest.Placement.Provider = "AZURE"

	params := structs.Map(registerAgentTOServiceRequest)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
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
	result.IgnoreUpgrade = true
	return result, nil
}

func (c *Client) createOCCMAzure(occmDetails createOCCMDetails, proxyCertificates []string, clientID string) (OCCMMResult, error) {
	log.Print("createOCCMAzure")
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

	if !strings.Contains(occmDetails.VnetID, "/") {
		log.Print("Compose vnetID...")
		resourceGroup := occmDetails.ResourceGroup
		if occmDetails.VnetResourceGroup != "" {
			resourceGroup = occmDetails.VnetResourceGroup
		}
		registerAgentTOService.Placement.Network = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s", occmDetails.SubscriptionID, resourceGroup, occmDetails.VnetID)
	} else {
		registerAgentTOService.Placement.Network = occmDetails.VnetID
	}
	log.Print("createOCCMAzure VnetID: ")
	log.Print(registerAgentTOService.Placement.Network)
	if !strings.Contains(occmDetails.SubnetID, "/") {
		log.Print("Compose subnetID...")
		registerAgentTOService.Placement.Subnet = fmt.Sprintf("%s/subnets/%s", registerAgentTOService.Placement.Network, occmDetails.SubnetID)
	} else {
		registerAgentTOService.Placement.Subnet = occmDetails.SubnetID
	}
	log.Print("createOCCMAzure SubnetID: ")
	log.Print(registerAgentTOService.Placement.Subnet)
	userData, newClientID, err := c.getCustomData(registerAgentTOService, proxyCertificates, clientID)
	if err != nil {
		return OCCMMResult{}, err
	}

	log.Print(userData)
	log.Printf("deployAzureVM %s client_id %s", occmDetails.Name, newClientID)

	c.UserData = userData
	var result OCCMMResult
	result.ClientID = newClientID
	result.AccountID = c.AccountID

	principalID, err := c.CallDeployAzureVM(occmDetails)
	if err != nil {
		return OCCMMResult{}, err
	}

	result.PrincipalID = principalID
	log.Print("Sleep for 2 minutes")
	time.Sleep(time.Duration(120) * time.Second)

	retries := 26
	for {
		occmResp, err := c.checkOCCMStatus(newClientID)
		if err != nil {
			return OCCMMResult{}, err
		}
		if occmResp.Status == "active" {
			break
		} else {
			if retries == 0 {
				log.Print("Taking too long for status to be active")
				return OCCMMResult{}, fmt.Errorf("Taking too long for OCCM agent to be active or not properly setup")
			}
			time.Sleep(time.Duration(30) * time.Second)
			retries--
		}
	}

	return result, nil
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

func (c *Client) deleteOCCMAzure(request deleteOCCMDetails, clientID string) error {

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
		occmResp, err := c.checkOCCMStatus(clientID)
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

	if err := c.callOCCMDelete(clientID); err != nil {
		return err
	}

	return nil
}
