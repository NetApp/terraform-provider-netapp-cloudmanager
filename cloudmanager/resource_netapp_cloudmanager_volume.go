package cloudmanager

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCVOVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOVolumeCreate,
		Read:   resourceCVOVolumeRead,
		Delete: resourceCVOVolumeDelete,
		Exists: resourceCVOVolumeExists,
		Update: resourceCVOVolumeUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceVolumeCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"working_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"working_environment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"svm_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"aggregate_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Required: true,
			},
			"snapshot_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"enable_thin_provisioning": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_compression": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_deduplication": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"export_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"export_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"export_policy_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"export_policy_nfs_version": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"export_policy_rule_access_control": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"readonly", "readwrite", "none"}, false),
			},
			"export_policy_rule_super_user": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"throughput": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"provider_volume_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"S3", "Blob", "cloudStorage", "none"}, false),
			},
			"tiering_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "snapshot_only", "auto", "all"}, false),
				Default:      "auto",
			},
			"volume_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"nfs", "cifs", "iscsi"}, false),
				Default:      "nfs",
			},
			"share_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"igroups": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"os_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"initiator": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alias": {
							Type:     schema.TypeString,
							Required: true,
						},
						"iqn": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"snapshot_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule_type": {
										Type:         schema.TypeString,
										ValidateFunc: validation.StringInSlice([]string{"5min", "8hour", "hourly", "daily", "weekly", "monthly"}, true),
										Required:     true,
										ForceNew:     true,
									},
									"retention": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceCVOVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %s", d.Get("name").(string))

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	var svm string
	var workingEnvironmentType string
	var createAggregateifNotExists bool
	volume := volumeRequest{}
	quote := quoteRequest{}
	// quote volume
	quote.Name = d.Get("name").(string)
	quote.Size.Size = d.Get("size").(float64)
	quote.Size.Unit = d.Get("unit").(string)
	quote.SnapshotPolicyName = d.Get("snapshot_policy_name").(string)
	quote.ProviderVolumeType = d.Get("provider_volume_type").(string)
	if v, ok := d.GetOk("enable_deduplication"); ok {
		quote.EnableDeduplication = v.(bool)
		volume.EnableDeduplication = v.(bool)
	}
	if v, ok := d.GetOk("enable_thin_provisioning"); ok {
		quote.EnableThinProvisioning = v.(bool)
		volume.EnableThinProvisioning = v.(bool)
	}
	if v, ok := d.GetOk("enable_compression"); ok {
		quote.EnableCompression = v.(bool)
		volume.EnableCompression = v.(bool)
	}
	quote.VerifyNameUniqueness = true // hard code to always true
	if v, ok := d.GetOk("iops"); ok {
		quote.Iops = v.(int)
	}
	if v, ok := d.GetOk("throughput"); ok {
		quote.Throughput = v.(int)
	}
	if v, ok := d.GetOk("aggregate_name"); ok {
		quote.AggregateName = v.(string)
		volume.AggregateName = v.(string)
		createAggregateifNotExists = false
	} else {
		createAggregateifNotExists = true
	}
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	}

	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	volume.WorkingEnvironmentID = weInfo.PublicID
	volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	if svm == "" {
		if weInfo.SvmName != "" {
			svm = weInfo.SvmName
		} else {
			svm = "svm_" + weInfo.Name
		}
	}
	volume.SvmName = svm
	workingEnvironmentType = weInfo.WorkingEnvironmentType
	volume.WorkingEnvironmentType = workingEnvironmentType

	if workingEnvironmentType != "ON_PREM" {
		// Check if snapshot_nolicy_name exists
		if !client.findSnapshotPolicy(weInfo.PublicID, quote.SnapshotPolicyName, clientID) {
			// If snapshot_policy_name does not exist, create the snapshot policy
			if v, ok := d.GetOk("snapshot_policy"); ok {
				policy := v.(*schema.Set)
				if policy.Len() > 0 {
					err := client.createSnapshotPolicy(weInfo.PublicID, quote.SnapshotPolicyName, policy, clientID)
					if err != nil {
						return err
					}
				}
			}
		}
		quote.WorkingEnvironmentType = workingEnvironmentType
		quote.WorkingEnvironmentID = weInfo.PublicID
		quote.SvmName = svm
		if v, ok := d.GetOk("capacity_tier"); ok {
			if v.(string) != "none" {
				quote.CapacityTier = v.(string)
				volume.CapacityTier = v.(string)
				if v, ok = d.GetOk("tiering_policy"); ok {
					if v.(string) != "none" {
						quote.TieringPolicy = v.(string)
					}
				}
			}
		}
		// exmaple of the export policy info in volume creation from the UI timeline. Be aware that it might be out of date as time changes.
		// Note that for update, export policy info has different structure.
		// "exportPolicyInfo": {
		// 	"policyType": "custom",
		// 	"rules": [
		// 	  {
		// 		"nfsVersion": [
		// 		  "nfs4",
		// 		  "nfs3"
		// 		],
		// 		"superuser": false,
		// 		"ruleAccessControl": "readonly",
		// 		"ips": [
		// 		  "0.0.0.0"
		// 		],
		// 		"index": 1
		// 	  }
		// 	]
		//   }
		if d.HasChange("export_policy_name") || d.HasChange("export_policy_type") || d.HasChange("export_policy_ip") || d.HasChange("export_policy_nfs_version") || d.HasChange("export_policy_rule_access_control") || d.HasChange("export_policy_rule_super_user") {
			var exportPolicyTypeOK, exportPolicyIPOK, exportPolicyNfsVersionOK, exportPolicyRuleAccessControlOK, exportPolicyRuleSuperUserOK bool
			var nfsVersion, policyIps []string
			if v, ok := d.GetOk("export_policy_name"); ok {
				volume.ExportPolicyInfo.Name = v.(string)
			}
			if v, ok := d.GetOk("export_policy_ip"); ok {
				ips := make([]string, 0, len(v.([]interface{})))
				for _, x := range v.([]interface{}) {
					ips = append(ips, x.(string))
				}
				policyIps = ips
				exportPolicyIPOK = true
			}
			if v, ok := d.GetOk("export_policy_nfs_version"); ok {
				nfs := make([]string, 0, v.(*schema.Set).Len())
				for _, x := range v.(*schema.Set).List() {
					nfs = append(nfs, x.(string))
				}
				nfsVersion = nfs
				exportPolicyNfsVersionOK = true
			}
			if _, ok := d.GetOk("export_policy_rule_access_control"); ok {
				exportPolicyRuleAccessControlOK = true
			}
			if _, ok := d.GetOkExists("export_policy_rule_super_user"); ok {
				exportPolicyRuleSuperUserOK = true
			}
			if v, ok := d.GetOk("export_policy_type"); ok {
				volume.ExportPolicyInfo.PolicyType = v.(string)
				exportPolicyTypeOK = true
			}
			if !exportPolicyTypeOK || !exportPolicyIPOK || !exportPolicyNfsVersionOK || !exportPolicyRuleAccessControlOK || !exportPolicyRuleSuperUserOK {
				return fmt.Errorf("export_policy_type, export_policy_ip, export_policy_nfs_version, export_policy_rule_access_control and export_policy_rule_super_user are required for export policy")
			}
			var rules []ExportPolicyRule
			rules = make([]ExportPolicyRule, len(policyIps))
			for i, x := range policyIps {
				rules[i] = ExportPolicyRule{}
				eachRule := make([]string, 1)
				eachRule[0] = x
				rules[i].Ips = eachRule
				rules[i].NfsVersion = nfsVersion
				rules[i].Superuser = d.Get("export_policy_rule_super_user").(bool)
				rules[i].RuleAccessControl = d.Get("export_policy_rule_access_control").(string)
				rules[i].Index = int32(i + 1)
			}
			volume.ExportPolicyInfo.Rules = rules
		}

		response, err := client.quoteVolume(quote, clientID)
		if err != nil {
			log.Printf("Error quoting volume")
			return err
		}
		volume.NewAggregate = response["newAggregate"].(bool)
		volume.AggregateName = response["aggregateName"].(string)
		volume.NumOfDisks = response["numOfDisks"].(float64)
	}

	volume.ProviderVolumeType = d.Get("provider_volume_type").(string)
	volume.Name = d.Get("name").(string)
	volume.SnapshotPolicyName = d.Get("snapshot_policy_name").(string)
	volume.Size.Size = d.Get("size").(float64)
	volume.Size.Unit = d.Get("unit").(string)
	volumeProtocol := d.Get("volume_protocol").(string)
	if v, ok := d.GetOk("comment"); ok {
		volume.Comment = v.(string)
	}
	if v, ok := d.GetOk("iops"); ok {
		volume.Iops = v.(int)
	}
	if v, ok := d.GetOk("throughput"); ok {
		volume.Throughput = v.(int)
	}
	if o, ok := d.GetOk("tags"); ok {
		tags := make([]volumeTag, 0)
		for k, v := range o.(map[string]interface{}) {
			tag := volumeTag{}
			tag.TagKey = k
			tag.TagValue = v.(string)
			tags = append(tags, tag)
		}
		volume.VolumeTags = tags
	}
	if volumeProtocol == "cifs" {
		exist, err := client.checkCifsExists(workingEnvironmentType, volume.WorkingEnvironmentID, volume.SvmName, clientID)
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
	} else if volumeProtocol == "iscsi" {
		isNewIgroup, _, err := createIscsiVolumeHelper(d, meta)
		if err != nil {
			return err
		}
		if v, ok := d.GetOk("os_name"); ok {
			volume.IscsiInfo.OsName = v.(string)
		}
		if isNewIgroup {
			log.Print("Need to create igroup")
			igroups := d.Get("igroups").(*schema.Set)
			if igroups.Len() > 1 {
				return fmt.Errorf("can not create more than one new igroup")
			}
			if _, ok := d.GetOk("initiator"); !ok {
				return fmt.Errorf("initiator is required when creating new igroup")
			}
			volume.IscsiInfo.IgroupCreationRequest.IgroupName = igroups.List()[0].(string)
			if v, ok := d.GetOk("initiator"); ok {
				ini := v.(*schema.Set)
				if ini.Len() > 0 {
					initiators := make([]string, 0, ini.Len())
					for _, v := range expandInitiator(ini) {
						initiators = append(initiators, v.Iqn)
					}
					volume.IscsiInfo.IgroupCreationRequest.Initiators = initiators
				}
			}
		} else {
			if v, ok := d.GetOk("igroups"); ok {
				igroups := make([]string, 0, v.(*schema.Set).Len())
				for _, x := range v.(*schema.Set).List() {
					igroups = append(igroups, x.(string))
				}
				volume.IscsiInfo.Igroups = igroups
			}
		}
	}
	volume.WorkingEnvironmentType = workingEnvironmentType
	err = client.createVolume(volume, createAggregateifNotExists, clientID)
	if err != nil {
		log.Print("Error creating volume")
		return err
	}
	res, err := client.getVolume(volume, clientID)
	if err != nil {
		log.Print("Error reading volume after creation")
		return err
	}
	for _, volume := range res {
		if volume.SvmName == svm && volume.Name == d.Get("name") {
			d.SetId(volume.ID)
			break
		}
	}

	return resourceCVOVolumeRead(d, meta)
}

func resourceCVOVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching volume: %s", d.Get("name").(string))

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	}

	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	volume.WorkingEnvironmentID = weInfo.PublicID
	volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	if svm == "" {
		if weInfo.SvmName != "" {
			svm = weInfo.SvmName
		} else {
			svm = "svm_" + weInfo.Name
		}
	}
	volume.SvmName = svm

	res, err := client.getVolume(volume, clientID)
	if err != nil {
		log.Print("Error reading volume")
		return err
	}
	for _, volume := range res {
		if volume.ID == d.Id() {
			if _, ok := d.GetOk("aggregate_name"); ok {
				d.Set("aggregate_name", volume.AggregateName)
			}
			if _, ok := d.GetOk("snapshot_policy_name"); ok {
				d.Set("snapshot_policy_name", volume.SnapshotPolicyName)
			}
			if _, ok := d.GetOk("enable_thin_provisioning"); ok {
				d.Set("enable_thin_provisioning", volume.EnableThinProvisioning)
			}
			// setting this two in create are not working right now.
			// if _, ok := d.GetOk("enable_deduplication"); ok {
			// 	d.Set("enable_deduplication", volume.EnableDeduplication)
			// }
			// if _, ok := d.GetOk("enable_compression"); ok {
			// 	d.Set("enable_compression", volume.EnableCompression)
			// }
			if _, ok := d.GetOk("export_policy_ip"); ok {
				d.Set("export_policy_ip", volume.ExportPolicyInfo.Ips)
			}
			if _, ok := d.GetOk("export_policy_nfs_version"); ok {
				d.Set("export_policy_nfs_version", volume.ExportPolicyInfo.NfsVersion)
			}
			if _, ok := d.GetOk("export_policy_type"); ok {
				d.Set("export_policy_type", volume.ExportPolicyInfo.PolicyType)
			}
			if v, ok := d.GetOk("provider_volume_type"); ok {
				d.Set("provider_volume_type", v.(string))
			}
			if _, ok := d.GetOk("comment"); ok {
				d.Set("comment", volume.Comment)
			}
			if v, ok := d.GetOk("capacity_tier"); ok {
				if v.(string) != "none" {
					d.Set("capacity_tier", volume.CapacityTier)
					if v, ok = d.GetOk("tiering_policy"); ok {
						if v.(string) != "none" {
							d.Set("tiering_policy", volume.TieringPolicy)
						}
					}
				}
			}
			if _, ok := d.GetOk("export_policy_name"); ok {
				d.Set("export_policy_name", volume.ExportPolicyInfo.Name)
			}
			if d.Get("unit") != "GB" {
				d.Set("size", convertSizeUnit(volume.Size.Size, volume.Size.Unit, d.Get("unit").(string)))
				d.Set("unit", d.Get("unit").(string))
			} else {
				d.Set("size", volume.Size.Size)
				d.Set("unit", volume.Size.Unit)
			}
			if d.Get("volume_protocol") == "cifs" {
				if _, ok := d.GetOk("share_name"); ok {
					if len(volume.ShareInfo) > 0 {
						d.Set("share_name", volume.ShareInfo[0].ShareName)
					}
				}
				if _, ok := d.GetOk("permission"); ok {
					if len(volume.ShareInfo) > 0 {
						if len(volume.ShareInfo[0].AccessControlList) > 0 {
							d.Set("permission", volume.ShareInfo[0].AccessControlList[0].Permission)
						}
					}
				}
				if _, ok := d.GetOk("users"); ok {
					if len(volume.ShareInfo) > 0 {
						if len(volume.ShareInfo[0].AccessControlList) > 0 {
							d.Set("users", volume.ShareInfo[0].AccessControlList[0].Users)
						}
					}
				}
			}
			return nil
		}
	}

	return fmt.Errorf("error reading volume: volume doesn't exist")
}

func resourceCVOVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume: %s", d.Get("name").(string))
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	}

	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	volume.WorkingEnvironmentID = weInfo.PublicID
	volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	if svm == "" {
		if weInfo.SvmName != "" {
			svm = weInfo.SvmName
		} else {
			svm = "svm_" + weInfo.Name
		}
	}
	volume.SvmName = svm

	volume.Name = d.Get("name").(string)

	err = client.deleteVolume(volume, clientID)
	if err != nil {
		log.Print("Error deleting volume")
		return err
	}
	return nil
}

func resourceCVOVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume: %s", d.Get("name").(string))
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := volumeRequest{}
	volume.Name = d.Get("name").(string)
	volume.ID = d.Id()

	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return false, fmt.Errorf("cannot find working environment")
	}
	volume.WorkingEnvironmentID = weInfo.PublicID
	volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType

	res, err := client.getVolumeByID(volume, clientID)
	if err != nil {
		log.Print("Error reading volume")
		return false, err
	}

	if res.ID != d.Id() {
		d.SetId("")
		return false, nil
	}
	return true, nil
}

func resourceCVOVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume: %s", d.Get("name").(string))
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	volume.Name = d.Get("name").(string)
	if d.HasChange("export_policy_ip") || d.HasChange("export_policy_nfs_version") || d.HasChange("export_policy_rule_super_user") || d.HasChange("export_policy_rule_access_control") {
		var exportPolicyTypeOK, exportPolicyIPOK, exportPolicyNfsVersionOK, exportPolicyRuleAccessControlOK, exportPolicyRuleSuperUserOK bool
		if v, ok := d.GetOk("export_policy_name"); ok {
			volume.ExportPolicyInfo.Name = v.(string)
		}
		if v, ok := d.GetOk("export_policy_nfs_version"); ok {
			nfsVersions := make([]string, 0, v.(*schema.Set).Len())
			for _, x := range v.(*schema.Set).List() {
				nfsVersions = append(nfsVersions, x.(string))
			}
			volume.ExportPolicyInfo.NfsVersion = nfsVersions
			exportPolicyNfsVersionOK = true
		}
		if v, ok := d.GetOk("export_policy_type"); ok {
			volume.ExportPolicyInfo.PolicyType = v.(string)
			exportPolicyTypeOK = true
		}
		if v, ok := d.GetOk("export_policy_ip"); ok {
			ips := make([]string, 0, len(v.([]interface{})))
			for _, x := range v.([]interface{}) {
				ips = append(ips, x.(string))
			}
			volume.ExportPolicyInfo.Ips = ips
			exportPolicyIPOK = true
		}
		if _, ok := d.GetOkExists("export_policy_rule_super_user"); ok {
			exportPolicyRuleSuperUserOK = true
		}
		if _, ok := d.GetOk("export_policy_rule_access_control"); ok {
			exportPolicyRuleAccessControlOK = true
		}
		// example of export poliy info for update.
		// "exportPolicyInfo": {
		// 	"ips": [
		// 	  "0.0.0.0",
		// 	  "0.0.0.1"
		// 	],
		// 	"nfsVersion": [
		// 	  "nfs4",
		// 	  "nfs3"
		// 	],
		// 	"policyType": "custom",
		// 	"rules": [
		// 	  {
		// 		"nfsVersion": [
		// 		  "nfs4",
		// 		  "nfs3"
		// 		],
		// 		"superuser": false,
		// 		"ruleAccessControl": "readonly",
		// 		"ips": [
		// 		  "10.0.0.0"
		// 		],
		// 		"index": 1
		// 	  },
		// 	  {
		// 		"nfsVersion": [
		// 		  "nfs4",
		// 		  "nfs3"
		// 		],
		// 		"superuser": false,
		// 		"ruleAccessControl": "readonly",
		// 		"ips": [
		// 		  "10.0.0.1"
		// 		],
		// 		"index": 2
		// 	  }
		// 	]
		//   }
		if !exportPolicyTypeOK || !exportPolicyIPOK || !exportPolicyNfsVersionOK || !exportPolicyRuleAccessControlOK || !exportPolicyRuleSuperUserOK {
			return fmt.Errorf("export_policy_type, export_policy_ip, export_policy_nfs_version, export_policy_rule_access_control and export_policy_rule_super_user are required for export policy")
		}
		var rules []ExportPolicyRule
		rules = make([]ExportPolicyRule, len(volume.ExportPolicyInfo.Ips))
		for i, x := range volume.ExportPolicyInfo.Ips {
			rules[i] = ExportPolicyRule{}
			eachRule := make([]string, 1)
			eachRule[0] = x
			rules[i].Ips = eachRule
			rules[i].NfsVersion = volume.ExportPolicyInfo.NfsVersion
			rules[i].Superuser = d.Get("export_policy_rule_super_user").(bool)
			rules[i].RuleAccessControl = d.Get("export_policy_rule_access_control").(string)
			rules[i].Index = int32(i + 1)
		}
		volume.ExportPolicyInfo.Rules = rules

	}

	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	}
	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	volume.WorkingEnvironmentID = weInfo.PublicID
	volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	if svm == "" {
		svm = weInfo.SvmName
	}
	volume.SvmName = svm

	if d.HasChange("export_policy_name") {
		volume.ExportPolicyInfo.Name = d.Get("export_policy_name").(string)
	}
	if v, ok := d.GetOk("export_policy_nfs_version"); ok {
		nfs := make([]string, 0, v.(*schema.Set).Len())
		for _, x := range v.(*schema.Set).List() {
			nfs = append(nfs, x.(string))
		}
		volume.ExportPolicyInfo.NfsVersion = nfs
	}
	if d.HasChange("permission") || d.HasChange("users") {
		volume.ShareInfoUpdate.ShareName = d.Get("share_name").(string)
		volume.ShareInfoUpdate.AccessControlList = make([]accessControlList, 1)
		volume.ShareInfoUpdate.AccessControlList[0].Permission = d.Get("permission").(string)
		users := make([]string, 0, d.Get("users").(*schema.Set).Len())
		for _, x := range d.Get("users").(*schema.Set).List() {
			users = append(users, x.(string))
		}
		volume.ShareInfoUpdate.AccessControlList[0].Users = users
	}
	if d.HasChange("snapshot_policy_name") {
		volume.SnapshotPolicyName = d.Get("snapshot_policy_name").(string)
	}
	if d.HasChange("tiering_policy") {
		volume.TieringPolicy = d.Get("tiering_policy").(string)
	}

	err = client.updateVolume(volume, clientID)
	if err != nil {
		log.Print("Error updating volume")
		return err
	}

	return resourceCVOVolumeRead(d, meta)
}

func resourceVolumeCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	// Check supported modification: Use volume name as an indication to know if this is a creation or modification
	if !(diff.HasChange("name")) {
		changeableParams := []string{"volume_protocol", "export_policy_type", "export_policy_ip", "export_policy_name", "export_policy_nfs_version",
			"share_name", "permission", "users", "tiering_policy", "snapshot_policy_name", "export_policy_rule_access_control", "export_policy_rule_super_user"}
		changedKeys := diff.GetChangedKeysPrefix("")
		for _, key := range changedKeys {
			found := false
			for _, changeable := range changeableParams {
				if strings.Contains(key, changeable) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("change %s is not allowed", key)
			}
		}
	}

	if diff.HasChange("volume_protocol") {
		currentVolumeType, expectVolumeType := diff.GetChange("volume_protocol")
		if currentVolumeType.(string) == "" {
			if expectVolumeType.(string) == "nfs" {
				if _, ok := diff.GetOk("export_policy_type"); !ok {
					return fmt.Errorf("export_policy_type is required when volume type is nfs")
				}
				if _, ok := diff.GetOk("export_policy_ip"); !ok {
					return fmt.Errorf("export_policy_ip is required when volume type is nfs")
				}
				if _, ok := diff.GetOk("export_policy_nfs_version"); !ok {
					return fmt.Errorf("export_policy_nfs_version is required when volume type is nfs")
				}
			} else if expectVolumeType.(string) == "cifs" {
				if _, ok := diff.GetOk("share_name"); !ok {
					return fmt.Errorf("share_name is required when volume type is cifs")
				}
				if _, ok := diff.GetOk("permission"); !ok {
					return fmt.Errorf("permission is required when volume type is cifs")
				}
				if _, ok := diff.GetOk("users"); !ok {
					return fmt.Errorf("users is required when volume type is cifs")
				}
			} else if expectVolumeType.(string) == "iscsi" {
				if _, ok := diff.GetOk("igroups"); !ok {
					return fmt.Errorf("igroups is required when volume type is iscsi")
				}
				if _, ok := diff.GetOk("os_name"); !ok {
					return fmt.Errorf("os_name is required when volume type is iscsi")
				}
			}
		} else {
			return fmt.Errorf("volume type can not be changed")
		}
	}
	providerVolumeType := diff.Get("provider_volume_type")
	if _, ok := diff.GetOk("iops"); !ok && (providerVolumeType == "io1" || providerVolumeType == "gp3") {
		return fmt.Errorf("iops is required when provider_volume_type is io1 or gp3")
	}
	if _, ok := diff.GetOk("throughput"); !ok && providerVolumeType == "gp3" {
		return fmt.Errorf("throughput is required when provider_volume_type is gp3")
	}
	capacityTier := diff.Get("capacity_tier")
	if _, ok := diff.GetOk("tiering_policy"); !ok && capacityTier == "S3" {
		return fmt.Errorf("tiering policy is required when capacity tier is S3")
	}
	return nil
}

