package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/structs"
	"github.com/hashicorp/terraform/helper/schema"
)

// GiBToBytes converting GB to bytes
const GiBToBytes = 1024 * 1024 * 1024

// TiBToGiB converting TiB to GiB
const TiBToGiB = 1024

type exportPolicyRule struct {
	AllowedClients string `structs:"allowedClients"`
	Nfsv3          bool   `structs:"nfsv3"`
	Nfsv4          bool   `structs:"nfsv4"`
	RuleIndex      int    `structs:"ruleIndex"`
	UnixReadOnly   bool   `structs:"unixReadOnly"`
	UnixReadWrite  bool   `structs:"unixReadWrite"`
}

type exportPolicyRuleResponse struct {
	AllowedClients string `json:"allowedClients"`
	RuleIndex      int    `json:"ruleIndex"`
	UnixReadOnly   bool   `json:"unixReadOnly"`
	UnixReadWrite  bool   `json:"unixReadWrite"`
	Nfsv3          nfs    `json:"nfsv3"`
	Nfsv4          nfs    `json:"nfsv4"`
}

type nfs struct {
	Checked bool `structs:"checked"`
}

type exportPolicy struct {
	Rules []exportPolicyRule `structs:"rules"`
}

type gcpVolumeRequest struct {
	Name                   string             `structs:"name,omitempty"`
	Region                 string             `structs:"region,omitempty"`
	VolumePath             string             `structs:"volumePath,omitempty"`
	ProtocolTypes          []string           `structs:"protocolTypes,omitempty"`
	Network                string             `structs:"networkName,omitempty"`
	Size                   float64            `structs:"quotaInBytes,omitempty"`
	ServiceLevel           string             `structs:"serviceLevel,omitempty"`
	SnapshotPolicy         snapshotPolicy     `structs:"snapshotPolicy,omitempty"`
	ExportPolicy           []exportPolicyRule `structs:"rules"`
	VolumeID               string             `structs:"volumeId,omitempty"`
	WorkingEnvironmentName string             `structs:"workingEnvironmentName"`
}

type gcpVolumeResponse struct {
	Name                  string         `json:"name,omitempty"`
	Region                string         `json:"region,omitempty"`
	CreationToken         string         `json:"creationToken,omitempty"`
	ProtocolTypes         []string       `json:"protocolTypes,omitempty"`
	Network               string         `json:"network,omitempty"`
	Size                  int            `json:"quotaInBytes,omitempty"`
	ServiceLevel          string         `json:"serviceLevel,omitempty"`
	SnapshotPolicy        snapshotPolicy `json:"snapshotPolicy,omitempty"`
	ExportPolicy          rules          `json:"exportPolicy,omitempty"`
	VolumeID              string         `json:"volumeId,omitempty"`
	LifeCycleState        string         `json:"lifeCycleState"`
	LifeCycleStateDetails string         `json:"lifeCycleStateDetails"`
	Zone                  string         `json:"zone,omitempty"`
	StorageClass          string         `json:"storageClass,omitempty"`
	TypeDP                bool           `json:"isDataProtection,omitempty"`
	MountPoints           []mountPoints  `json:"mountPoints,omitempty"`
}

// Get volume API returns different struct of export policy from create volume API.
type rules struct {
	Rules []exportPolicyRuleResponse `json:"rules"`
}

type mountPoints struct {
	Export       string `structs:"export"`
	Server       string `structs:"server"`
	ProtocolType string `structs:"protocolType"`
}

type snapshotPolicy struct {
	Enabled         bool            `structs:"enabled"`
	DailySchedule   dailySchedule   `structs:"dailySchedule"`
	HourlySchedule  hourlySchedule  `structs:"hourlySchedule"`
	MonthlySchedule monthlySchedule `structs:"monthlySchedule"`
	WeeklySchedule  weeklySchedule  `structs:"weeklySchedule"`
}

type dailySchedule struct {
	Hour            int `structs:"hour"`
	Minute          int `structs:"minute"`
	SnapshotsToKeep int `structs:"snapshotsToKeep"`
}

type hourlySchedule struct {
	Minute          int `structs:"minute"`
	SnapshotsToKeep int `structs:"snapshotsToKeep"`
}

type monthlySchedule struct {
	DaysOfMonth     string `structs:"daysOfMonth"`
	Hour            int    `structs:"hour"`
	Minute          int    `structs:"minute"`
	SnapshotsToKeep int    `structs:"snapshotsToKeep"`
}

