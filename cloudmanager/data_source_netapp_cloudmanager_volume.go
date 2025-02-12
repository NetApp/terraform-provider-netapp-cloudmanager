package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceCVOVolume() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCVOVolumeRead,
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
				Optional: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Optional: true,
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
			"throughput": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"provider_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
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
			"mount_point": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceCVOVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching volume: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := volumeRequest{}
	weInfo, err := client.getWorkingEnvironmentDetail(d, clientID, true, "")
	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	volume.WorkingEnvironmentID = weInfo.PublicID
	volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	volume.SvmName = weInfo.SvmName
	d.Set("working_environment_id", volume.WorkingEnvironmentID)
	d.Set("svm_name", volume.SvmName)
	res, err := client.getVolume(volume, clientID, true, "")
	if err != nil {
		log.Print("Error reading volume")
		return err
	}
	for _, volume := range res {
		if volume.Name == d.Get("name").(string) {
			d.SetId(volume.ID)
			err = d.Set("aggregate_name", volume.AggregateName)
			if err != nil {
				return fmt.Errorf("error setting aggregate_name: %s", err.Error())
			}
			err = d.Set("snapshot_policy_name", volume.SnapshotPolicyName)
			if err != nil {
				return fmt.Errorf("error setting snapshot_policy_name: %s", err.Error())
			}
			err = d.Set("enable_thin_provisioning", volume.EnableThinProvisioning)
			if err != nil {
				return fmt.Errorf("error setting enable_thin_provisioning: %s", err.Error())
			}
			err = d.Set("enable_deduplication", volume.EnableDeduplication)
			if err != nil {
				return fmt.Errorf("error setting enable_deduplication: %s", err.Error())
			}
			err = d.Set("enable_compression", volume.EnableCompression)
			if err != nil {
				return fmt.Errorf("error setting enable_compression: %s", err.Error())
			}
			err = d.Set("provider_volume_type", volume.ProviderVolumeType)
			if err != nil {
				return fmt.Errorf("error setting provider_volume_type: %s", err.Error())
			}
			err = d.Set("capacity_tier", volume.CapacityTier)
			if err != nil {
				return fmt.Errorf("error setting capacity_tier: %s", err.Error())
			}
			err = d.Set("tiering_policy", volume.TieringPolicy)
			if err != nil {
				return fmt.Errorf("error setting tiering_policy: %s", err.Error())
			}
			err = d.Set("size", convertSizeUnit(volume.Size.Size, volume.Size.Unit, "GB"))
			if err != nil {
				return fmt.Errorf("error setting size: %s", err.Error())
			}
			err = d.Set("unit", volume.Size.Unit)
			if err != nil {
				return fmt.Errorf("error setting unit: %s", err.Error())
			}
			if len(volume.ShareInfo) > 0 {
				err = d.Set("share_name", volume.ShareInfo[0].ShareName)
				if err != nil {
					return fmt.Errorf("error setting share_name: %s", err.Error())
				}
				if len(volume.ShareInfo[0].AccessControlList) > 0 {
					err = d.Set("permission", volume.ShareInfo[0].AccessControlList[0].Permission)
					if err != nil {
						return fmt.Errorf("error setting permission: %s", err.Error())
					}
				}
				if len(volume.ShareInfo[0].AccessControlList) > 0 {
					err = d.Set("users", volume.ShareInfo[0].AccessControlList[0].Users)
					if err != nil {
						return fmt.Errorf("error setting users: %s", err.Error())
					}
				}
				err = d.Set("volume_protocol", "cifs")
				if err != nil {
					return fmt.Errorf("error setting volume_protocol: %s", err.Error())
				}
			} else if volume.IscsiEnabled {
				err = d.Set("volume_protocol", "iscsi")
				if err != nil {
					return fmt.Errorf("error setting volume_protocol: %s", err.Error())
				}
			} else {
				err = d.Set("volume_protocol", "nfs")
				if err != nil {
					return fmt.Errorf("error setting volume_protocol: %s", err.Error())
				}
				err = d.Set("export_policy_name", volume.ExportPolicyInfo.Name)
				if err != nil {
					return fmt.Errorf("error setting export_policy_name: %s", err.Error())
				}
				err = d.Set("export_policy_ip", volume.ExportPolicyInfo.Ips)
				if err != nil {
					return fmt.Errorf("error setting export_policy_ip: %s", err.Error())
				}
				err = d.Set("export_policy_nfs_version", volume.ExportPolicyInfo.NfsVersion)
				if err != nil {
					return fmt.Errorf("error setting export_policy_nfs_version: %s", err.Error())
				}
				err = d.Set("export_policy_type", volume.ExportPolicyInfo.PolicyType)
				if err != nil {
					return fmt.Errorf("error setting export_policy_type: %s", err.Error())
				}
				err = d.Set("mount_point", volume.MountPoint)
				if err != nil {
					return fmt.Errorf("error setting mount_point: %s", err.Error())
				}
			}
			return nil
		}
	}
	return fmt.Errorf("error reading volume: volume doesn't exist")
}
