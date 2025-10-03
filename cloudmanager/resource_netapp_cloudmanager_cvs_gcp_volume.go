package cloudmanager

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCVSGCPVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVSGCPVolumeCreate,
		Read:   resourceCVSGCPVolumeRead,
		Delete: resourceCVSGCPVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_types": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"size_unit": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"gb"}, true),
			},
			"service_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "medium",
				ValidateFunc: validation.StringInSlice([]string{"low", "medium", "high"}, true),
			},
			"volume_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"snapshot_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"daily_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
								},
							},
						},
						"hourly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
								},
							},
						},
						"monthly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days_of_month": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Default:  "1",
									},
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
								},
							},
						},
						"weekly_schedule": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Default:  "Sunday",
									},
									"hour": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
									"minute": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
									"snapshots_to_keep": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  0,
									},
								},
							},
						},
					},
				},
			},
			"export_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_clients": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"rule_index": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  true,
									},
									"unix_read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},
									"unix_read_write": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},
									"nfsv3": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},
									"nfsv4": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
										Default:  false,
									},
								},
							},
						},
					},
				},
			},
			"working_environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCVSGCPVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	volume := gcpVolumeRequest{}

	volume.Name = d.Get("name").(string)
	volume.Region = d.Get("region").(string)
	volume.Network = d.Get("network").(string)
	clientID := d.Get("client_id").(string)
	volume.WorkingEnvironmentName = d.Get("working_environment_name").(string)
	protocols := d.Get("protocol_types")
	for _, protocol := range protocols.([]interface{}) {
		volume.ProtocolTypes = append(volume.ProtocolTypes, protocol.(string))
	}

	// // size in 1 GiB increments, api takes in bytes only
	// volume.Size = d.Get("size").(int) * GiBToBytes
	volume.Size = math.Round(convertSizeUnit(d.Get("size").(float64), d.Get("size_unit").(string), "B")*10) / 10
	if v, ok := d.GetOk("service_level"); ok {
		volume.ServiceLevel = v.(string)
	}

	if v, ok := d.GetOk("snapshot_policy"); ok {
		if len(v.([]interface{})) > 0 {
			policy := v.([]interface{})[0].(map[string]interface{})
			volume.SnapshotPolicy = expandSnapshotPolicy(policy)
		}
	}

	if v, ok := d.GetOk("export_policy"); ok {
		policy := v.(*schema.Set)
		if policy.Len() > 0 {
			volume.ExportPolicy = expandExportPolicy(policy)
		}
	}

	if v, ok := d.GetOk("volume_path"); ok {
		volume.VolumePath = v.(string)
	}

	info := cvsInfo{}
	info.AccountName = d.Get("account").(string)

	res, err := client.createGCPVolume(volume, info, clientID)
	if err != nil {
		return err
	}
	d.SetId(res.VolumeID)

	return resourceCVSGCPVolumeRead(d, meta)
}

func resourceCVSGCPVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := gcpVolumeRequest{}
	info := cvsInfo{}
	volume.Name = d.Get("name").(string)
	volume.WorkingEnvironmentName = d.Get("working_environment_name").(string)
	volume.Region = d.Get("region").(string)
	volume.VolumeID = d.Id()
	info.AccountName = d.Get("account").(string)

	res, err := client.getGCPVolume(volume, info, clientID)
	if err != nil {
		log.Print("Error reading volume")
		return err
	}

	if err := d.Set("service_level", res.ServiceLevel); err != nil {
		return fmt.Errorf("error reading volume service_level: %s", err)
	}
	if err := d.Set("protocol_types", res.ProtocolTypes); err != nil {
		return fmt.Errorf("error reading volume protocol_types: %s", err)
	}
	if err := d.Set("volume_path", res.CreationToken); err != nil {
		return fmt.Errorf("error reading volume path or Creation Token: %s", err)
	}
	network := res.Network
	index := strings.Index(network, "networks/")
	if index > -1 {
		network = network[index+len("networks/"):]
	}
	if err := d.Set("network", network); err != nil {
		return fmt.Errorf("error reading volume network: %s", err)
	}
	if err := d.Set("region", res.Region); err != nil {
		return fmt.Errorf("error reading volume region: %s", err)
	}
	snapshotPolicy := flattenSnapshotPolicy(res.SnapshotPolicy)
	exportPolicy := flattenExportPolicy(res.ExportPolicy)
	if err := d.Set("snapshot_policy", snapshotPolicy); err != nil {
		return fmt.Errorf("error reading volume snapshot_policy: %s", err)
	}
	if len(res.ExportPolicy.Rules) > 0 {
		if err := d.Set("export_policy", exportPolicy); err != nil {
			return fmt.Errorf("error reading volume export_policy: %s", err)
		}
	} else {
		a := schema.NewSet(schema.HashString, []interface{}{})
		if err := d.Set("export_policy", a); err != nil {
			return fmt.Errorf("error reading volume export_policy: %s", err)
		}
	}
	return nil
}

func resourceCVSGCPVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := gcpVolumeRequest{}
	info := cvsInfo{}
	volume.WorkingEnvironmentName = d.Get("working_environment_name").(string)
	volume.Region = d.Get("region").(string)
	info.AccountName = d.Get("account").(string)
	volume.VolumeID = d.Id()

	err := client.deleteGCPVolume(volume, info, clientID)
	if err != nil {
		log.Print("Error deleting volume")
		return err
	}
	d.SetId("")
	return nil
}
