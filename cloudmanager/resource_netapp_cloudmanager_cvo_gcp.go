package cloudmanager

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCVOGCP() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOGCPCreate,
		Read:   resourceCVOGCPRead,
		Delete: resourceCVOGCPDelete,
		Update: resourceCVOGCPUpdate,
		Exists: resourceCVOGCPExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceCVOGCPCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gcp_service_account": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_encryption_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "GCP",
				ValidateFunc: validation.StringInSlice([]string{"GCP", "NONE"}, false),
			},
			"gcp_encryption_parameters": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gcp_volume_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "pd-ssd",
				ValidateFunc: validation.StringInSlice([]string{"pd-balanced", "pd-standard", "pd-ssd"}, false),
			},
			"gcp_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"gcp_volume_size_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "TB",
				ValidateFunc: validation.StringInSlice([]string{"GB", "TB"}, false),
			},
			"ontap_version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "latest",
			},
			"use_latest_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "gcp-cot-standard-paygo",
				ValidateFunc: validation.StringInSlice(GCPLicenseTypes, false),
			},
			"capacity_package_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Essential", "Professional", "Freemium"}, false),
			},
			"provided_license": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "n1-standard-8",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "default",
			},
			"network_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "default",
			},
			"svm_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"tier_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "standard",
				ValidateFunc: validation.StringInSlice([]string{"standard", "nearline", "coldline"}, false),
			},
			"nss_account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"writing_speed_state": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "cloudStorage",
				ValidateFunc: validation.StringInSlice([]string{"cloudStorage", "NONE"}, false),
			},
			"gcp_label": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"label_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				// ValidateFunc: func(val interface{}, )
			},
			"firewall_rule": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_ha": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"platform_serial_number_node1": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"platform_serial_number_node2": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node1_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node2_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mediator_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc0_node_and_data_connectivity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc1_cluster_connectivity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc2_ha_connectivity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc3_data_replication": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet0_node_and_data_connectivity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet1_cluster_connectivity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet2_ha_connectivity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet3_data_replication": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc0_firewall_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc1_firewall_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc2_firewall_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc3_firewall_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"svm_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"upgrade_ontap_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCVOGCPCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating CVO GCP: %#v", d)

	client := meta.(*Client)

	cvoDetails := createCVOGCPDetails{}

	cvoDetails.Name = d.Get("name").(string)
	cvoDetails.Region = d.Get("zone").(string)
	cvoDetails.GCPServiceAccount = d.Get("gcp_service_account").(string)
	cvoDetails.DataEncryptionType = d.Get("data_encryption_type").(string)
	cvoDetails.WorkspaceID = d.Get("workspace_id").(string)
	cvoDetails.GCPVolumeType = d.Get("gcp_volume_type").(string)
	cvoDetails.SvmPassword = d.Get("svm_password").(string)
	capacityTier := d.Get("capacity_tier").(string)
	if capacityTier == "cloudStorage" {
		cvoDetails.CapacityTier = capacityTier
		cvoDetails.TierLevel = d.Get("tier_level").(string)
	}
	cvoDetails.GCPVolumeSize.Size = d.Get("gcp_volume_size").(int)
	cvoDetails.GCPVolumeSize.Unit = d.Get("gcp_volume_size_unit").(string)
	cvoDetails.VsaMetadata.OntapVersion = d.Get("ontap_version").(string)
	cvoDetails.VsaMetadata.UseLatestVersion = d.Get("use_latest_version").(bool)
	cvoDetails.VsaMetadata.LicenseType = d.Get("license_type").(string)
	cvoDetails.VpcID = d.Get("vpc_id").(string)
	cvoDetails.Project = d.Get("project_id").(string)
	cvoDetails.VsaMetadata.InstanceType = d.Get("instance_type").(string)
	subnetID := d.Get("subnet_id").(string)
	if c, ok := d.GetOk("gcp_label"); ok {
		labels := c.(*schema.Set)
		if labels.Len() > 0 {
			cvoDetails.GCPLabels = expandGCPLabels(labels)
		}
	}

	var networkProjectID string
	if c, ok := d.GetOk("network_project_id"); ok {
		networkProjectID = c.(string)
	} else {
		networkProjectID = cvoDetails.Project
	}

	hasSelfLink := strings.HasPrefix(subnetID, "https://www.googleapis.com/compute/") || strings.HasPrefix(subnetID, "projects/")
	if hasSelfLink != true {
		cvoDetails.SubnetID = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], subnetID)
	} else {
		cvoDetails.SubnetID = subnetID
	}
	cvoDetails.SubnetPath = cvoDetails.SubnetID

	if c, ok := d.GetOk("firewall_rule"); ok {
		cvoDetails.FirewallRule = c.(string)
	}

	if c, ok := d.GetOk("writing_speed_state"); ok {
		cvoDetails.WritingSpeedState = c.(string)
	}

	if c, ok := d.GetOk("nss_account"); ok {
		cvoDetails.NssAccount = c.(string)
	}

	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	if c, ok := d.GetOk("capacity_package_name"); ok {
		cvoDetails.VsaMetadata.CapacityPackageName = c.(string)
	}

	if c, ok := d.GetOk("provided_license"); ok {
		cvoDetails.VsaMetadata.ProvidedLicense = c.(string)
	}

	if c, ok := d.GetOk("serial_number"); ok {
		cvoDetails.SerialNumber = c.(string)
	}

	if c, ok := d.GetOk("gcp_encryption_parameters"); ok {
		cvoDetails.GcpEncryptionParameters.Key = c.(string)
	}

	cvoDetails.IsHA = d.Get("is_ha").(bool)

	if cvoDetails.IsHA == true {
		if c, ok := d.GetOk("platform_serial_number_node1"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode1 = c.(string)
		}
		if c, ok := d.GetOk("platform_serial_number_node2"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode2 = c.(string)
		}
		if c, ok := d.GetOk("node1_zone"); ok {
			cvoDetails.HAParams.Node1Zone = c.(string)
		}
		if c, ok := d.GetOk("node2_zone"); ok {
			cvoDetails.HAParams.Node2Zone = c.(string)
		}
		if c, ok := d.GetOk("mediator_zone"); ok {
			cvoDetails.HAParams.MediatorZone = c.(string)
		}
		if c, ok := d.GetOk("vpc0_node_and_data_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC0NodeAndDataConnectivity = c.(string)
		}
		if c, ok := d.GetOk("vpc1_cluster_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC1ClusterConnectivity = c.(string)
		}
		if c, ok := d.GetOk("vpc2_ha_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC2HAConnectivity = c.(string)
		}
		if c, ok := d.GetOk("vpc3_data_replication"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC3DataReplication = c.(string)
		}
		if c, ok := d.GetOk("subnet0_node_and_data_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], c.(string))
			}
			cvoDetails.HAParams.Subnet0NodeAndDataConnectivity = c.(string)
		}
		if c, ok := d.GetOk("subnet1_cluster_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], c.(string))
			}
			cvoDetails.HAParams.Subnet1ClusterConnectivity = c.(string)
		}
		if c, ok := d.GetOk("subnet2_ha_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], c.(string))
			}
			cvoDetails.HAParams.Subnet2HAConnectivity = c.(string)
		}
		if c, ok := d.GetOk("subnet3_data_replication"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if hasSelfLink != true {
				c = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], c.(string))
			}
			cvoDetails.HAParams.Subnet3DataReplication = c.(string)
		}
		if c, ok := d.GetOk("vpc0_firewall_rule_name"); ok {
			cvoDetails.HAParams.VPC0FirewallRuleName = c.(string)
		}
		if c, ok := d.GetOk("vpc1_firewall_rule_name"); ok {
			cvoDetails.HAParams.VPC1FirewallRuleName = c.(string)
		}
		if c, ok := d.GetOk("vpc2_firewall_rule_name"); ok {
			cvoDetails.HAParams.VPC2FirewallRuleName = c.(string)
		}
		if c, ok := d.GetOk("vpc3_firewall_rule_name"); ok {
			cvoDetails.HAParams.VPC3FirewallRuleName = c.(string)
		}
	} else if cvoDetails.WritingSpeedState == "" {
		cvoDetails.WritingSpeedState = "NORMAL"
		d.Set("writing_speed_state", "NORMAL")
	}

	err := validateCVOGCPParams(cvoDetails)
	if err != nil {
		log.Print("Error validating parameters")
		return err
	}

	res, err := client.createCVOGCP(cvoDetails)
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(res.PublicID)
	d.Set("svm_name", res.SvmName)

	log.Printf("Created cvo: %v", res)

	return resourceCVOGCPRead(d, meta)
}

func resourceCVOGCPRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()

	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	_, err := client.getWorkingEnvironmentInfo(id)
	if err != nil {
		log.Print("Error getting cvo")
		return err
	}

	return nil
}

func resourceCVOGCPDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting CVO: %#v", d)

	client := meta.(*Client)

	id := d.Id()
	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	isHA := d.Get("is_ha").(bool)

	deleteErr := client.deleteCVOGCP(id, isHA)
	if deleteErr != nil {
		log.Print("Error deleting cvo")
		return deleteErr
	}

	return nil
}

func resourceCVOGCPUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating CVO: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)

	// check if svm_password is changed
	if d.HasChange("svm_password") {
		respErr := updateCVOSVMPassword(d, meta)
		if respErr != nil {
			return respErr
		}
	}

	// check if license_type and instance type are changed
	if d.HasChange("instance_type") || d.HasChange("license_type") {
		respErr := updateCVOLicenseInstanceType(d, meta)
		if respErr != nil {
			return respErr
		}
	}

	// check if tier_level is changed
	if d.HasChange("tier_level") && d.Get("capacity_tier").(string) == "cloudStorage" {
		respErr := updateCVOTierLevel(d, meta)
		if respErr != nil {
			return respErr
		}
	}

	// check if gcp_label has changes
	if d.HasChange("gcp_label") {
		respErr := updateCVOUserTags(d, meta, "gcp_label")
		if respErr != nil {
			return respErr
		}
		return resourceCVOGCPRead(d, meta)
	}
	// upgrade ontap version
	upgradeErr := client.checkAndDoUpgradeOntapVersion(d)
	if upgradeErr != nil {
		return upgradeErr
	}

	return nil
}

