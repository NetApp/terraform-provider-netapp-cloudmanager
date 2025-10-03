package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// only list what is needed
type workingEnvironmentInfo struct {
	Name                   string      `json:"name"`
	PublicID               string      `json:"publicId"`
	CloudProviderName      string      `json:"cloudProviderName"`
	ProviderName           string      `json:"providerName"`
	IsHA                   bool        `json:"isHA"`
	WorkingEnvironmentType string      `json:"workingEnvironmentType"`
	SvmName                string      `json:"svmName"`
	Svms                   interface{} `json:"svms"`
}

type workingEnvironmentResult struct {
	VsaWorkingEnvironment       []workingEnvironmentInfo `json:"vsaWorkingEnvironments"`
	OnPremWorkingEnvironments   []workingEnvironmentInfo `json:"onPremWorkingEnvironments"`
	AzureVsaWorkingEnvironments []workingEnvironmentInfo `json:"azureVsaWorkingEnvironments"`
	GcpVsaWorkingEnvironments   []workingEnvironmentInfo `json:"gcpVsaWorkingEnvironments"`
}

type workingEnvironmentOntapClusterPropertiesResponse struct {
	ActionsRequired                interface{}            `json:"actionsRequired"`
	ActiveActions                  interface{}            `json:"activeActions"`
	AwsProperties                  interface{}            `json:"awsProperties"` // aws
	CapacityFeatures               interface{}            `json:"capacityFeatures"`
	CbsProperties                  interface{}            `json:"cbsProperties"`
	CloudSyncProperties            interface{}            `json:"cloudSyncProperties"` // aws
	CloudProviderName              string                 `json:"cloudProviderName"`
	ComplianceProperties           interface{}            `json:"complianceProperties"`
	CreatorUserEmail               string                 `json:"creatorUserEmail"`
	CronJobSchedules               interface{}            `json:"cronJobSchedules"` // aws
	EncryptionProperties           interface{}            `json:"encryptionProperties"`
	FpolicyProperties              interface{}            `json:"fpolicyProperties"`
	HAProperties                   haProperties           `json:"haProperties"`
	InterClusterLifs               interface{}            `json:"interClusterLifs"` // aws
	IsHA                           bool                   `json:"isHA"`
	LicensesInformation            interface{}            `json:"licensesInformation"`
	MonitoringProperties           interface{}            `json:"monitoringProperties"`
	Name                           string                 `json:"name"`
	OntapClusterProperties         ontapClusterProperties `json:"ontapClusterProperties"`
	ProviderProperties             providerProperties     `json:"providerProperties"`
	PublicID                       string                 `json:"publicId"`
	ReplicationProperties          interface{}            `json:"replicationProperties"`
	ReservedSize                   interface{}            `json:"reservedSize"`
	SaasProperties                 interface{}            `json:"saasProperties"`
	Schedules                      interface{}            `json:"schedules"`
	SnapshotPolicies               []cvoSnapshotPolicy    `json:"snapshotPolicies"`
	Status                         cvoStatus              `json:"status"`
	SupportRegistrationInformation []interface{}          `json:"supportRegistrationInformation"`
	SupportRegistrationProperties  interface{}            `json:"supportRegistrationProperties"`
	SupportedFeatures              interface{}            `json:"supportedFeatures"`
	SvmName                        string                 `json:"svmName"`
	Svms                           interface{}            `json:"svms"`
	TenantID                       string                 `json:"tenantId"`
	WorkingEnvironmentType         string                 `json:"workingEnvironmentType"`
}

type ontapClusterProperties struct {
	AggregateCount                   int                   `json:"aggregateCount"`
	VolumeCount                      int                   `json:"volumeCount"`
	FlashCache                       bool                  `json:"flashCache"`
	SpaceReportingLogical            bool                  `json:"spaceReportingLogical"`
	KeystoneSubscription             bool                  `json:"keystoneSubscription"`
	BroadcastDomainInfo              []broadcastDomainInfo `json:"broadcastDomainInfo"`
	CanConfigureCapacityTier         bool                  `json:"canConfigureCapacityTier"`
	CapacityTierInfo                 capacityTierInfo      `json:"capacityTierInfo"`
	ClusterName                      string                `json:"clusterName"`
	ClusterUUID                      string                `json:"clusterUuid"`
	CreationTime                     interface{}           `json:"creationTime"`
	Evaluation                       bool                  `json:"evaluation"`
	IsSpaceReportingLogical          bool                  `json:"isSpaceReportingLogical"`
	LastModifiedOffbox               interface{}           `json:"lastModifiedOffbox"`
	LicensePackageName               interface{}           `json:"licensePackageName"`
	LicenseType                      licenseType           `json:"licenseType"`
	Nodes                            []node                `json:"nodes"`
	OffboxTarget                     bool                  `json:"offboxTarget"`
	OntapVersion                     string                `json:"ontapVersion"`
	SystemManagerURL                 string                `json:"systemManagerUrl"`
	UpgradeVersions                  []upgradeVersion      `json:"upgradeVersions"`
	UsedCapacity                     capacityLimit         `json:"usedCapacity"`
	UserName                         string                `json:"userName"`
	VscanFileOperationDefaultProfile string                `json:"vscanFileOperationDefaultProfile"`
	WormEnabled                      bool                  `json:"wormEnabled"`
	WritingSpeedState                string                `json:"writingSpeedState"`
}

// common fields of GCP properties and Azure properties
type providerProperties struct {
	RegionName   string `json:"regionName"`
	InstanceType string `json:"instanceType"`
	NumOfNics    int    `json:"numOfNics"`
}

// type gcpProperties struct {
// 	Name           string      `json:"name"`
// 	RegionName     string      `json:"regionName"`
// 	ZoneName       []string    `json:"zoneName"`
// 	InstanceType   string      `json:"instanceType"`
// 	SubnetCidr     string      `json:"subnetCidr"`
// 	NumOfNics      int         `json:"numOfNics"`
// 	Labels         interface{} `json:"labels"`
// 	ProjectName    string      `json:"projectName"`
// 	DeploymentName string      `json:"deploymentName"`
// }

type haProperties struct {
	FailoverMode             interface{}   `json:"failoverMode"`
	MediatorStatus           interface{}   `json:"mediatorStatus"`
	MediatorVersionInfo      interface{}   `json:"mediatorVersionInfo"`
	MediatorVersionsToUpdate []interface{} `json:"mediatorVersionsToUpdate"`
	RouteTables              []string      `json:"routeTables"`
}

type broadcastDomainInfo struct {
	BroadcastDomain string `json:"broadcastDomain"`
	IPSpace         string `json:"ipSpace"`
	Mtu             int    `json:"mtu"`
}

type capacityTierInfo struct {
	CapacityTierUsedSize capacityLimit `json:"capacityTierUsedSize"`
	S3BucketName         string        `json:"s3BucketName"`
	TierLevel            string        `json:"tierLevel"`
}

type node struct {
	CloudProviderID      string      `json:"cloudProviderId"`
	Health               bool        `json:"health"`
	InTakeover           bool        `json:"inTakeover"`
	Lifs                 []lif       `json:"lifs"`
	Name                 string      `json:"name"`
	PlatformLicense      interface{} `json:"platformLicense"`
	PlatformSerialNumber interface{} `json:"platformSerialNumber"`
	SerialNumber         string      `json:"serialNumber"`
	SystemID             string      `json:"systemId"`
}

