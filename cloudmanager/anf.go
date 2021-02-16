package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/structs"
)

type anfVolumeRequest struct {
	Size                   float64  `structs:"quotaInBytes"`
	Name                   string   `structs:"name"`
	VolumePath             string   `structs:"volumePath"`
	ProtocolTypes          []string `structs:"protocolTypes"`
	ServiceLevel           string   `structs:"serviceLevel"`
	SubnetName             string   `structs:"subnetName"`
	VirtualNetworkName     string   `structs:"virtualNetworkName"`
	Location               string   `structs:"location"`
	Rules                  []rule   `structs:"rules"`
	WorkingEnvironmentName string   `structs:"workingEnvironmentName"`
}

// VolumePath is returned as creationToken
type anfVolumeResponse struct {
	Size                   float64                   `json:"quotaInBytes"`
	Name                   string                    `json:"name"`
	VolumePath             string                    `json:"creationToken"`
	ProtocolTypes          []string                  `json:"protocolTypes"`
	ServiceLevel           string                    `json:"serviceLevel"`
	SubnetName             string                    `json:"subnet"`
	Location               string                    `json:"location"`
	Rules                  map[string][]ruleResponse `json:"exportPolicy"`
	WorkingEnvironmentName string                    `json:"workingEnvironmentName"`
}

type ruleResponse struct {
	AllowedClients string `json:"allowedClients"`
	Cifs           bool   `json:"cifs"`
	Nfsv3          bool   `json:"nfsv3"`
	Nfsv41         bool   `json:"nfsv41"`
	RuleIndex      int    `json:"ruleIndex"`
	UnixReadOnly   bool   `json:"unixReadOnly"`
	UnixReadWrite  bool   `json:"unixReadWrite"`
}

type rule struct {
	AllowedClients string `structs:"allowedClients"`
	Cifs           bool   `structs:"cifs"`
	Nfsv3          bool   `structs:"nfsv3"`
	Nfsv41         bool   `structs:"nfsv41"`
	RuleIndex      int    `structs:"ruleIndex"`
	UnixReadOnly   bool   `structs:"unixReadOnly"`
	UnixReadWrite  bool   `structs:"unixReadWrite"`
}

// AWS,ANF and GCP share this struct.
type cvsInfo struct {
	AccountName        string `structs:"accountName"`
	AccountID          string `structs:"accountID"`
	CredentialsID      string `structs:"credentialsID"`
	SubscriptionName   string `structs:"subscriptionName"`
	NetAppAccountName  string `structs:"netAppAccountName"`
	ResourceGroupsName string `structs:"resourceGroupsName"`
	CapacityPools      string `structs:"capacityPools"`
	VirtualNetworkName string `structs:"virtualNetworkName"`
	SubnetName         string `structs:"subnetName"`
}

func (c *Client) getAccountByName(name string) (string, error) {
	log.Print("getAccount")

	baseURL := "/tenancy/account"
	hostType := "CloudManagerHost"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getAccount request failed ", statusCode)
		return "", err
	}
	responseError := apiResponseChecker(statusCode, response, "getAccount")
	if responseError != nil {
		return "", responseError
	}

	var results []accountIDResult
	if err := json.Unmarshal(response, &results); err != nil {
		log.Print("Failed to unmarshall response from getAccount ", err)
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no account exists")
	}
	if name != "" {
		for _, result := range results {
			if name == result.AccountName {
				return result.AccountID, nil
			}
		}
		return "", fmt.Errorf("account: %s not found", name)
	}

	// return the first account if name is not provided.
	return results[0].AccountID, nil
}

func (c *Client) createANFVolume(vol anfVolumeRequest, info cvsInfo) error {
	baseURL, err := c.getCVSAPIRoot(info.AccountName, vol.WorkingEnvironmentName)
	if err != nil {
		return err
	}
	subscription, err := c.getSubscription(baseURL, info.SubscriptionName)
	if err != nil {
		return err
	}
	subnet, err := c.getSubnetID(fmt.Sprintf("%s/subscriptions/%s", baseURL, subscription), info.VirtualNetworkName, info.SubnetName, vol.Location)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/netAppAccounts/%s/capacityPools/%s/volumes", baseURL, subscription, info.ResourceGroupsName, info.NetAppAccountName, info.CapacityPools)
	hostType := "CVSHost"
	param := structs.Map(vol)
	param["subnetId"] = subnet
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType)
	if err != nil {
		log.Print("createANFVolume request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createANFVolume")
	if responseError != nil {
		return responseError
	}

	return nil
}

func (c *Client) getANFVolume(vol anfVolumeRequest, info cvsInfo) (anfVolumeResponse, error) {
	baseURL, err := c.getCVSAPIRoot(info.AccountName, vol.WorkingEnvironmentName)
	if err != nil {
		return anfVolumeResponse{}, err
	}
	subscription, err := c.getSubscription(baseURL, info.SubscriptionName)
	if err != nil {
		return anfVolumeResponse{}, err
	}
	subnet, err := c.getSubnetID(fmt.Sprintf("%s/subscriptions/%s", baseURL, subscription), info.VirtualNetworkName, info.SubnetName, vol.Location)
	if err != nil {
		return anfVolumeResponse{}, err
	}
	baseURL = fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/netAppAccounts/%s/capacityPools/%s/volumes/%s", baseURL, subscription, info.ResourceGroupsName, info.NetAppAccountName, info.CapacityPools, vol.Name)
	hostType := "CVSHost"
	param := structs.Map(vol)
	param["subnetId"] = subnet
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getANFVolume request failed ", statusCode)
		return anfVolumeResponse{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getANFVolume")
	if responseError != nil {
		return anfVolumeResponse{}, responseError
	}
	var result anfVolumeResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getANFVolume ", err)
		return anfVolumeResponse{}, err
	}

	return result, nil
}

