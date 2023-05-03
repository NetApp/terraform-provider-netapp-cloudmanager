package cloudmanager

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/fatih/structs"
)

type cbsRequest struct {
	Provider string       `structs:"provider"`
	Region   string       `structs:"region"`
	Aws      awsDetails   `structs:"aws,omitempty"`
	Azure    azureDetails `structs:"azure,omitempty"`
	Gcp      gcpDetails   `structs:"gcp,omitempty"`
	// Sgws
	// ontap-s3
	Bucket                  string       `structs:"bucket,omitempty"`
	IPSpace                 string       `structs:"ip-space,omitempty"`
	BackupPolicy            backupPolicy `structs:"backup-policy"`
	AutoBackupEnabled       bool         `structs:"auto-backup-enabled,omitempty"`
	MaxTransferRate         int          `structs:"max-transfer-rate,omitempty"`
	ExportExistingSnapshots bool         `structs:"export-existing-snapshots,omitempty"`

	WorkingEnvironmentID string `structs:"workingEnvironmentId"`
	AccountID            string `structs:"account-id"`
}

type cbsVolumeRequest struct {
	Volume []cbsVolume `structs:"volume,omitempty"`
}

type cbsVolume struct {
	VolumeID     string       `structs:"volume-id"`
	Mode         string       `structs:"mode,omitempty"`
	BackupPolicy backupPolicy `structs:"backup-policy,omitempty"`
}

type awsDetails struct {
	AccountID           string          `structs:"account-id,omitempty"`
	AccessKey           string          `structs:"access-key,omitempty"`
	SecretPassword      string          `structs:"secret-password,omitempty"`
	KmsKeyID            string          `structs:"kms-key-id,omitempty"`
	PrivateEndpoint     privateEndpoint `structs:"private-endpoint,omitempty"`
	ArchiveStorageClass string          `structs:"archive-storage-class,omitempty"`
}

type azureDetails struct {
	ResourceGroup   string          `structs:"resource-group,omitempty"`
	StorageAccount  string          `structs:"storage-account,omitempty"`
	Subscription    string          `structs:"subscription,omitempty"`
	PrivateEndpoint privateEndpoint `structs:"private-endpoint,omitempty"`
	KeyVault        keyVault        `structs:"key-vault,omitempty"`
}

type gcpDetails struct {
	ProjectID      string `structs:"project-id,omitempty"`
	AccessKey      string `structs:"access-key,omitempty"`
	SecretPassword string `structs:"secret-password,omitempty"`
	Kms            kms    `structs:"kms,omitempty"`
}
type privateEndpoint struct {
	ID string `structs:"id"`
}

type backupPolicy struct {
	Name            string        `structs:"name"`
	Rule            []ruleDetails `structs:"rule,omitempty"`
	ArchiveAfteDays string        `structs:"archive-after-days,omitempty"`
	ObjectLock      string        `structs:"object-lock,omitempty"`
	// SgwsArchival sgwsArchival `structs:"sgws-archival"``
}

type ruleDetails struct {
	Label     string `structs:"label"`
	Retention string `structs:"retention"`
}

type keyVault struct {
	KeyVaultID string `structs:"key-vault-id"`
	KeyName    string `structs:"key-name"`
}

type kms struct {
	KeyRingID   string `structs:"key-ring-id"`
	CryptoKeyID string `structs:"crypto-key-id"`
}

type cbsAPICallResult struct {
	ID string `json:"job-id"`
}

type cbsJobResult struct {
	Job []cbsJobDetails `json:"job"`
}

type cbsJobDetails struct {
	ID                   string      `json:"id"`
	WorkingEnvironmentID string      `json:"working-environment-id"`
	JobType              string      `json:"type"`
	JobStatus            string      `json:"status"`
	JobError             string      `json:"error"`
	JobTime              int         `json:"time"`
	Data                 dataDetails `json:"data"`
}

type dataDetails struct {
	MultiVolumeBackup multiVolumeBackupDetail `json:"multi-volume-backup"`
	Restore           restoreDetails          `json:"restore"`
}

type multiVolumeBackupDetail struct {
	Volume []jobVolumeDetails `json:"volume"`
}

type jobVolumeDetails struct {
	ID           string `json:"id"`
	VolumeStatus string `json:"status"`
	VolumeError  string `json:"error"`
}

