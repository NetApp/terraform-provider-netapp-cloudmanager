package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/structs"
)

// AWSLicenseTypes is the AWS License types
var AWSLicenseTypes = []string{"cot-standard-paygo", "cot-premium-paygo", "cot-explore-paygo", "cot-premium-byol",
	"ha-cot-standard-paygo", "ha-cot-premium-paygo", "ha-cot-premium-byol", "ha-cot-explore-paygo", "capacity-paygo",
	"ha-capacity-paygo"}

// createCVOAWSDetails the users input for creating a CVO
type createCVOAWSDetails struct {
	Name                        string                  `structs:"name"`
	DataEncryptionType          string                  `structs:"dataEncryptionType"`
	WorkspaceID                 string                  `structs:"tenantId,omitempty"`
	Region                      string                  `structs:"region"`
	VpcID                       string                  `structs:"vpcId,omitempty"`
	SvmPassword                 string                  `structs:"svmPassword"`
	SvmName                     string                  `structs:"svmName,omitempty"`
	VsaMetadata                 vsaMetadata             `structs:"vsaMetadata"`
	EbsVolumeSize               ebsVolumeSize           `structs:"ebsVolumeSize"`
	EbsVolumeType               string                  `structs:"ebsVolumeType"`
	SubnetID                    string                  `structs:"subnetId"`
	CapacityTier                string                  `structs:"capacityTier,omitempty"`
	TierLevel                   string                  `structs:"tierLevel,omitempty"`
	NssAccount                  string                  `structs:"nssAccount,omitempty"`
	WritingSpeedState           string                  `structs:"writingSpeedState,omitempty"`
	ClusterKeyPairName          string                  `structs:"clusterKeyPairName,omitempty"`
	IOPS                        int                     `structs:"iops,omitempty"`
	Throughput                  int                     `structs:"throughput,omitempty"`
	OptimizedNetworkUtilization bool                    `structs:"optimizedNetworkUtilization"`
	InstanceTenancy             string                  `structs:"instanceTenancy"`
	InstanceProfileName         string                  `structs:"instanceProfileName,omitempty"`
	SecurityGroupID             string                  `structs:"securityGroupId,omitempty"`
	WormRequest                 wormRequest             `structs:"wormRequest,omitempty"`
	CloudProviderAccount        string                  `structs:"cloudProviderAccount,omitempty"`
	BackupVolumesToCbs          bool                    `structs:"backupVolumesToCbs"`
	EnableCompliance            bool                    `structs:"enableCompliance"`
	EnableMonitoring            bool                    `structs:"enableMonitoring"`
	AwsEncryptionParameters     awsEncryptionParameters `structs:"awsEncryptionParameters,omitempty"`
	AwsTags                     []userTags              `structs:"awsTags,omitempty"`
	IsHA                        bool
	HAParams                    haParamsAWS `structs:"haParams,omitempty"`
}

// haParamsAWS the input for requesting a CVO
type haParamsAWS struct {
	ClusterFloatingIP           string   `structs:"clusterFloatingIP,omitempty"`
	DataFloatingIP              string   `structs:"dataFloatingIP,omitempty"`
	DataFloatingIP2             string   `structs:"dataFloatingIP2,omitempty"`
	SvmFloatingIP               string   `structs:"svmFloatingIP,omitempty"`
	FailoverMode                string   `structs:"failoverMode,omitempty"`
	Node1SubnetID               string   `structs:"node1SubnetId,omitempty"`
	Node2SubnetID               string   `structs:"node2SubnetId,omitempty"`
	MediatorSubnetID            string   `structs:"mediatorSubnetId,omitempty"`
	MediatorKeyPairName         string   `structs:"mediatorKeyPairName,omitempty"`
	PlatformSerialNumberNode1   string   `structs:"platformSerialNumberNode1,omitempty"`
	PlatformSerialNumberNode2   string   `structs:"platformSerialNumberNode2,omitempty"`
	MediatorInstanceProfileName string   `structs:"mediatorInstanceProfileName,omitempty"`
	MediatorAssignPublicIP      bool     `structs:"mediatorAssignPublicIP,omitempty"`
	MediatorSecurityGroupID     string   `structs:"mediatorSecurityGroupId,omitempty"`
	RouteTableIds               []string `structs:"routeTableIds,omitempty"`
	AssumeRoleArn               string   `structs:"assumeRoleArn,omitempty"`
}