type weeklySchedule struct {
	Day             string `structs:"day"`
	Hour            int    `structs:"hour"`
	Minute          int    `structs:"minute"`
	SnapshotsToKeep int    `structs:"snapshotsToKeep"`
}

func (c *Client) createGCPVolume(vol gcpVolumeRequest, info cvsInfo, clientID string) (gcpVolumeResponse, error) {
	baseURL, err := c.getCVSAPIRoot(info.AccountName, vol.WorkingEnvironmentName, clientID)
	if err != nil {
		return gcpVolumeResponse{}, err
	}
	baseURL = fmt.Sprintf("%s/locations/%s/volumes", baseURL, vol.Region)
	hostType := "CVSHost"
	param := structs.Map(vol)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, param, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createGCPVolume request failed ", statusCode)
		return gcpVolumeResponse{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "createGCPVolume")
	if responseError != nil {
		return gcpVolumeResponse{}, responseError
	}
	var result gcpVolumeResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createGCPVolume", err)
		return gcpVolumeResponse{}, err
	}
	return result, nil
}

func (c *Client) deleteGCPVolume(vol gcpVolumeRequest, info cvsInfo, clientID string) error {
	baseURL, err := c.getCVSAPIRoot(info.AccountName, vol.WorkingEnvironmentName, clientID)
	if err != nil {
		return err
	}
	baseURL = fmt.Sprintf("%s/locations/%s/volumes/%s", baseURL, vol.Region, vol.VolumeID)
	hostType := "CVSHost"
	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("deleteGCPVolume request failed ", statusCode)
		return err
	}
	responseError := apiResponseChecker(statusCode, response, "deleteGCPVolume")
	if responseError != nil {
		return responseError
	}
	return nil
}

func (c *Client) getGCPVolume(vol gcpVolumeRequest, info cvsInfo, clientID string) (gcpVolumeResponse, error) {
	baseURL, err := c.getCVSAPIRoot(info.AccountName, vol.WorkingEnvironmentName, clientID)
	if err != nil {
		return gcpVolumeResponse{}, err
	}
	baseURL = fmt.Sprintf("%s/locations/%s/volumes/%s", baseURL, vol.Region, vol.VolumeID)
	hostType := "CVSHost"
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getGCPVolume request failed ", statusCode)
		return gcpVolumeResponse{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getGCPVolume")
	if responseError != nil {
		return gcpVolumeResponse{}, responseError
	}

	var result gcpVolumeResponse
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getGCPVolume ", err)
		return gcpVolumeResponse{}, err
	}

	return result, nil
}

