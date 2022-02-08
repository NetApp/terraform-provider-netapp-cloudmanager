package cloudmanager

import (
	"log"
	"math"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCVSANFVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVSAzureVolumeCreate,
		Read:   resourceCVSAzureVolumeRead,
		Delete: resourceCVSAzureVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceVolumeCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"working_environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_groups": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"netapp_account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"capacity_pool": {
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
				ValidateFunc: validation.StringInSlice([]string{"gb"}, false),
			},
			"volume_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_types": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_network": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"export_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_clients": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"cifs": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"nfsv3": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"nfsv41": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"unix_read_only": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"unix_read_write": {
										Type:     schema.TypeBool,
										Optional: true,
										ForceNew: true,
									},
									"rule_index": {
										Type:     schema.TypeInt,
										Optional: true,
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

func resourceCVSAzureVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating volume: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := anfVolumeRequest{}
	volume.Name = d.Get("name").(string)
	volume.Location = d.Get("location").(string)
	volume.ServiceLevel = d.Get("service_level").(string)
	volume.Size = math.Round(convertSizeUnit(d.Get("size").(float64), d.Get("size_unit").(string), "B")*10) / 10
	volume.SubnetName = d.Get("subnet").(string)
	volume.VolumePath = d.Get("volume_path").(string)
	volume.VirtualNetworkName = d.Get("virtual_network").(string)
	volume.WorkingEnvironmentName = d.Get("working_environment_name").(string)
	if v, ok := d.GetOk("protocol_types"); ok {
		protocolTypes := make([]string, 0, v.(*schema.Set).Len())
		for _, x := range v.(*schema.Set).List() {
			protocolTypes = append(protocolTypes, x.(string))
		}
		volume.ProtocolTypes = protocolTypes
	}
	if v, ok := d.GetOk("export_policy"); ok {
		if len(v.([]interface{})) > 0 {
			rules := make([]rule, 0, len(v.([]interface{})))
			for _, v1 := range v.([]interface{}) {
				v2 := v1.(map[string]interface{})
				ruleList := v2["rule"].([]interface{})
				for _, v3 := range ruleList {
					rule := rule{}
					ruleMap := v3.(map[string]interface{})
					rule.AllowedClients = ruleMap["allowed_clients"].(string)
					rule.UnixReadOnly = ruleMap["unix_read_only"].(bool)
					rule.RuleIndex = ruleMap["rule_index"].(int)
					rules = append(rules, rule)
				}
			}
			volume.Rules = rules
		}
	}
	info := cvsInfo{}
	info.AccountName = d.Get("account").(string)
	info.SubscriptionName = d.Get("subscription").(string)
	info.VirtualNetworkName = d.Get("virtual_network").(string)
	info.SubnetName = d.Get("subnet").(string)
	info.ResourceGroupsName = d.Get("resource_groups").(string)
	info.NetAppAccountName = d.Get("netapp_account").(string)
	info.CapacityPools = d.Get("capacity_pool").(string)
	err := client.createANFVolume(volume, info, clientID)
	if err != nil {
		return err
	}
	// volume Id is not returned, so use name.
	d.SetId(volume.Name)

	return resourceCVSAzureVolumeRead(d, meta)
}

func resourceCVSAzureVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	volume := anfVolumeRequest{}
	info := cvsInfo{}
	volume.Name = d.Get("name").(string)
	volume.WorkingEnvironmentName = d.Get("working_environment_name").(string)
	volume.Location = d.Get("location").(string)
	info.AccountName = d.Get("account").(string)
	info.SubscriptionName = d.Get("subscription").(string)
	info.VirtualNetworkName = d.Get("virtual_network").(string)
	info.SubnetName = d.Get("subnet").(string)
	info.ResourceGroupsName = d.Get("resource_groups").(string)
	info.NetAppAccountName = d.Get("netapp_account").(string)
	info.CapacityPools = d.Get("capacity_pool").(string)
	result, err := client.getANFVolume(volume, info, clientID)
	if err != nil {
		log.Print("Error reading volume")
		return err
	}
	d.Set("size", math.Round(convertSizeUnit(result.Size, "B", d.Get("size_unit").(string))*10)/10)
	d.Set("volume_path", result.VolumePath)
	d.Set("protocol_types", result.ProtocolTypes)
	d.Set("service_level", result.ServiceLevel)
	d.Set("location", result.Location)
	// subnet is returned as empty string in get volume API.
	//d.Set("subnet", result.SubnetName)

	rules := make([]map[string]interface{}, 1)
	rule := make(map[string]interface{})
	ruleList := make([]map[string]interface{}, len(result.Rules["rules"]))
	for _, ruleContent := range result.Rules["rules"] {
		ruleDict := make(map[string]interface{})
		ruleDict["allowed_clients"] = ruleContent.AllowedClients
		ruleDict["cifs"] = ruleContent.Cifs
		ruleDict["nfsv3"] = ruleContent.Nfsv3
		ruleDict["nfsv41"] = ruleContent.Nfsv41
		ruleDict["ruleIndex"] = ruleContent.RuleIndex
		ruleDict["unixReadOnly"] = ruleContent.UnixReadOnly
		ruleDict["unixReadWrite"] = ruleContent.UnixReadWrite
		ruleList = append(ruleList, ruleDict)
	}
	rule["rule"] = ruleList
	rules[0] = rule
	d.Set("export_policy", rules)

	return nil
}

func resourceCVSAzureVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	info := cvsInfo{}
	volume := anfVolumeRequest{}
	volume.Name = d.Get("name").(string)
	volume.WorkingEnvironmentName = d.Get("working_environment_name").(string)
	info.AccountName = d.Get("account").(string)
	info.SubscriptionName = d.Get("subscription").(string)
	info.VirtualNetworkName = d.Get("virtual_network").(string)
	info.SubnetName = d.Get("subnet").(string)
	info.ResourceGroupsName = d.Get("resource_groups").(string)
	info.NetAppAccountName = d.Get("netapp_account").(string)
	info.CapacityPools = d.Get("capacity_pool").(string)
	err := client.deleteANFVolume(volume, info, clientID)
	if err != nil {
		return err
	}

	return nil
}
