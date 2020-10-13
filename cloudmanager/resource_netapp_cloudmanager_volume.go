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
				Default:  true,
			},
			"enable_compression": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"enable_deduplication": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
				Type:     schema.TypeSet,
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
			"iops": {
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
		},
	}
}

func resourceCVOVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	var svm string
	var workingEnvironmentType string
	var cloudProvider string
	volume := volumeRequest{}
	quote := quoteRequest{}
	// quote volume
	quote.Name = d.Get("name").(string)
	quote.Size.Size = d.Get("size").(float64)
	quote.Size.Unit = d.Get("unit").(string)
	quote.SnapshotPolicyName = d.Get("snapshot_policy_name").(string)
	quote.ProviderVolumeType = d.Get("provider_volume_type").(string)
	quote.EnableDeduplication = d.Get("enable_deduplication").(bool)
	quote.EnableThinProvisioning = d.Get("enable_thin_provisioning").(bool)
	quote.VerifyNameUniqueness = true // hard code to always true
	if v, ok := d.GetOk("iops"); ok {
		quote.Iops = v.(int)
	}
	if v, ok := d.GetOk("aggregate_name"); ok {
		quote.AggregateName = v.(string)
	}
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	}
	if v, ok := d.GetOk("working_environment_id"); ok {
		volume.WorkingEnvironmentID = v.(string)
		quote.WorkingEnvironmentID = v.(string)
		weInfo, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return nil
		}
		workingEnvironmentType = weInfo.WorkingEnvironmentType
		weInfo, err = client.findWorkingEnvironmentByName(weInfo.Name)
		if err != nil {
			return err
		}
		if svm == "" {
			svm = weInfo.SvmName
		}
		quote.SvmName = svm
		volume.SvmName = svm
		cloudProvider = strings.ToLower(weInfo.CloudProviderName)
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		weInfo, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		if svm == "" {
			svm = weInfo.SvmName
		}
		volume.WorkingEnvironmentID = weInfo.PublicID
		quote.WorkingEnvironmentID = weInfo.PublicID
		quote.SvmName = svm
		volume.SvmName = svm
		workingEnvironmentType = weInfo.WorkingEnvironmentType
		cloudProvider = strings.ToLower(weInfo.CloudProviderName)
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	quote.WorkingEnvironmentType = workingEnvironmentType
	volume.WorkingEnvironmentType = workingEnvironmentType
	if v, ok := d.GetOk("capacity_tier"); ok {
		if v.(string) != "none" {
			quote.CapacityTier = v.(string)
			if v, ok = d.GetOk("tiering_policy"); ok {
				if v.(string) != "none" {
					quote.TieringPolicy = v.(string)
				}
			}
		}
	} else {
		if cloudProvider == "aws" {
			quote.CapacityTier = "S3"
		} else if cloudProvider == "azure" {
			quote.CapacityTier = "Blob"
		} else if cloudProvider == "gcp" {
			quote.CapacityTier = "cloudStorage"
		}
	}
	response, err := client.quoteVolume(quote)
	if err != nil {
		log.Printf("Error quoting volume")
		return err
	}
	volume.NewAggregate = response["newAggregate"].(bool)
	volume.AggregateName = response["aggregateName"].(string)
	volume.NumOfDisks = response["numOfDisks"].(float64)
	volume.ProviderVolumeType = d.Get("provider_volume_type").(string)
	volume.Name = d.Get("name").(string)
	volume.SnapshotPolicyName = d.Get("snapshot_policy_name").(string)
	volume.EnableThinProvisioning = d.Get("enable_thin_provisioning").(bool)
	volume.EnableCompression = d.Get("enable_compression").(bool)
	volume.EnableDeduplication = d.Get("enable_deduplication").(bool)
	volume.Size.Size = d.Get("size").(float64)
	volume.Size.Unit = d.Get("unit").(string)
	volume_protocol := d.Get("volume_protocol").(string)
	if v, ok := d.GetOk("export_policy_name"); ok {
		volume.ExportPolicyInfo.Name = v.(string)
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
	if v, ok := d.GetOk("capacity_tier"); ok {
		if v.(string) != "none" {
			volume.CapacityTier = v.(string)
		}
	} else {
		if cloudProvider == "aws" {
			volume.CapacityTier = "S3"
		} else if cloudProvider == "azure" {
			volume.CapacityTier = "Blob"
		} else if cloudProvider == "gcp" {
			volume.CapacityTier = "cloudStorage"
		}
	}
	if v, ok := d.GetOk("iops"); ok {
		volume.Iops = v.(int)
	}
	if volume_protocol == "cifs" {
		exist, err := client.checkCifsExists(volume.WorkingEnvironmentID, volume.SvmName)
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
	} else if volume_protocol == "iscsi" {
		isNewIgroup, _, err := createIscsiVolumeHelper(d, meta)
		if err != nil {
			return err
		}
		if v, ok := d.GetOk("os_name"); ok {
			volume.IscsiInfo.OsName = v.(string)
		}
		if isNewIgroup {
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
	err = client.createVolume(volume)
	if err != nil {
		log.Print("Error creating volume")
		return err
	}
	res, err := client.getVolume(volume)
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
	log.Printf("Fetching volume: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	}
	if v, ok := d.GetOk("working_environment_id"); ok {
		volume.WorkingEnvironmentID = v.(string)
		weInfo, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return nil
		}
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
		weInfo, err = client.findWorkingEnvironmentByName(weInfo.Name)
		if err != nil {
			return err
		}
		if svm == "" {
			svm = weInfo.SvmName
		}
		volume.SvmName = svm
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		weInfo, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		volume.WorkingEnvironmentID = weInfo.PublicID
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
		if svm == "" {
			svm = weInfo.SvmName
		}
		volume.SvmName = svm
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	res, err := client.getVolume(volume)
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
			if _, ok := d.GetOk("enable_deduplication"); ok {
				d.Set("enable_deduplication", volume.EnableDeduplication)
			}
			if _, ok := d.GetOk("enable_compression"); ok {
				d.Set("enable_compression", volume.EnableCompression)
			}
			if _, ok := d.GetOk("export_policy_ip"); ok {
				d.Set("export_policy_ip", volume.ExportPolicyInfo.Ips)
			}
			if _, ok := d.GetOk("export_policy_nfs_version"); ok {
				d.Set("export_policy_nfs_version", volume.ExportPolicyInfo.NfsVersion)
			}
			if _, ok := d.GetOk("export_policy_type"); ok {
				d.Set("export_policy_type", volume.ExportPolicyInfo.PolicyType)
			}
			if _, ok := d.GetOk("provider_volume_type"); ok {
				d.Set("provider_volume_type", volume.ProviderVolumeType)
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

	return fmt.Errorf("Error reading volume: volume doesn't exist")
}

func resourceCVOVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	}
	if v, ok := d.GetOk("working_environment_id"); ok {
		volume.WorkingEnvironmentID = v.(string)
		weInfo, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return err
		}
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
		weInfo, err = client.findWorkingEnvironmentByName(weInfo.Name)
		if err != nil {
			return err
		}
		if svm == "" {
			svm = weInfo.SvmName
		}
		volume.SvmName = svm
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		weInfo, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		volume.WorkingEnvironmentID = weInfo.PublicID
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
		if svm == "" {
			svm = weInfo.SvmName
		}
		volume.SvmName = svm
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}

	volume.Name = d.Get("name").(string)

	err := client.deleteVolume(volume)
	if err != nil {
		log.Print("Error deleting volume")
		return err
	}
	return nil
}

func resourceCVOVolumeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of volume: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	volume := volumeRequest{}
	volume.Name = d.Get("name").(string)
	volume.ID = d.Id()
	if v, ok := d.GetOk("working_environment_id"); ok {
		volume.WorkingEnvironmentID = v.(string)
		weInfo, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return false, err
		}
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		weInfo, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return false, err
		}
		volume.WorkingEnvironmentID = weInfo.PublicID
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	} else {
		return false, fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	res, err := client.getVolumeByID(volume)
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
	log.Printf("Updating volume: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	volume.Name = d.Get("name").(string)
	volume.ExportPolicyInfo.PolicyType = d.Get("export_policy_type").(string)
	if v, ok := d.GetOk("export_policy_ip"); ok {
		ips := make([]string, 0, v.(*schema.Set).Len())
		for _, x := range v.(*schema.Set).List() {
			ips = append(ips, x.(string))
		}
		volume.ExportPolicyInfo.Ips = ips
	}
	if v, ok := d.GetOk("working_environment_id"); ok {
		volume.WorkingEnvironmentID = v.(string)
		weInfo, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return err
		}
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
		weInfo, err = client.findWorkingEnvironmentByName(weInfo.Name)
		if err != nil {
			return err
		}
		if svm == "" {
			svm = weInfo.SvmName
		}
		volume.SvmName = svm
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		weInfo, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return err
		}
		volume.WorkingEnvironmentID = weInfo.PublicID
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
		if svm == "" {
			svm = weInfo.SvmName
		}
		volume.SvmName = svm
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
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
		volume.ShareInfoUpdate.AccessControlList[0].Permission = d.Get("permission").(string)
		users := make([]string, 0, d.Get("Users").(*schema.Set).Len())
		for _, x := range d.Get("Users").(*schema.Set).List() {
			users = append(users, x.(string))
		}
		volume.ShareInfoUpdate.AccessControlList[0].Users = users
	}
	err := client.updateVolume(volume)
	if err != nil {
		log.Print("Error updating volume")
		return err
	}

	return resourceCVOVolumeRead(d, meta)
}

func resourceVolumeCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
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
	provider_volume_type := diff.Get("provider_volume_type")
	if _, ok := diff.GetOk("iops"); !ok && provider_volume_type == "io1" {
		return fmt.Errorf("iops is required when provider_volume_type is io1")
	}
	capacityTier := diff.Get("capacity_tier")
	if _, ok := diff.GetOk("tiering_policy"); !ok && capacityTier == "S3" {
		return fmt.Errorf("tiering policy is required when capacity tier is S3")
	}
	return nil
}

func createIscsiVolumeHelper(d *schema.ResourceData, meta interface{}) (bool, bool, error) {
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	igroup := igroup{}

	var workingEnvironmentType string
	var workingEnvironmentId string
	var isNewIgroup bool
	var isNewInitiator bool
	var svm string
	if v, ok := d.GetOk("igroup_name"); ok {
		igroup.IgroupName = v.(string)
	}
	if v, ok := d.GetOk("working_environment_id"); ok {
		igroup.WorkingEnvironmentId = v.(string)
		workingEnvironmentId = v.(string)
		we_info, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return false, false, err
		}
		workingEnvironmentType = we_info.WorkingEnvironmentType
		we_info, err = client.findWorkingEnvironmentByName(we_info.Name)
		if err != nil {
			return false, false, err
		}
		if svm == "" {
			svm = we_info.SvmName
		}
		igroup.SvmName = svm
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		we_info, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return false, false, err
		}
		igroup.WorkingEnvironmentId = we_info.PublicID

		workingEnvironmentId = we_info.PublicID
		workingEnvironmentType = we_info.WorkingEnvironmentType
		if svm == "" {
			svm = we_info.SvmName
		}
		igroup.SvmName = svm
	} else {
		return false, false, fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	igroup.WorkingEnvironmentType = workingEnvironmentType
	res, err := client.getIgroups(igroup)
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
	if isNewIgroup {
		var initiators []initiator
		if v, ok := d.GetOk("initiator"); ok {
			initiators = expandInitiator(v.(*schema.Set))
		}
		getAll := initiator{}
		getAll.WorkingEnvironmentId = workingEnvironmentId
		getAll.WorkingEnvironmentType = workingEnvironmentType
		res, err := client.getInitiator(getAll)
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
				client.createInitiator(expectIni)
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
