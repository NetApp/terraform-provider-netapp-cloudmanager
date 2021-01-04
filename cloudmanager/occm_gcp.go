package cloudmanager

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
	"gopkg.in/yaml.v2"
)

// OCCMMGCPResult the response name  of a occm
type OCCMMGCPResult struct {
	Name string `json:"name"`
}

// GCP template
type tcontent struct {
	Resources []tresource `yaml:"resources"`
}

type tresource struct {
	Name       string     `yaml:"name"`
	Properties properties `yaml:"properties"`
	Type       string     `yaml:"type"`
	Metadata   metaData   `yaml:"metadata,omitempty"`
}

type properties struct {
	Name              string             `yaml:"name,omitempty"`
	SizeGb            int                `yaml:"sizeGb,omitempty"`
	SourceImage       string             `yaml:"sourceImage,omitempty"`
	Type              string             `yaml:"type,omitempty"`
	Tags              tags               `yaml:"tags,omitempty"`
	Disks             []pdisk            `yaml:"disks,omitempty"`
	MachineType       string             `yaml:"machineType,omitempty"`
	Zone              string             `yaml:"zone"`
	Metadata          pmetadata          `yaml:"metadata,omitempty"`
	NetworkInterfaces []networkInterface `yaml:"networkInterfaces,omitempty"`
	ServiceAccounts   []serviceAccount   `yaml:"serviceAccounts,omitempty"`
}

type metaData struct {
	DependsOn []string `yaml:"dependsOn"`
}

type tags struct {
	Items []string `yaml:"items"`
}

type pdisk struct {
	AutoDelete bool   `yaml:"autoDelete"`
	Boot       bool   `yaml:"boot"`
	DeviceName string `yaml:"deviceName"`
	Name       string `yaml:"name"`
	Source     string `yaml:"source"`
	Type       string `yaml:"type"`
}

type pmetadata struct {
	Items []item `yaml:"items"`
}

type item struct {
	Key, Value interface{}
}

type networkInterface struct {
	AccessConfigs []accessConfig `yaml:"accessConfigs"`
	Kind          string         `yaml:"kind"`
	Subnetwork    string         `yaml:"subnetwork"`
}

type accessConfig struct {
	Kind        string `yaml:"kind"`
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	NetworkTier string `yaml:"networkTier"`
}

type serviceAccount struct {
	Email  string   `yaml:"email"`
	Scopes []string `yaml:"scopes"`
}

func (c *Client) getCustomDataForGCP(registerAgentTOService registerAgentTOServiceRequest, proxyCertificates []string) (string, error) {
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

	userDataRespone.ProxySettings.ProxyCertificates = proxyCertificates
	rawUserData, _ := json.Marshal(userDataRespone)
	userData := string(rawUserData)
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

func (c *Client) deployGCPVM(occmDetails createOCCMDetails, proxyCertificates []string) (OCCMMResult, error) {
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

	userData, err := c.getCustomDataForGCP(registerAgentTOService, proxyCertificates)
	if err != nil {
		return OCCMMResult{}, err
	}

	c.UserData = userData
	var result OCCMMResult
	result.ClientID = c.ClientID
	result.AccountID = c.AccountID

	gcpCustomData := base64.StdEncoding.EncodeToString([]byte(userData))

	gcpSaScopes := []string{"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/compute.readonly",
		"https://www.googleapis.com/auth/ndev.cloudman",
		"https://www.googleapis.com/auth/ndev.cloudman.readonly"}

	content := tcontent{}

	// first resource
	t := tresource{}
	t.Name = fmt.Sprintf("%s-vm", occmDetails.Name)
	if occmDetails.FirewallTags == true {
		t.Properties.Tags.Items = []string{"firewall-tag-bvsu", "http-server", "https-server"}
	}
	deviceName := fmt.Sprintf("%s-vm-disk-boot", occmDetails.Name)
	t.Properties.Disks = []pdisk{
		{AutoDelete: true,
			Boot:       true,
			DeviceName: deviceName,
			Name:       deviceName,
			Source:     fmt.Sprintf("\\\"$(ref.%s.selfLink)\\\"", deviceName),
			Type:       "PERSISTENT",
		},
	}
	t.Properties.MachineType = fmt.Sprintf("zones/%s/machineTypes/%s", occmDetails.Zone, occmDetails.MachineType)
	t.Properties.Metadata.Items = []item{
		{Key: "serial-port-enable", Value: 1},
		{Key: "customData", Value: gcpCustomData}}

	var accessConfigs []accessConfig
	if occmDetails.AssociatePublicIP == true {
		accessConfigs = []accessConfig{{Kind: "compute#accessConfig", Name: "External NAT", Type: "ONE_TO_ONE_NAT", NetworkTier: "PREMIUM"}}
	}
	var projectID string
	if occmDetails.NetworkProjectID != "" {
		projectID = occmDetails.NetworkProjectID
	} else {
		projectID = occmDetails.GCPProject
	}
	t.Properties.NetworkInterfaces = []networkInterface{
		{AccessConfigs: accessConfigs,
			Kind:       "compute#networkInterface",
			Subnetwork: fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", projectID, occmDetails.Region, occmDetails.SubnetID),
		},
	}
	t.Properties.ServiceAccounts = []serviceAccount{{Email: occmDetails.ServiceAccountEmail, Scopes: gcpSaScopes}}
	t.Properties.Zone = occmDetails.Zone
	t.Type = "compute.v1.instance"
	t.Metadata.DependsOn = []string{deviceName}

	// the resource which the first resource depends on
	td := tresource{}
	td.Name = deviceName
	td.Properties.Name = deviceName
	td.Properties.SizeGb = 100
	td.Properties.SourceImage = fmt.Sprintf("projects/%s/global/images/family/%s", c.GCPImageProject, c.GCPImageFamily)
	td.Properties.Type = fmt.Sprintf("zones/%s/diskTypes/pd-ssd", occmDetails.Zone)
	td.Properties.Zone = occmDetails.Zone
	td.Type = "compute.v1.disks"

	content.Resources = []tresource{t, td}
	data, err := yaml.Marshal(&content)
	if err != nil {
		return OCCMMResult{}, fmt.Errorf("error: %v", err)
	}
	mydata := string(data)
	c.GCPDeploymentTemplate = fmt.Sprintf(`{
		"name": "%s%s",
		"target": {
		"config": {
		"content": "%s"
		}
	}
	}`, occmDetails.Name, occmDetails.GCPCommonSuffixName, mydata)

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
