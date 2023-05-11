package cloudmanager

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCBS() *schema.Resource {
	return &schema.Resource{
		Create: resourceCBSCreate,
		Read:   resourceCBSRead,
		Delete: resourceCBSDelete,
		// Update: resourceCBSUpdate,
		// Exists: resourceCBSExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		// CustomizeDiff: resourceCBSCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"working_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"working_environment_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cloud_provider": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE", "GCP"}, false),
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"aws_cbs_parameters": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"access_key": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"secret_password": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kms_key_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_endpoint_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"archive_storage_class": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"azure_cbs_parameters": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_group": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"storage_account": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subscription": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_endpoint_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_vault_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"key_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			// commenting out as gcp is not tested yet
			// "gcp_cbs_parameters": {
			// 	Type:     schema.TypeSet,
			// 	MaxItems: 1,
			// 	Optional: true,
			// 	ForceNew: true,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"project_id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"access_key": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"secret_password": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"kms_key_ring_id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 			"kms_crypto_key_id": {
			// 				Type:     schema.TypeString,
			// 				Optional: true,
			// 			},
			// 		},
			// 	},
			// },
			"bucket": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_space": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_policy": {
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"policy_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"retention": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
						"archive_after_days": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"object_lock": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"auto_backup_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"max_transfer_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"export_existing_snapshots": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"volume_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"backup_policy": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"policy_rules": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rule": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"label": {
																Type:     schema.TypeString,
																Optional: true,
																ForceNew: true,
															},
															"retention": {
																Type:     schema.TypeString,
																Optional: true,
																ForceNew: true,
															},
														},
													},
												},
											},
										},
									},
									"archive_after_days": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"object_lock": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCBSCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Enabling cloud backup: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	createCBSRequest := cbsRequest{}
	createCBSVolumeRequest := &cbsVolumeRequest{}
	var volumesIDNameMap = map[string]map[string]string{}

	workingEnv, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	createCBSRequest.WorkingEnvironmentID = workingEnv.PublicID
	createCBSRequest.AccountID = d.Get("account_id").(string)
	createCBSRequest.Region = d.Get("region").(string)

	if a, ok := d.GetOk("cloud_provider"); ok {
		createCBSRequest.Provider = a.(string)
	}
	if a, ok := d.GetOk("region"); ok {
		createCBSRequest.Region = a.(string)
	}
	if a, ok := d.GetOk("bucket"); ok {
		createCBSRequest.Bucket = a.(string)
	}
	if a, ok := d.GetOk("ip_space"); ok {
		createCBSRequest.IPSpace = a.(string)
	}
	if a, ok := d.GetOk("auto_backup_enabled"); ok {
		createCBSRequest.AutoBackupEnabled = a.(bool)
	}
	if a, ok := d.GetOk("max_transfer_rate"); ok {
		createCBSRequest.MaxTransferRate = a.(int)
	}
	if a, ok := d.GetOk("export_existing_snapshots"); ok {
		createCBSRequest.ExportExistingSnapshots = a.(bool)
	}
	if v, ok := d.GetOk("backup_policy"); ok {
		backupPolicy := v.(*schema.Set)
		createCBSRequest.BackupPolicy = expandBackupPolicy(backupPolicy)
	}
	// AWS
	if v, ok := d.GetOk("aws_cbs_parameters"); ok {
		aws := v.(*schema.Set)
		createCBSRequest.Aws = expandAws(aws)
	}
	// AZURE
	if v, ok := d.GetOk("azure_cbs_parameters"); ok {
		azure := v.(*schema.Set)
		createCBSRequest.Azure = expandAzure(azure)
	}
	// commenting out as gcp is not tested yet
	// GCP
	// if v, ok := d.GetOk("gcp_cbs_parameters"); ok {
	// 	gcp := v.(*schema.Set)
	// 	createCBSRequest.Gcp = expandGcp(gcp)
	// }

	if volumes, ok := d.GetOk("volumes"); ok {
		volumesSet := volumes.([]interface{})
		volumesConfigs := make([]cbsVolume, 0, len(volumesSet))
		for _, x := range volumesSet {
			volume := cbsVolume{}
			volumeConfig := x.(map[string]interface{})
			volumeName := volumeConfig["volume_name"].(string)
			volumeRequest := volumeRequest{}
			volumeRequest.Name = volumeName
			volumeRequest.WorkingEnvironmentID = createCBSRequest.WorkingEnvironmentID
			getVolmeDetails, err := client.getVolume(volumeRequest, clientID)
			if err != nil {
				log.Print("Error getting volumes details ", createCBSVolumeRequest.Volume)
				return err
			}
			for _, vol := range getVolmeDetails {
				if vol.Name == volumeName {
					volume.VolumeID = vol.ID
					volumesIDNameMap[vol.ID] = make(map[string]string)
					volumesIDNameMap[vol.ID]["name"] = volumeName
					break
				}
			}
			if volume.VolumeID == "" {
				return fmt.Errorf("error retrieving volumes details: volume %s does not exist", volumeName)
			}
			if m, ok := volumeConfig["mode"]; ok {
				volume.Mode = m.(string)
			}
			if b, ok := volumeConfig["backup_policy"]; ok {
				backupPolicy := b.(*schema.Set)
				volume.BackupPolicy = expandBackupPolicy(backupPolicy)
			}
			volumesConfigs = append(volumesConfigs, volume)
		}
		createCBSVolumeRequest.Volume = volumesConfigs

	}

	res, err := client.createCBS(createCBSRequest, clientID)
	if err != nil {
		log.Print("Error enabling cloud backup on the working environment ", createCBSRequest.WorkingEnvironmentID)
		return err
	}

	d.SetId(res.ID)

	if len(createCBSVolumeRequest.Volume) != 0 {
		res, err = client.enableBackupForSingleORMultipleVolumes(createCBSRequest, *createCBSVolumeRequest, clientID, volumesIDNameMap)
		if err != nil {
			log.Print("Error enabling cloud backup on the volumes ", createCBSVolumeRequest.Volume)
			return err
		}
	}

	log.Printf("Enabled backup cloud: %v", res)
	return resourceCBSRead(d, meta)
}