type upgradeVersion struct {
	ImageVersion      string `json:"imageVersion"`
	LastModified      int    `json:"lastModified"`
	AutoUpdateAllowed bool   `json:"autoUpdateAllowed"`
}

type licenseType struct {
	CapacityLimit capacityLimit `json:"capacityLimit"`
	Name          string        `json:"name"`
}

type capacityLimit struct {
	Size float64 `json:"size"`
	Unit string  `json:"unit"`
}

type lif struct {
	DataProtocols []string `json:"dataProtocols"`
	IP            string   `json:"ip"`
	LifType       string   `json:"lifType"`
	Netmask       string   `json:"netmask"`
	NodeName      string   `json:"nodeName"`
	PrivateIP     bool     `json:"privateIp"`
}

type cvoStatus struct {
	ExtendedFailureReason interface{}   `json:"extendedFailureReason"`
	FailureCauses         failureCauses `json:"failureCauses"`
	Message               string        `json:"message"`
	Status                string        `json:"status"`
}

type failureCauses struct {
	InvalidCloudProviderCredentials bool `json:"invalidCloudProviderCredentials"`
	InvalidOntapCredentials         bool `json:"invalidOntapCredentials"`
	NoCloudProviderConnection       bool `json:"noCloudProviderConnection"`
}

type svm struct {
	Name              string   `structs:"name"`
	State             string   `structs:"state"`
	Language          string   `structs:"language"`
	AllowedAggregates []string `structs:"allowAggregates"`
	Ver3Enabled       bool     `structs:"ver3Enabled"`
	Ver4Enabled       bool     `structs:"ver4Enabled"`
}

type configValuesUpdateRequest struct {
	GcpBlockProjectSSHKeys bool `structs:"gcpBlockProjectSshKeys"`
	GcpSerialPortEnable    bool `structs:"gcpSerialPortEnable"`
	GcpEnableOsLogin       bool `structs:"gcpEnableOsLogin"`
	GcpEnableOsLoginSk     bool `structs:"gcpEnableOsLoginSk"`
}

type configValuesResponse struct {
	GcpInstanceMetadataItems gcpInstanceMetadata `json:"gcpInstanceMetadataItems"`
}
type gcpInstanceMetadata struct {
	BlockProjectSSHKeys bool `json:"blockProjectSshKeys"`
	SerialPortEnable    bool `json:"serialPortEnable"`
	EnableOsLogin       bool `json:"enableOsLogin"`
	EnableOsLoginSk     bool `json:"enableOsLoginSk"`
}

// userTags the input for requesting a CVO
type userTags struct {
	TagKey   string `structs:"tagKey"`
	TagValue string `structs:"tagValue,omitempty"`
}

// modifyUserTagsRequest the input for requesting tags modification
type modifyUserTagsRequest struct {
	Tags []userTags `structs:"tags"`
}

// setPasswordRequest the input for for setting password
type setPasswordRequest struct {
	Password string `structs:"password"`
}

// licenseAndInstanceTypeModificationRequest the input for license and instance type modification
type licenseAndInstanceTypeModificationRequest struct {
	InstanceType string `structs:"instanceType"`
	LicenseType  string `structs:"licenseType"`
}

// changeTierLevelRequest the input for tier level change
type changeTierLevelRequest struct {
	Level string `structs:"level"`
}

// changeWritingSpeedStateRequest the input for writing speed state change
type changeWritingSpeedStateRequest struct {
	WritingSpeedState string `structs:"writingSpeedState"`
}

// upgradeOntapVersionRequest
type upgradeOntapVersionRequest struct {
	UpdateType      string `structs:"updateType"`
	UpdateParameter string `structs:"updateParameter"`
}

// set config flag
type setFlagRequest struct {
	Value     bool   `structs:"value"`
	ValueType string `structs:"valueType"`
}

// svmNameModificationRequest
type svmNameModificationRequest struct {
	SvmNewName string `structs:"svmNewName"`
	SvmName    string `structs:"svmName"`
}

// snapshotPolicy
type cvoSnapshotPolicy struct {
	Name        string           `json:"name"`
	Schedules   []policySchedule `json:"schedules"`
	Description string           `json:"description"`
}

type policySchedule struct {
	Frequency string `json:"frequency"`
	Retention int    `json:"retention"`
}

// vmInstance
type vmInstance struct {
	NetworkInterfaces []networkInterfaces `json:"networkInterfaces"`
}

// networkInterfaces
type networkInterfaces struct {
	AccessConfigs []accessConfigs `json:"accessConfigs"`
}

// accessConfigs
type accessConfigs struct {
	NatIP string `json:"natIP"`
}

// check deployment mode and validate the input
func (c *Client) checkDeploymentMode(d *schema.ResourceData, clientID string) (bool, string, error) {
	isSaaS := true
	// get the deployment_mode
	deploymentMode := "Standard"
	if deployment, ok := d.GetOk("deployment_mode"); ok {
		deploymentMode = deployment.(string)
	}
	if acnt, ok := d.GetOk("tenant_id"); ok {
		accessTokenResult, err := c.getAccessToken()
		if err != nil {
			c.revertDeploymentModeParameters(d, clientID)
			return false, "", err
		}
		c.Token = accessTokenResult.Token
		c.AccountID = acnt.(string)
		// check if the tenant_id is SaaS
		account, err := c.getAccountDetails(clientID)
		if err != nil {
			c.revertDeploymentModeParameters(d, clientID)
			log.Print("not able to get account details")
			return false, "", err
		}
		isSaaS = account.IsSaas
	}
	connectorIP := ""
	if connector, ok := d.GetOk("connector_ip"); ok {
		connectorIP = connector.(string)
	}
	if deploymentMode == "Restricted" {
		if c.AccountID == "" {
			c.revertDeploymentModeParameters(d, clientID)
			return false, "", fmt.Errorf("tenant_id is required for deployment_mode Restricted")
		}
		if isSaaS {
			c.revertDeploymentModeParameters(d, clientID)
			return false, "", fmt.Errorf("the tenant_id %s is not a Restricted mode account", c.AccountID)
		}
		if connectorIP == "" {
			c.revertDeploymentModeParameters(d, clientID)
			return false, "", fmt.Errorf("connector_ip is required for deployment_mode Restricted")
		}
	} else if deploymentMode == "Standard" {
		if connectorIP != "" {
			c.revertDeploymentModeParameters(d, clientID)
			return false, "", fmt.Errorf("connector_ip is not required for deployment_mode Standard")
		}
		if !isSaaS {
			c.revertDeploymentModeParameters(d, clientID)
			return false, "", fmt.Errorf("the tenant_id %s is not in Standard mode account", c.AccountID)
		}
	}
	log.Printf("=== deployment_mode: %s ===", deploymentMode)
	return isSaaS, connectorIP, nil
}