type restoreDetails struct {
	RestoreType string         `json:"type"`
	Source      sourceDetails  `json:"source"`
	Target      targetDetails  `json:"target"`
	Batch       []batchDetails `json:"batch"`
}

type sourceDetails struct {
	WorkingEnvironmentID   string `json:"working-environment-id"`
	WorkingEnvironmentName string `json:"working-environment-name"`
	Bucket                 string `json:"bucket"`
	VolumeID               string `json:"volume-id"`
	VolumeName             string `json:"volume-name"`
	Snapshot               string `json:"snapshot"`
}

type targetDetails struct {
	WorkingEnvironmentID   string `json:"working-environment-id"`
	WorkingEnvironmentName string `json:"working-environment-name"`
	Svm                    string `json:"svm"`
	VolumeName             string `json:"volume-name"`
	VolumeSize             int    `json:"volume-size"`
	Path                   string `json:"path"`
}

type batchDetails struct {
	ID          string        `json:"id"`
	BatchStatus string        `json:"status"`
	BatchError  string        `json:"error"`
	BatchTime   int           `json:"time"`
	BatchFile   []fileDetails `json:"file"`
}

type fileDetails struct {
	Inode     int    `json:"inode"`
	FilePaht  string `json:"path"`
	FileType  string `json:"type"`
	FileSize  int    `json:"size"`
	FileMtime int    `json:"mtime"`
}

// cbsStatusResult for creating a cbs
type cbsStatusResult struct {
	Name                    string             `json:"name"`
	ID                      string             `json:"id"`
	Region                  string             `json:"region"`
	Status                  string             `json:"status"`
	OntapVersion            string             `json:"ontap-version"`
	BackupEnablementStatus  string             `json:"backup-enablement-status"`
	CBSType                 string             `json:"type"`
	CloudProvider           string             `json:"provider"`
	ProviderAccountID       string             `json:"provider-account-id"`
	ProviderAccountName     string             `json:"provider-account-name"`
	Bucket                  string             `json:"bucket"`
	ArchiveStorageClass     string             `json:"archive-storage-class"`
	ResourceGroup           string             `json:"resource-group"`
	StorageAccount          string             `json:"storage-account"`
	StorageServer           string             `json:"storage-server"`
	UsedCapacityGb          string             `json:"used-capacity-gb"`
	ChargingCapacity        string             `json:"charging-capacity"`
	LogicalUsedSize         string             `json:"logical-used-size"`
	BackedUpVolumeCount     string             `json:"backed-up-volume-count"`
	TotalVolumesCount       string             `json:"total-volumes-count"`
	BackupPolicyCount       string             `json:"backup-policy-count"`
	FailedBackupVolumeCount string             `json:"failed-backup-volume-count"`
	CatalogEnabled          bool               `json:"catalog-enabled"`
	AutoBackupEnabled       bool               `json:"auto-backup-enabled"`
	BackupPolicy            backupPolicyResult `json:"backup-policy"`
	PrivateEndpointRequired bool               `json:"private-endpoint-required"`
	License                 licenseResult      `json:"license"`
	IPSpace                 string             `json:"ip-space"`
	ProviderAccessKey       string             `json:"provider-access-key"`
	DeleteYearlySnapshots   bool               `json:"delete-yearly-snapshots"`
	ExportExistingSnapshots bool               `json:"export-existing-snapshots"`
}

type backupPolicyResult struct {
	Name            string       `json:"name"`
	Rules           []ruleResult `json:"rule"`
	ArchiveAfteDays string       `json:"archive-after-days"`
}

type ruleResult struct {
	Label      string `json:"label"`
	Rentention string `json:"retention"`
}

type licenseResult struct {
	FreeTrialEnd int  `json:"free-trial-end"`
	Eligible     bool `json:"eligible"`
}

