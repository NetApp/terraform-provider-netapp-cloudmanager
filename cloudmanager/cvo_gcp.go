package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

// GCPLicenseTypes is the GCP License types
var GCPLicenseTypes = []string{"gcp-cot-standard-paygo", "gcp-cot-explore-paygo", "gcp-cot-premium-paygo",
	"gcp-cot-premium-byol", "gcp-ha-cot-standard-paygo", "gcp-ha-cot-premium-paygo", "gcp-ha-cot-explore-paygo",
	"gcp-ha-cot-premium-byol", "capacity-paygo", "ha-capacity-paygo"}

// createCVOGCPDetails the users input for creating a CVO
type createCVOGCPDetails struct {
	Name                    string                  `structs:"name"`
	DataEncryptionType      string                  `structs:"dataEncryptionType"`
	WorkspaceID             string                  `structs:"tenantId,omitempty"`
	Region                  string                  `structs:"region"`
	GCPServiceAccount       string                  `structs:"gcpServiceAccount,omitempty"`
	VpcID                   string                  `structs:"vpcId"`
	SvmPassword             string                  `structs:"svmPassword"`
	SvmName                 string                  `structs:"svmName,omitempty"`
	VsaMetadata             vsaMetadata             `structs:"vsaMetadata"`
	GCPVolumeSize           diskSize                `structs:"gcpVolumeSize"`
	GCPVolumeType           string                  `structs:"gcpVolumeType"`
	SubnetID                string                  `structs:"subnetId"`
	SubnetPath              string                  `structs:"subnetPath"`
	Project                 string                  `structs:"project"`
	CapacityTier            string                  `structs:"capacityTier,omitempty"`
	TierLevel               string                  `structs:"tierLevel,omitempty"`
	NssAccount              string                  `structs:"nssAccount,omitempty"`
	WormRequest             wormRequest             `structs:"wormRequest,omitempty"`
	WritingSpeedState       string                  `structs:"writingSpeedState,omitempty"`
	SerialNumber            string                  `structs:"serialNumber,omitempty"`
	GCPLabels               []gcpLabels             `structs:"gcpLabels,omitempty"`
	GcpEncryptionParameters gcpEncryptionParameters `structs:"gcpEncryptionParameters,omitempty"`
	FirewallRule            string                  `structs:"firewallRule,omitempty"`
	BackupVolumesToCbs      bool                    `structs:"backupVolumesToCbs"`
	EnableCompliance        bool                    `structs:"enableCompliance"`
	IsHA                    bool
	HAParams                haParamsGCP `structs:"haParams,omitempty"`
	FlashCache              bool        `structs:"flashCache"`
}

// gcpLabels the input for requesting a CVO
type gcpLabels struct {
	LabelKey   string `structs:"labelKey"`
	LabelValue string `structs:"labelValue,omitempty"`
}

// gcpSVMs the input for adding SVMs to a CVO HA
type gcpSVM struct {
	SvmName string `structs:"svmName"`
	// more parameters might be added in the future
}

// gcpEncryptionParameters the input for requesting a CVO
type gcpEncryptionParameters struct {
	Key string `structs:"key,omitempty"`
}

// haParamsGCP the input for requesting a CVO
type haParamsGCP struct {
	PlatformSerialNumberNode1      string `structs:"platformSerialNumberNode1,omitempty"`
	PlatformSerialNumberNode2      string `structs:"platformSerialNumberNode2,omitempty"`
	Node1Zone                      string `structs:"node1Zone,omitempty"`
	Node2Zone                      string `structs:"node2Zone,omitempty"`
	MediatorZone                   string `structs:"mediatorZone,omitempty"`
	VPC0NodeAndDataConnectivity    string `structs:"vpc0NodeAndDataConnectivity,omitempty"`
	VPC1ClusterConnectivity        string `structs:"vpc1ClusterConnectivity,omitempty"`
	VPC2HAConnectivity             string `structs:"vpc2HAConnectivity,omitempty"`
	VPC3DataReplication            string `structs:"vpc3DataReplication,omitempty"`
	Subnet0NodeAndDataConnectivity string `structs:"subnet0NodeAndDataConnectivity,omitempty"`
	Subnet1ClusterConnectivity     string `structs:"subnet1ClusterConnectivity,omitempty"`
	Subnet2HAConnectivity          string `structs:"subnet2HAConnectivity,omitempty"`
	Subnet3DataReplication         string `structs:"subnet3DataReplication,omitempty"`
	VPC0FirewallRuleName           string `structs:"vpc0FirewallRuleName,omitempty"`
	VPC1FirewallRuleName           string `structs:"vpc1FirewallRuleName,omitempty"`
	VPC2FirewallRuleName           string `structs:"vpc2FirewallRuleName,omitempty"`
	VPC3FirewallRuleName           string `structs:"vpc3FirewallRuleName,omitempty"`
}

