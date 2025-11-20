package cloudmanager

import (
	"fmt"
	"log"
	"regexp"
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
				Optional: true,
				ForceNew: true,
			},
			"gcp_service_account": {
				Type:     schema.TypeString,
				Optional: true,
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
				Default:      "capacity-paygo",
				ValidateFunc: validation.StringInSlice(GCPLicenseTypes, false),
			},
			"capacity_package_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Essential", "Professional", "Freemium", "Edge", "Optimized"}, false),
			},
			"provided_license": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Log the values to understand what we're getting
					log.Printf("DiffSuppressFunc - key: %s, old: %s, new: %s", k, old, new)
					suppress := old == "n1-standard-8" && new == ""
					log.Printf("DiffSuppressFunc - suppressing: %t", suppress)
					return suppress
				},
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
			"svm": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"svm_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tier_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "nearline", "coldline"}, false),
			},
			"saas_subscription_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"worm_retention_period_length": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"worm_retention_period_unit": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"years", "months", "days", "hours", "minutes", "seconds"}, true),
				Optional:     true,
				ForceNew:     true,
			},
			"nss_account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"writing_speed_state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL", "HIGH"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloudStorage"}, false),
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
			"firewall_tag_name_rule": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewall_ip_ranges": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_volumes_to_cbs": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_compliance": {
				Type:     schema.TypeBool,
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
			"vpc0_firewall_rule_tag_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc1_firewall_rule_tag_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc2_firewall_rule_tag_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc3_firewall_rule_tag_name": {
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
			},
			"flash_cache": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"upgrade_ontap_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"retries": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
			"connector_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Restricted"}, false),
				Default:      "Standard",
			},
		},
	}
}

func resourceCVOGCPCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating CVO GCP: %#v", d)

	client := meta.(*Client)
	cvoDetails := createCVOGCPDetails{}

	clientID := d.Get("client_id").(string)

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	client.Retries = d.Get("retries").(int)

	cvoDetails.Name = d.Get("name").(string)
	log.Print("Create cvo name ", cvoDetails.Name)
	if c, ok := d.GetOk("gcp_service_account"); ok {
		cvoDetails.GCPServiceAccount = c.(string)
	}
	cvoDetails.DataEncryptionType = d.Get("data_encryption_type").(string)
	cvoDetails.WorkspaceID = d.Get("workspace_id").(string)
	cvoDetails.GCPVolumeType = d.Get("gcp_volume_type").(string)
	cvoDetails.SvmPassword = d.Get("svm_password").(string)
	if c, ok := d.GetOk("svm_name"); ok {
		cvoDetails.SvmName = c.(string)
	}
	if c, ok := d.GetOk("capacity_tier"); ok {
		cvoDetails.CapacityTier = c.(string)
	}
	if c, ok := d.GetOk("tier_level"); ok {
		cvoDetails.TierLevel = c.(string)
	}
	if c, ok := d.GetOk("saas_subscription_id"); ok {
		cvoDetails.SaasSubscriptionID = c.(string)
	}
	cvoDetails.GCPVolumeSize.Size = d.Get("gcp_volume_size").(int)
	cvoDetails.GCPVolumeSize.Unit = d.Get("gcp_volume_size_unit").(string)
	cvoDetails.VsaMetadata.OntapVersion = d.Get("ontap_version").(string)
	cvoDetails.VsaMetadata.UseLatestVersion = d.Get("use_latest_version").(bool)
	cvoDetails.VsaMetadata.LicenseType = d.Get("license_type").(string)

	if c, ok := d.GetOk("capacity_package_name"); ok {
		cvoDetails.VsaMetadata.CapacityPackageName = c.(string)
	} else {
		// by Capacity - set default capacity package name
		if strings.HasSuffix(cvoDetails.VsaMetadata.LicenseType, "capacity-paygo") {
			cvoDetails.VsaMetadata.CapacityPackageName = "Essential"
		}
	}

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

	if c, ok := d.GetOk("backup_volumes_to_cbs"); ok {
		cvoDetails.BackupVolumesToCbs = c.(bool)
	}

	if c, ok := d.GetOk("enable_compliance"); ok {
		cvoDetails.EnableCompliance = c.(bool)
	}
	// In both single and HA case, flash_cache only can be set when the selected instance_type
	if c, ok := d.GetOk("flash_cache"); ok {
		cvoDetails.FlashCache = c.(bool)
		match, _ := regexp.MatchString("^n2-standard-(16|32|48|64)$", cvoDetails.VsaMetadata.InstanceType)
		if !match {
			return fmt.Errorf("instance_type has to be one of n2-standard-16,32,48,64")
		}
	}

	if c, ok := d.GetOk("zone"); ok {
		cvoDetails.Region = c.(string)
	}
	cvoDetails.IsHA = d.Get("is_ha").(bool)
	if !cvoDetails.IsHA {
		if cvoDetails.Region == "" {
			return fmt.Errorf("zone is required in single GCP CVO")
		}
	} else {
		if c, ok := d.GetOk("node1_zone"); ok {
			cvoDetails.HAParams.Node1Zone = c.(string)
			if cvoDetails.Region == "" {
				cvoDetails.Region = cvoDetails.HAParams.Node1Zone
			}
		}
	}
	var networkProjectID string
	if c, ok := d.GetOk("network_project_id"); ok {
		networkProjectID = c.(string)
	} else {
		networkProjectID = cvoDetails.Project
	}

	hasSelfLink := strings.HasPrefix(subnetID, "https://www.googleapis.com/compute/") || strings.HasPrefix(subnetID, "projects/")
	if !hasSelfLink {
		cvoDetails.SubnetID = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], subnetID)
	} else {
		cvoDetails.SubnetID = subnetID
	}
	cvoDetails.SubnetPath = cvoDetails.SubnetID

	if c, ok := d.GetOk("firewall_rule"); ok {
		cvoDetails.FirewallRule = c.(string)
	}
	if c, ok := d.GetOk("firewall_tag_name_rule"); ok {
		cvoDetails.FirewallTagNameRule = c.(string)
	}
	if c, ok := d.GetOk("firewall_ip_ranges"); ok {
		cvoDetails.FirewallIPRanges = c.(bool)
	}
	if c, ok := d.GetOk("writing_speed_state"); ok {
		cvoDetails.WritingSpeedState = strings.ToUpper(c.(string))
	} else {
		cvoDetails.WritingSpeedState = "NORMAL"
	}

	if c, ok := d.GetOk("nss_account"); ok {
		cvoDetails.NssAccount = c.(string)
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

	if c, ok := d.GetOk("worm_retention_period_length"); ok {
		cvoDetails.WormRequest.RetentionPeriod.Length = c.(int)
	}
	if c, ok := d.GetOk("worm_retention_period_unit"); ok {
		cvoDetails.WormRequest.RetentionPeriod.Unit = c.(string)
	}

	// initialize the svmList for GCP CVO HA SVMs adding
	svmList := []gcpSVM{}

	if cvoDetails.IsHA {
		if cvoDetails.VsaMetadata.LicenseType == "capacity-paygo" {
			log.Print("Set licenseType as default value ha-capacity-paygo")
			cvoDetails.VsaMetadata.LicenseType = "ha-capacity-paygo"
		}
		if c, ok := d.GetOk("platform_serial_number_node1"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode1 = c.(string)
		}
		if c, ok := d.GetOk("platform_serial_number_node2"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode2 = c.(string)
		}

		if c, ok := d.GetOk("node2_zone"); ok {
			cvoDetails.HAParams.Node2Zone = c.(string)
		}
		if c, ok := d.GetOk("mediator_zone"); ok {
			cvoDetails.HAParams.MediatorZone = c.(string)
		}
		if c, ok := d.GetOk("vpc0_node_and_data_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC0NodeAndDataConnectivity = c.(string)
		}
		if c, ok := d.GetOk("vpc1_cluster_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC1ClusterConnectivity = c.(string)
		}
		if c, ok := d.GetOk("vpc2_ha_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC2HAConnectivity = c.(string)
		}
		if c, ok := d.GetOk("vpc3_data_replication"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
				c = fmt.Sprintf("projects/%s/global/networks/%s", networkProjectID, c.(string))
			}
			cvoDetails.HAParams.VPC3DataReplication = c.(string)
		}
		if c, ok := d.GetOk("subnet0_node_and_data_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
				c = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], c.(string))
			}
			cvoDetails.HAParams.Subnet0NodeAndDataConnectivity = c.(string)
		}
		if c, ok := d.GetOk("subnet1_cluster_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
				c = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], c.(string))
			}
			cvoDetails.HAParams.Subnet1ClusterConnectivity = c.(string)
		}
		if c, ok := d.GetOk("subnet2_ha_connectivity"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
				c = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], c.(string))
			}
			cvoDetails.HAParams.Subnet2HAConnectivity = c.(string)
		}
		if c, ok := d.GetOk("subnet3_data_replication"); ok {
			hasSelfLink := strings.HasPrefix(c.(string), "https://www.googleapis.com/compute/") || strings.HasPrefix(c.(string), "projects/")
			if !hasSelfLink {
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
		if c, ok := d.GetOk("vpc0_firewall_rule_tag_name"); ok {
			cvoDetails.HAParams.VPC0FirewallRuleTagName = c.(string)
		}
		if c, ok := d.GetOk("vpc1_firewall_rule_tag_name"); ok {
			cvoDetails.HAParams.VPC1FirewallRuleTagName = c.(string)
		}
		if c, ok := d.GetOk("vpc2_firewall_rule_tag_name"); ok {
			cvoDetails.HAParams.VPC2FirewallRuleTagName = c.(string)
		}
		if c, ok := d.GetOk("vpc3_firewall_rule_tag_name"); ok {
			cvoDetails.HAParams.VPC3FirewallRuleTagName = c.(string)
		}
		if c, ok := d.GetOk("svm"); ok {
			svms := c.(*schema.Set)
			svmList = expandGCPSVMs(svms)
		}
	}

	err = validateCVOGCPParams(cvoDetails)
	if err != nil {
		log.Print("Error validating parameters")
		return err
	}

	res, err := client.createCVOGCP(cvoDetails, clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("Error creating instance")
		return err
	}
	log.Printf("createCVOGCP %s result %#v  client_id %s", cvoDetails.Name, res, clientID)
	d.SetId(res.PublicID)
	d.Set("svm_name", res.SvmName)
	d.Set("writing_speed_state", res.OntapClusterProperties.WritingSpeedState)
	log.Printf("Created cvo: %v", res)

	// Add SVMs on GCP CVO HA
	for _, svm := range svmList {
		err := client.addSVMtoCVO(res.PublicID, clientID, svm.SvmName, isSaas, connectorIP)
		if err != nil {
			log.Printf("Error adding SVM %v: %v", svm.SvmName, err)
			return err
		}
	}

	return resourceCVOGCPRead(d, meta)
}

func resourceCVOGCPRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()

	connectorIP := ""
	clientID := d.Get("client_id").(string)

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	resp, err := client.getCVOProperties(id, clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("Error reading cvo")
		return err
	}
	d.Set("svm_name", resp.SvmName)
	d.Set("writing_speed_state", resp.OntapClusterProperties.WritingSpeedState)
	d.Set("instance_type", resp.ProviderProperties.InstanceType)

	return nil
}

func resourceCVOGCPDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting CVO: %#v", d)

	client := meta.(*Client)

	id := d.Id()
	clientID := d.Get("client_id").(string)

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	isHA := d.Get("is_ha").(bool)
	deleteErr := client.deleteCVOGCP(id, isHA, clientID, isSaas, connectorIP)
	if deleteErr != nil {
		log.Print("Error deleting cvo")
		return deleteErr
	}

	return nil
}

func resourceCVOGCPUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating CVO: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	// check if svm_password is changed
	if d.HasChange("svm_password") {
		respErr := updateCVOSVMPassword(d, meta, clientID, isSaas, connectorIP)
		if respErr != nil {
			return respErr
		}
	}

	//  check if svm_name is changed
	if d.HasChange("svm_name") {
		svmName, svmNewName := d.GetChange("svm_name")
		respErr := client.updateCVOSVMName(d, clientID, svmName.(string), svmNewName.(string), isSaas, connectorIP)
		if respErr != nil {
			return respErr
		}
	}

	// check if svm list changes
	if d.Get("is_ha").(bool) && d.HasChange("svm") {
		respErr := client.updateCVOSVMs(d, clientID, isSaas, connectorIP)
		if respErr != nil {
			return respErr
		}
	}

	instanceType := d.Get("instance_type").(string)
	// In both single and HA case, flash_cache only can be set when the selected instance_type
	if _, ok := d.GetOk("flash_cache"); ok {
		match, _ := regexp.MatchString("^n2-standard-(16|32|48|64)$", instanceType)
		if !match {
			return fmt.Errorf("instance_type has to be one of n2-standard-16,32,48,64")
		}
		if d.Get("is_ha").(bool) && d.Get("writing_speed_state").(string) == "" {
			return fmt.Errorf("in HA, writing_speed_state has to be set when flash_cache is set")
		}
	}

	// check if license_type and instance type are changed
	if d.HasChange("instance_type") || d.HasChange("license_type") {
		respErr := updateCVOLicenseInstanceType(d, meta, clientID, isSaas, connectorIP)
		if respErr != nil {
			return respErr
		}
	}

	// check if tier_level is changed
	if d.HasChange("tier_level") {
		respErr := updateCVOTierLevel(d, meta, clientID, isSaas, connectorIP)
		if respErr != nil {
			return respErr
		}
	}

	// check if writing_speed_state is changed
	if d.HasChange("writing_speed_state") {
		currentWritingSpeedState, expectWritingSpeedState := d.GetChange("writing_speed_state")
		if currentWritingSpeedState.(string) == "" && strings.ToUpper(expectWritingSpeedState.(string)) == "NORMAL" {
			d.Set("writing_speed_state", expectWritingSpeedState.(string))
			log.Print("writing_speed_state: default value is NORMAL. No change call is needed.")
			return nil
		}
		respErr := updateCVOWritingSpeedState(d, meta, clientID, isSaas, connectorIP)
		if respErr != nil {
			return respErr
		}
		return nil
	}

	// check if gcp_label has changes
	if d.HasChange("gcp_label") {
		respErr := updateCVOUserTags(d, meta, "gcp_label", clientID, isSaas, connectorIP)
		if respErr != nil {
			return respErr
		}
		return resourceCVOGCPRead(d, meta)
	}
	// upgrade ontap version
	upgradeErr := client.checkAndDoUpgradeOntapVersion(d, clientID, isSaas, connectorIP)
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
	clientID := d.Get("client_id").(string)

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return false, err
	}

	name := d.Get("name").(string)
	resID, err := client.findWorkingEnvironmentByName(name, clientID, isSaas, connectorIP)
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
