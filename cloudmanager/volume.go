package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

type volumeRequest struct {
	Name                      string                 `structs:"name"`
	WorkingEnvironmentID      string                 `structs:"workingEnvironmentId"`
	SvmName                   string                 `structs:"svmName"`
	AggregateName             string                 `structs:"aggregateName"`
	Size                      size                   `structs:"size"`
	SnapshotPolicyName        string                 `structs:"snapshotPolicyName,omitempty"`
	EnableThinProvisioning    bool                   `structs:"enableThinProvisioning"`
	EnableCompression         bool                   `structs:"enableCompression"`
	EnableDeduplication       bool                   `structs:"enableDeduplication"`
	ExportPolicyInfo          ExportPolicyInfo       `structs:"exportPolicyInfo,omitempty"`
	ID                        string                 `structs:"uuid"`
	NewAggregate              bool                   `structs:"newAggregate"`
	CapacityTier              string                 `structs:"capacityTier,omitempty"`
	ProviderVolumeType        string                 `structs:"providerVolumeType"`
	TieringPolicy             string                 `structs:"tieringPolicy,omitempty"`
	NumOfDisks                float64                `structs:"maxNumOfDisksApprovedToAdd,omitempty"`
	AutoVsaCapacityManagement bool                   `structs:"autoVsaCapacityManagement"`
	DiskSize                  size                   `structs:"diskSize,omitempty"`
	Iops                      int                    `structs:"iops,omitempty"`
	Throughput                int                    `structs:"throughput,omitempty"`
	WorkingEnvironmentType    string                 `structs:"workingEnvironmentType,omitempty"`
	ShareInfo                 shareInfoRequest       `structs:"shareInfo,omitempty"`
	ShareInfoUpdate           shareInfoUpdateRequest `structs:"shareInfo,omitempty"`
	IscsiInfo                 iscsiInfo              `structs:"iscsiInfo,omitempty"`
	FileSystemID              string                 `structs:"fileSystemId,omitempty"`
	TenantID                  string                 `structs:"tenantId,omitempty"`
	EnableStorageEfficiency   bool                   `structs:"enableStorageEfficiency,omitempty"`
	VolumeTags                []volumeTag            `structs:"volumeTags,omitempty"`
	VolumeFSXTags             []volumeTag            `structs:"awsTags,omitempty"`
	Comment                   string                 `structs:"comment,omitempty"`
	ConnectorIP               string
}

type volumeResponse struct {
	Name                   string                   `json:"name"`
	SvmName                string                   `json:"svmName"`
	AggregateName          string                   `json:"aggregateName"`
	Size                   size                     `json:"size"`
	SnapshotPolicyName     string                   `json:"snapshotPolicy"`
	EnableThinProvisioning bool                     `json:"thinProvisioning"`
	EnableCompression      bool                     `json:"compression"`
	EnableDeduplication    bool                     `json:"deduplication"`
	ExportPolicyInfo       ExportPolicyInfoResponse `json:"exportPolicyInfo"`
	ID                     string                   `json:"uuid"`
	CapacityTier           string                   `json:"capacityTier,omitempty"`
	TieringPolicy          string                   `json:"tieringPolicy,omitempty"`
	ProviderVolumeType     string                   `json:"providerVolumeType"`
	ShareInfo              []shareInfoResponse      `json:"shareInfo"`
	MountPoint             string                   `json:"mountPoint"`
	IscsiEnabled           bool                     `json:"iscsiEnabled"`
	Comment                string                   `json:"comment"`
}

// ExportPolicyInfo describes the export policy section.
type ExportPolicyInfo struct {
	Name       string             `structs:"name,omitempty"`
	PolicyType string             `structs:"policyType,omitempty"`
	Ips        []string           `structs:"ips,omitempty"`
	NfsVersion []string           `structs:"nfsVersion,omitempty"`
	Rules      []ExportPolicyRule `structs:"rules,omitempty"`
}

// ExportPolicyInfoResponse describes the export policy section in API response.
type ExportPolicyInfoResponse struct {
	Name       string             `json:"name"`
	PolicyType string             `json:"policyType"`
	Ips        []string           `json:"ips"`
	NfsVersion []string           `json:"nfsVersion"`
	Rules      []ExportPolicyRule `json:"rules"`
}

// ExportPolicyRule describes the export policy rule section.
type ExportPolicyRule struct {
	// Protocols         []string `structs:"protocols"`
	// Clients           []string `structs:"clients"`
	// RoRule            []string `structs:"ro_rule"`
	// RwRule            []string `structs:"rw_rule"`
	Superuser         bool     `structs:"superuser"`
	Index             int32    `structs:"index,omitempty"`
	RuleAccessControl string   `structs:"ruleAccessControl"`
	Ips               []string `structs:"ips"`
	NfsVersion        []string `structs:"nfsVersion,omitempty"`
}