// revert the configuration of deployment_mode, tenant_id, and connector_ip to the previous state
func (c *Client) revertDeploymentModeParameters(d *schema.ResourceData, clientID string) error {
	if d.HasChange("deployment_mode") {
		previous, _ := d.GetChange("deployment_mode")
		d.Set("deployment_mode", previous)
	}
	if d.HasChange("tenant_id") {
		previous, _ := d.GetChange("tenant_id")
		d.Set("tenant_id", previous)
	}
	if d.HasChange("connector_ip") {
		previous, _ := d.GetChange("connector_ip")
		d.Set("connector_ip", previous)
	}
	return nil
}

// Check HTTP response code, return error if HTTP request is not successful.
func apiResponseChecker(statusCode int, response []byte, funcName string) error {

	if statusCode >= 300 || statusCode < 200 {
		log.Printf("%s request failed: %v", funcName, string(response))
		return fmt.Errorf("code: %d, message: %s", statusCode, string(response))
	}

	return nil

}

func (c *Client) checkTaskStatus(id string, clientID string) (int, string, error) {

	log.Printf("checkTaskStatus: %s", id)

	baseURL := fmt.Sprintf("/occm/api/audit/activeTask/%s", id)

	hostType := "CloudManagerHost"

	var statusCode int
	var response []byte
	networkRetries := 9
	initialWaitTime := 1 * time.Second
	for i := 0; i < networkRetries; i++ {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
		if err != nil || code == 504 {
			if i < networkRetries {
				if code == 504 {
					waitTime := initialWaitTime * (1 << i) // Exponential backoff
					time.Sleep(waitTime)
				} else {
					time.Sleep(initialWaitTime)
				}
				log.Printf("checkTaskStatus id=%s code=%v error=%v Retries %v client=%s", id, code, err, networkRetries, clientID)
			} else {
				log.Printf("checkTaskStatus request failed after %v retries: %v, %v", i+1, code, err)
				return 0, "", err
			}
		} else {
			log.Printf("checkTaskStatus get request %s response code %v clientID %s", id, code, clientID)
			statusCode = code
			response = result
			break
		}
	}

	responseError := apiResponseChecker(statusCode, response, "checkTaskStatus")
	if responseError != nil {
		return 0, "", responseError
	}

	var result cvoStatusResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkTaskStatus ", err)
		return 0, "", err
	}

	return result.Status, result.Error, nil
}

func (c *Client) checkTaskStatusForNotSaas(id string, clientID string, connectorIP string) (int, string, error) {

	log.Printf("checkTaskStatus: %s", id)

	baseURL := fmt.Sprintf("/occm/api/audit/activeTask/%s", id)

	hostType := "http://" + connectorIP

	var statusCode int
	var response []byte
	networkRetries := 9
	initialWaitTime := 1 * time.Second
	for i := 0; i < networkRetries; i++ {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
		if err != nil || code == 504 {
			if i < networkRetries {
				if code == 504 {
					waitTime := initialWaitTime * (1 << i) // Exponential backoff
					time.Sleep(waitTime)
				} else {
					time.Sleep(initialWaitTime)
				}
				log.Printf("checkTaskStatus id=%s code=%v error=%v Retries %v client=%s", id, code, err, networkRetries, clientID)
			} else {
				log.Printf("checkTaskStatus request failed after %v retries: %v, %v", i+1, code, err)
				return 0, "", err
			}
		} else {
			log.Printf("checkTaskStatus get request %s response code %v clientID %s", id, code, clientID)
			statusCode = code
			response = result
			break
		}
	}

	responseError := apiResponseChecker(statusCode, response, "checkTaskStatus")
	if responseError != nil {
		return 0, "", responseError
	}

	var result cvoStatusResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkTaskStatus ", err)
		return 0, "", err
	}

	return result.Status, result.Error, nil
}

func (c *Client) waitOnCompletion(id string, actionName string, task string, retries int, waitInterval int, clientID string) error {
	for {
		cvoStatus, failureErrorMessage, err := c.checkTaskStatus(id, clientID)
		if err != nil {
			return err
		}
		if cvoStatus == 1 {
			return nil
		} else if cvoStatus == -1 {
			return fmt.Errorf("failed to %s %s, error: %s", task, actionName, failureErrorMessage)
		} else if cvoStatus == 0 {
			if retries == 0 {
				log.Print("Taking too long to ", task, actionName)
				return fmt.Errorf("taking too long for %s to %s or not properly setup", actionName, task)
			}
			log.Printf("Sleep for %d seconds", waitInterval)
			time.Sleep(time.Duration(waitInterval) * time.Second)
			retries--
		}

	}
}

func (c *Client) waitOnCompletionForNotSaas(id string, actionName string, task string, retries int, waitInterval int, clientID string, connectorIP string) error {
	for {
		cvoStatus, failureErrorMessage, err := c.checkTaskStatusForNotSaas(id, clientID, connectorIP)
		if err != nil {
			return err
		}
		if cvoStatus == 1 {
			return nil
		} else if cvoStatus == -1 {
			return fmt.Errorf("failed to %s %s, error: %s", task, actionName, failureErrorMessage)
		} else if cvoStatus == 0 {
			if retries == 0 {
				log.Print("Taking too long to ", task, actionName)
				return fmt.Errorf("taking too long for %s to %s or not properly setup", actionName, task)
			}
			log.Printf("Sleep for %d seconds", waitInterval)
			time.Sleep(time.Duration(waitInterval) * time.Second)
			retries--
		}

	}
}

// get working environment information by working environment id
// response: publicId, name, isHA, providerName, workingEnvironmentType, ...
func (c *Client) getWorkingEnvironmentInfo(id string, clientID string, isSaas bool, connectorIP string) (workingEnvironmentInfo, error) {
	baseURL := fmt.Sprintf("/occm/api/ontaps/working-environments/%s", id)
	hostType := "CloudManagerHost"
	if !isSaas {
		hostType = "http://" + connectorIP
	}

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		log.Print("Call API ", baseURL)
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(1 * time.Second)
				networkRetries--
			} else {
				log.Printf("getWorkingEnvironmentInfo: ID %s request failed. Err: %v", id, err)
				return workingEnvironmentInfo{}, err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}
	responseError := apiResponseChecker(statusCode, response, "getWorkingEnvironmentInfo")
	if responseError != nil {
		log.Printf("apiResponseChecker error %v", responseError)
		return workingEnvironmentInfo{}, responseError
	}

	var result workingEnvironmentInfo
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}

	result.CloudProviderName = result.ProviderName
	return result, nil
}

func findWE(name string, weList []workingEnvironmentInfo) (workingEnvironmentInfo, error) {

	for i := range weList {
		if weList[i].Name == name {
			log.Printf("Found working environment: %v", weList[i])
			return weList[i], nil
		}
	}
	return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment %s in the list", name)
}

func findWEForID(id string, weList []workingEnvironmentInfo) (workingEnvironmentInfo, error) {

	for i := range weList {
		if weList[i].PublicID == id {
			log.Printf("Found working environment: %v", weList[i])
			return weList[i], nil
		}
	}
	return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment %s in the list", id)
}