func checkIfLabelMissing(cLabels *schema.Set, eLabels *schema.Set) error {
	// check if current gcp_labels in future gcp_labels
	for _, currentLabel := range cLabels.List() {
		found := false
		clabel := currentLabel.(map[string]interface{})
		ckey := clabel["label_key"].(string)
		for _, expectLabel := range eLabels.List() {
			elabel := expectLabel.(map[string]interface{})
			ekey := elabel["label_key"].(string)
			if ekey == ckey {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("label key %s in gcp_label cannot be removed", ckey)
		}
	}
	return nil
}

func checkLabelDiff(diff *schema.ResourceDiff) error {
	// gcp_label only can be added and modified on the label values
	if diff.HasChange("gcp_label") {
		currentLabel, expectLabel := diff.GetChange("gcp_label")

		if currentLabel != nil {
			if expectLabel == nil {
				return fmt.Errorf("gcp_label deletion is not supported")
			}
			cLabels := currentLabel.(*schema.Set)
			eLabels := expectLabel.(*schema.Set)
			if cLabels.Len() > eLabels.Len() {
				return fmt.Errorf("gcp_label deletion is not supported")
			}

			respErr := checkUserTagKeyUnique(eLabels, "label_key")
			if respErr != nil {
				return respErr
			}
			return checkIfLabelMissing(cLabels, eLabels)
		}
	}
	return nil
}

func resourceCVOGCPCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	respErr := checkLabelDiff(diff)
	if respErr != nil {
		return respErr
	}
	return nil
}

func resourceCVOGCPExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	client.ClientID = d.Get("client_id").(string)
	name := d.Get("name").(string)

	resID, err := client.findWorkingEnvironmentByName(name)
	if err != nil {
		log.Print("Error getting cvo")
		return false, err
	}

	if resID.PublicID != id {
		d.SetId("")
		return false, err
	}

	return true, nil
}
