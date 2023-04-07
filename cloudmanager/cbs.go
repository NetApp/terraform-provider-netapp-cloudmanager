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
	ProjectID      string `struct:"project-id,omitempty"`
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
	Label      string `structs:"label"`
	Retentioin string `structs:"retention"`
}

type keyVault struct {
	KeyVaultID string `structs:"key-vault-id"`
	KeyName    string `structs:"key-name"`
}

type kms struct {
	KeyRingID   string `structs:"key-ring-id"`
	CryptoKeyID string `structs:"crypto-key-id"`
}

type cbsCreateResult struct {
	ID string `json:"job-id"`
}

// fsxStatusResult for creating a fsx
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

func (c *Client) createCBS(cbs cbsRequest, clientID string) (cbsCreateResult, error) {
	log.Print("createCBS...")

	creationRetryCount := 10
	creationWaitTime := 60

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createFSX request, failed to get AccessToken")
		return cbsCreateResult{}, err
	}
	c.Token = accessTokenResult.Token
	hostType := "CloudManagerHost"
	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v3/backup/working-environment/%s", cbs.AccountID, cbs.WorkingEnvironmentID)
	params := structs.Map(cbs)

	log.Printf("\tparams: %+v", params)
	statusCode, response, _, err := c.CallAPIMethod("POST", baseURL, params, c.Token, hostType, clientID)
	if err != nil {
		log.Print("createCBS request failed ", statusCode)
		return cbsCreateResult{}, err
	}

	responseError := apiResponseChecker(statusCode, response, "createCBS")
	if responseError != nil {
		return cbsCreateResult{}, responseError
	}
	var result cbsCreateResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from createCBS ", err)
		return cbsCreateResult{}, err
	}
	log.Print("cbsCreate result:", result)
	err = c.waitOnCompletionCBS(result.ID, cbs, "CBS", "create", creationRetryCount, creationWaitTime, clientID)
	if err != nil {
		return cbsCreateResult{}, err
	}

	return result, nil
}

func (c *Client) getCBS(cbs cbsRequest, clientID string) (cbsWEResult, error) {
	log.Print("getCBS...")

	accessTokenResult, err := c.getAccessToken()
	if err != nil {
		log.Print("in createFSX request, failed to get AccessToken")
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
	return result, nil
}

func (c *Client) unRegisterWE(cbs cbsRequest, clientID string) error {
	log.Print("unregister working environment...")

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

	return nil
}

func (c *Client) checkTaskStatusCBS(id string, accountID string, workingEnvironmentID string, clientID string) (cbsStatusResult, error) {
	log.Printf("checkTaskStatusCBS: %s, %s, %s", id, accountID, workingEnvironmentID)

	baseURL := fmt.Sprintf("/account/%s/providers/cloudmanager_cbs/api/v1/backup/working-environment/%s/status", accountID, workingEnvironmentID)

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
				log.Print("checkTaskStatusCBS request failed ", code)
				return cbsStatusResult{}, err
			}
		} else {
			statusCode = code
			response = result
			break
		}
	}

	responseError := apiResponseChecker(statusCode, response, "checkTaskStatusCBS")
	if responseError != nil {
		return cbsStatusResult{}, responseError
	}

	var result cbsStatusResult
	if err := json.Unmarshal(response, &result); err != nil {
		log.Print("Failed to unmarshall response from checkTaskStatusCBS ", err)
		return cbsStatusResult{}, err
	}

	return result, nil
}

func (c *Client) waitOnCompletionCBS(id string, cbs cbsRequest, actionName string, task string, retries int, waitInterval int, clientID string) error {
	for {
		cbsStatus, err := c.checkTaskStatusCBS(id, cbs.AccountID, cbs.WorkingEnvironmentID, clientID)
		if err != nil {
			return err
		}
		// if cbsStatus.BackupEnablementStatus == "ON" && cbsStatus.Status == "ON" {
		if cbsStatus.BackupEnablementStatus == "ON" {
			return nil
		} else if retries == 0 {
			log.Printf("Taking too long to %s %s weID %s backup status %s", task, actionName, cbs.WorkingEnvironmentID, cbsStatus.BackupEnablementStatus)
			return fmt.Errorf("taking too long for %s to %s or not properly setup", actionName, task)
		}
		log.Printf("\tcheck status %+v", cbsStatus)
		log.Printf("Sleep for %d seconds - we %s status %s backup enablement status %s", waitInterval, cbsStatus.ID, cbsStatus.Status, cbsStatus.BackupEnablementStatus)
		time.Sleep(time.Duration(waitInterval) * time.Second)
		retries--
	}
}

// validateCBSParams validates params
func validateCBSParams(cbsDetails cbsRequest) error {
	if cbsDetails.Provider == "AWS" && (cbsDetails.Aws.AccountID == "" || cbsDetails.Aws.ArchiveStorageClass == "") {
		return fmt.Errorf("account_id and archive_storage_class are required for AWS CBS")
	}
	if cbsDetails.Provider == "AZURE" && cbsDetails.Azure.Subscription == "" {
		return fmt.Errorf("subscription aisrequired for AZURE CBS")
	}
	if cbsDetails.Provider == "GCP" && cbsDetails.Gcp.ProjectID == "" {
		return fmt.Errorf("project_id aisrequired for GCP CBS")
	}
	return nil
}
