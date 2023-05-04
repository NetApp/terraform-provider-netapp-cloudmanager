package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/structs"
)

// AzureLicenseTypes is the Azure License types
var AzureLicenseTypes = []string{"azure-cot-standard-paygo", "azure-cot-premium-paygo", "azure-cot-premium-byol",
	"azure-cot-explore-paygo", "azure-ha-cot-standard-paygo", "azure-ha-cot-premium-paygo", "azure-ha-cot-premium-byol",
	"capacity-paygo", "ha-capacity-paygo"}

// createCVOAzureDetails the users input for creating a CVO
type createCVOAzureDetails struct {
	Name                        string                    `structs:"name"`
	DataEncryptionType          string                    `structs:"dataEncryptionType"`
	WorkspaceID                 string                    `structs:"tenantId,omitempty"`
	Region                      string                    `structs:"region"`
	SubscriptionID              string                    `structs:"subscriptionId"`
	SaasSubscriptionID          string                    `structs:"saasSubscriptionId"`
	VnetID                      string                    `structs:"vnetId,omitempty"`
	SvmPassword                 string                    `structs:"svmPassword"`
	SvmName                     string                    `structs:"svmName,omitempty"`
	VsaMetadata                 vsaMetadata               `structs:"vsaMetadata"`
	DiskSize                    diskSize                  `structs:"diskSize"`
	StorageType                 string                    `structs:"storageType"`
	SubnetID                    string                    `structs:"subnetId"`
	Cidr                        string                    `structs:"cidr"`
	CapacityTier                string                    `structs:"capacityTier,omitempty"`
	TierLevel                   string                    `structs:"tierLevel,omitempty"`
	AvailabilityZone            int                       `structs:"availabilityZone,omitempty"`
	AzureEncryptionParameters   azureEncryptionParameters `structs:"azureEncryptionParameters,omitempty"`
	NssAccount                  string                    `structs:"nssAccount,omitempty"`
	WormRequest                 wormRequest               `structs:"wormRequest,omitempty"`
	WritingSpeedState           string                    `structs:"writingSpeedState,omitempty"`
	OptimizedNetworkUtilization bool                      `structs:"optimizedNetworkUtilization"`
	SecurityGroupID             string                    `structs:"securityGroupId,omitempty"`
	CloudProviderAccount        string                    `structs:"cloudProviderAccount,omitempty"`
	BackupVolumesToCbs          bool                      `structs:"backupVolumesToCbs"`
	EnableCompliance            bool                      `structs:"enableCompliance"`
	EnableMonitoring            bool                      `structs:"enableMonitoring"`
	AllowDeployInExistingRg     bool                      `structs:"allowDeployInExistingRg,omitempty"`
	AzureTags                   []userTags                `structs:"azureTags,omitempty"`
	IsHA                        bool
	ResourceGroup               string `structs:"resourceGroup,omitempty"`
	VnetResourceGroup           string
	VnetForInternal             string
	SerialNumber                string        `structs:"serialNumber,omitempty"`
	HAParams                    haParamsAzure `structs:"haParams,omitempty"`
}

type azureEncryptionParameters struct {
	Key                  string `structs:"key"`
	VaultName            string `structs:"vaultName"`
	UserAssignedIdentity string `structs:"userAssignedIdentity,omitempty"`
}

// haParamsAzure the input for requesting a CVO
type haParamsAzure struct {
	PlatformSerialNumberNode1 string `structs:"platformSerialNumberNode1,omitempty"`
	PlatformSerialNumberNode2 string `structs:"platformSerialNumberNode2,omitempty"`
	AvailabilityZoneNode1     int    `structs:"availabilityZoneNode1,omitempty"`
	AvailabilityZoneNode2     int    `structs:"availabilityZoneNode2,omitempty"`
	EnableHTTPS               bool   `structs:"enableHttps"`
}

// cvoListAzure the users input for getting cvo
type cvoListAzure struct {
	CVO []cvoResult `json:"azureVsaWorkingEnvironments"`
}

// accountForNSSResult the users input for creating a cvo
type accountForNSSResult struct {
	NssAccounts []nssAccountResult `json:"nssAccounts"`
}

// nssAccountResult the users input for creating a cvo
type nssAccountResult struct {
	PublicID string `json:"publicId"`
}

func (c *Client) getNSS(clientID string) (string, error) {

	log.Print("getNSS")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createCVO request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := "/occm/api/accounts"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getNSS request failed ", statusCode)
		return "", err
	}

	log.Print("getNSS ")
	log.Print(string(response))

	responseError := apiResponseChecker(statusCode, response, "getNSS")
	if responseError != nil {
		return "", responseError
	}

	var result accountForNSSResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getNSS ", err)
		return "", err
	}

	log.Print("getNSS ", result)

	if len(result.NssAccounts) == 0 {
		return "", fmt.Errorf("could not find any NSS account")
	}

	return result.NssAccounts[0].PublicID, nil
}