// ebsVolumeSize the input for requesting a CVO
type ebsVolumeSize struct {
	Size int    `structs:"size"`
	Unit string `structs:"unit"`
}

// vsaMetadata the input for requesting a CVO
type vsaMetadata struct {
	OntapVersion         string `structs:"ontapVersion"`
	UseLatestVersion     bool   `structs:"useLatestVersion"`
	LicenseType          string `structs:"licenseType"`
	InstanceType         string `structs:"instanceType,omitempty"`
	PlatformSerialNumber string `structs:"platformSerialNumber,omitempty"`
	CapacityPackageName  string `structs:"capacityPackageName,omitempty"`
	ProvidedLicense      string `structs:"providedLicense,omitempty"`
}

// awsEncryptionParameters the input for requesting a CVO
type awsEncryptionParameters struct {
	KmsKeyID  string `structs:"kmsKeyId,omitempty"`
	KmsKeyArn string `structs:"kmsKeyArn,omitempty"`
}

// cvoList the users input for getting cvo
type cvoList struct {
	CVO []cvoResult `json:"vsaWorkingEnvironments"`
}

// cvoResult the users input for creating a cvo
type cvoResult struct {
	PublicID string `json:"publicId"`
	SvmName  string `json:"svmName"`
}

// tenantResult the users input for creating a cvo
type tenantResult struct {
	PublicID string `json:"publicId"`
}

// cvoStatusResult the users input for creating a cvo
type cvoStatusResult struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

// wormRequest
type wormRequest struct {
	RetentionPeriod wormRetentionPeriod `structs:"retentionPeriod"`
}

// wormRetentionPeriod
type wormRetentionPeriod struct {
	Length int    `structs:"length"`
	Unit   string `structs:"unit"`
}

func (c *Client) getTenant(clientID string) (string, error) {

	log.Print("getTenant client=", clientID)

	baseURL := "/occm/api/tenants"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getTenant request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getTenant")
	if responseError != nil {
		return "", responseError
	}

	var result []tenantResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getTenant ", err)
		return "", err
	}

	return result[0].PublicID, nil
}

func (c *Client) getCVOAWS(id string, clientID string) (string, error) {

	log.Print("getCVOAWS")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCVOAWS request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := "/occm/api/working-environments"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getCVOAWS request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getCVOAWS")
	if responseError != nil {
		return "", responseError
	}

	var result cvoList
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getCVOAWS ", err)
		return "", err
	}

	for _, cvoID := range result.CVO {
		if cvoID.PublicID == id {
			return cvoID.PublicID, nil
		}
	}

	return "", nil
}

