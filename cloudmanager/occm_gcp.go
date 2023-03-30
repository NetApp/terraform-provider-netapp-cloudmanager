package cloudmanager

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
	Name       string            `yaml:"name"`
	Properties properties        `yaml:"properties"`
	Type       string            `yaml:"type"`
	Metadata   metaData          `yaml:"metadata,omitempty"`
	Labels     map[string]string `yaml:"labels,omitempty"`
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
	Labels            map[string]string  `yaml:"labels,omitempty"`
}

type label struct {
	Key   string `yaml:"key,omitempty"`
	Value string `yaml:"value,omitempty"`
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

func (c *Client) getCustomDataForGCP(registerAgentTOService registerAgentTOServiceRequest, proxyCertificates []string, clientID string) (string, string, error) {
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

	userDataRespone, err := c.registerAgentTOServiceForGCP(registerAgentTOService, clientID)
	if err != nil {
		return "", "", err
	}

	newClientID := userDataRespone.ClientID
	c.AccountID = userDataRespone.AccountID

	userDataRespone.ProxySettings.ProxyCertificates = proxyCertificates
	rawUserData, _ := json.Marshal(userDataRespone)
	userData := string(rawUserData)
	log.Print("userData ", userData)

	return userData, newClientID, nil
}

func (c *Client) registerAgentTOServiceForGCP(registerAgentTOServiceRequest registerAgentTOServiceRequest, clientID string) (createUserData, error) {

	baseURL := "/agents-mgmt/connector-setup"
	hostType := "CloudManagerHost"

	registerAgentTOServiceRequest.Placement.Provider = "GCP"
	log.Print(c.Token)

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

func (c *Client) deployGCPVM(occmDetails createOCCMDetails, proxyCertificates []string, clientID string, retries int) (OCCMMResult, error) {
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

	userData, newClientID, err := c.getCustomDataForGCP(registerAgentTOService, proxyCertificates, clientID)
	if err != nil {
		return OCCMMResult{}, err
	}

	c.UserData = userData
	var result OCCMMResult
	log.Printf("Set result clientid=%s", newClientID)
	result.ClientID = newClientID
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
	if len(occmDetails.Tags) > 0 {
		t.Properties.Tags.Items = append(t.Properties.Tags.Items, occmDetails.Tags...)
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
	subnetID, err := convertSubnetID(projectID, occmDetails, occmDetails.SubnetID)
	if err != nil {
		return OCCMMResult{}, err
	}
	t.Properties.NetworkInterfaces = []networkInterface{
		{AccessConfigs: accessConfigs,
			Kind:       "compute#networkInterface",
			Subnetwork: subnetID,
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

	if occmDetails.Labels != nil {
		t.Properties.Labels = occmDetails.Labels
		td.Properties.Labels = occmDetails.Labels
	}

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
	log.Printf("deployGCPVM: call depolyments api base client=%s", newClientID)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, nil, "", hostType, newClientID)
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

func (c *Client) getdeployGCPVM(occmDetails createOCCMDetails, id string, clientID string) (string, error) {

	log.Print("getdeployGCPVM")

	baseURL := fmt.Sprintf("/deploymentmanager/v2/projects/%s/global/deployments/%s%s", occmDetails.GCPProject, occmDetails.Name, occmDetails.GCPCommonSuffixName)
	hostType := "GCPDeploymentManager"

	log.Print("GET")
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, "", hostType, clientID)
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

func (c *Client) getDisk(occmDetails createOCCMDetails, clientID string) (map[string]interface{}, error) {
	hostType := "GCPCompute"
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/disks/%s-vm-disk-boot", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, "", hostType, clientID)
	if err != nil {
		log.Printf("getDisk request failed: %s", err.Error())
		return nil, err
	}

	responseError := apiResponseChecker(statusCode, response, "getDisk")
	if responseError != nil {
		return nil, responseError
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getDisk")
		return nil, err
	}
	return result, nil
}

func (c *Client) getVMInstance(occmDetails createOCCMDetails, clientID string) (map[string]interface{}, error) {

	log.Print("getVMInstance")

	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"

	log.Print("GET")
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, "", hostType, clientID)
	if err != nil {
		log.Printf("getVMInstance request failed: %s", err.Error())
		return nil, err
	}

	responseError := apiResponseChecker(statusCode, response, "getVMInstance")
	if responseError != nil {
		return nil, responseError
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVMInstance")
		return nil, err
	}

	return result, nil
}

// This function is not used because can't get the update instance API from GCP working. Receive the following error: Boot disk must be the first disk attached to the instance.
// Although there is only one disk exists all the time, I can't figure it out to make it work.
func (c *Client) updateVMInstance(occmDetails createOCCMDetails, clientID string, updatePropertities map[string]interface{}) error {
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"
	statusCode, response, _, err := c.CallAPIMethod("PUT", baseURL, updatePropertities, "", hostType, clientID)

	if err != nil {
		log.Print("updateVMInstance request failed")
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "getVMInstance")
	if responseError != nil {
		return responseError
	}

	return nil

}

func (c *Client) setVMLabels(occmDetails createOCCMDetails, labels map[string]interface{}, clientID string) error {
	log.Print("setVMLabels")

	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm/setLabels", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, labels, "", hostType, clientID)
	if err != nil {
		log.Printf("setVMLabels request failed: %s", err.Error())
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "setVMLabels")
	if responseError != nil {
		return responseError
	}
	return nil
}

func (c *Client) setDiskLabels(occmDetails createOCCMDetails, labels map[string]interface{}, clientID string) error {
	log.Print("setDiskLabels")

	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/disks/%s-vm-disk-boot/setLabels", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, labels, "", hostType, clientID)
	if err != nil {
		log.Printf("setDiskLabels request failed: %s", err.Error())
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "setDiskLabels")
	if responseError != nil {
		return responseError
	}
	return nil
}

func (c *Client) setVMInstaceTags(occmDetails createOCCMDetails, fingerprint string, clientID string) error {
	log.Print("setVMInstaceTags")

	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm/setTags", occmDetails.GCPProject, occmDetails.Region, occmDetails.Name)
	hostType := "GCPDeploymentManager"
	body := make(map[string]interface{})
	body["items"] = occmDetails.Tags
	body["fingerprint"] = fingerprint
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, body, "", hostType, clientID)
	if err != nil {
		log.Print("setVMInstaceTags request failed")
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "setVMInstaceTags")
	if responseError != nil {
		return responseError
	}
	return nil
}

func (c *Client) deleteOCCMGCP(request deleteOCCMDetails, clientID string) error {

	log.Printf("deleteOCCMGCP %s client %s", request.Name, clientID)

	baseURL := fmt.Sprintf("/deploymentmanager/v2/projects/%s/global/deployments/%s%s", request.Project, request.Name, request.GCPCommonSuffixName)
	hostType := "GCPDeploymentManager"

	log.Print("DELETE")
	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, "", hostType, clientID)
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

func convertSubnetID(projectID string, occmDetails createOCCMDetails, input string) (string, error) {
	if !strings.Contains(input, "/") {
		return fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", projectID, occmDetails.Region, occmDetails.SubnetID), nil
	}
	return input, nil

}