func createIscsiVolumeHelper(d *schema.ResourceData, meta interface{}) (bool, bool, error) {
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	igroup := igroup{}

	var workingEnvironmentType string
	var workingEnvironmentID string
	var isNewIgroup bool
	var isNewInitiator bool
	var svm string
	if v, ok := d.GetOk("igroup_name"); ok {
		igroup.IgroupName = v.(string)
	}

	workingEnvDetail, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		return false, false, fmt.Errorf("cannot find working environment")
	}
	igroup.WorkingEnvironmentID = workingEnvDetail.PublicID
	workingEnvironmentID = workingEnvDetail.PublicID
	workingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
	if svm == "" {
		if workingEnvDetail.SvmName != "" {
			svm = workingEnvDetail.SvmName
		} else {
			svm = "svm_" + workingEnvDetail.Name
		}
	}
	igroup.SvmName = svm

	igroup.WorkingEnvironmentType = workingEnvironmentType
	res, err := client.getIgroups(igroup, clientID)
	if err != nil {
		log.Print("Error reading igroups")
		return false, false, err
	}
	for _, expectIg := range d.Get("igroups").(*schema.Set).List() {
		findIgroup := false
		for _, currentIg := range res {
			if currentIg.IgroupName == expectIg.(string) && isNewIgroup {
				return false, false, fmt.Errorf("igroups can not contain both existed and new igroup")
			}
			if currentIg.IgroupName == expectIg.(string) {
				findIgroup = true
				break
			}
		}
		if !findIgroup {
			isNewIgroup = true
		}
	}
	if isNewIgroup && workingEnvironmentType != "ON_PREM" {
		var initiators []initiator
		if v, ok := d.GetOk("initiator"); ok {
			initiators = expandInitiator(v.(*schema.Set))
		}
		getAll := initiator{}
		getAll.WorkingEnvironmentID = workingEnvironmentID
		getAll.WorkingEnvironmentType = workingEnvironmentType
		res, err := client.getInitiator(getAll, clientID)
		if err != nil {
			return false, false, err
		}
		// Check initiators does not contain both existed and new initiator.
		for _, expectIni := range initiators {
			findInitiator := false
			for _, currentIni := range res {
				if expectIni.Iqn == currentIni.Iqn && isNewInitiator {
					return false, false, fmt.Errorf("initiators can not contain both existed and new initiator")
				}
				if expectIni.Iqn == currentIni.Iqn {
					findInitiator = true
					break
				}
			}
			if !findInitiator {
				isNewInitiator = true
			}
		}
		if isNewInitiator {
			for _, expectIni := range initiators {
				client.createInitiator(expectIni, clientID)
			}
		}
	}
	return isNewIgroup, isNewInitiator, nil
}

func expandInitiator(set *schema.Set) []initiator {
	var initiators []initiator
	for _, v := range set.List() {
		initiator := initiator{}
		ini := v.(map[string]interface{})
		initiator.AliasName = ini["alias"].(string)
		initiator.Iqn = ini["iqn"].(string)
		initiators = append(initiators, initiator)
	}
	return initiators
}