func expandSnapshotPolicy(data map[string]interface{}) snapshotPolicy {
	snapshotPolicy := snapshotPolicy{}

	if v, ok := data["enabled"]; ok {
		snapshotPolicy.Enabled = v.(bool)
	}

	if v, ok := data["daily_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			dailySchedule := v.([]interface{})[0].(map[string]interface{})
			if hour, ok := dailySchedule["hour"]; ok {
				snapshotPolicy.DailySchedule.Hour = hour.(int)
			}
			if minute, ok := dailySchedule["minute"]; ok {
				snapshotPolicy.DailySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := dailySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.DailySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["hourly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			hourlySchedule := v.([]interface{})[0].(map[string]interface{})
			if minute, ok := hourlySchedule["minute"]; ok {
				snapshotPolicy.HourlySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := hourlySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.HourlySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["monthly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			monthlySchedule := v.([]interface{})[0].(map[string]interface{})
			if daysOfMonth, ok := monthlySchedule["days_of_month"]; ok {
				snapshotPolicy.MonthlySchedule.DaysOfMonth = daysOfMonth.(string)
			}
			if hour, ok := monthlySchedule["hour"]; ok {
				snapshotPolicy.MonthlySchedule.Hour = hour.(int)
			}
			if minute, ok := monthlySchedule["minute"]; ok {
				snapshotPolicy.MonthlySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := monthlySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.MonthlySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	if v, ok := data["weekly_schedule"]; ok {
		if len(v.([]interface{})) > 0 {
			weeklySchedule := v.([]interface{})[0].(map[string]interface{})
			if day, ok := weeklySchedule["day"]; ok {
				snapshotPolicy.WeeklySchedule.Day = day.(string)
			}
			if hour, ok := weeklySchedule["hour"]; ok {
				snapshotPolicy.WeeklySchedule.Hour = hour.(int)
			}
			if minute, ok := weeklySchedule["minute"]; ok {
				snapshotPolicy.WeeklySchedule.Minute = minute.(int)
			}
			if snapshotsToKeep, ok := weeklySchedule["snapshots_to_keep"]; ok {
				snapshotPolicy.WeeklySchedule.SnapshotsToKeep = snapshotsToKeep.(int)
			}
		}
	}
	return snapshotPolicy
}

// expandExportPolicy converts set to exportPolicy struct
func expandExportPolicy(set *schema.Set) []exportPolicyRule {
	exportPolicy := exportPolicy{}
	for _, v := range set.List() {
		log.Printf("here here here here : %#v", v)
		rules := v.(map[string]interface{})
		ruleSet := rules["rule"].(*schema.Set).List()
		ruleConfigs := make([]exportPolicyRule, 0, len(ruleSet))
		for _, x := range ruleSet {
			exportPolicyRule := exportPolicyRule{}
			ruleConfig := x.(map[string]interface{})
			exportPolicyRule.AllowedClients = ruleConfig["allowed_clients"].(string)
			exportPolicyRule.RuleIndex = ruleConfig["rule_index"].(int)
			exportPolicyRule.UnixReadOnly = ruleConfig["unix_read_only"].(bool)
			exportPolicyRule.UnixReadWrite = ruleConfig["unix_read_write"].(bool)
			exportPolicyRule.Nfsv3 = ruleConfig["nfsv3"].(bool)
			exportPolicyRule.Nfsv4 = ruleConfig["nfsv4"].(bool)
			ruleConfigs = append(ruleConfigs, exportPolicyRule)
			exportPolicy.Rules = append(exportPolicy.Rules, exportPolicyRule)
		}
	}
	return exportPolicy.Rules
}

// flattenExportPolicy converts exportPolicy struct to []map[string]interface{}
func flattenExportPolicy(v rules) interface{} {
	rules := make([]map[string]interface{}, 0, len(v.Rules))
	for _, exportPolicyRule := range v.Rules {
		ruleMap := make(map[string]interface{})
		ruleMap["allowed_clients"] = exportPolicyRule.AllowedClients
		ruleMap["unix_read_only"] = exportPolicyRule.UnixReadOnly
		ruleMap["unix_read_write"] = exportPolicyRule.UnixReadWrite
		ruleMap["rule_index"] = exportPolicyRule.RuleIndex
		ruleMap["nfsv3"] = exportPolicyRule.Nfsv3.Checked
		ruleMap["nfsv4"] = exportPolicyRule.Nfsv4.Checked
		rules = append(rules, ruleMap)
	}
	result := make([]map[string]interface{}, 1)
	result[0] = make(map[string]interface{})
	result[0]["rule"] = rules
	return result
}

// flattenSnapshotPolicy converts snapshotPolicy struct to []map[string]interface{}
func flattenSnapshotPolicy(v snapshotPolicy) interface{} {
	flattened := make([]map[string]interface{}, 1)
	sp := make(map[string]interface{})
	sp["enabled"] = v.Enabled
	hourly := make([]map[string]interface{}, 1)
	hourly[0] = make(map[string]interface{})
	hourly[0]["minute"] = v.HourlySchedule.Minute
	hourly[0]["snapshots_to_keep"] = v.HourlySchedule.SnapshotsToKeep
	daily := make([]map[string]interface{}, 1)
	daily[0] = make(map[string]interface{})
	daily[0]["hour"] = v.DailySchedule.Hour
	daily[0]["minute"] = v.DailySchedule.Minute
	daily[0]["snapshots_to_keep"] = v.DailySchedule.SnapshotsToKeep
	monthly := make([]map[string]interface{}, 1)
	monthly[0] = make(map[string]interface{})
	monthly[0]["days_of_month"] = v.MonthlySchedule.DaysOfMonth
	monthly[0]["hour"] = v.MonthlySchedule.Hour
	monthly[0]["minute"] = v.MonthlySchedule.Minute
	monthly[0]["snapshots_to_keep"] = v.MonthlySchedule.SnapshotsToKeep
	weekly := make([]map[string]interface{}, 1)
	weekly[0] = make(map[string]interface{})
	weekly[0]["day"] = v.WeeklySchedule.Day
	weekly[0]["hour"] = v.WeeklySchedule.Hour
	weekly[0]["minute"] = v.WeeklySchedule.Minute
	weekly[0]["snapshots_to_keep"] = v.WeeklySchedule.SnapshotsToKeep
	sp["daily_schedule"] = daily
	sp["hourly_schedule"] = hourly
	sp["weekly_schedule"] = weekly
	sp["monthly_schedule"] = monthly
	flattened[0] = sp
	return flattened
}