func (c *Client) findWorkingEnvironmentByName(name string, clientID string, isSaas bool, connectorIP string) (workingEnvironmentInfo, error) {
	// check working environment exists or not
	baseURL := fmt.Sprintf("/occm/api/working-environments/exists/%s", name)
	hostType := "CloudManagerHost"
	if !isSaas {
		hostType = "http://" + connectorIP
	}

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("findWorkingEnvironmentByName request failed. (check exists) ", statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "findWorkingEnvironmentByName")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	// get working environment information
	baseURL = "/occm/api/working-environments"
	statusCode, response, _, err = c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("findWorkingEnvironmentByName %s request failed (%d)", name, statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError = apiResponseChecker(statusCode, response, "findWorkingEnvironmentByName")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var workingEnvironments workingEnvironmentResult
	if err := json.Unmarshal(response, &workingEnvironments); err != nil {
		log.Print("Failed to unmarshall response from findWorkingEnvironmentByName")
		return workingEnvironmentInfo{}, err
	}

	var workingEnvironment workingEnvironmentInfo
	workingEnvironment, err = findWE(name, workingEnvironments.VsaWorkingEnvironment)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.OnPremWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.AzureVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWE(name, workingEnvironments.GcpVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}

	log.Printf("Cannot find the working environment %s", name)

	return workingEnvironmentInfo{}, err
}

// get WE directly from REST API using a given ID
func (c *Client) findWorkingEnvironmentByID(id string, clientID string, isSaas bool, connectorIP string) (workingEnvironmentInfo, error) {

	workingEnvInfo, err := c.getWorkingEnvironmentInfo(id, clientID, isSaas, connectorIP)
	if err != nil {
		return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by working_environment_id %s", id)
	}
	workingEnvDetail, err := c.findWorkingEnvironmentByName(workingEnvInfo.Name, clientID, isSaas, connectorIP)

	if err != nil {
		return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by working_environment_name %s", workingEnvInfo.Name)
	}
	return workingEnvDetail, nil
}

func (c *Client) getFSXWorkingEnvironmentInfo(tenantID string, id string, clientID string, isSaaS bool, connectorIP string) (workingEnvironmentInfo, error) {
	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s/%s", tenantID, id)
	hostType := "CloudManagerHost"
	if !isSaaS {
		hostType = "http://" + connectorIP
	}

	var result workingEnvironmentInfo

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("getFSXWorkingEnvironmentInfo %s request failed (%d)", id, statusCode)
		log.Printf("error: %#v", err)
		return workingEnvironmentInfo{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getFSXWorkingEnvironmentInfo")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var system map[string]interface{}
	if err := json.Unmarshal(response, &system); err != nil {
		log.Print("Failed to unmarshall response from getFSXWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}
	result.Name = system["name"].(string)

	baseURL = fmt.Sprintf("/occm/api/fsx/working-environments/%s/svms", id)
	statusCode, response, _, err = c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("getFSXWorkingEnvironmentInfo %s request failed (%d)", id, statusCode)
		return workingEnvironmentInfo{}, err
	}
	responseError = apiResponseChecker(statusCode, response, "getFSXWorkingEnvironmentInfo")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}
	var info []map[string]interface{}
	if err := json.Unmarshal(response, &info); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}
	//assume there is only one svm in fsx
	result.SvmName = info[0]["name"].(string)

	return result, nil
}

func (c *Client) getAPIRoot(workingEnvironmentID string, clientID string, isSaas bool, connectorIP string) (string, string, error) {

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", "", err
		}
		c.Token = accesTokenResult.Token
	}

	// fsx working environment starts with "fs-" prefix.
	if strings.HasPrefix(workingEnvironmentID, "fs-") {
		return "/occm/api/fsx", "", nil
	}
	// onPrem working environment starts with "OnPrem" prefix.
	if strings.HasPrefix(workingEnvironmentID, "OnPrem") {
		return "/occm/api/onprem", "", nil
	}
	workingEnvDetail, err := c.getWorkingEnvironmentInfo(workingEnvironmentID, clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("Cannot get working environment information.")
		return "", "", err
	}
	log.Printf("Working environment %v", workingEnvDetail)

	var baseURL string
	if workingEnvDetail.CloudProviderName != "Amazon" {
		if workingEnvDetail.IsHA {
			baseURL = fmt.Sprintf("/occm/api/%s/ha", strings.ToLower(workingEnvDetail.CloudProviderName))
		} else {
			baseURL = fmt.Sprintf("/occm/api/%s/vsa", strings.ToLower(workingEnvDetail.CloudProviderName))
		}
	} else {
		if workingEnvDetail.IsHA {
			baseURL = "/occm/api/aws/ha"
		} else {
			baseURL = "/occm/api/vsa"
		}
	}
	log.Printf("API root = %s", baseURL)
	return baseURL, workingEnvDetail.CloudProviderName, nil
}

func getAPIRootForWorkingEnvironment(isHA bool, workingEnvironmentID string) string {

	var baseURL string

	if workingEnvironmentID == "" {
		if isHA {
			baseURL = "/occm/api/gcp/ha/working-environments"
		} else {
			baseURL = "/occm/api/gcp/vsa/working-environments"
		}
	} else {
		if isHA {
			baseURL = fmt.Sprintf("/occm/api/gcp/ha/working-environments/%s", workingEnvironmentID)
		} else {
			baseURL = fmt.Sprintf("/occm/api/gcp/vsa/working-environments/%s", workingEnvironmentID)
		}
	}

	log.Printf("API root = %s", baseURL)
	return baseURL
}

// read working environment information and return the details
func (c *Client) getWorkingEnvironmentDetail(d *schema.ResourceData, clientID string, isSaas bool, connectorIP string) (workingEnvironmentInfo, error) {
	var workingEnvDetail workingEnvironmentInfo
	var err error

	if a, ok := d.GetOk("file_system_id"); ok {
		workingEnvDetail, err = c.getFSXWorkingEnvironmentInfo(d.Get("tenant_id").(string), a.(string), clientID, isSaas, connectorIP)

		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by working_environment_id %s", a.(string))
		}
		return workingEnvDetail, nil
	}

	if a, ok := d.GetOk("working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		workingEnvDetail, err = c.findWorkingEnvironmentByID(WorkingEnvironmentID, clientID, isSaas, connectorIP)

		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by working_environment_id %s", WorkingEnvironmentID)
		}
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string), clientID, isSaas, connectorIP)

		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", workingEnvDetail.PublicID, a.(string))
	} else {
		return workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by working_environment_id or working_environment_name")
	}
	return workingEnvDetail, nil
}

func (c *Client) getFSXSVM(id string, clientID string, isSaas bool, connectorIP string) (string, error) {

	log.Print("getFSXSVM")

	baseURL := fmt.Sprintf("/occm/api/fsx/working-environments/%s/svms", id)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getFSXSVM request failed ", statusCode)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getFSXSVM")
	if responseError != nil {
		return "", responseError
	}

	var result []fsxSVMResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getFSXSVM ", err)
		return "", err
	}

	if len(result) == 0 {
		return "", fmt.Errorf("no SVM found for %s", id)
	}

	return result[0].Name, nil
}

func (c *Client) getAWSFSXByName(name string, tenantID string, clientID string, isSaas bool, connectorIP string) (string, error) {

	log.Print("getAWSFSXByName")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getAWSFSXByName request, failed to get AccessToken")
		return "", err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/fsx-ontap/working-environments/%s", tenantID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getAWSFSXByName request failed ", statusCode, err)
		return "", err
	}

	responseError := apiResponseChecker(statusCode, response, "getAWSFSXByName")
	if responseError != nil {
		return "", responseError
	}

	var result []fsxResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getAWSFSXByName ", err)
		return "", err
	}

	for _, fsxID := range result {
		if fsxID.Name == name {
			return fsxID.ID, nil
		}
	}

	return "", nil
}