type cbsWEResult struct {
	Name                    string             `json:"name"`
	ID                      string             `json:"id"`
	Region                  string             `json:"region"`
	Status                  string             `json:"status"`
	OntapVersion            string             `json:"ontap-version"`
	BackupEnablementStatus  string             `json:"backup-enablement-status"`
	CBSType                 string             `json:"type"`
	CloudProvider           string             `json:"provider"`
	ProviderAccountID       string             `json:"provider-account-id"`
	ProviderAccountName     string             `json:"provider-account-name"`
	Bucket                  string             `json:"bucket"`
	ArchiveStorageClass     string             `json:"archive-storage-class"`
	ResourceGroup           string             `json:"resource-group"`
	StorageAccount          string             `json:"storage-account"`
	StorageServer           string             `json:"storage-server"`
	UsedCapacityGb          string             `json:"used-capacity-gb"`
	ChargingCapacity        string             `json:"charging-capacity"`
	LogicalUsedSize         string             `json:"logical-used-size"`
	BackedUpVolumeCount     string             `json:"backed-up-volume-count"`
	TotalVolumesCount       string             `json:"total-volumes-count"`
	BackupPolicyCount       string             `json:"backup-policy-count"`
	FailedBackupVolumeCount string             `json:"failed-backup-volume-count"`
	CatalogEnabled          bool               `json:"catalog-enabled"`
	AutoBackupEnabled       bool               `json:"auto-backup-enabled"`
	ObjectLock              string             `json:"object-lock"`
	MaxTransferRate         int                `json:"max-transfer-rate"`
	BackupPolicy            backupPolicyResult `json:"backup-policy"`
	PrivateEndpointRequired bool               `json:"private-endpoint-required"`
	License                 licenseResult      `json:"license"`
	IPSpace                 string             `json:"ip-space"`
	ProviderAccessKey       string             `json:"provider-access-key"`
	DeleteYearlySnapshots   bool               `json:"delete-yearly-snapshots"`
	ExportExistingSnapshots bool               `json:"export-existing-snapshots"`
	StorageGridID           string             `json:"storage-grid-id"`
	SgwsArchival            interface{}        `json:"sgws-archival"`
	RemoteMccID             string             `json:"remote-mcc-id"`
}

type cbsGetVolumeResult struct {
	Volume []cbsVolumeResult `json:"volume"`
}

type cbsVolumeResult struct {
	Name          string `json:"name"`
	ID            string `json:"file-system-id"`
	SnapshotCount string `json:"snapshot-count"`
}

type cbsGetSnapshotVolumeResult struct {
	Snapshot []cbsSnapshotVolumeResult `json:"snapshot"`
}

type cbsSnapshotVolumeResult struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// Create working environment cloud backup
func (c *Client) createCBS(cbs cbsRequest, clientID string) (cbsAPICallResult, error) {
	log.Print("createCBS...")

	creationRetryCount := 60
	creationWaitTime := 10

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createCBS request, failed to get AccessToken")
		return cbsAPICallResult{}, err
	}
	c.Token = accessTokenResult.Token
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v3/backup/working-environment/%s", cbs.AccountID, cbs.WorkingEnvironmentID)
	params := structs.Map(cbs)

	log.Printf("\tparams: %+v", params)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createCBS request failed ", statusCode)
		return cbsAPICallResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createCBS")
	if responseError != nil {
		return cbsAPICallResult{}, responseError
	}
	var result cbsAPICallResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createCBS ", err)
		return cbsAPICallResult{}, err
	}
	log.Print("cbsCreate result:", result)
	_, err = c.waitOnJobCompletionCBS(result.ID, cbs, "CBS", "create", creationRetryCount, creationWaitTime, clientID)
	if err != nil {
		return cbsAPICallResult{}, err
	}

	return result, nil
}

// enbale backup for single or multiple volumes
func (c *Client) enableBackupForSingleORMultipleVolumes(cbs cbsRequest, cbsVolume cbsVolumeRequest, clientID string, volumesIDNameMap map[string]map[string]string) (cbsAPICallResult, error) {
	log.Print("enableBackupForSingleORMultipleVolumes...")

	creationRetryCount := 60
	creationWaitTime := 10

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in enableBackupForSingleORMultipleVolumes request, failed to get AccessToken")
		return cbsAPICallResult{}, err
	}
	c.Token = accessTokenResult.Token
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v2/backup/working-environment/%s/volume", cbs.AccountID, cbs.WorkingEnvironmentID)
	params := structs.Map(cbsVolume)

	log.Printf("\tparams: %+v", params)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("enableBackupForSingleORMultipleVolumes request failed ", statusCode)
		return cbsAPICallResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "enableBackupForSingleORMultipleVolumes")
	if responseError != nil {
		return cbsAPICallResult{}, responseError
	}
	var result cbsAPICallResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from enableBackupForSingleORMultipleVolumes ", err)
		return cbsAPICallResult{}, err
	}
	log.Print("enableBackupForSingleORMultipleVolumes result:", result)
	cbsJobStatus, err := c.waitOnJobCompletionCBS(result.ID, cbs, "backup for Volumes", "enable", creationRetryCount, creationWaitTime, clientID)
	if err != nil {
		errVolumeMessage := ""
		noOfVolsFailed := 0
		for _, eachVolume := range cbsJobStatus[0].Data.MultiVolumeBackup.Volume {
			if eachVolume.VolumeStatus == "FAILED" {
				errVolumeMessage += "for volume " + volumesIDNameMap[eachVolume.ID]["name"] + ", " + eachVolume.VolumeError + "\n"
				noOfVolsFailed++
			}
		}
		if errVolumeMessage != "" {
			return cbsAPICallResult{}, fmt.Errorf("%s. Following %v volume have failed:\n%v", err, noOfVolsFailed, errVolumeMessage)
		}
		return cbsAPICallResult{}, err
	}

	return result, nil
}

