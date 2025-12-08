package cloudmanager

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
)

// OCCMMGCPResult the response name  of a occm
type OCCMMGCPResult struct {
	Name string `json:"name"`
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

func (c *Client) deployGCPVM(occmDetails createOCCMDetails, proxyCertificates []string, clientID string) (OCCMMResult, error) {
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

	// Get GCP token for API calls
	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return OCCMMResult{}, err
	}

	// Step 1: Create the disk first
	err = c.createGCPDisk(occmDetails, token, newClientID)
	if err != nil {
		return OCCMMResult{}, err
	}

	// Step 2: Create the VM instance
	err = c.createGCPInstance(occmDetails, token, newClientID, gcpCustomData)
	if err != nil {
		return OCCMMResult{}, err
	}

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
				return OCCMMResult{}, fmt.Errorf("taking too long for OCCM agent to be active or not properly setup")
			}
			time.Sleep(time.Duration(30) * time.Second)
			retries--
		}
	}

	return result, nil
}

func (c *Client) getdeployGCPVM(occmDetails createOCCMDetails, id string, clientID string) (string, error) {
	log.Print("getdeployGCPVM")
	log.Printf("Expected ID parameter: %s", id)

	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return "", err
	}

	instanceName := fmt.Sprintf("%s-vm", occmDetails.Name)
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s", occmDetails.GCPProject, occmDetails.Zone, instanceName)
	hostType := "GCPCompute"

	log.Printf("Making GET request to: %s", baseURL)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, token, hostType, clientID)
	if err != nil {
		log.Printf("getdeployGCPVM request failed: %v", err)
		return "", err
	}

	// Handle 404 - instance not found
	// if statusCode == 404 {
	// 	log.Printf("Instance not found (404): %s", instanceName)

	// 	// Try to check if this is a legacy deployment manager instance
	// 	// by looking for any instance with the base name
	// 	legacyInstanceName := occmDetails.Name
	// 	legacyURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s", occmDetails.GCPProject, occmDetails.Zone, legacyInstanceName)
	// 	log.Printf("Trying legacy instance name: %s", legacyInstanceName)

	// 	legacyStatusCode, legacyResponse, _, legacyErr := c.CallAPIMethod("GET", legacyURL, nil, token, hostType, clientID)
	// 	if legacyErr == nil && legacyStatusCode == 200 {
	// 		log.Printf("Found legacy instance with name: %s", legacyInstanceName)
	// 		var legacyResult map[string]interface{}
	// 		if err := json.Unmarshal(legacyResponse, &legacyResult); err == nil {
	// 			if name, ok := legacyResult["name"].(string); ok && name == legacyInstanceName {
	// 				log.Printf("Legacy instance found: %s, returning expected ID: %s", name, id)
	// 				return id, nil
	// 			}
	// 		}
	// 	}

	// 	return "", nil
	// }

	responseError := apiResponseChecker(statusCode, response, "getdeployGCPVM")
	if responseError != nil {
		log.Printf("getdeployGCPVM response error: %v", responseError)
		return "", responseError
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getdeployGCPVM")
		return "", err
	}

	if name, ok := result["name"].(string); ok && name == instanceName {
		// Instance exists, return the expected ID that was passed in
		log.Printf("Instance found: %s, returning expected ID: %s", name, id)
		return id, nil
	}

	return "", nil
}

func (c *Client) getDisk(occmDetails createOCCMDetails, clientID string) (map[string]interface{}, error) {
	log.Print("getDisk")

	hostType := "GCPCompute"
	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return nil, err
	}
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/disks/%s-vm-disk-boot", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, token, hostType, clientID)
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

	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return nil, err
	}
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"

	log.Print("GET")
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, token, hostType, clientID)
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
	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return err
	}
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"
	statusCode, response, _, err := c.CallAPIMethod("PUT", baseURL, updatePropertities, token, hostType, clientID)

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
	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return err
	}
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm/setLabels", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, labels, token, hostType, clientID)
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
	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return err
	}
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/disks/%s-vm-disk-boot/setLabels", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPCompute"
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, labels, token, hostType, clientID)
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
	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return err
	}
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s-vm/setTags", occmDetails.GCPProject, occmDetails.Zone, occmDetails.Name)
	hostType := "GCPDeploymentManager"
	body := make(map[string]interface{})
	body["items"] = occmDetails.Tags
	body["fingerprint"] = fingerprint
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, body, token, hostType, clientID)
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

	token, err := c.getGCPToken(c.GCPServiceAccountKey)
	if err != nil {
		return err
	}

	// Step 1: Delete the VM instance
	err = c.deleteGCPInstance(request, token, clientID)
	if err != nil {
		log.Printf("Warning: Failed to delete instance: %v", err)
		// Continue with deletion process even if instance deletion fails
	}

	// We specify "autoDelete": true in line 576, as a result the disk will be deleted automatically when the instance is deleted.
	// Keep the code here commented out in case we need to delete the disk manually in the future.
	// // Step 2: Delete the disk
	// err = c.deleteGCPDisk(request, token, clientID)
	// if err != nil {
	// 	log.Printf("Warning: Failed to delete disk: %v", err)
	// 	// Continue with deletion process even if disk deletion fails
	// }

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
				return fmt.Errorf("taking too long for instance to finish terminating")
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

