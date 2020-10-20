package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

// AWSLicenseTypes is the AWS License types
var AWSLicenseTypes = []string{"cot-standard-paygo", "cot-premium-paygo", "cot-explore-paygo", "cot-premium-byol", "ha-cot-standard-paygo", "ha-cot-premium-paygo", "ha-cot-premium-byol", "ha-cot-explore-paygo"}

// createCVOAWSDetails the users input for creating a CVO
type createCVOAWSDetails struct {
	Name                        string                  `structs:"name"`
	DataEncryptionType          string                  `structs:"dataEncryptionType"`
	WorkspaceID                 string                  `structs:"tenantId,omitempty"`
	Region                      string                  `structs:"region"`
	VpcID                       string                  `structs:"vpcId,omitempty"`
	SvmPassword                 string                  `structs:"svmPassword"`
	VsaMetadata                 vsaMetadata             `structs:"vsaMetadata"`
	EbsVolumeSize               ebsVolumeSize           `structs:"ebsVolumeSize"`
	EbsVolumeType               string                  `structs:"ebsVolumeType"`
	SubnetID                    string                  `structs:"subnetId"`
	CapacityTier                string                  `structs:"capacityTier,omitempty"`
	TierLevel                   string                  `structs:"tierLevel,omitempty"`
	NssAccount                  string                  `structs:"nssAccount,omitempty"`
	WritingSpeedState           string                  `structs:"writingSpeedState,omitempty"`
	IOPS                        int                     `structs:"iops,omitempty"`
	OptimizedNetworkUtilization bool                    `structs:"optimizedNetworkUtilization"`
	InstanceTenancy             string                  `structs:"instanceTenancy"`
	InstanceProfileName         string                  `structs:"instanceProfileName,omitempty"`
	SecurityGroupID             string                  `structs:"securityGroupId,omitempty"`
	CloudProviderAccount        string                  `structs:"cloudProviderAccount,omitempty"`
	BackupVolumesToCbs          bool                    `structs:"backupVolumesToCbs"`
	EnableCompliance            bool                    `structs:"enableCompliance"`
	EnableMonitoring            bool                    `structs:"enableMonitoring"`
	AwsEncryptionParameters     awsEncryptionParameters `structs:"awsEncryptionParameters,omitempty"`
	AwsTags                     []awsTags               `structs:"awsTags,omitempty"`
	IsHA                        bool
	HAParams                    haParamsAWS `structs:"haParams,omitempty"`
}

// haParamsAWS the input for requesting a CVO
type haParamsAWS struct {
	ClusterFloatingIP      string   `structs:"clusterFloatingIP,omitempty"`
	DataFloatingIP         string   `structs:"dataFloatingIP,omitempty"`
	DataFloatingIP2        string   `structs:"dataFloatingIP2,omitempty"`
	SvmFloatingIP          string   `structs:"svmFloatingIP,omitempty"`
	FailoverMode           string   `structs:"failoverMode,omitempty"`
	Node1SubnetID          string   `structs:"node1SubnetId,omitempty"`
	Node2SubnetID          string   `structs:"node2SubnetId,omitempty"`
	MediatorSubnetID       string   `structs:"mediatorSubnetId,omitempty"`
	MediatorKeyPairName    string   `structs:"mediatorKeyPairName,omitempty"`
	MediatorAssignPublicIP bool     `structs:"mediatorAssignPublicIP,omitempty"`
	RouteTableIds          []string `structs:"routeTableIds,omitempty"`
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
}

// awsEncryptionParameters the input for requesting a CVO
type awsEncryptionParameters struct {
	KmsKeyID string `structs:"kmsKeyId,omitempty"`
}

// awsTags the input for requesting a CVO
type awsTags struct {
	TagKey   string `structs:"tagKey"`
	TagValue string `structs:"tagValue,omitempty"`
}

// deleteCVODetails the users input for deleting a cvo
type deleteCVODetails struct {
	InstanceID string
	Region     string
}

