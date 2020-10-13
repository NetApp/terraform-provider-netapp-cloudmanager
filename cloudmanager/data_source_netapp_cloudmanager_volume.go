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
		},
	}
}

func dataSourceCVOVolumeRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching volume: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	volume := volumeRequest{}
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
		volume.SvmName = weInfo.SvmName
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		weInfo, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		volume.WorkingEnvironmentID = weInfo.PublicID
		volume.SvmName = weInfo.SvmName
		volume.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	res, err := client.getVolume(volume)
	if err != nil {
		log.Print("Error reading volume")
		return err
	}
	for _, volume := range res {
		d.SetId(volume.ID)
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
		if _, ok := d.GetOk("capacity_tier"); ok {
			d.Set("capacity_tier", volume.CapacityTier)
		}
		if _, ok := d.GetOk("tiering_policy"); ok {
			d.Set("tiering_policy", volume.TieringPolicy)
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

	return fmt.Errorf("Error reading volume: volume doesn't exist")
}