// read working environemnt information and return the details
func (c *Client) getWorkingEnvironmentDetailForSnapMirror(d *schema.ResourceData, clientID string, isSaas bool, connectorIP string) (workingEnvironmentInfo, workingEnvironmentInfo, error) {
	var sourceWorkingEnvDetail workingEnvironmentInfo
	var destWorkingEnvDetail workingEnvironmentInfo
	var err error

	if a, ok := d.GetOk("source_working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		// fsx working environment starts with "fs-" prefix.
		if strings.HasPrefix(WorkingEnvironmentID, "fs-") {
			if b, ok := d.GetOk("tenant_id"); ok {
				tenantID := b.(string)
				id, err := c.getAWSFSX(WorkingEnvironmentID, tenantID, isSaas, connectorIP)
				if err != nil {
					log.Print("Error getting AWS FSX")
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}

				if id != WorkingEnvironmentID {
					log.Print("Error getting AWS FSX")
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("could not find source working environment ID %v", WorkingEnvironmentID)
				}
				sourceWorkingEnvDetail.PublicID = WorkingEnvironmentID
				svmName, err := c.getFSXSVM(WorkingEnvironmentID, clientID, isSaas, connectorIP)
				if err != nil {
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}
				sourceWorkingEnvDetail.SvmName = svmName
			} else {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find FSX working environment by destination_working_environment_id %s, need tenant_id", WorkingEnvironmentID)
			}
		} else {
			sourceWorkingEnvDetail, err = c.findWorkingEnvironmentForID(WorkingEnvironmentID, clientID, isSaas, connectorIP)
			if err != nil {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by source_working_environment_id %s", WorkingEnvironmentID)
			}
		}
	} else if a, ok = d.GetOk("source_working_environment_name"); ok {
		sourceWorkingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string), clientID, isSaas, connectorIP)
		if sourceWorkingEnvDetail.PublicID == "" {
			if b, ok := d.GetOk("tenant_id"); ok {
				workingEnvironmentName := a.(string)
				tenantID := b.(string)
				WorkingEnvironmentID, err := c.getAWSFSXByName(workingEnvironmentName, tenantID, clientID, isSaas, connectorIP)
				if err != nil {
					log.Print("Error getting AWS FSX: ", err)
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}
				sourceWorkingEnvDetail.PublicID = WorkingEnvironmentID
				if sourceWorkingEnvDetail.PublicID != "" {
					svmName, err := c.getFSXSVM(WorkingEnvironmentID, clientID, isSaas, connectorIP)
					if err != nil {
						return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
					}
					sourceWorkingEnvDetail.SvmName = svmName
				}
			}
		}
		if err != nil && sourceWorkingEnvDetail.PublicID == "" {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by source_working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", sourceWorkingEnvDetail.PublicID, a.(string))
	} else {
		return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by source_working_environment_id or source_working_environment_name")
	}

	if a, ok := d.GetOk("destination_working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		// fsx working environment starts with "fs-" prefix.
		if strings.HasPrefix(WorkingEnvironmentID, "fs-") {
			if b, ok := d.GetOk("tenant_id"); ok {
				tenantID := b.(string)
				id, err := c.getAWSFSX(WorkingEnvironmentID, tenantID, isSaas, connectorIP)
				if err != nil {
					log.Print("Error getting AWS FSX")
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}
				if id != WorkingEnvironmentID {
					log.Print("Error getting AWS FSX")
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("could not find destination working environment ID %v", WorkingEnvironmentID)
				}
				destWorkingEnvDetail.PublicID = WorkingEnvironmentID
				svmName, err := c.getFSXSVM(WorkingEnvironmentID, clientID, isSaas, connectorIP)
				if err != nil {
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}
				destWorkingEnvDetail.SvmName = svmName
			} else {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find FSX working environment by destination_working_environment_id %s, need tenant_id", WorkingEnvironmentID)
			}
		} else {
			destWorkingEnvDetail, err = c.findWorkingEnvironmentForID(WorkingEnvironmentID, clientID, isSaas, connectorIP)
			if err != nil {
				return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by destination_working_environment_id %s", WorkingEnvironmentID)
			}
			log.Print("findWorkingEnvironmentForID", destWorkingEnvDetail)
		}
	} else if a, ok = d.GetOk("destination_working_environment_name"); ok {
		destWorkingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string), clientID, isSaas, connectorIP)
		log.Printf("Get environment id %v by %v", destWorkingEnvDetail.PublicID, a.(string))
		if destWorkingEnvDetail.PublicID == "" {
			if b, ok := d.GetOk("tenant_id"); ok {
				workingEnvironmentName := a.(string)
				tenantID := b.(string)
				WorkingEnvironmentID, err := c.getAWSFSXByName(workingEnvironmentName, tenantID, clientID, isSaas, connectorIP)
				if err != nil {
					log.Print("Error getting AWS FSX: ", err)
					return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
				}
				if destWorkingEnvDetail.PublicID != "" {
					destWorkingEnvDetail.PublicID = WorkingEnvironmentID
					svmName, err := c.getFSXSVM(WorkingEnvironmentID, clientID, isSaas, connectorIP)
					if err != nil {
						return workingEnvironmentInfo{}, workingEnvironmentInfo{}, err
					}
					destWorkingEnvDetail.SvmName = svmName
				}
			}
		}
		if err != nil && destWorkingEnvDetail.PublicID == "" {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by destination_working_environment_name %s", a.(string))
		}
	} else {
		return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("cannot find working environment by destination_working_environment_id or destination_working_environment_name")
	}
	return sourceWorkingEnvDetail, destWorkingEnvDetail, nil
}