func (c *Client) getCVOAzure(id string, clientID string) (string, error) {

	log.Print("getCVOAzure")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCVOAzure request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := "/occm/api/working-environments"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getCVOAzure request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getCVOAzure")
	if responseError != nil {
		return "", responseError
	}

	log.Print(string(response))

	var result cvoListAzure
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getCVOAzure ", err)
		return "", err
	}

	for _, cvoID := range result.CVO {
		if cvoID.PublicID == id {
			return cvoID.PublicID, nil
		}
	}

	return "", nil
}

func (c *Client) createCVOAzure(cvoDetails createCVOAzureDetails, clientID string) (cvoResult, error) {

	log.Print("createCVO")

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

	if cvoDetails.NssAccount == "" && !strings.HasPrefix(cvoDetails.SerialNumber, "Eval-") && (cvoDetails.VsaMetadata.LicenseType == "azure-cot-premium-byol" || cvoDetails.VsaMetadata.LicenseType == "azure-ha-cot-premium-byol") {
		nssAccount, err := c.getNSS(clientID)
		if err != nil {
			log.Print("getNSS request failed ", err)
			return cvoResult{}, err
		}
		log.Print("getNSS result ", nssAccount)
		cvoDetails.NssAccount = nssAccount
	}

	if cvoDetails.Cidr == "" {
		var rg string
		if cvoDetails.VnetResourceGroup != "" {
			rg = cvoDetails.VnetResourceGroup
		} else {
			rg = cvoDetails.ResourceGroup
		}
		cidr, err := c.CallVNetGetCidr(cvoDetails.SubscriptionID, rg, cvoDetails.VnetForInternal)
		if err != nil {
			log.Print("CallVNetGetCidr request failed")
			return cvoResult{}, err
		}
		cvoDetails.Cidr = cidr
		log.Print("cidr result ", cvoDetails.Cidr)
	}

	var baseURL string
	var CreationRetries int

	log.Print("retries ", c.Retries)
	if !cvoDetails.IsHA {
		baseURL = "/occm/api/azure/vsa/working-environments"
		CreationRetries = c.Retries
	} else if cvoDetails.IsHA {
		baseURL = "/occm/api/azure/ha/working-environments"
		CreationRetries = c.Retries + 30
	}

	log.Print(baseURL)

	hostType := "CloudManagerHost"
	params := structs.Map(cvoDetails)

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createCVO request failed ", statusCode)
		return cvoResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createCVO")
	if responseError != nil {
		return cvoResult{}, responseError
	}

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "create", CreationRetries, 60, clientID)
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

func (c *Client) deleteCVOAzure(id string, isHA bool, clientID string) error {

	log.Print("deleteCVO")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteCVO request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	var baseURL string

	if !isHA {
		baseURL = fmt.Sprintf("/occm/api/azure/vsa/working-environments/%s", id)
	} else if isHA {
		baseURL = fmt.Sprintf("/occm/api/azure/ha/working-environments/%s", id)
	}

	hostType := "CloudManagerHost"

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("deleteCVO request failed ", statusCode)
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

// validateCVOParams validates params
func validateCVOAzureParams(cvoDetails createCVOAzureDetails) error {
	if cvoDetails.VsaMetadata.UseLatestVersion && cvoDetails.VsaMetadata.OntapVersion != "latest" {
		return fmt.Errorf("ontap_version parameter not required when having use_latest_version as true")
	}

	// by Node byol license
	if cvoDetails.IsHA && cvoDetails.VsaMetadata.LicenseType == "azure-ha-cot-premium-byol" {
		if cvoDetails.HAParams.PlatformSerialNumberNode1 == "" || cvoDetails.HAParams.PlatformSerialNumberNode2 == "" {
			return fmt.Errorf("both platform_serial_number_node1 and platform_serial_number_node2 parameters are required when having ha type as true and license_type as azure-ha-cot-premium-byol")
		}
	}

	if !cvoDetails.IsHA && (cvoDetails.HAParams.PlatformSerialNumberNode1 != "" || cvoDetails.HAParams.PlatformSerialNumberNode2 != "") {
		return fmt.Errorf("both platform_serial_number_node1 and platform_serial_number_node2 parameters are not required when having ha type as false")
	}

	// by Node byol license
	if cvoDetails.VsaMetadata.LicenseType == "azure-cot-premium-byol" {
		if cvoDetails.SerialNumber == "" {
			return fmt.Errorf("serial_number parameter is required when having license_type as azure-cot-premium-byol")
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