type exportPolicyInfoResponse struct {
	Name       string   `json:"name"`
	PolicyType string   `json:"policyType"`
	Ips        []string `json:"ips"`
	NfsVersion []string `json:"nfsVersion"`
	// Rules      exportPolicyRule `json:"rules"`
}

// type exportPolicyRule struct {
// 	Protocols []string `structs:"protocols"`
// 	Clients   []string `structs:"clients"`
// 	RoRule    []string `structs:"ro_rule"`
// 	RwRule    []string `structs:"rw_rule"`
// 	Superuser []string `structs:"superuser"`
// 	Index     int32    `structs:"index"`
// }

type size struct {
	Size float64 `structs:"size"`
	Unit string  `structs:"unit"`
}

type quoteRequest struct {
	Size                   size             `structs:"size"`
	WorkingEnvironmentID   string           `structs:"workingEnvironmentId"`
	SvmName                string           `structs:"svmName"`
	AggregateName          string           `structs:"aggregateName,omitempty"`
	EnableThinProvisioning bool             `structs:"enableThinProvisioning"`
	EnableCompression      bool             `structs:"enableCompression"`
	EnableDeduplication    bool             `structs:"enableDeduplication"`
	ReplicationFlow        bool             `structs:"replicationFlow"`
	ExportPolicyInfo       ExportPolicyInfo `structs:"exportPolicyInfo,omitempty"`
	SnapshotPolicyName     string           `structs:"snapshotPolicyName"`
	Name                   string           `structs:"name"`
	CapacityTier           string           `structs:"capacityTier,omitempty"`
	ProviderVolumeType     string           `structs:"providerVolumeType"`
	TieringPolicy          string           `structs:"tieringPolicy,omitempty"`
	VerifyNameUniqueness   bool             `structs:"verifyNameUniqueness"`
	Iops                   int              `structs:"iops,omitempty"`
	Throughput             int              `structs:"throughput,omitempty"`
	WorkingEnvironmentType string           `structs:"workingEnvironmentType"`
	ConnectorIP            string
}

type shareInfoRequest struct {
	ShareName     string        `structs:"shareName,omitempty"`
	AccessControl accessControl `structs:"accessControl,omitempty"`
}

type accessControl struct {
	Permission string   `structs:"permission,omitempty"`
	Users      []string `structs:"users,omitempty"`
}

// cifs volume response has different strcuture comparing to cifs create.
type shareInfoResponse struct {
	ShareName         string                      `json:"shareName"`
	AccessControlList []accessControlListResponse `json:"accessControlList"`
}

type accessControlListResponse struct {
	Permission string   `json:"permission"`
	Users      []string `json:"users"`
}

type shareInfoUpdateRequest struct {
	ShareName         string              `structs:"shareName,omitempty"`
	AccessControlList []accessControlList `structs:"accessControlList,omitempty"`
}

type accessControlList struct {
	Permission string   `structs:"permission,omitempty"`
	Users      []string `structs:"users,omitempty"`
}

type iscsiInfo struct {
	OsName                string `structs:"osName,omitempty"`
	IgroupCreationRequest struct {
		Initiators []string `structs:"initiators,omitempty"`
		IgroupName string   `structs:"igroupName,omitempty"`
	} `structs:"igroupCreationRequest,omitempty"`
	Igroups []string `structs:"igroups,omitempty"`
}

type initiator struct {
	AliasName              string `structs:"aliasName,omitempty"`
	Iqn                    string `structs:"iqn,omitempty"`
	WorkingEnvironmentID   string `structs:"workingEnvironmentId,omitempty"`
	SvmName                string `structs:"svmName,omitempty"`
	WorkingEnvironmentType string `structs:"workingEnvironmentType,omitempty"`
}

type volumeTag struct {
	TagKey   string `structs:"tagKey"`
	TagValue string `structs:"tagValue"`
}

type createSnapshotPolicyRequest struct {
	SnapshotPolicyName   string        `structs:"snapshotPolicyName"`
	Schedules            []scheduleReq `structs:"schedules"`
	WorkingEnvironmentID string        `structs:"workingEnvironmentId"`
}

type scheduleReq struct {
	ScheduleType string `structs:"scheduleType"`
	Retention    int    `structs:"retention"`
}