// get all WE from REST API and then using a given ID get the WE
func (c *Client) findWorkingEnvironmentForID(id string, clientID string, isSaas bool, connectorIP string) (workingEnvironmentInfo, error) {
	hostType := "CloudManagerHost"
	if !isSaas {
		hostType = "http://" + connectorIP
	}

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	baseURL := "/occm/api/working-environments"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("findWorkingEnvironmentForId %s request failed (%d)", id, statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "findWorkingEnvironmentForId")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var workingEnvironments workingEnvironmentResult
	if err := json.Unmarshal(response, &workingEnvironments); err != nil {
		log.Print("Failed to unmarshall response from findWorkingEnvironmentForId")
		return workingEnvironmentInfo{}, err
	}

	var workingEnvironment workingEnvironmentInfo
	workingEnvironment, err = findWEForID(id, workingEnvironments.VsaWorkingEnvironment)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWEForID(id, workingEnvironments.OnPremWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWEForID(id, workingEnvironments.AzureVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}
	workingEnvironment, err = findWEForID(id, workingEnvironments.GcpVsaWorkingEnvironments)
	if err == nil {
		return workingEnvironment, nil
	}

	log.Printf("Cannot find the working environment %s", id)

	return workingEnvironmentInfo{}, err
}

// get working environment properties
func (c *Client) getWorkingEnvironmentProperties(apiRoot string, id string, field string, clientID string, isSaas bool, connectorIP string) (workingEnvironmentOntapClusterPropertiesResponse, error) {
	hostType := "CloudManagerHost"
	if !isSaas {
		hostType = "http://" + connectorIP
	}
	baseURL := fmt.Sprintf("%s/working-environments/%s?fields=%s", apiRoot, id, field)
	log.Printf("Call %s", baseURL)

	var statusCode int
	var response []byte
	networkRetries := 10
	httpRetries := 10 // retry if API returns non 200 status code
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
		if err != nil {
			log.Printf("network error %v", err)
			if networkRetries > 0 {
				time.Sleep(2 * time.Second)
				networkRetries--
			} else {
				log.Printf("getWorkingEnvironmentProperties %s request failed (%d) %s", baseURL, statusCode, err)
				return workingEnvironmentOntapClusterPropertiesResponse{}, err
			}
		} else {
			if code > 400 { //client error or server error
				log.Printf("http status code %v", code)
				log.Printf("http response %s", result)
				if httpRetries > 0 {
					time.Sleep(10 * time.Second)
					httpRetries--
					continue
				}
			}
			statusCode = code
			response = result
			break
		}
	}
	responseError := apiResponseChecker(statusCode, response, "getWorkingEnvironmentProperties")
	if responseError != nil {
		return workingEnvironmentOntapClusterPropertiesResponse{}, responseError
	}

	var result workingEnvironmentOntapClusterPropertiesResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentProperties ", err)
		return workingEnvironmentOntapClusterPropertiesResponse{}, err
	}
	log.Printf("Get cvo properities result %+v", result)
	return result, nil
}

// customized check diff user-tags (aws_tag, azure_tag, gcp_label)
func checkUserTagDiff(diff *schema.ResourceDiff, tagName string, keyName string) error {
	if diff.HasChange(tagName) {
		_, expectTag := diff.GetChange(tagName)
		etags := expectTag.(*schema.Set)
		if etags.Len() > 0 {
			log.Println("etags len: ", etags.Len())
			// check each of the tag_key in the list is unique
			respErr := checkUserTagKeyUnique(etags, keyName)
			if respErr != nil {
				return respErr
			}
		}
	}
	return nil
}

// check each of the tag_key or label_key in the list is unique
func checkUserTagKeyUnique(etags *schema.Set, keyName string) error {
	m := make(map[string]bool)
	for _, v := range etags.List() {
		tag := v.(map[string]interface{})
		tkey := tag[keyName].(string)
		if _, ok := m[tkey]; !ok {
			m[tkey] = true
		} else {
			return fmt.Errorf("%s %s is not unique", keyName, tkey)
		}
	}
	return nil
}

// expandUserTags converts set to userTags struct
func expandUserTags(set *schema.Set) []userTags {
	tags := []userTags{}

	for _, v := range set.List() {
		tag := v.(map[string]interface{})
		userTag := userTags{}
		userTag.TagKey = tag["tag_key"].(string)
		userTag.TagValue = tag["tag_value"].(string)
		tags = append(tags, userTag)
	}
	return tags
}

func (c *Client) callCMUpdateAPI(method string, request interface{}, baseURL string, id string, functionName string, clientID string, isSaas bool, connectorIP string) error {
	apiRoot, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return err
	}
	baseURL = apiRoot + baseURL

	hostType := "CloudManagerHost"
	if !isSaas {
		hostType = "http://" + connectorIP
	}

	params := structs.Map(request)

	if c.Token == "" {
		accessTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Printf("in %s request, failed to get AccessToken", functionName)
			return err
		}
		c.Token = accessTokenResult.Token
	}

	statusCode, response, _, err := c.CallAPIMethod(method, baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Printf("%s request failed: %d", functionName, statusCode)
		log.Print("call api response: ", response)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, functionName)
	if responseError != nil {
		return responseError
	}
	return nil
}

// modify CVO SVM name
func (c *Client) updateCVOSVMName(d *schema.ResourceData, clientID string, svmName string, svmNewName string, isSaas bool, connectorIP string) error {
	var request svmNameModificationRequest
	// Update svm name
	id := d.Id()
	request.SvmName = svmName
	request.SvmNewName = svmNewName
	baseURL := fmt.Sprintf("/working-environments/%s/svm", id)
	log.Printf("Modify %s SVM %s with %s", id, svmName, svmNewName)
	updateErr := c.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOSVMName", clientID, isSaas, connectorIP)
	if updateErr != nil {
		return updateErr
	}
	log.Printf("\tUpdated %s on svm_name", id)
	return nil
}

// update CVO user-tags
func updateCVOUserTags(d *schema.ResourceData, meta interface{}, tagName string, clientID string, isSaas bool, connectorIP string) error {
	client := meta.(*Client)
	var request modifyUserTagsRequest
	if c, ok := d.GetOk(tagName); ok {
		tags := c.(*schema.Set)
		if tags.Len() > 0 {
			if tagName == "gcp_label" {
				request.Tags = expandGCPLabelsToUserTags(tags)
			} else {
				request.Tags = expandUserTags(tags)
			}
			log.Print("Update user-tags: ", request.Tags)
		}
	}
	// Update tags
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/user-tags", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOUserTags", clientID, isSaas, connectorIP)
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s %s: %v", id, tagName, request.Tags)
	return nil
}

// set the cluster password of a specific cloud volumes ONTAP
func updateCVOSVMPassword(d *schema.ResourceData, meta interface{}, clientID string, isSaas bool, connectorIP string) error {
	client := meta.(*Client)
	var request setPasswordRequest
	request.Password = d.Get("svm_password").(string)
	// Update password
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/set-password", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOSVMPassword", clientID, isSaas, connectorIP)
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s svm_password", id)
	return nil
}

// update SVMs on GCP CVO HA
func (c *Client) updateCVOSVMs(d *schema.ResourceData, clientID string, isSaas bool, connectorIP string) error {
	id := d.Id()
	currentSVMs, expectSVMs := d.GetChange("svm")
	cSVMs := expandGCPSVMs(currentSVMs.(*schema.Set))
	eSVMs := expandGCPSVMs(expectSVMs.(*schema.Set))

	// expectList will be used to keep the new SVMs which will be added later.
	// currentList will be used to keep the SVMs which will be removed later.
	// But rather than adding and deleting the SVMs, try to do rename/update them first.
	expectList := make(map[string]bool)
	currentList := make(map[int]string)

	for _, svm := range eSVMs {
		expectList[svm.SvmName] = true
	}

	i := 0
	for _, svm := range cSVMs {
		svmName := svm.SvmName
		if _, ok := expectList[svmName]; !ok {
			currentList[i] = svmName
			i++
		} else {
			delete(expectList, svmName)
		}
	}
	log.Printf("eList: %#v", expectList)
	log.Printf("cList: %#v", currentList)

	j := 0
	for svmName := range expectList {
		if len(currentList) > 0 {
			// update SVM
			respErr := c.updateCVOSVMName(d, clientID, currentList[j], svmName, isSaas, connectorIP)
			if respErr != nil {
				return respErr
			}
			delete(currentList, j)
			j++
		} else {
			// add SVM
			respErr := c.addSVMtoCVO(id, clientID, svmName, isSaas, connectorIP)
			if respErr != nil {
				log.Printf("Error adding SVM %v: %v", svmName, respErr)
				return respErr
			}
		}
	}
	if len(currentList) > 0 {
		for _, svmName := range currentList {
			// delete SVM
			respErr := c.deleteSVMfromCVO(id, clientID, svmName, isSaas, connectorIP)
			if respErr != nil {
				log.Printf("Error deleting SVM %v: %v", svmName, respErr)
				return respErr
			}
		}
	}
	return nil
}