func (c *Client) getGCPToken(gcpServiceAccountKey string) (string, error) {
	var token string
	log.Printf("getGCPToken...")
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/compute.readonly",
		"https://www.googleapis.com/auth/ndev.cloudman",
		"https://www.googleapis.com/auth/ndev.cloudman.readonly",
		"https://www.googleapis.com/auth/devstorage.full_control",
		"https://www.googleapis.com/auth/devstorage.read_write",
	}
	if gcpServiceAccountKey != "" {
		var c = struct {
			Email      string `json:"client_email"`
			PrivateKey string `json:"private_key"`
		}{}
		json.Unmarshal([]byte(gcpServiceAccountKey), &c)
		config := &jwt.Config{
			Email:      c.Email,
			PrivateKey: []byte(c.PrivateKey),
			Scopes:     scopes,
			TokenURL:   google.JWTTokenURL,
		}
		gcpToken, err := config.TokenSource(context.Background()).Token()
		if err != nil {
			return "", err
		}
		token = gcpToken.AccessToken
	} else {
		// find default application credential
		ctx := context.Background()
		credential, err := google.FindDefaultCredentials(ctx, scopes...)
		if err != nil {
			return "", fmt.Errorf("cannot get credentials: %v", err)
		}
		t, err := credential.TokenSource.Token()
		if err != nil {
			return "", fmt.Errorf("getGCPToken failed on get token from credential: %v", err)
		}
		token = t.AccessToken
	}

	return token, nil
}

// createGCPDisk creates a disk using GCP Compute API
func (c *Client) createGCPDisk(occmDetails createOCCMDetails, token, clientID string) error {
	log.Print("Creating GCP disk...")

	// Validate that required fields are not empty
	if occmDetails.Zone == "" {
		return fmt.Errorf("zone is required but not provided")
	}
	if occmDetails.GCPProject == "" {
		return fmt.Errorf("GCP project is required but not provided")
	}

	deviceName := fmt.Sprintf("%s-vm-disk-boot", occmDetails.Name)

	diskBody := map[string]interface{}{
		"name":        deviceName,
		"sizeGb":      100,
		"sourceImage": fmt.Sprintf("projects/%s/global/images/family/%s", c.GCPImageProject, c.GCPImageFamily),
		"type":        fmt.Sprintf("zones/%s/diskTypes/pd-ssd", occmDetails.Zone),
		"zone":        occmDetails.Zone,
	}

	// Add labels if specified
	if occmDetails.Labels != nil {
		labels := map[string]string{}
		for key, value := range occmDetails.Labels {
			labels[key] = value
		}
		diskBody["labels"] = labels
	}

	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/disks", occmDetails.GCPProject, occmDetails.Zone)
	hostType := "GCPCompute"

	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, diskBody, token, hostType, clientID)
	if err != nil {
		log.Printf("createGCPDisk request failed: %s", err.Error())
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "createGCPDisk")
	if responseError != nil {
		return responseError
	}

	log.Printf("Disk creation initiated: %s", deviceName)

	// Wait for disk to be ready
	return c.waitForDiskReady(occmDetails, token, clientID, deviceName)
}