func (c *Client) createCVOGCP(cvoDetails createCVOGCPDetails, clientID string) (cvoResult, error) {
	log.Printf("\n\ncreateCVO %s client_id %s", cvoDetails.Name, clientID)

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createCVO request, failed to get AccessToken")
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

	if cvoDetails.NssAccount == "" && (cvoDetails.VsaMetadata.LicenseType == "gcp-cot-premium-byol" || cvoDetails.VsaMetadata.LicenseType == "gcp-ha-cot-premium-byol") && !strings.HasPrefix(cvoDetails.SerialNumber, "Eval-") {
		nssAccount, err := c.getNSS(clientID)
		if err != nil {
			log.Print("getNSS request failed ", err)
			return cvoResult{}, err
		}
		log.Printf("%s %s getNSS result %s", cvoDetails.Name, clientID, nssAccount)
		cvoDetails.NssAccount = nssAccount
	}

	baseURL := getAPIRootForWorkingEnvironment(cvoDetails.IsHA, "")

	hostType := "CloudManagerHost"
	params := structs.Map(cvoDetails)

	log.Printf("Create GCP CVO: %#v", params)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createCVO request failed ", statusCode)
		return cvoResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createCVO")
	if responseError != nil {
		return cvoResult{}, responseError
	}

	log.Print("retries ", c.Retries)
	var CreationRetries int
	if !cvoDetails.IsHA {
		CreationRetries = c.Retries
	} else {
		CreationRetries = c.Retries + 30
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "create", CreationRetries, 60, clientID)
	if err != nil {
		return cvoResult{}, err
	}

	var result cvoResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Printf("Response: %#v", response)
		log.Print("Failed to unmarshall response from createCVO ", err)
		return cvoResult{}, err
	}

	return result, nil
}