func (c *Client) waitOnCompletionCVOUpdate(id string, retryCount int, waitInterval int, clientID string, isSaas bool, connectorIP string) error {
	// check upgrade status
	log.Print("Check CVO update status")
	// check upgrade status
	apiRoot, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return fmt.Errorf("cannot get root API")
	}

	for {
		cvoResp, err := c.getWorkingEnvironmentProperties(apiRoot, id, "status,ontapClusterProperties", clientID, isSaas, connectorIP)
		if err != nil {
			return err
		}
		if cvoResp.Status.Status != "UPDATING" {
			log.Print("CVO update is done")
			return nil
		}
		if retryCount <= 0 {
			log.Print("Taking too long for status to be active")
			return fmt.Errorf("taking too long for CVO to be active or not properly setup")
		}
		log.Printf("Update status %s...(%d)", cvoResp.Status.Status, retryCount)
		time.Sleep(time.Duration(waitInterval) * time.Second)
		retryCount--
	}
}

func (c *Client) getCVOProperties(id string, clientID string, isSaas bool, connectorIP string) (workingEnvironmentOntapClusterPropertiesResponse, error) {
	apiRoot, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return workingEnvironmentOntapClusterPropertiesResponse{}, fmt.Errorf("cannot get root API")
	}
	cvoResp, err := c.getWorkingEnvironmentProperties(apiRoot, id, "*", clientID, isSaas, connectorIP)
	if err != nil {
		return workingEnvironmentOntapClusterPropertiesResponse{}, err
	}
	return cvoResp, nil
}

// set the license_type and instance type of a specific cloud volumes ONTAP
func updateCVOLicenseInstanceType(d *schema.ResourceData, meta interface{}, clientID string, isSaas bool, connectorIP string) error {
	client := meta.(*Client)
	var request licenseAndInstanceTypeModificationRequest
	if c, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = c.(string)
	}
	if c, ok := d.GetOk("license_type"); ok {
		request.LicenseType = c.(string)
	}
	if request.LicenseType == "capacity-paygo" && d.Get("is_ha").(bool) {
		log.Print("Set licenseType as default value ha-capacity-paygo")
		request.LicenseType = "ha-capacity-paygo"
	}
	// Update license type and instance type
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/license-instance-type", id)
	log.Printf("Update license and instance type: %#v", request)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOLicenseInstanceType", clientID, isSaas, connectorIP)
	if updateErr != nil {
		return updateErr
	}
	// check upgrade status
	retryCount := 65
	if d.Get("is_ha").(bool) {
		retryCount = retryCount * 2
	}
	err := client.waitOnCompletionCVOUpdate(id, retryCount, 60, clientID, isSaas, connectorIP)
	if err != nil {
		return fmt.Errorf("update CVO failed %v", err)
	}
	log.Printf("Updated %s license and instance type: %v", id, request)
	return nil
}

// update tier_level of a specific cloud volumes ONTAP
func updateCVOTierLevel(d *schema.ResourceData, meta interface{}, clientID string, isSaas bool, connectorIP string) error {
	client := meta.(*Client)
	var request changeTierLevelRequest
	if c, ok := d.GetOk("tier_level"); ok {
		request.Level = c.(string)
	}
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/change-tier-level", id)
	updateErr := client.callCMUpdateAPI("POST", request, baseURL, id, "updateCVOTierLevel", clientID, isSaas, connectorIP)
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s tier_level: %v", id, request)
	return nil
}

// update writing_speed_state of a specific CVO
func updateCVOWritingSpeedState(d *schema.ResourceData, meta interface{}, clientID string, isSaas bool, connectorIP string) error {
	client := meta.(*Client)
	var request changeWritingSpeedStateRequest
	if c, ok := d.GetOk("writing_speed_state"); ok {
		request.WritingSpeedState = strings.ToUpper(c.(string))
	}
	log.Printf("writing_speed_state value %s", request.WritingSpeedState)
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/writing-speed", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOWritingSpeedState", clientID, isSaas, connectorIP)
	if updateErr != nil {
		return updateErr
	}

	// check upgrade status
	retryCount := 10
	if d.Get("is_ha").(bool) {
		retryCount = retryCount * 2
	}

	err := client.waitOnCompletionCVOUpdate(id, retryCount, 60, clientID, isSaas, connectorIP)
	if err != nil {
		return fmt.Errorf("update CVO failed %v", err)
	}
	log.Printf("Updated %s writing_speed_state: %v", id, request)
	return nil
}

func (c *Client) waitOnCompletionOntapImageUpgrade(apiRoot string, id string, targetVersion string, retryCount int, waitInterval int, clientID string, isSaas bool, connectorIP string) error {
	// check upgrade status
	log.Print("Check CVO ontap image upgrade status")

	for {
		cvoResp, err := c.getWorkingEnvironmentProperties(apiRoot, id, "status,ontapClusterProperties", clientID, isSaas, connectorIP)
		if err != nil {
			return err
		}
		if cvoResp.Status.Status != "UPDATING" && cvoResp.OntapClusterProperties.OntapVersion != "" {
			if strings.Contains(targetVersion, cvoResp.OntapClusterProperties.OntapVersion) {
				log.Print("ONTAP image upgrade is done")
				return nil
			}
			log.Printf("Update ontap image failed on checking version (%s, %s)", cvoResp.OntapClusterProperties.OntapVersion, targetVersion)
			return fmt.Errorf("update ontap version failed. Current version %s", cvoResp.OntapClusterProperties.OntapVersion)
		}
		if retryCount <= 0 {
			log.Print("Taking too long for status to be active")
			return fmt.Errorf("taking too long for CVO to be active or not properly setup")
		}
		log.Printf("Update %s status %s...(%d)", targetVersion, cvoResp.Status.Status, retryCount)
		time.Sleep(time.Duration(waitInterval) * time.Second)
		retryCount--
	}
}

// check if ontap_version is the list of upgrade available versions
func (c *Client) upgradeOntapVersionAvailable(apiRoot string, id string, ontapVersion string, clientID string, isSaas bool, connectorIP string) (string, error) {
	log.Print("upgradeOntapVersionAvailable: Check if target version is in the upgrade version list")

	var upgradeOntapVersions []upgradeVersion

	WEProperties, err := c.getWorkingEnvironmentProperties(apiRoot, id, "ontapClusterProperties.fields(upgradeVersions)", clientID, isSaas, connectorIP)
	if err != nil {
		return "", fmt.Errorf("upgradeOntapVersionAvailable %s not able to get the properties %v", id, err)
	}
	log.Printf("Get current ontap version: %s", WEProperties.OntapClusterProperties.OntapVersion)

	upgradeOntapVersions = WEProperties.OntapClusterProperties.UpgradeVersions

	if upgradeOntapVersions != nil {
		for _, ugVersion := range upgradeOntapVersions {
			version := ugVersion.ImageVersion
			if strings.Contains(ontapVersion, version) {
				return version, nil
			}
		}
		return "", fmt.Errorf("working environment %s: ontap version %s is not in the upgrade versions list (%+v)", id, ontapVersion, upgradeOntapVersions)
	}
	return "", fmt.Errorf("working environment %s: no upgrade version availble", id)
}

