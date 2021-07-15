package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

type apiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// only list what is needed
type workingEnvironmentInfo struct {
	Name                   string `json:"name"`
	PublicID               string `json:"publicId"`
	CloudProviderName      string `json:"cloudProviderName"`
	IsHA                   bool   `json:"isHA"`
	WorkingEnvironmentType string `json:"workingEnvironmentType"`
	SvmName                string `json:"svmName"`
}

type workingEnvironmentResult struct {
	VsaWorkingEnvironment       []workingEnvironmentInfo `json:"vsaWorkingEnvironments"`
	OnPremWorkingEnvironments   []workingEnvironmentInfo `json:"onPremWorkingEnvironments"`
	AzureVsaWorkingEnvironments []workingEnvironmentInfo `json:"azureVsaWorkingEnvironments"`
	GcpVsaWorkingEnvironments   []workingEnvironmentInfo `json:"gcpVsaWorkingEnvironments"`
}

// userTags the input for requesting a CVO
type userTags struct {
	TagKey   string `structs:"tagKey"`
	TagValue string `structs:"tagValue,omitempty"`
}

// modifyUserTagsRequest the input for requesting tags modificaiton
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

// Check HTTP response code, return error if HTTP request is not successed.
func apiResponseChecker(statusCode int, response []byte, funcName string) error {

	if statusCode >= 300 || statusCode < 200 {
		log.Printf("%s request failed", funcName)
		return fmt.Errorf("code: %d, message: %s", statusCode, string(response))
	}

	return nil

}

func (c *Client) checkTaskStatus(id string) (int, string, error) {

	log.Printf("checkTaskStatus: %s", id)

	baseURL := fmt.Sprintf("/occm/api/audit/activeTask/%s", id)

	hostType := "CloudManagerHost"

	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(1 * time.Second)
				networkRetries--
			} else {
				log.Print("checkTaskStatus request failed ", code)
				return 0, "", err
			}
		} else {
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

func (c *Client) waitOnCompletion(id string, actionName string, task string, retries int, waitInterval int) error {
	for {
		cvoStatus, failureErrorMessage, err := c.checkTaskStatus(id)
		if err != nil {
			return err
		}
		if cvoStatus == 1 {
			return nil
		} else if cvoStatus == -1 {
			return fmt.Errorf("Failed to %s %s, error: %s", task, actionName, failureErrorMessage)
		} else if cvoStatus == 0 {
			if retries == 0 {
				log.Print("Taking too long to ", task, actionName)
				return fmt.Errorf("Taking too long for %s to %s or not properly setup", actionName, task)
			}
			log.Printf("Sleep for %d seconds", waitInterval)
			time.Sleep(time.Duration(waitInterval) * time.Second)
			retries--
		}

	}
}

// get working environment information by working environment id
// response: publicId, name, isHA, cloudProvider, workingEnvironmentType
func (c *Client) getWorkingEnvironmentInfo(id string) (workingEnvironmentInfo, error) {
	baseURL := fmt.Sprintf("/occm/api/working-environments/%s", id)
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("getWorkingEnvironmentInfo %s request failed (%d)", id, statusCode)
		return workingEnvironmentInfo{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getWorkingEnvironmentInfo")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var result workingEnvironmentInfo
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getWorkingEnvironmentInfo ", err)
		return workingEnvironmentInfo{}, err
	}

	return result, nil
}

func findWE(name string, weList []workingEnvironmentInfo) (workingEnvironmentInfo, error) {

	for i := range weList {
		if weList[i].Name == name {
			log.Printf("Found working environment: %v", weList[i])
			return weList[i], nil
		}
	}
	return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment %s in the list", name)
}

func findWEForID(id string, weList []workingEnvironmentInfo) (workingEnvironmentInfo, error) {

	for i := range weList {
		if weList[i].PublicID == id {
			log.Printf("Found working environment: %v", weList[i])
			return weList[i], nil
		}
	}
	return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment %s in the list", id)
}

func (c *Client) findWorkingEnvironmentByName(name string) (workingEnvironmentInfo, error) {
	// check working environment exists or not
	baseURL := fmt.Sprintf("/occm/api/working-environments/exists/%s", name)
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
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
	statusCode, response, _, err = c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
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
func (c *Client) findWorkingEnvironmentByID(id string) (workingEnvironmentInfo, error) {
	workingEnvInfo, err := c.getWorkingEnvironmentInfo(id)
	if err != nil {
		return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_id %s", id)
	}
	workingEnvDetail, err := c.findWorkingEnvironmentByName(workingEnvInfo.Name)
	if err != nil {
		return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_name %s", workingEnvInfo.Name)
	}
	return workingEnvDetail, nil
}

func (c *Client) getAPIRoot(workingEnvironmentID string) (string, string, error) {

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			log.Print("Not able to get the access token.")
			return "", "", err
		}
		c.Token = accesTokenResult.Token
	}
	workingEnvDetail, err := c.getWorkingEnvironmentInfo(workingEnvironmentID)
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

func (c *Client) getAPIRootForWorkingEnvironment(isHA bool, workingEnvironmentID string) string {

	var baseURL string

	if workingEnvironmentID == "" {
		if isHA == true {
			baseURL = "/occm/api/gcp/ha/working-environments"
		} else {
			baseURL = "/occm/api/gcp/vsa/working-environments"
		}
	} else {
		if isHA == true {
			baseURL = fmt.Sprintf("/occm/api/gcp/ha/working-environments/%s", workingEnvironmentID)
		} else {
			baseURL = fmt.Sprintf("/occm/api/gcp/vsa/working-environments/%s", workingEnvironmentID)
		}
	}

	log.Printf("API root = %s", baseURL)
	return baseURL
}

// read working environemnt information and return the details
func (c *Client) getWorkingEnvironmentDetail(d *schema.ResourceData) (workingEnvironmentInfo, error) {
	var workingEnvDetail workingEnvironmentInfo
	var err error

	if a, ok := d.GetOk("working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		workingEnvDetail, err = c.findWorkingEnvironmentByID(WorkingEnvironmentID)
		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_id %s", WorkingEnvironmentID)
		}
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", workingEnvDetail.PublicID, a.(string))
	} else {
		return workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by working_enviroment_id or working_environment_name")
	}
	return workingEnvDetail, nil
}