func (c *Client) createCVOAWS(cvoDetails createCVOAWSDetails, clientID string) (cvoResult, error) {

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

	if cvoDetails.NssAccount == "" {
		if cvoDetails.VsaMetadata.PlatformSerialNumber != "" {
			if !strings.HasPrefix(cvoDetails.VsaMetadata.PlatformSerialNumber, "Eval-") && cvoDetails.VsaMetadata.LicenseType == "cot-premium-byol" {
				nssAccount, err := c.getNSS(clientID)
				if err != nil {
					log.Print("getNSS request failed ", err)
					return cvoResult{}, err
				}
				log.Print("getNSS result ", nssAccount)
				cvoDetails.NssAccount = nssAccount
			}
		} else if cvoDetails.HAParams.PlatformSerialNumberNode1 != "" && cvoDetails.HAParams.PlatformSerialNumberNode2 != "" {
			if !strings.HasPrefix(cvoDetails.HAParams.PlatformSerialNumberNode1, "Eval-") && !strings.HasPrefix(cvoDetails.HAParams.PlatformSerialNumberNode2, "Eval-") && cvoDetails.VsaMetadata.LicenseType == "ha-cot-premium-byol" {
				nssAccount, err := c.getNSS(clientID)
				if err != nil {
					log.Print("getNSS request failed ", err)
					return cvoResult{}, err
				}
				log.Print("getNSS result ", nssAccount)
				cvoDetails.NssAccount = nssAccount
			}
		}
	}

	if cvoDetails.VpcID == "" {
		vpcID, err := c.CallVPCGet(cvoDetails.SubnetID, cvoDetails.Region)
		if err != nil {
			log.Print("CallVPCGet request failed")
			return cvoResult{}, err
		}
		log.Print("vpcID result ", vpcID)
		cvoDetails.VpcID = vpcID
	}

	var baseURL string
	var CreationRetries int

	log.Print("retries ", c.Retries)
	if !cvoDetails.IsHA {
		baseURL = "/occm/api/vsa/working-environments"
		CreationRetries = c.Retries
	} else if cvoDetails.IsHA {
		baseURL = "/occm/api/aws/ha/working-environments"
		CreationRetries = c.Retries + 30
	}

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

func (c *Client) deleteCVO(id string, isHA bool, clientID string) error {

	log.Print("deleteCVO")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteCVO request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	var baseURL string

	if !isHA {
		baseURL = fmt.Sprintf("/occm/api/vsa/working-environments/%s", id)
	} else if isHA {
		baseURL = fmt.Sprintf("/occm/api/aws/ha/working-environments/%s", id)
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
func validateCVOParams(cvoDetails createCVOAWSDetails) error {
	if cvoDetails.VsaMetadata.UseLatestVersion == true && cvoDetails.VsaMetadata.OntapVersion != "latest" {
		return fmt.Errorf("ontap_version parameter not required when having use_latest_version as true")
	}

	// by Node byol license for existing customers
	if cvoDetails.VsaMetadata.PlatformSerialNumber != "" && cvoDetails.VsaMetadata.LicenseType != "cot-premium-byol" {
		return fmt.Errorf("platform_serial_number parameter required only when having license_type as cot-premium-byol")
	}

	// by Node byol license for existing customers
	if cvoDetails.IsHA == true && cvoDetails.VsaMetadata.LicenseType == "ha-cot-premium-byol" {
		if cvoDetails.HAParams.PlatformSerialNumberNode1 == "" || cvoDetails.HAParams.PlatformSerialNumberNode2 == "" {
			return fmt.Errorf("both platform_serial_number_node1 and platform_serial_number_node2 parameters are required when having ha type as true and license_type as ha-cot-premium-byol")
		}
	}

	// by Node byol license for existing customers
	if cvoDetails.IsHA == false && (cvoDetails.HAParams.PlatformSerialNumberNode1 != "" || cvoDetails.HAParams.PlatformSerialNumberNode2 != "") {
		return fmt.Errorf("both platform_serial_number_node1 and platform_serial_number_node2 parameters are only required when having ha type as true and license_type as ha-cot-premium-byol")
	}

	if (cvoDetails.EbsVolumeType == "io1" || cvoDetails.EbsVolumeType == "gp3") && cvoDetails.IOPS == 0 {
		return fmt.Errorf("iops parameter is required when having ebs_volume_type as io1 or gp3")
	}

	if cvoDetails.IOPS != 0 && cvoDetails.EbsVolumeType != "io1" && cvoDetails.EbsVolumeType != "gp3" {
		return fmt.Errorf("iops parameter is not supported when ebs_volume_type is not io1 or gp3")
	}

	if cvoDetails.EbsVolumeType == "gp3" && cvoDetails.Throughput == 0 {
		return fmt.Errorf("throughput parameter required when ebs_volume_type is gp3")
	}

	if cvoDetails.IsHA == true && cvoDetails.SubnetID != "" {
		return fmt.Errorf("subnet_id not required when having ha as true")
	}

	// by Capacity license
	if cvoDetails.VsaMetadata.CapacityPackageName != "" {
		log.Print("Verify cvo parameter capacity_package_name is not empty")
		if cvoDetails.IsHA == true && cvoDetails.VsaMetadata.LicenseType != "ha-capacity-paygo" {
			return fmt.Errorf("license_type must be ha-capacity-paygo")
		}
		if cvoDetails.IsHA == false && cvoDetails.VsaMetadata.LicenseType != "capacity-paygo" {
			return fmt.Errorf("license_type must be capacity-paygo")
		}
	}

	if strings.HasSuffix(cvoDetails.VsaMetadata.LicenseType, "capacity-paygo") && cvoDetails.VsaMetadata.CapacityPackageName == "" {
		return fmt.Errorf("capacity_package_name is required with capacity based license type")
	}

	return nil
}
