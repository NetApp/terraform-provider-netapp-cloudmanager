package cloudmanager

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
)

// OCCMMGCPResult the response name  of a occm
type OCCMMGCPResult struct {
	Name string `json:"name"`
}

func (c *Client) getCustomDataForGCP(registerAgentTOService registerAgentTOServiceRequest) (string, error) {
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

	userDataRespone, err := c.registerAgentTOServiceForGCP(registerAgentTOService)
	if err != nil {
		return "", err
	}

	c.ClientID = userDataRespone.ClientID
	c.AccountID = userDataRespone.AccountID

	userData := fmt.Sprintf(`{"instanceName":"%s","company":"%s","clientId":"%s","clientSecret":"%s","systemId":"%s","tenancyAccountId":"%s"}`, userDataRespone.Name, userDataRespone.Company, userDataRespone.ClientID, userDataRespone.ClientSecret, userDataRespone.UUID, userDataRespone.AccountID)
	log.Print("userData ", userData)

	return userData, nil
}

func (c *Client) registerAgentTOServiceForGCP(registerAgentTOServiceRequest registerAgentTOServiceRequest) (createUserData, error) {

	baseURL := "/agents-mgmt/connector-setup"
	hostType := "CloudManagerHost"

	registerAgentTOServiceRequest.Placement.Provider = "GCP"
	log.Print(c.Token)

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

func (c *Client) deployGCPVM(occmDetails createOCCMDetails) (OCCMMResult, error) {

	var registerAgentTOService registerAgentTOServiceRequest
	registerAgentTOService.Name = occmDetails.Name
	registerAgentTOService.Placement.Region = occmDetails.Region
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

	registerAgentTOService.Placement.Subnet = occmDetails.SubnetID

	userData, err := c.getCustomDataForGCP(registerAgentTOService)
	if err != nil {
		return OCCMMResult{}, err
	}

	c.UserData = userData
	var result OCCMMResult
	result.ClientID = c.ClientID
	result.AccountID = c.AccountID

	gcpCustomData := base64.StdEncoding.EncodeToString([]byte(userData))
	gcpSaScopes := `- https://www.googleapis.com/auth/cloud-platform\n      - https://www.googleapis.com/auth/compute\n      - https://www.googleapis.com/auth/compute.readonly\n      - https://www.googleapis.com/auth/ndev.cloudman\n      - https://www.googleapis.com/auth/ndev.cloudman.readonly`

	var tags string
	var accessConfig string
	if occmDetails.FirewallTags == true {
		tags = `    tags:\n      items:\n      - firewall-tag-bvsu\n      - http-server\n      - https-server\n`
	}

	if occmDetails.AssociatePublicIP == true {
		accessConfig = `\n      - kind: compute#accessConfig\n        name: External NAT\n        type: ONE_TO_ONE_NAT\n        networkTier: PREMIUM`
	} else {
		accessConfig = ` []`
	}

	if occmDetails.NetworkProjectID != "" {
		c.GCPDeploymentTemplate = fmt.Sprintf(`{
			"name": "%s%s",
			"target": {
			"config": {
			"content": "resources:\n- name: %s-vm\n  properties:\n%s    disks:\n    - autoDelete: true\n      boot: true\n      deviceName: %s-vm-disk-boot\n      name: %s-vm-disk-boot\n      source: \"$(ref.%s-vm-disk-boot.selfLink)\"\n      type: PERSISTENT\n    machineType: zones/%s/machineTypes/%s\n    metadata:\n      items:\n      - key: serial-port-enable\n        value: 1\n      - key: customData\n        value: %s\n    networkInterfaces:\n    - accessConfigs:%s\n      kind: compute#networkInterface\n      subnetwork: projects/%s/regions/%s/subnetworks/%s\n    serviceAccounts:\n    - email: %s\n      scopes:\n      %s\n    zone: %s\n  type: compute.v1.instance\n  metadata:\n    dependsOn:\n    - %s-vm-disk-boot\n- name: %s-vm-disk-boot\n  properties:\n    name: %s-vm-disk-boot\n    sizeGb: 100\n    sourceImage: projects/%s/global/images/family/%s\n    type: zones/%s/diskTypes/pd-ssd\n    zone: %s\n  type: compute.v1.disks"
			}
		}
		}`, occmDetails.Name, occmDetails.GCPCommonSuffixName, occmDetails.Name, tags, occmDetails.Name, occmDetails.Name, occmDetails.Name, occmDetails.Zone, occmDetails.MachineType, gcpCustomData, accessConfig, c.Project, occmDetails.Region, occmDetails.NetworkProjectID, occmDetails.ServiceAccountEmail, gcpSaScopes, occmDetails.Zone, occmDetails.Name, occmDetails.Name, occmDetails.Name, c.GCPImageProject, c.GCPImageFamily, occmDetails.Zone, occmDetails.Zone)
	} else {
		c.GCPDeploymentTemplate = fmt.Sprintf(`{
			"name": "%s%s",
			"target": {
			"config": {
			"content": "resources:\n- name: %s-vm\n  properties:\n%s    disks:\n    - autoDelete: true\n      boot: true\n      deviceName: %s-vm-disk-boot\n      name: %s-vm-disk-boot\n      source: \"$(ref.%s-vm-disk-boot.selfLink)\"\n      type: PERSISTENT\n    machineType: zones/%s/machineTypes/%s\n    metadata:\n      items:\n      - key: serial-port-enable\n        value: 1\n      - key: customData\n        value: %s\n    networkInterfaces:\n    - accessConfigs:%s\n      kind: compute#networkInterface\n      subnetwork: projects/%s/regions/%s/subnetworks/%s\n    serviceAccounts:\n    - email: %s\n      scopes:\n      %s\n    zone: %s\n  type: compute.v1.instance\n  metadata:\n    dependsOn:\n    - %s-vm-disk-boot\n- name: %s-vm-disk-boot\n  properties:\n    name: %s-vm-disk-boot\n    sizeGb: 100\n    sourceImage: projects/%s/global/images/family/%s\n    type: zones/%s/diskTypes/pd-ssd\n    zone: %s\n  type: compute.v1.disks"
			}
		}
		}`, occmDetails.Name, occmDetails.GCPCommonSuffixName, occmDetails.Name, tags, occmDetails.Name, occmDetails.Name, occmDetails.Name, occmDetails.Zone, occmDetails.MachineType, gcpCustomData, accessConfig, c.Project, occmDetails.Region, occmDetails.SubnetID, occmDetails.ServiceAccountEmail, gcpSaScopes, occmDetails.Zone, occmDetails.Name, occmDetails.Name, occmDetails.Name, c.GCPImageProject, c.GCPImageFamily, occmDetails.Zone, occmDetails.Zone)
	}

	baseURL := fmt.Sprintf("/deploymentmanager/v2/projects/%s/global/deployments", occmDetails.GCPProject)
	hostType := "GCPDeploymentManager"

	log.Print("POST")
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, nil, "", hostType)
	if err != nil {
		log.Print("deployGCPVM request failed")
		return OCCMMResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "deployGCPVM")
	if responseError != nil {
		return OCCMMResult{}, responseError
	}

	log.Print("Sleep for 2 minutes")
	time.Sleep(time.Duration(120) * time.Second)

	retries := 16
	for {
		occmResp, err := c.checkOCCMStatus()
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

func (c *Client) getdeployGCPVM(occmDetails createOCCMDetails, id string) (string, error) {

	log.Print("getdeployGCPVM")

	baseURL := fmt.Sprintf("/deploymentmanager/v2/projects/%s/global/deployments/%s%s", occmDetails.GCPProject, occmDetails.Name, occmDetails.GCPCommonSuffixName)
	hostType := "GCPDeploymentManager"

	log.Print("GET")
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, "", hostType)
	if err != nil {
		log.Print("getdeployGCPVM request failed")
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getdeployGCPVM")
	if responseError != nil {
		return "", responseError
	}

	var result OCCMMGCPResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolumeByID")
		return "", err
	}

	if result.Name == id {
		return result.Name, nil
	}

	return "", nil
}

func (c *Client) deleteOCCMGCP(request deleteOCCMDetails) error {

	log.Print("deleteOCCMGCP")

	baseURL := fmt.Sprintf("/deploymentmanager/v2/projects/%s/global/deployments/%s%s", request.Project, request.Name, request.GCPCommonSuffixName)
	hostType := "GCPDeploymentManager"

	log.Print("DELETE")
	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, "", hostType)
	if err != nil {
		log.Print("deleteOCCMGCP request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteOCCMGCP")
	if responseError != nil {
		return responseError
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