func resourceCBSRead(d *schema.ResourceData, meta interface{}) error {
	log.Print("Fetching backup cloud...")

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	readCBSRequest := cbsRequest{}
	workingEnv, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	readCBSRequest.WorkingEnvironmentID = workingEnv.PublicID
	readCBSRequest.AccountID = d.Get("account_id").(string)
	res, err := client.getCBS(readCBSRequest, clientID)
	if err != nil {
		log.Print("Error retrieving WE backup details")
		return err
	}
	if res.ID == workingEnv.PublicID {
		return nil
	}
	return fmt.Errorf("error retrieving cloud backup: cloud backup does not exist")
}

func resourceCBSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Disabling cloud backup...")
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)

	workingEnv, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}

	// delete snapshot copies (volumes)
	disableWECBSRequest := cbsRequest{}
	disableWECBSRequest.WorkingEnvironmentID = workingEnv.PublicID
	disableWECBSRequest.AccountID = d.Get("account_id").(string)
	volumes, err := client.getCBSVolume(disableWECBSRequest, clientID)
	if err != nil {
		log.Print("Error retrieving volume backup details")
		return err
	}
	if len(volumes) != 0 {
		for _, eachVolume := range volumes {
			if count, _ := strconv.Atoi(eachVolume.SnapshotCount); count > 0 {
				err = client.deleteSnapshotCopiesVolume(disableWECBSRequest, clientID, eachVolume.ID)
				if err != nil {
					log.Print("Error deleting snapshot copies for volumes: ", volumes)
					return err
				}
			}
		}
	}
	// delete snapshot copies (working environment)
	err = client.deleteSnapshotCopiesWE(disableWECBSRequest, clientID)
	if err != nil {
		log.Print("Error delete snapshot copies working environment", disableWECBSRequest.WorkingEnvironmentID)
		return err
	}
	// unregister a working environment
	log.Print("Unregistering working environment: ", disableWECBSRequest.WorkingEnvironmentID)
	err = client.unRegisterWE(disableWECBSRequest, clientID)
	if err != nil {
		log.Print("Error unregistering working environment", disableWECBSRequest.WorkingEnvironmentID)
		return err
	}
	return nil
}