// Read working environment cloud backup details
func (c *Client) getCBS(cbs cbsRequest, clientID string) (cbsWEResult, error) {
	log.Print("getCBS...")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCBS request, failed to get AccessToken")
		return cbsWEResult{}, err
	}
	c.Token = accessTokenResult.Token
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v1/backup/working-environment/%s", cbs.AccountID, cbs.WorkingEnvironmentID)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getCBS request failed ", statusCode)
		return cbsWEResult{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getCBS")
	if responseError != nil {
		return cbsWEResult{}, responseError
	}
	var result cbsWEResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getCBS ", err)
		return result, err
	}
	log.Printf("\tget CBS: %+v", result)
	if result.BackupEnablementStatus == "ON" {
		return result, nil
	}
	return cbsWEResult{}, fmt.Errorf("working environment %s backup status is %s", cbs.WorkingEnvironmentID, result.BackupEnablementStatus)
}

// Read volume for working environment cloud backup details
func (c *Client) getCBSVolume(cbs cbsRequest, clientID string) ([]cbsVolumeResult, error) {
	log.Print("getCBSVolume...")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in getCBSVolume request, failed to get AccessToken")
		return []cbsVolumeResult{}, err
	}
	c.Token = accessTokenResult.Token
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v1/backup/working-environment/%s/volume", cbs.AccountID, cbs.WorkingEnvironmentID)
	statusCode, response, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("getCBSVolume request failed ", statusCode)
		return []cbsVolumeResult{}, err
	}
	responseError := apiResponseChecker(statusCode, response, "getCBSVolume")
	if responseError != nil {
		return []cbsVolumeResult{}, responseError
	}
	var result cbsGetVolumeResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from getCBSVolume ", err)
		return []cbsVolumeResult{}, err
	}
	if result.Volume == nil {
		return []cbsVolumeResult{}, nil
	}
	log.Printf("\tget Volume: %+v", result)
	return result.Volume, nil
}

// delete snapshot copies (volume)
func (c *Client) deleteSnapshotCopiesVolume(cbs cbsRequest, clientID string, volumeID string) error {
	log.Print("delete snapshot copies volume...")

	jobRetryCount := 60
	jobWaitTime := 10

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteSnapshotCopiesVolume request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v1/backup/working-environment/%s/volume/%s/snapshot", cbs.AccountID, cbs.WorkingEnvironmentID, volumeID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("deleteSnapshotCopiesVolume request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "deleteSnapshotCopiesVolume")
	if responseError != nil {
		return responseError
	}

	var result cbsAPICallResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from deleteSnapshotCopiesVolume ", err)
		return err
	}
	log.Print("deleteSnapshotCopiesVolume result:", result)
	_, err = c.waitOnJobCompletionCBS(result.ID, cbs, "CBS", "deleteSnapshotCopiesVolume", jobRetryCount, jobWaitTime, clientID)
	if err != nil {
		return err
	}
	return nil
}