// read working environemnt information and return the details
func (c *Client) getWorkingEnvironmentDetailForSnapMirror(d *schema.ResourceData) (workingEnvironmentInfo, workingEnvironmentInfo, error) {
	var sourceWorkingEnvDetail workingEnvironmentInfo
	var destWorkingEnvDetail workingEnvironmentInfo
	var err error

	if a, ok := d.GetOk("source_working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		sourceWorkingEnvDetail, err = c.findWorkingEnvironmentForID(WorkingEnvironmentID)
		if err != nil {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by source_working_environment_id %s", WorkingEnvironmentID)
		}
	} else if a, ok = d.GetOk("source_working_environment_name"); ok {
		sourceWorkingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by source_working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", sourceWorkingEnvDetail.PublicID, a.(string))
	} else {
		return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by source_working_environment_id or source_working_environment_name")
	}

	if a, ok := d.GetOk("destination_working_environment_id"); ok {
		WorkingEnvironmentID := a.(string)
		destWorkingEnvDetail, err = c.findWorkingEnvironmentForID(WorkingEnvironmentID)
		if err != nil {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by destination_working_environment_id %s", WorkingEnvironmentID)
		}
		log.Print("findWorkingEnvironmentForID", destWorkingEnvDetail)
	} else if a, ok = d.GetOk("destination_working_environment_name"); ok {
		destWorkingEnvDetail, err = c.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by destination_working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", destWorkingEnvDetail.PublicID, a.(string))
	} else {
		return workingEnvironmentInfo{}, workingEnvironmentInfo{}, fmt.Errorf("Cannot find working environment by destination_working_environment_id or destination_working_environment_name")
	}
	return sourceWorkingEnvDetail, destWorkingEnvDetail, nil
}

// get all WE from REST API and then using a given ID get the WE
func (c *Client) findWorkingEnvironmentForID(id string) (workingEnvironmentInfo, error) {
	hostType := "CloudManagerHost"

	if c.Token == "" {
		accesTokenResult, err := c.getAccessToken()
		if err != nil {
			return workingEnvironmentInfo{}, err
		}
		c.Token = accesTokenResult.Token
	}
	baseURL := "/occm/api/working-environments"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType)
	if err != nil {
		log.Printf("findWorkingEnvironmentForId %s request failed (%d)", id, statusCode)
		return workingEnvironmentInfo{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "findWorkingEnvironmentByName")
	if responseError != nil {
		return workingEnvironmentInfo{}, responseError
	}

	var workingEnvironments workingEnvironmentResult
	if err := json.Unmarshal(response, &workingEnvironments); err != nil {
		log.Print("Failed to unmarshall response from findWorkingEnvironmentByName")
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
			return fmt.Errorf("Tag key %s is not unique", tkey)
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

func (c *Client) callCMUpdateAPI(method string, request interface{}, baseURL string, id string, functionName string) error {
	apiRoot, _, err := c.getAPIRoot(id)
	baseURL = apiRoot + baseURL

	hostType := "CloudManagerHost"
	params := structs.Map(request)

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Printf("in %s request, failed to get AccessToken", functionName)
		return err
	}
	c.Token = accessTokenResult.Token

	statusCode, response, _, err := c.CallAPIMethod(method, baseURL, params, c.Token, hostType)
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

// update CVO user-tags
func updateCVOUserTags(d *schema.ResourceData, meta interface{}, tagName string) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
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
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOUserTags")
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s %s: %v", id, tagName, request.Tags)
	return nil
}

// set the cluster password of a specific cloud volumes ONTAP
func updateCVOSVMPassword(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var request setPasswordRequest
	request.Password = d.Get("svm_password").(string)
	// Update password
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/set-password", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOSVMPassword")
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s svm_password", id)
	return nil
}

// set the license_type and instance type of a specific cloud volumes ONTAP
func updateCVOLicenseInstanceType(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var request licenseAndInstanceTypeModificationRequest
	if c, ok := d.GetOk("instance_type"); ok {
		request.InstanceType = c.(string)
	}
	if c, ok := d.GetOk("license_type"); ok {
		request.LicenseType = c.(string)
	}

	// Update license type and instance type
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/license-instance-type", id)
	updateErr := client.callCMUpdateAPI("PUT", request, baseURL, id, "updateCVOLicenseInstanceType")
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s license and instance type: %v", id, request)
	return nil
}

// update tier_level of a specific cloud volumes ONTAP
func updateCVOTierLevel(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var request changeTierLevelRequest
	if c, ok := d.GetOk("tier_level"); ok {
		request.Level = c.(string)
	}
	id := d.Id()
	baseURL := fmt.Sprintf("/working-environments/%s/change-tier-level", id)
	updateErr := client.callCMUpdateAPI("POST", request, baseURL, id, "updateCVOTierLevel")
	if updateErr != nil {
		return updateErr
	}
	log.Printf("Updated %s tier_level: %v", id, request)
	return nil
}
