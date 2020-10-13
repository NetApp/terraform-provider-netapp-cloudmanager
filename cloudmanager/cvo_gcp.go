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
var GCPLicenseTypes = []string{"gcp-cot-standard-paygo", "gcp-cot-explore-paygo", "gcp-cot-premium-paygo", "gcp-cot-premium-byol"}

// createCVOGCPDetails the users input for creating a CVO
type createCVOGCPDetails struct {
	Name               string      `structs:"name"`
	DataEncryptionType string      `structs:"dataEncryptionType"`
	WorkspaceID        string      `structs:"tenantId,omitempty"`
	Region             string      `structs:"region"`
	GCPServiceAccount  string      `structs:"gcpServiceAccount"`
	VpcID              string      `structs:"vpcId"`
	SvmPassword        string      `structs:"svmPassword"`
	VsaMetadata        vsaMetadata `structs:"vsaMetadata"`
	GCPVolumeSize      diskSize    `structs:"gcpVolumeSize"`
	GCPVolumeType      string      `structs:"gcpVolumeType"`
	SubnetID           string      `structs:"subnetId"`
	SubnetPath         string      `structs:"subnetPath"`
	Project            string      `structs:"project"`
	CapacityTier       string      `structs:"capacityTier"`
	TierLevel          string      `structs:"tierLevel"`
	NssAccount         string      `structs:"nssAccount,omitempty"`
	WritingSpeedState  string      `structs:"writingSpeedState,omitempty"`
	SerialNumber       string      `structs:"serialNumber,omitempty"`
	GCPLabels          []gcpLabels `structs:"gcpLabels,omitempty"`
	FirewallRule       string      `structs:"firewallRule,omitempty"`
}

// gcpLabels the input for requesting a CVO
type gcpLabels struct {
	LabelKey   string `structs:"labelKey"`
	LabelValue string `structs:"labelValue,omitempty"`
}

func (c *Client) createCVOGCP(cvoDetails createCVOGCPDetails) (cvoResult, error) {

	log.Print("createCVO")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createCVO request, failed to get AccessToken")
		return cvoResult{}, err
	}
	c.Token = accessTokenResult.Token

	if cvoDetails.WorkspaceID == "" {
		tenantID, err := c.getTenant()
		if err != nil {
			log.Print("getTenant request failed ", err)
			return cvoResult{}, err
		}
		log.Print("tenant result ", tenantID)
		cvoDetails.WorkspaceID = tenantID
	}

	if cvoDetails.NssAccount == "" && cvoDetails.VsaMetadata.LicenseType == "gcp-cot-premium-byol" && !strings.HasPrefix(cvoDetails.SerialNumber, "Eval-") {
		nssAccount, err := c.getNSS()
		if err != nil {
			log.Print("getNSS request failed ", err)
			return cvoResult{}, err
		}
		log.Print("getNSS result ", nssAccount)
		cvoDetails.NssAccount = nssAccount
	}

	baseURL := "/occm/api/gcp/vsa/working-environments"

	hostType := "CloudManagerHost"
	params := structs.Map(cvoDetails)

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType)
	if err != nil {
		log.Print("createCVO request failed ", statusCode)
		return cvoResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createCVO")
	if responseError != nil {
		return cvoResult{}, responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "create", 60, 60)
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

func (c *Client) deleteCVOGCP(id string) error {

	log.Print("deleteCVO")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteCVO request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/occm/api/gcp/vsa/working-environments/%s", id)

	hostType := "CloudManagerHost"

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("deleteCVO request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteCVO")
	if responseError != nil {
		return responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "delete", 40, 60)
	if err != nil {
		return err
	}

	return nil
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

// validateCVOGCPParams validates params
func validateCVOGCPParams(cvoDetails createCVOGCPDetails) error {
	if cvoDetails.VsaMetadata.UseLatestVersion == true && cvoDetails.VsaMetadata.OntapVersion != "latest" {
		return fmt.Errorf("ontap_version parameter not required when having use_latest_version as true")
	}

	if cvoDetails.VsaMetadata.LicenseType == "gcp-cot-premium-byol" {
		if cvoDetails.SerialNumber == "" {
			return fmt.Errorf("serial_number parameter is required when having license_type as gcp-cot-premium-byol")
		}
	}

	return nil
}