// cvoList the users input for getting cvo
type cvoList struct {
	CVO []cvoResult `json:"vsaWorkingEnvironments"`
}

// cvoResult the users input for creating a cvo
type cvoResult struct {
	PublicID string `json:"publicId"`
}

// tenantResult the users input for creating a cvo
type tenantResult struct {
	PublicID string `json:"publicId"`
}

// cvoStatusResult the users input for creating a cvo
type cvoStatusResult struct {
	Status int `json:"status"`
}

func (c *Client) getTenant() (string, error) {

	log.Print("getTenant")

	baseURL := "/occm/api/tenants"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
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

func (c *Client) getCVOAWSByID(id string) error {

	log.Print("getCVOAWSByID")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCVOAWSByID request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/occm/api/working-environments/%s", id)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Print("getCVOAWSByID request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "getCVOAWSByID")
	if responseError != nil {
		return responseError
	}

	return nil
}

func (c *Client) getCVOAWS(id string) (string, error) {

	log.Print("getCVOAWS")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCVOAWS request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := "/occm/api/working-environments"

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
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

func (c *Client) createCVOAWS(cvoDetails createCVOAWSDetails) (cvoResult, error) {

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

	if cvoDetails.NssAccount == "" && !strings.HasPrefix(cvoDetails.VsaMetadata.PlatformSerialNumber, "Eval-") && (cvoDetails.VsaMetadata.LicenseType == "cot-premium-byol" || cvoDetails.VsaMetadata.LicenseType == "ha-cot-premium-byol") {
		nssAccount, err := c.getNSS()
		if err != nil {
			log.Print("getNSS request failed ", err)
			return cvoResult{}, err
		}
		log.Print("getNSS result ", nssAccount)
		cvoDetails.NssAccount = nssAccount
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
	var creationWaitTime int

	if cvoDetails.IsHA == false {
		baseURL = "/occm/api/vsa/working-environments"
		creationWaitTime = 60
	} else if cvoDetails.IsHA == true {
		baseURL = "/occm/api/aws/ha/working-environments"
		creationWaitTime = 90
	}

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

	err = c.waitOnCompletion(onCloudRequestID, "CVO", "create", creationWaitTime, 60)
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

func (c *Client) deleteCVO(id string, isHA bool) error {

	log.Print("deleteCVO")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteCVO request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	var baseURL string

	if isHA == false {
		baseURL = fmt.Sprintf("/occm/api/vsa/working-environments/%s", id)
	} else if isHA == true {
		baseURL = fmt.Sprintf("/occm/api/aws/ha/working-environments/%s", id)
	}

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

// expandAWSTags converts set to awsEncryptionParameters struct
func expandAWSTags(set *schema.Set) []awsTags {
	tags := []awsTags{}

	for _, v := range set.List() {
		tag := v.(map[string]interface{})
		awsTag := awsTags{}
		awsTag.TagKey = tag["tag_key"].(string)
		awsTag.TagValue = tag["tag_value"].(string)
		tags = append(tags, awsTag)
	}
	return tags
}

// validateCVOParams validates params
func validateCVOParams(cvoDetails createCVOAWSDetails) error {
	if cvoDetails.VsaMetadata.UseLatestVersion == true && cvoDetails.VsaMetadata.OntapVersion != "latest" {
		return fmt.Errorf("ontap_version parameter not required when having use_latest_version as true")
	}

	if cvoDetails.VsaMetadata.PlatformSerialNumber != "" && cvoDetails.VsaMetadata.LicenseType != "cot-premium-byol" {
		return fmt.Errorf("platform_serial_number parameter required only when having license_type as cot-premium-byol")
	}

	if (cvoDetails.IOPS == 0 && cvoDetails.EbsVolumeType == "io1") || (cvoDetails.IOPS != 0 && cvoDetails.EbsVolumeType != "io1") {
		return fmt.Errorf("iops parameter required when having ebs_volume_type as io1")
	}
	return nil
}