func (c *Client) createVolume(vol volumeRequest, createAggregateIfNotFound bool, clientID string, isSaas bool, connectorIP string) error {
	var id string
	if vol.FileSystemID != "" {
		id = vol.FileSystemID
	} else {
		id = vol.WorkingEnvironmentID
	}
	baseURL, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return err
	}
	if vol.FileSystemID != "" || vol.WorkingEnvironmentType == "ON_PREM" {
		baseURL = fmt.Sprintf("%s/volumes", baseURL)
	} else {
		baseURL = fmt.Sprintf("%s/volumes?createAggregateIfNotFound=%s", baseURL, strconv.FormatBool(createAggregateIfNotFound))
	}

	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	param := structs.Map(vol)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createVolume request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createVolume")
	if responseError != nil {
		return responseError
	}

	if isSaas {
		err = c.waitOnCompletion(onCloudRequestID, "volume", "create", 40, 10, clientID)
	} else {
		err = c.waitOnCompletionForNotSaas(onCloudRequestID, "volume", "create", 40, 10, clientID, connectorIP)
	}

	return err
}

func (c *Client) deleteVolume(vol volumeRequest, clientID string, isSaas bool, connectorIP string) error {
	var id string
	if vol.FileSystemID != "" {
		id = vol.FileSystemID
	} else {
		id = vol.WorkingEnvironmentID
	}
	baseURL, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/volumes/%s/%s/%s", baseURL, id, vol.SvmName, vol.Name)

	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("deleteVolume request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "deleteVolume")
	if responseError != nil {
		return responseError
	}

	log.Print("Wait for volume deletion.")

	if isSaas {
		err = c.waitOnCompletion(onCloudRequestID, "volume", "delete", 60, 10, clientID)
	} else {
		err = c.waitOnCompletionForNotSaas(onCloudRequestID, "volume", "delete", 60, 10, clientID, connectorIP)
	}

	return err
}

func (c *Client) getVolume(vol volumeRequest, clientID string, isSaas bool, connectorIP string) ([]volumeResponse, error) {
	var result []volumeResponse

	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	var id string
	if vol.FileSystemID != "" {
		id = vol.FileSystemID
	} else {
		id = vol.WorkingEnvironmentID
	}
	baseURL, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return result, err
	}
	if vol.FileSystemID != "" {
		baseURL = fmt.Sprintf("%s/volumes?fileSystemId=%s", baseURL, id)
	} else {
		baseURL = fmt.Sprintf("%s/volumes?workingEnvironmentId=%s", baseURL, id)
	}

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getVolume request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getVolume")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolume ", err)
		return result, err
	}

	return result, nil
}

func (c *Client) getVolumeForOnPrem(vol volumeRequest, clientID string) ([]volumeResponse, error) {
	var result []volumeResponse
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/occm/api/onprem/volumes?workingEnvironmentId=%s", vol.WorkingEnvironmentID)

	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getVolumeForOnPrem request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getVolumeForOnPrem")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getVolumeForOnPrem ", err)
		return result, err
	}

	return result, nil
}

func (c *Client) getVolumeByID(request volumeRequest, clientID string, isSaas bool, connectorIP string) (volumeResponse, error) {
	res, err := c.getVolume(request, clientID, isSaas, connectorIP)
	if err != nil {
		return volumeResponse{}, err
	}
	for _, vol := range res {
		if vol.ID == request.ID {
			return vol, nil
		}
	}
	return volumeResponse{}, fmt.Errorf("error fetching volume: volume doesn't exist")
}

func (c *Client) updateVolume(request volumeRequest, clientID string, isSaas bool, connectorIP string) error {
	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	var id string
	if request.FileSystemID != "" {
		id = request.FileSystemID
	} else {
		id = request.WorkingEnvironmentID
	}
	baseURL, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/volumes/%s/%s/%s", baseURL, id, request.SvmName, request.Name)
	params := structs.Map(request)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("PUT", baseURL, params, c.Token, hostType, clientID)

	responseError := apiResponseChecker(statusCode, response, "updateVolume")
	if responseError != nil {
		return responseError
	}

	if err != nil {
		log.Print("updateVolume request failed ", statusCode)
		return err
	}

	if isSaas {
		err = c.waitOnCompletion(onCloudRequestID, "volume", "update", 40, 10, clientID)
	} else {
		err = c.waitOnCompletionForNotSaas(onCloudRequestID, "volume", "update", 40, 10, clientID, connectorIP)
	}

	return err
}

func (c *Client) quoteVolume(request quoteRequest, clientID string, isSaas bool, connectorIP string) (map[string]interface{}, error) {
	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID, clientID, isSaas, connectorIP)
	if err != nil {
		return nil, err
	}
	baseURL = fmt.Sprintf("%s/volumes/quote", baseURL)
	params := structs.Map(request)

	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("quoteVolume request failed ", statusCode)
		return nil, err
	}
	responseError := apiResponseChecker(statusCode, response, "quoteVolume")
	if responseError != nil {
		return nil, responseError
	}
	var result map[string]interface{}
	json.Unmarshal(response, &result)
	return result, nil

}

