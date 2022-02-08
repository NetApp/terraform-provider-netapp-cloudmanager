package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceFsxVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceFSXVolumeCreate,
		Read:   resourceFSXVolumeRead,
		Delete: resourceFSXVolumeDelete,
		Update: resourceFSXVolumeUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceFSXVolumeCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"svm_name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"tiering_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"none", "snapshot_only", "auto", "all"}, false),
			},
			"volume_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"nfs", "cifs"}, false),
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
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_storage_efficiency": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceFSXVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	var svm string
	volume := volumeRequest{}

	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		log.Printf("Cannot find working environment: %#v", err)
		return fmt.Errorf("Cannot find working environment: %#v", err)
	}
	if svm == "" {
		svm = weInfo.SvmName
	}
	volume.SvmName = svm

	volume.FileSystemID = d.Get("file_system_id").(string)
	volume.TenantID = d.Get("tenant_id").(string)
	volume.EnableCompression = false
	volume.EnableThinProvisioning = true
	volume.EnableStorageEfficiency = d.Get("enable_storage_efficiency").(bool)
	err = client.setCommonAttributes(d, &volume, clientID)
	if err != nil {
		return err
	}

	err = client.createVolume(volume, false, clientID)
	if err != nil {
		log.Printf("Error creating volume: %#v", err)
		return err
	}
	res, err := client.getVolume(volume, clientID)
	if err != nil {
		log.Printf("Error reading volume after creation: %#v", err)
		return err
	}
	for _, volume := range res {
		if volume.SvmName == svm && volume.Name == d.Get("name") {
			d.SetId(volume.ID)
			break
		}
	}
	return resourceFSXVolumeRead(d, meta)
}

func resourceFSXVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching volume: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	} else {
		weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
		if err != nil {
			log.Printf("Cannot find working environment: %#v", err)
			return fmt.Errorf("Cannot find working environment: %#v", err)
		}
		if svm == "" {
			svm = weInfo.SvmName
		}
	}
	volume.SvmName = svm
	volume.Name = d.Get("name").(string)
	volume.FileSystemID = d.Get("file_system_id").(string)
	res, err := client.getVolume(volume, clientID)
	if err != nil {
		log.Printf("Error reading volume: %#v", err)
		return err
	}
	for _, volume := range res {
		if volume.ID == d.Id() {
			if _, ok := d.GetOk("enable_storage_efficiency"); ok {
				// EnableDeduplication reflects  enable_storage_efficiency for fsx volume
				d.Set("enable_storage_efficiency", volume.EnableDeduplication)
			}
			if v, ok := d.GetOk("tiering_policy"); ok {
				if v.(string) != "none" {
					d.Set("tiering_policy", volume.TieringPolicy)
				}
			}
			if _, ok := d.GetOk("snapshot_policy_name"); ok {
				d.Set("snapshot_policy_name", volume.SnapshotPolicyName)
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

func resourceFSXVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting volume: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := volumeRequest{}
	var svm string
	if v, ok := d.GetOk("svm_name"); ok {
		svm = v.(string)
	} else {
		weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
		if err != nil {
			log.Printf("Cannot find working environment: %#v", err)
			return fmt.Errorf("Cannot find working environment: %#v", err)
		}
		svm = weInfo.SvmName
	}
	volume.FileSystemID = d.Get("file_system_id").(string)
	volume.SvmName = svm

	volume.Name = d.Get("name").(string)

	err := client.deleteVolume(volume, clientID)
	if err != nil {
		log.Printf("Error deleting volume: %#v", err)
		return err
	}
	return nil
}

func resourceFSXVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating volume: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
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

	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID)
	if err != nil {
		log.Printf("Cannot find working environment: %#v", err)
		return fmt.Errorf("Cannot find working environment: %#v", err)
	}
	if svm == "" {
		svm = weInfo.SvmName
	}
	volume.SvmName = svm
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
	err = client.updateVolume(volume, clientID)
	if err != nil {
		log.Printf("Error updating volume: %#v", err)
		return err
	}

	return resourceFSXVolumeRead(d, meta)
}

func resourceFSXVolumeCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	if diff.HasChange("volume_protocol") {
		currentVolumeType, expectedVolumeType := diff.GetChange("volume_protocol")
		if currentVolumeType.(string) == "" {
			if expectedVolumeType.(string) == "nfs" {
				if _, ok := diff.GetOk("export_policy_type"); !ok {
					return fmt.Errorf("export_policy_type is required when volume type is nfs")
				}
				if _, ok := diff.GetOk("export_policy_ip"); !ok {
					return fmt.Errorf("export_policy_ip is required when volume type is nfs")
				}
				if _, ok := diff.GetOk("export_policy_nfs_version"); !ok {
					return fmt.Errorf("export_policy_nfs_version is required when volume type is nfs")
				}
			} else if expectedVolumeType.(string) == "cifs" {
				if _, ok := diff.GetOk("share_name"); !ok {
					return fmt.Errorf("share_name is required when volume type is cifs")
				}
				if _, ok := diff.GetOk("permission"); !ok {
					return fmt.Errorf("permission is required when volume type is cifs")
				}
				if _, ok := diff.GetOk("users"); !ok {
					return fmt.Errorf("users is required when volume type is cifs")
				}
			}
		} else {
			return fmt.Errorf("volume type %s can not be changed to %s", currentVolumeType.(string), expectedVolumeType.(string))
		}
	}
	return nil
}