// createGCPInstance creates a VM instance using GCP Compute API
func (c *Client) createGCPInstance(occmDetails createOCCMDetails, token, clientID, gcpCustomData string) error {
	log.Print("Creating GCP instance...")

	// Validate that required fields are not empty
	if occmDetails.Zone == "" {
		return fmt.Errorf("zone is required but not provided")
	}
	if occmDetails.GCPProject == "" {
		return fmt.Errorf("GCP project is required but not provided")
	}

	deviceName := fmt.Sprintf("%s-vm-disk-boot", occmDetails.Name)
	instanceName := fmt.Sprintf("%s-vm", occmDetails.Name)

	gcpSaScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/compute.readonly",
		"https://www.googleapis.com/auth/ndev.cloudman",
		"https://www.googleapis.com/auth/ndev.cloudman.readonly",
	}

	// Build tags
	var tags []string
	if occmDetails.FirewallTags {
		tags = []string{"firewall-tag-bvsu", "http-server", "https-server"}
	}
	if len(occmDetails.Tags) > 0 {
		tags = append(tags, occmDetails.Tags...)
	}

	// Build access configs for public IP
	var accessConfigs []map[string]interface{}
	if occmDetails.AssociatePublicIP {
		accessConfigs = []map[string]interface{}{
			{
				"kind":        "compute#accessConfig",
				"name":        "External NAT",
				"type":        "ONE_TO_ONE_NAT",
				"networkTier": "PREMIUM",
			},
		}
	}

	// Determine project ID for network
	var projectID string
	if occmDetails.NetworkProjectID != "" {
		projectID = occmDetails.NetworkProjectID
	} else {
		projectID = occmDetails.GCPProject
	}

	subnetID, err := convertSubnetID(projectID, occmDetails, occmDetails.SubnetID)
	if err != nil {
		return err
	}

	instanceBody := map[string]interface{}{
		"name":        instanceName,
		"machineType": fmt.Sprintf("zones/%s/machineTypes/%s", occmDetails.Zone, occmDetails.MachineType),
		"zone":        occmDetails.Zone,
		"disks": []map[string]interface{}{
			{
				"autoDelete": true,
				"boot":       true,
				"deviceName": deviceName,
				"source":     fmt.Sprintf("projects/%s/zones/%s/disks/%s", occmDetails.GCPProject, occmDetails.Zone, deviceName),
				"type":       "PERSISTENT",
			},
		},
		"metadata": map[string]interface{}{
			"items": []map[string]interface{}{
				{"key": "serial-port-enable", "value": "1"},
				{"key": "customData", "value": gcpCustomData},
			},
		},
		"networkInterfaces": []map[string]interface{}{
			{
				"accessConfigs": accessConfigs,
				"kind":          "compute#networkInterface",
				"subnetwork":    subnetID,
			},
		},
		"serviceAccounts": []map[string]interface{}{
			{
				"email":  occmDetails.ServiceAccountEmail,
				"scopes": gcpSaScopes,
			},
		},
	}

	// Add tags if specified
	if len(tags) > 0 {
		instanceBody["tags"] = map[string]interface{}{
			"items": tags,
		}
	}

	// Add labels if specified
	if occmDetails.Labels != nil {
		labels := map[string]string{}
		for key, value := range occmDetails.Labels {
			labels[key] = value
		}
		instanceBody["labels"] = labels
	}

	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances", occmDetails.GCPProject, occmDetails.Zone)
	hostType := "GCPCompute"

	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, instanceBody, token, hostType, clientID)
	if err != nil {
		log.Printf("createGCPInstance request failed: %s", err.Error())
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "createGCPInstance")
	if responseError != nil {
		return responseError
	}

	log.Printf("Instance creation initiated: %s", instanceName)

	return nil

	// // Wait for instance to be ready
	// return c.waitForInstanceReady(occmDetails, token, clientID, instanceName)
}

// deleteGCPInstance deletes a VM instance using GCP Compute API
func (c *Client) deleteGCPInstance(request deleteOCCMDetails, token, clientID string) error {
	log.Print("Deleting GCP instance...")

	instanceName := fmt.Sprintf("%s-vm", request.Name)
	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/instances/%s", request.Project, request.Region, instanceName)
	hostType := "GCPCompute"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, token, hostType, clientID)
	if err != nil {
		log.Printf("deleteGCPInstance request failed: %s", err.Error())
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteGCPInstance")
	if responseError != nil {
		return responseError
	}

	log.Printf("Instance deletion initiated: %s", instanceName)

	return nil

	// // Wait for instance to be deleted
	// return c.waitForInstanceDeleted(request, token, clientID, instanceName)
}

// waitForDiskReady waits for a GCP disk to be in READY state
func (c *Client) waitForDiskReady(occmDetails createOCCMDetails, token, clientID, diskName string) error {
	log.Printf("Waiting for disk to be ready: %s", diskName)

	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/disks/%s", occmDetails.GCPProject, occmDetails.Zone, diskName)
	hostType := "GCPCompute"

	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, token, hostType, clientID)
		if err != nil {
			log.Printf("Error checking disk status: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if statusCode == 200 {
			var diskStatus map[string]interface{}
			if err := json.Unmarshal(response, &diskStatus); err != nil {
				log.Printf("Failed to unmarshal disk status: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			status, ok := diskStatus["status"].(string)
			if !ok {
				log.Printf("Disk status not found in response")
				time.Sleep(5 * time.Second)
				continue
			}

			log.Printf("Disk %s status: %s", diskName, status)
			if status == "READY" {
				log.Printf("Disk %s is ready", diskName)
				return nil
			}
		}

		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for disk %s to be ready", diskName)
}

// We specify "autoDelete": true in line 576, as a result the disk will be deleted automatically when the instance is deleted.
// Keep the code here commented out in case we need to delete the disk manually in the future.
// // deleteGCPDisk deletes a disk using GCP Compute API
// func (c *Client) deleteGCPDisk(request deleteOCCMDetails, token, clientID string) error {
// 	log.Print("Deleting GCP disk...")

// 	diskName := fmt.Sprintf("%s-vm-disk-boot", request.Name)
// 	baseURL := fmt.Sprintf("/compute/v1/projects/%s/zones/%s/disks/%s", request.Project, request.Region, diskName)
// 	hostType := "GCPCompute"

// 	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, token, hostType, clientID)
// 	if err != nil {
// 		log.Printf("deleteGCPDisk request failed: %s", err.Error())
// 		return err
// 	}

// 	responseError := apiResponseChecker(statusCode, response, "deleteGCPDisk")
// 	if responseError != nil {
// 		return responseError
// 	}

// 	log.Printf("Disk deletion initiated: %s", diskName)

// 	// Wait for disk to be deleted
// 	return c.waitForDiskDeleted(request, token, clientID, diskName)
// }