func (c *Client) createInitiator(request initiator, clientID string, isSaas bool, connectorIP string) error {
	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID, clientID, isSaas, connectorIP)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/volumes/initiator", baseURL)
	params := structs.Map(request)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createInitiator request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createInitiator")
	if responseError != nil {
		return responseError
	}
	return nil
}

func (c *Client) getInitiator(request initiator, clientID string, isSaas bool, connectorIP string) ([]initiator, error) {
	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID, clientID, isSaas, connectorIP)
	var result []initiator
	if err != nil {
		return result, err
	}
	baseURL = fmt.Sprintf("%s/volumes/initiator", baseURL)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createInitiator request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "createInitiator")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getInitiator ", err)
		return result, err
	}
	return result, nil
}

type igroup struct {
	IgroupName             string   `json:"igroupName"`
	OsType                 string   `json:"osType"`
	PortsetName            string   `json:"portsetName"`
	IgroupType             string   `json:"igroupType"`
	Initiators             []string `json:"initiators"`
	WorkingEnvironmentID   string   `structs:"workingEnvironmentId"`
	SvmName                string   `structs:"svmName"`
	WorkingEnvironmentType string   `structs:"workingEnvironmentType,omitempty"`
}

func (c *Client) getIgroups(request igroup, clientID string, isSaas bool, connectorIP string) ([]igroup, error) {
	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	baseURL, _, err := c.getAPIRoot(request.WorkingEnvironmentID, clientID, isSaas, connectorIP)
	var result []igroup
	if err != nil {
		return result, err
	}
	if request.WorkingEnvironmentType == "ON_PREM" {
		log.Print("get igroup onPrem")
		baseURL = fmt.Sprintf("/occm/api/ontaps/working-environments/%s/volumes/%s/igroups", request.WorkingEnvironmentID, request.SvmName)
	} else {
		baseURL = fmt.Sprintf("%s/volumes/igroups/%s/%s", baseURL, request.WorkingEnvironmentID, request.SvmName)
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getIgroups request failed ", statusCode)
		return result, err
	}
	responseError := apiResponseChecker(statusCode, response, "getIgroups")
	if responseError != nil {
		return result, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getIgroups ", err)
		return result, err
	}
	return result, nil
}

func (c *Client) checkCifsExists(workingEnvironmentType string, id string, svm string, clientID string, isSaas bool, connectorIP string) (bool, error) {
	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	baseURL, _, err := c.getAPIRoot(id, clientID, isSaas, connectorIP)
	var result []map[string]interface{}
	if err != nil {
		return false, err
	}
	if workingEnvironmentType == "ON_PREM" {
		baseURL = fmt.Sprintf("%s/working-environments/%s/cifs?vserver=%s", baseURL, id, svm)
	} else {
		baseURL = fmt.Sprintf("%s/working-environments/%s/cifs?svm=%s", baseURL, id, svm)
	}
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("chkeckCifsExists request failed ", statusCode)
		return false, err
	}
	responseError := apiResponseChecker(statusCode, response, "checkCifsExists")
	if responseError != nil {
		return false, responseError
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from chkeckCifsExists ", err)
		return false, err
	}
	if len(result) > 0 {
		return true, nil
	}
	return false, nil
}

func convertSizeUnit(size float64, from string, to string) float64 {
	if strings.ToUpper(from) == "GB" && strings.ToUpper(to) == "GB" {
		return size
	}
	if strings.ToUpper(from) == "GB" && strings.ToUpper(to) == "TB" {
		size = size / 1024
	}
	if strings.ToUpper(from) == "GB" && strings.ToUpper(to) == "B" {
		size = size * 1073741824
	}
	if strings.ToUpper(from) == "B" && strings.ToUpper(to) == "GB" {
		size = size / 1073741824
	}
	return size
}