// delete snapshot copies (working environment)
func (c *Client) deleteSnapshotCopiesWE(cbs cbsRequest, clientID string) error {
	log.Print("delete snapshot copies working environment...")

	jobRetryCount := 60
	jobWaitTime := 10

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in deleteSnapshotCopiesWE request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v2/backup/working-environment/%s/snapshot", cbs.AccountID, cbs.WorkingEnvironmentID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("deleteSnapshotCopiesWE request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "unRegisterWE")
	if responseError != nil {
		return responseError
	}

	var result cbsAPICallResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from deleteSnapshotCopiesWE ", err)
		return err
	}
	log.Print("deleteSnapshotCopiesWE result:", result)
	_, err = c.waitOnJobCompletionCBS(result.ID, cbs, "CBS", "deleteSnapshotCopiesWE", jobRetryCount, jobWaitTime, clientID)
	if err != nil {
		return err
	}
	return nil
}

// unRegisterWE: unregister working environment
func (c *Client) unRegisterWE(cbs cbsRequest, clientID string) error {
	log.Print("unregister working environment...")

	jobRetryCount := 60
	jobWaitTime := 10
	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in unRegisterWE request, failed to get AccessToken")
		return err
	}
	c.Token = accessTokenResult.Token

	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v1/backup/working-environment/%s", cbs.AccountID, cbs.WorkingEnvironmentID)

	hostType := "CloudManagerHost"

	statusCode, response, _, err := c.CallAPIMethod("DELETE", baseURL, nil, c.Token, hostType, clientID)
	if err != nil {
		log.Print("unRegisterWE request failed ", statusCode)
		return err
	}

	responseError := apiResponseChecker(statusCode, response, "unRegisterWE")
	if responseError != nil {
		return responseError
	}
	var result cbsAPICallResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from unRegisterWE ", err)
		return err
	}

	log.Print("unRegisterWE result:", result)
	_, err = c.waitOnJobCompletionCBS(result.ID, cbs, "CBS", "unRegisterWE", jobRetryCount, jobWaitTime, clientID)
	if err != nil {
		return err
	}
	return nil
}

// checkJobStatusCBS: retrieve job status
func (c *Client) checkJobStatusCBS(jobID string, accountID string, workingEnvironmentID string, clientID string) ([]cbsJobDetails, error) {
	log.Printf("checkJobStatusCBS: job-id:%s, act:%s, weid:%s", jobID, accountID, workingEnvironmentID)

	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v1/job/%s", accountID, jobID)

	hostType := "CloudManagerHost"

	var statusCode int
	var response []byte
	networkRetries := 3
	for {
		code, result, _, err := c.CallAPIMethod("GET", baseURL, nil, c.Token, hostType, clientID)
		if err != nil {
			if networkRetries > 0 {
				time.Sleep(1 * time.Second)
				networkRetries--
			} else {
				log.Print("checkJobStatusCBS request failed ", code)
				return nil, err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}

	responseError := apiResponseChecker(statusCode, response, "checkJobStatusCBS")
	if responseError != nil {
		return nil, responseError
	}

	var result cbsJobResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkJobStatusCBS ", err)
		return nil, err
	}

	return result.Job, nil
}

// waitOnJobCompletionCBS: check job completed or not
func (c *Client) waitOnJobCompletionCBS(id string, cbs cbsRequest, actionName string, task string, retries int, waitInterval int, clientID string) ([]cbsJobDetails, error) {
	for {
		cbsJobStatus, err := c.checkJobStatusCBS(id, cbs.AccountID, cbs.WorkingEnvironmentID, clientID)
		if err != nil {
			return cbsJobStatus, err
		}
		if cbsJobStatus[0].JobStatus == "FAILED" {
			return cbsJobStatus, fmt.Errorf("cbs jobID %s WE %s %s %s status FAILED: %s", id, cbs.WorkingEnvironmentID, task, actionName, cbsJobStatus[0].JobError)
		} else if cbsJobStatus[0].JobStatus == "COMPLETED" {
			return cbsJobStatus, nil
		} else if retries == 0 {
			log.Printf("Taking too long to %s %s jobID %s backup status %s", task, actionName, id, cbsJobStatus[0].JobStatus)
			return cbsJobStatus, fmt.Errorf("taking too long for %s %s or not properly setup", actionName, task)
		}
		log.Printf("\tcheck job status %+v", cbsJobStatus)
		log.Printf("Sleep for %d seconds - jobID %s we %s job status %s", waitInterval, id, cbs.WorkingEnvironmentID, cbsJobStatus[0].JobStatus)
		time.Sleep(time.Duration(waitInterval) * time.Second)
		retries--
	}
}