func (c *Client) setOCCMConfig(request configValuesUpdateRequest, clientID string, isSaas bool, occmDetails createOCCMDetails) error {
	log.Print("setOCCMConfig: set OCCM configuration")

	connectorIP := ""
	if !isSaas {
		vm, err := c.getVMInstance(occmDetails, clientID)
		if err != nil {
			log.Print("Error creating instance")
			return err
		}
		var vmInstance vmInstance
		vmjsonbody, err := json.Marshal(vm)
		if err != nil {
			log.Print("Failed to marshall response from getVMInstance ", err)
			return err
		}
		if err := json.Unmarshal(vmjsonbody, &vmInstance); err != nil {
			log.Print("Failed to unmarshall response from getVMInstance ", err)
			return err
		}

		if len(vmInstance.NetworkInterfaces) > 0 {
			if len(vmInstance.NetworkInterfaces[0].AccessConfigs) > 0 {
				connectorIP = vmInstance.NetworkInterfaces[0].AccessConfigs[0].NatIP
			}
		}
	}

	hostType := "CloudManagerHost"
	if !isSaas {
		hostType = "http://" + connectorIP
	}

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return err
		}
		c.Token = accesTokenResult.Token
	}

	baseURL := "/occm/api/occm/config"
	params := structs.Map(request)
	log.Printf("\tparams: %+v", params)
	statusCode, response, _, err := c.CallAPIMethod("PUT", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("setOCCMConfig request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "setOCCMConfig")
	if responseError != nil {
		return responseError
	}

	var result configValuesResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from setOCCMConfig ", err)
		return err
	}
	log.Printf("\tsetOCCMConfig result: %+v", result.GcpInstanceMetadataItems)
	return nil
}

func (c *Client) setConfigFlag(request setFlagRequest, keyPath string, clientID string, isSaas bool, connectorIP string) error {
	log.Print("setConfigFlag: set flag to allow ONTAP image upgrade")

	hostType := "CloudManagerHost"
	if !isSaas {
		hostType = "http://" + connectorIP
	}

	baseURL := fmt.Sprintf("/occm/api/occm/config/%s", keyPath)
	params := structs.Map(request)
	statusCode, response, _, err := c.CallAPIMethod("PUT", baseURL, params, c.Token, hostType, clientID)

	responseError := apiResponseChecker(statusCode, response, "setUpgradeCheckingBypass")
	if responseError != nil {
		return responseError
	}

	if err != nil {
		log.Print("setUpgradeCheckingBypass request failed ", statusCode)
		return err
	}
	return nil

}

// upgrade CVO ontap version
func (c *Client) upgradeCVOOntapImage(apiRoot string, id string, ontapVersion string, isHa bool, clientID string, isSaas bool, connectorIP string) error {
	// set config flag to skip the upgrade check
	var setFlag setFlagRequest
	setFlag.Value = true
	setFlag.ValueType = "BOOLEAN"

	log.Printf("Set config flag")
	setFlagErr := c.setConfigFlag(setFlag, "skip-eligibility-paygo-upgrade", clientID, isSaas, connectorIP)
	if setFlagErr != nil {
		log.Printf("upgradeCVOOntapVersion failed on setConfigFlag call %v", setFlagErr)
		return setFlagErr
	}

	// upgrade image
	var request upgradeOntapVersionRequest
	request.UpdateType = "OCCM_PROVIDED"
	request.UpdateParameter = ontapVersion

	baseURL := fmt.Sprintf("/working-environments/%s/update-image", id)
	log.Printf("upgradeCVOOntapVersion - %s %v", baseURL, request)
	updateErr := c.callCMUpdateAPI("POST", request, baseURL, id, "upgradeCVOOntapVersion", clientID, isSaas, connectorIP)
	if updateErr != nil {
		log.Printf("upgradeCVOOntapVersion failed on API call %v", updateErr)
		return updateErr
	}

	// check upgrade status
	retryCount := 65
	if isHa {
		retryCount = retryCount * 2
	}
	err := c.waitOnCompletionOntapImageUpgrade(apiRoot, id, ontapVersion, retryCount, 60, clientID, isSaas, connectorIP)
	if err != nil {
		return fmt.Errorf("upgrade ontap image %s failed %v", ontapVersion, err)
	}
	log.Printf("Updated %s ontap_version: %v", id, request)
	return nil
}

func (c *Client) doUpgradeCVOOntapVersion(id string, isHA bool, ontapVersion string, clientID string, isSaas bool, connectorIP string) error {
	// only when the upgrade_ontap_version is true, use_latest_version is false and the ontap_version is not "latest"
	log.Print("Check CVO ontap image upgrade status ... ")
	apiRoot, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return fmt.Errorf("cannot get root API")
	}

	upgradeVersion, err := c.upgradeOntapVersionAvailable(apiRoot, id, ontapVersion, clientID, isSaas, connectorIP)
	if err != nil {
		return err
	}

	return c.upgradeCVOOntapImage(apiRoot, id, upgradeVersion, isHA, clientID, isSaas, connectorIP)
}

func checkOntapVersionChangeWithoutUpgrade(d *schema.ResourceData) error {
	var wrongChange = false
	if d.HasChange("ontap_version") {
		currentVersion, _ := d.GetChange("ontap_version")
		d.Set("ontap_version", currentVersion)
		wrongChange = true
	}
	if d.HasChange("use_latest_version") {
		current, _ := d.GetChange("use_latest_version")
		d.Set("use_latest_version", current)
		wrongChange = true
	}
	if wrongChange {
		return fmt.Errorf("upgrade_ontap_version is not turned on. The change will not be done")
	}
	log.Printf("No ontap version upgrade")
	return nil
}

func (c *Client) checkAndDoUpgradeOntapVersion(d *schema.ResourceData, clientID string, isSaas bool, connectorIP string) error {
	upgradeOntapVersion := d.Get("upgrade_ontap_version").(bool)
	if upgradeOntapVersion {
		ontapVersion := d.Get("ontap_version").(string)
		log.Printf("Check if need upgrade - ontapVersion %s", ontapVersion)

		if ontapVersion == "latest" {
			return fmt.Errorf("ontap_version only can be upgraded with the specific ontap_version not \"latest\"")
		}
		if d.Get("use_latest_version").(bool) {
			return fmt.Errorf("ontap_version cannot be upgraded with \"use_latest_version\" true")
		}
		id := d.Id()
		respErr := c.doUpgradeCVOOntapVersion(id, d.Get("is_ha").(bool), ontapVersion, clientID, isSaas, connectorIP)
		if respErr != nil {
			currentVersion, _ := d.GetChange("ontap_version")
			d.Set("ontap_version", currentVersion)
			return respErr
		}
	} else {
		respErr := checkOntapVersionChangeWithoutUpgrade(d)
		if respErr != nil {
			return respErr
		}
	}
	return nil
}