func (c *Client) deleteCVOGCP(id string, isHA bool, clientID string) error {

	log.Printf("deleteCVO: id %s client %s", id, clientID)

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteCVO request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := getAPIRootForWorkingEnvironment(isHA, id)

	hostType := "CloudManagerHost"

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("deleteCVO %s request failed %#v", id, statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteCVO")
	if responseError != nil {
		return responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "delete", 40, 60, clientID)
	if err != nil {
		return err
	}

	return nil
}

// This is used on GCP CVO HA only
func (c *Client) addSVMtoCVO(id string, clientID string, svmName string) error {
	log.Printf("addSVMtoCVO: id %s client %s svm %s", id, clientID, svmName)

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("In addSVMtoCVO request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	// GCP CVO SVM add and deletion only support HA
	baseURL := getAPIRootForWorkingEnvironment(true, id) + "/svm"
	hostType := "CloudManagerHost"

	var svm gcpSVM
	svm.SvmName = svmName
	params := structs.Map(svm)
	log.Printf("\taddSVMtoCVO params: %#v", params)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("addSVMtoCVO request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "addSVMtoCVO")
	if responseError != nil {
		return responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO_SVM", "add", 60, 60, clientID)

	return err
}

func (c *Client) deleteSVMfromCVO(id string, clientID string, svmName string) error {
	log.Printf("deleteSVMfromCVO: id %s client %s svm %s", id, clientID, svmName)

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("In deleteSVMfromCVO request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	// GCP CVO SVM add and deletion only support HA
	baseURL := getAPIRootForWorkingEnvironment(true, id)
	baseURL = fmt.Sprintf("%s/svm/%s", baseURL, svmName)
	log.Print("\tDelete svm url: ", baseURL)
	hostType := "CloudManagerHost"

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("deleteSVMfromCVO %s request failed %#v", id, statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteSVMfromCVO")
	if responseError != nil {
		return responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO_SVM", "delete", 40, 60, clientID)

	return err
}

// expandGCPLabels converts set to gcpLabels struct
func expandGCPLabels(set *schema.Set) []gcpLabels {
	labels := []gcpLabels{}

	for _, v := range set.List() {
		label := v.(map[string]interface{})
		gcpLabel := gcpLabels{}
		gcpLabel.LabelKey = label["label_key"].(string)
		gcpLabel.LabelValue = label["label_value"].(string)
		labels = append(labels, gcpLabel)
	}
	return labels
}

// expandGCPLabels converts set to gcpLabels struct
func expandGCPSVMs(set *schema.Set) []gcpSVM {
	svms := []gcpSVM{}

	for _, v := range set.List() {
		svm := v.(map[string]interface{})
		gcpSVM := gcpSVM{}
		gcpSVM.SvmName = svm["svm_name"].(string)
		svms = append(svms, gcpSVM)
	}
	return svms
}

// expandGCPLabelsToUseTags
func expandGCPLabelsToUserTags(set *schema.Set) []userTags {
	tags := []userTags{}

	for _, v := range set.List() {
		label := v.(map[string]interface{})
		userTag := userTags{}
		userTag.TagKey = label["label_key"].(string)
		userTag.TagValue = label["label_value"].(string)
		tags = append(tags, userTag)
	}
	return tags
}

// validateCVOGCPParams validates params
func validateCVOGCPParams(cvoDetails createCVOGCPDetails) error {
	if cvoDetails.VsaMetadata.UseLatestVersion && cvoDetails.VsaMetadata.OntapVersion != "latest" {
		return fmt.Errorf("ontap_version parameter not required when having use_latest_version as true")
	}

	// by Node byol license
	if cvoDetails.IsHA && cvoDetails.VsaMetadata.LicenseType == "gcp-ha-cot-premium-byol" {
		if cvoDetails.HAParams.PlatformSerialNumberNode1 == "" || cvoDetails.HAParams.PlatformSerialNumberNode2 == "" {
			return fmt.Errorf("both platform_serial_number_node1 and platform_serial_number_node2 parameters are required when having ha type as true and license_type as gcp-ha-cot-premium-byol")
		}
	}
	// by Node byol license
	if !cvoDetails.IsHA && (cvoDetails.HAParams.PlatformSerialNumberNode1 != "" || cvoDetails.HAParams.PlatformSerialNumberNode2 != "") {
		return fmt.Errorf("both platform_serial_number_node1 and platform_serial_number_node2 parameters are only required when having ha type as true and license_type as gcp-ha-cot-premium-byol")
	}

	// by Node byol license
	if cvoDetails.VsaMetadata.LicenseType == "gcp-cot-premium-byol" {
		if cvoDetails.SerialNumber == "" {
			return fmt.Errorf("serial_number parameter is required when having license_type as gcp-cot-premium-byol")
		}
	}

	// by Capacity license
	if cvoDetails.VsaMetadata.CapacityPackageName != "" {
		if cvoDetails.IsHA && cvoDetails.VsaMetadata.LicenseType != "ha-capacity-paygo" {
			return fmt.Errorf("license_type must be ha-capacity-paygo")
		}
		if !cvoDetails.IsHA && cvoDetails.VsaMetadata.LicenseType != "capacity-paygo" {
			return fmt.Errorf("license_type must be capacity-paygo")
		}
	}

	if strings.HasSuffix(cvoDetails.VsaMetadata.LicenseType, "capacity-paygo") && cvoDetails.VsaMetadata.CapacityPackageName == "" {
		return fmt.Errorf("capacity_package_name is required with capacity based license type")
	}
	return nil
}