func (c *Client) setCommonAttributes(WorkingEnvironmentType string, d *schema.ResourceData, volume *volumeRequest, clientID string) error {
	volume.Name = d.Get("name").(string)
	volume.Size.Size = d.Get("size").(float64)
	volume.Size.Unit = d.Get("unit").(string)
	volume.SnapshotPolicyName = d.Get("snapshot_policy_name").(string)
	var weid string
	if fsid, ok := d.GetOk("file_system_id"); ok {
		weid = fsid.(string)
	} else {
		weid = volume.WorkingEnvironmentID
	}

	if v, ok := d.GetOk("export_policy_type"); ok {
		volume.ExportPolicyInfo.PolicyType = v.(string)
	}
	if v, ok := d.GetOk("export_policy_ip"); ok {
		ips := make([]string, 0, v.(*schema.Set).Len())
		for _, x := range v.(*schema.Set).List() {
			ips = append(ips, x.(string))
		}
		volume.ExportPolicyInfo.Ips = ips
	}
	if v, ok := d.GetOk("export_policy_nfs_version"); ok {
		nfs := make([]string, 0, v.(*schema.Set).Len())
		for _, x := range v.(*schema.Set).List() {
			nfs = append(nfs, x.(string))
		}
		volume.ExportPolicyInfo.NfsVersion = nfs
	}

	volumeProtocol := d.Get("volume_protocol").(string)
	if volumeProtocol == "cifs" {

		exist, err := c.checkCifsExists(WorkingEnvironmentType, weid, volume.SvmName, clientID, true, "")
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("cifs has not been set up yet")
		}
		if v, ok := d.GetOk("share_name"); ok {
			volume.ShareInfo.ShareName = v.(string)
		}
		if v, ok := d.GetOk("permission"); ok {
			volume.ShareInfo.AccessControl.Permission = v.(string)
		}
		if v, ok := d.GetOk("users"); ok {
			users := make([]string, 0, v.(*schema.Set).Len())
			for _, x := range v.(*schema.Set).List() {
				users = append(users, x.(string))
			}
			volume.ShareInfo.AccessControl.Users = users
		}
	}
	return nil
}

// createSnapshotPolicy
func (c *Client) createSnapshotPolicy(workingEnviromentID string, snapshotPolicyName string, set *schema.Set, clientID string, isSaas bool, connectorIP string) error {
	log.Print("createSnapshotPolicy: ", snapshotPolicyName)
	snapshotPolicy := createSnapshotPolicyRequest{}
	snapshotPolicy.SnapshotPolicyName = snapshotPolicyName
	snapshotPolicy.WorkingEnvironmentID = workingEnviromentID
	for _, v := range set.List() {
		schedules := v.(map[string]interface{})
		scheduleSet := schedules["schedule"].([]interface{})
		scheduleConfigs := make([]scheduleReq, 0, len(scheduleSet))
		for _, x := range scheduleSet {
			snapshotPolicySchedule := scheduleReq{}
			scheduleConfig := x.(map[string]interface{})
			snapshotPolicySchedule.ScheduleType = scheduleConfig["schedule_type"].(string)
			snapshotPolicySchedule.Retention = scheduleConfig["retention"].(int)

			scheduleConfigs = append(scheduleConfigs, snapshotPolicySchedule)
		}
		snapshotPolicy.Schedules = scheduleConfigs
	}
	baseURL, _, err := c.getAPIRoot(snapshotPolicy.WorkingEnvironmentID, clientID, isSaas, connectorIP)

	hostType := ""
	if isSaas {
		hostType = "CloudManagerHost"
	} else {
		hostType = "http://" + connectorIP
	}

	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/working-environments/%s/snapshot-policy", baseURL, snapshotPolicy.WorkingEnvironmentID)
	param := structs.Map(snapshotPolicy)
	statusCode, response, onCloudRequestID, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createSnapshotPolicy request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "createSnapshotPolicy")
	if responseError != nil {
		return responseError
	}

	if isSaas {
		err = c.waitOnCompletion(onCloudRequestID, "snapshotPolicy", "create", 10, 10, clientID)

	} else {
		err = c.waitOnCompletionForNotSaas(onCloudRequestID, "snapshotPolicy", "create", 10, 10, clientID, connectorIP)
	}
	if err != nil {
		return err
	}

	if c.findSnapshotPolicy(workingEnviromentID, snapshotPolicyName, clientID, isSaas, connectorIP) {
		return nil
	}

	return fmt.Errorf("create snapshot policy failed")
}

// findSnapshotPolicy
func (c *Client) findSnapshotPolicy(workingEnviromentID string, snapshotPolicyName string, clientID string, isSaas bool, connectorIP string) bool {
	resp, err := c.getCVOProperties(workingEnviromentID, clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("cannot find working environment ", workingEnviromentID)
		return false
	}
	snapshotPolicies := resp.SnapshotPolicies
	for i := range snapshotPolicies {
		if snapshotPolicies[i].Name == snapshotPolicyName {
			log.Print("found snapshot policy: ", snapshotPolicyName)
			return true
		}
	}
	log.Print("cannot find snapshot policy ", snapshotPolicyName)
	return false
}