func (c *Client) getCVSWorkingEnvironment(accountID string, WorkingEnvironment string) (string, string, error) {
	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", "", err
		}
		c.Token = accesTokenResult.Token
	}

	baseURL := fmt.Sprintf("/cvs/accounts/%s/working-environments", accountID)
	hostType := "CVSHost"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getCVSWorkingEnvironment request failed ", statusCode)
		return "", "", err
	}
	responseError := apiResponseChecker(statusCode, response, "getCVSWorkingEnvironment")
	if responseError != nil {
		return "", "", responseError
	}
	var results []map[string]interface{}
	if err := json.Unmarshal(response, &results); err != nil {
		log.Print("Failed to unmarshall response from getCVSWorkingEnvironment ", err)
		return "", "", err
	}

	for _, result := range results {
		if strings.ToLower(result["name"].(string)) == strings.ToLower(WorkingEnvironment) {
			return result["credentialsId"].(string), result["provider"].(string), nil
		}
	}

	return "", "", fmt.Errorf(" working environment: %s doesn't exist", WorkingEnvironment)
}

func (c *Client) getCVSAPIRoot(accountName string, workingEnvironment string) (string, error) {
	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", err
		}
		c.Token = accesTokenResult.Token
	}
	accountID, err := c.getAccountByName(accountName)
	if err != nil {
		return "", err
	}
	credentialsID, provider, err := c.getCVSWorkingEnvironment(accountID, workingEnvironment)
	if err != nil {
		return "", err
	}

	if provider == "azure" {
		return fmt.Sprintf("/cvs/azure/accounts/%s/credentials/%s", accountID, credentialsID), nil
	} else if provider == "gcp" {
		return fmt.Sprintf("/cvs/gcp/accounts/%s/credentials/%s", accountID, credentialsID), nil
	} else if provider == "aws" {
		return fmt.Sprintf("/cvs/aws/accounts/%s/credentials/%s", accountID, credentialsID), nil
	} else {
		return "", fmt.Errorf("working environment's provider is not supported or not found")
	}
}

func (c *Client) getSubscription(baseURL string, subscription string) (string, error) {
	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", err
		}
		c.Token = accesTokenResult.Token
	}
	baseURL = fmt.Sprintf("%s/subscriptions", baseURL)
	hostType := "CVSHost"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getSubscriptions request failed ", statusCode)
		return "", err
	}
	responseError := apiResponseChecker(statusCode, response, "getSubscriptions")
	if responseError != nil {
		return "", responseError
	}

	var results []map[string]interface{}
	if err := json.Unmarshal(response, &results); err != nil {
		log.Print("Failed to unmarshall response from getSubscriptions ", err)
		return "", err
	}
	for _, result := range results {
		if strings.ToLower(subscription) == strings.ToLower(result["displayName"].(string)) {
			return result["subscriptionId"].(string), nil
		}
	}

	return "", fmt.Errorf("subscription: %s doesn't exist", subscription)

}

func (c *Client) getSubnetID(baseURL string, virtualNetwork string, subnet string, location string) (string, error) {
	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", err
		}
		c.Token = accesTokenResult.Token
	}
	baseURL = fmt.Sprintf("%s/virtualNetworks?location=%s", baseURL, location)
	hostType := "CVSHost"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getSubnetID request failed ", statusCode)
		return "", err
	}
	responseError := apiResponseChecker(statusCode, response, "getSubnetID")
	if responseError != nil {
		return "", responseError
	}

	var results []interface{}
	if err := json.Unmarshal(response, &results); err != nil {
		log.Print("Failed to unmarshall response from getSubnetID ", err)
		return "", err
	}
	for _, result := range results {
		if strings.ToLower(virtualNetwork) == strings.ToLower(result.(map[string]interface{})["name"].(string)) {
			subnetResults := result.(map[string]interface{})["subnets"].([]interface{})
			for _, subnetResult := range subnetResults {
				if strings.ToLower(subnet) == strings.ToLower(subnetResult.(map[string]interface{})["name"].(string)) {
					return subnetResult.(map[string]interface{})["subnetId"].(string), nil
				}
			}
		}
	}

	return "", fmt.Errorf("subnet: %s doesn't exist", subnet)

}

func (c *Client) deleteANFVolume(vol anfVolumeRequest, info cvsInfo) error {
	baseURL, err := c.getCVSAPIRoot(info.AccountName, vol.WorkingEnvironmentName)
	if err != nil {
		return err
	}
	subscription, err := c.getSubscription(baseURL, info.SubscriptionName)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/subscriptions/%s/resourceGroups/%s/netAppAccounts/%s/capacityPools/%s/volumes/%s", baseURL, subscription, info.ResourceGroupsName, info.NetAppAccountName, info.CapacityPools, vol.Name)
	hostType := "CVSHost"
	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("deleteANFVolume request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "deleteANFVolume")
	if responseError != nil {
		return responseError
	}

	return nil
}