func expandAws(awsParameterList *schema.Set) awsDetails {
	var params awsDetails

	for _, v := range awsParameterList.List() {
		paramSet := v.(map[string]interface{})
		if v, ok := paramSet["aws_account_id"]; ok {
			params.AccountID = v.(string)
		}
		if v, ok := paramSet["aws_access_key"]; ok {
			params.AccessKey = v.(string)
		}
		if v, ok := paramSet["secret_password"]; ok {
			params.SecretPassword = v.(string)
		}
		if v, ok := paramSet["kms_key_id"]; ok {
			params.KmsKeyID = v.(string)
		}
		if v, ok := paramSet["private_endpoint_id"]; ok {
			params.PrivateEndpoint.ID = v.(string)
		}
		if v, ok := paramSet["archive_storage_class"]; ok {
			params.ArchiveStorageClass = v.(string)
		}
	}
	return params
}

func expandAzure(azureParameterList *schema.Set) azureDetails {
	var params azureDetails

	for _, v := range azureParameterList.List() {
		paramSet := v.(map[string]interface{})
		if v, ok := paramSet["resource_group"]; ok {
			params.ResourceGroup = v.(string)
		}
		if v, ok := paramSet["storage_account"]; ok {
			params.StorageAccount = v.(string)
		}
		if v, ok := paramSet["subscription"]; ok {
			params.Subscription = v.(string)
		}
		if v, ok := paramSet["private_endpoint_id"]; ok {
			params.PrivateEndpoint.ID = v.(string)
		}
		if v, ok := paramSet["key_vault_id"]; ok {
			params.KeyVault.KeyVaultID = v.(string)
		}
		if v, ok := paramSet["key_name"]; ok {
			params.KeyVault.KeyName = v.(string)
		}
	}
	return params
}

// commenting out as gcp is not tested yet
// func expandGcp(gcpParameterList *schema.Set) gcpDetails {
// 	var params gcpDetails

// 	for _, v := range gcpParameterList.List() {
// 		paramSet := v.(map[string]interface{})
// 		if v, ok := paramSet["project_id"]; ok {
// 			params.ProjectID = v.(string)
// 		}
// 		if v, ok := paramSet["access_key"]; ok {
// 			params.AccessKey = v.(string)
// 		}
// 		if v, ok := paramSet["secret_password"]; ok {
// 			params.SecretPassword = v.(string)
// 		}
// 		if v, ok := paramSet["kms_key_ring_id"]; ok {
// 			params.Kms.KeyRingID = v.(string)
// 		}
// 		if v, ok := paramSet["kms_crypto_key_id"]; ok {
// 			params.Kms.CryptoKeyID = v.(string)
// 		}
// 	}
// 	return params
// }

func expandBackupPolicy(backupPolicyList *schema.Set) backupPolicy {
	var params backupPolicy
	for _, v := range backupPolicyList.List() {
		paramSet := v.(map[string]interface{})
		params.Name = paramSet["name"].(string)
		if v, ok := paramSet["archive_after_days"]; ok {
			params.ArchiveAfteDays = v.(string)
		}
		params.ObjectLock = paramSet["object_lock"].(string)
		// rule
		if v, ok := paramSet["policy_rules"]; ok {
			policyRules := v.(*schema.Set)
			for _, v := range policyRules.List() {
				rules := v.(map[string]interface{})
				ruleSet := rules["rule"].([]interface{})
				ruleConfigs := make([]ruleDetails, 0, len(ruleSet))
				for _, x := range ruleSet {
					rule := ruleDetails{}
					ruleConfig := x.(map[string]interface{})
					rule.Label = ruleConfig["label"].(string)
					rule.Retention = ruleConfig["retention"].(string)
					ruleConfigs = append(ruleConfigs, rule)
				}
				params.Rule = ruleConfigs
				log.Print("rules: ", params.Rule)
			}
		}
		// sgws-archival
	}
	log.Print("params:", params)
	return params
}
