package cloudmanager

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCVOAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOAWSCreate,
		Read:   resourceCVOAWSRead,
		Delete: resourceCVOAWSDelete,
		Update: resourceCVOAWSUpdate,
		Exists: resourceCVOAWSExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceCVOAWSCustomizeDiff,
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
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_encryption_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "AWS",
				ValidateFunc: validation.StringInSlice([]string{"AWS", "NONE"}, false),
			},
			"ebs_volume_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "gp2",
				ValidateFunc: validation.StringInSlice([]string{"gp3", "gp2", "io1", "sc1", "st1"}, false),
			},
			"ebs_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"ebs_volume_size_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "TB",
				ValidateFunc: validation.StringInSlice([]string{"GB", "TB"}, false),
			},
			"aws_encryption_kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"aws_encryption_kms_key_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ValidateFunc: validation.StringInSlice(AWSLicenseTypes, false),
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
				Default:  "m5.2xlarge",
			},
			"platform_serial_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"svm_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"tier_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "normal",
				ValidateFunc: validation.StringInSlice([]string{"normal", "ia", "ia-single", "intelligent"}, false),
			},
			"nss_account": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"writing_speed_state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL", "HIGH"}, true),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"S3", "NONE"}, false),
			},
			"instance_tenancy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "default",
				ValidateFunc: validation.StringInSlice([]string{"default", "dedicated"}, false),
			},
			"instance_profile_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cloud_provider_account": {
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
			"enable_monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"optimized_network_utilization": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"cluster_key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"aws_tag": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
			"failover_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"PrivateIP", "FloatingIP"}, false),
			},
			"mediator_assign_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"node1_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"node2_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mediator_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mediator_key_pair_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mediator_instance_profile_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_floating_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_floating_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"data_floating_ip2": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"svm_floating_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"route_table_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"upgrade_ontap_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"mediator_security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"assume_role_arn": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"retries": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  60,
			},
		},
	}
}

func resourceCVOAWSCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating CVO: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	client.Retries = d.Get("retries").(int)

	cvoDetails := createCVOAWSDetails{}

	cvoDetails.Name = d.Get("name").(string)
	cvoDetails.Region = d.Get("region").(string)
	cvoDetails.DataEncryptionType = d.Get("data_encryption_type").(string)
	cvoDetails.WorkspaceID = d.Get("workspace_id").(string)
	cvoDetails.EbsVolumeType = d.Get("ebs_volume_type").(string)
	cvoDetails.SvmPassword = d.Get("svm_password").(string)
	if c, ok := d.GetOk("svm_name"); ok {
		cvoDetails.SvmName = c.(string)
	}
	capacityTier := d.Get("capacity_tier").(string)
	if capacityTier == "S3" {
		cvoDetails.CapacityTier = capacityTier
		cvoDetails.TierLevel = d.Get("tier_level").(string)
	}
	cvoDetails.OptimizedNetworkUtilization = d.Get("optimized_network_utilization").(bool)
	cvoDetails.InstanceTenancy = d.Get("instance_tenancy").(string)

	if c, ok := d.GetOk("backup_volumes_to_cbs"); ok {
		cvoDetails.BackupVolumesToCbs = c.(bool)
	}

	if c, ok := d.GetOk("enable_compliance"); ok {
		cvoDetails.EnableCompliance = c.(bool)
	}

	if c, ok := d.GetOk("cluster_key_pair_name"); ok {
		cvoDetails.ClusterKeyPairName = c.(string)
	}

	cvoDetails.EnableMonitoring = d.Get("enable_monitoring").(bool)
	if c, ok := d.GetOk("aws_tag"); ok {
		tags := c.(*schema.Set)
		if tags.Len() > 0 {
			cvoDetails.AwsTags = expandUserTags(tags)
		}
	}

	if c, ok := d.GetOk("saas_subscription_id"); ok {
		cvoDetails.SaasSubscriptionID = c.(string)
	}

	cvoDetails.EbsVolumeSize.Size = d.Get("ebs_volume_size").(int)
	cvoDetails.EbsVolumeSize.Unit = d.Get("ebs_volume_size_unit").(string)
	cvoDetails.VsaMetadata.OntapVersion = d.Get("ontap_version").(string)
	cvoDetails.VsaMetadata.UseLatestVersion = d.Get("use_latest_version").(bool)
	cvoDetails.VsaMetadata.LicenseType = d.Get("license_type").(string)
	cvoDetails.VsaMetadata.InstanceType = d.Get("instance_type").(string)

	// by Capacity
	if c, ok := d.GetOk("capacity_package_name"); ok {
		cvoDetails.VsaMetadata.CapacityPackageName = c.(string)
	} else {
		// by Capacity - set default capacity package name
		if strings.HasSuffix(cvoDetails.VsaMetadata.LicenseType, "capacity-paygo") {
			cvoDetails.VsaMetadata.CapacityPackageName = "Essential"
		}
	}

	if cvoDetails.DataEncryptionType == "AWS" {
		// Only one of KMS key id or KMS arn should be specified
		if c, ok := d.GetOk("aws_encryption_kms_key_id"); ok {
			cvoDetails.AwsEncryptionParameters.KmsKeyID = c.(string)
		}

		if c, ok := d.GetOk("aws_encryption_kms_key_arn"); ok {
			cvoDetails.AwsEncryptionParameters.KmsKeyArn = c.(string)
		}
	}

	if c, ok := d.GetOk("vpc_id"); ok {
		cvoDetails.VpcID = c.(string)
	}

	if c, ok := d.GetOk("writing_speed_state"); ok {
		cvoDetails.WritingSpeedState = strings.ToUpper(c.(string))
	}

	if c, ok := d.GetOk("platform_serial_number"); ok {
		cvoDetails.VsaMetadata.PlatformSerialNumber = c.(string)
	}

	if c, ok := d.GetOk("nss_account"); ok {
		cvoDetails.NssAccount = c.(string)
	}

	if c, ok := d.GetOk("iops"); ok {
		cvoDetails.IOPS = c.(int)
	}

	if c, ok := d.GetOk("throughput"); ok {
		cvoDetails.Throughput = c.(int)
	}

	if c, ok := d.GetOk("instance_profile_name"); ok {
		cvoDetails.InstanceProfileName = c.(string)
	}

	if c, ok := d.GetOk("security_group_id"); ok {
		cvoDetails.SecurityGroupID = c.(string)
	}

	if c, ok := d.GetOk("cloud_provider_account"); ok {
		cvoDetails.CloudProviderAccount = c.(string)
	}

	if c, ok := d.GetOk("kms_key_id"); ok {
		cvoDetails.AwsEncryptionParameters.KmsKeyID = c.(string)
	}

	if c, ok := d.GetOk("provided_license"); ok {
		cvoDetails.VsaMetadata.ProvidedLicense = c.(string)
	}

	if c, ok := d.GetOk("worm_retention_period_length"); ok {
		cvoDetails.WormRequest.RetentionPeriod.Length = c.(int)
	}
	if c, ok := d.GetOk("worm_retention_period_unit"); ok {
		cvoDetails.WormRequest.RetentionPeriod.Unit = c.(string)
	}
	cvoDetails.IsHA = d.Get("is_ha").(bool)

	if cvoDetails.IsHA {
		if cvoDetails.VsaMetadata.LicenseType == "capacity-paygo" {
			log.Print("Set licenseType as default value ha-capacity-paygo")
			cvoDetails.VsaMetadata.LicenseType = "ha-capacity-paygo"
		}
		cvoDetails.HAParams.FailoverMode = d.Get("failover_mode").(string)
		cvoDetails.HAParams.Node1SubnetID = d.Get("node1_subnet_id").(string)
		cvoDetails.HAParams.Node2SubnetID = d.Get("node2_subnet_id").(string)
		cvoDetails.HAParams.MediatorSubnetID = d.Get("mediator_subnet_id").(string)
		cvoDetails.HAParams.MediatorKeyPairName = d.Get("mediator_key_pair_name").(string)

		if o, ok := d.GetOkExists("mediator_assign_public_ip"); ok {
			cvoDetails.HAParams.MediatorAssignPublicIP = o.(bool)
		}

		cvoDetails.HAParams.ClusterFloatingIP = d.Get("cluster_floating_ip").(string)
		cvoDetails.HAParams.DataFloatingIP = d.Get("data_floating_ip").(string)
		cvoDetails.HAParams.DataFloatingIP2 = d.Get("data_floating_ip2").(string)
		cvoDetails.HAParams.SvmFloatingIP = d.Get("svm_floating_ip").(string)
		routeTableIds := d.Get("route_table_ids")
		for _, routeTableID := range routeTableIds.(*schema.Set).List() {
			cvoDetails.HAParams.RouteTableIds = append(cvoDetails.HAParams.RouteTableIds, routeTableID.(string))
		}
		if c, ok := d.GetOk("mediator_instance_profile_name"); ok {
			cvoDetails.HAParams.MediatorInstanceProfileName = c.(string)
		}
		if c, ok := d.GetOk("platform_serial_number_node1"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode1 = c.(string)
		}

		if c, ok := d.GetOk("platform_serial_number_node2"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode2 = c.(string)
		}
		if c, ok := d.GetOk("mediator_security_group_id"); ok {
			cvoDetails.HAParams.MediatorSecurityGroupID = c.(string)
		}
		if c, ok := d.GetOk("assume_role_arn"); ok {
			cvoDetails.HAParams.AssumeRoleArn = c.(string)
		}
	} else {
		if c, ok := d.GetOk("subnet_id"); ok {
			cvoDetails.SubnetID = c.(string)
		}
	}

	err := validateCVOParams(cvoDetails)
	if err != nil {
		log.Print("Error validating parameters")
		return err
	}

	res, err := client.createCVOAWS(cvoDetails, clientID)
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(res.PublicID)
	d.Set("svm_name", res.SvmName)

	log.Printf("Created cvo: %v", res)

	return resourceCVOAWSRead(d, meta)
}

func resourceCVOAWSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()

	clientID := d.Get("client_id").(string)

	resp, err := client.getCVOProperties(id, clientID, true, "")
	if err != nil {
		log.Print("Error reading cvo")
		return err
	}
	d.Set("svm_name", resp.SvmName)
	if c, ok := d.GetOk("writing_speed_state"); ok {
		if strings.EqualFold(c.(string), resp.OntapClusterProperties.WritingSpeedState) {
			d.Set("writing_speed_state", c.(string))
		} else {
			d.Set("writing_speed_state", resp.OntapClusterProperties.WritingSpeedState)
		}
	}
	if _, ok := d.GetOk("route_table_ids"); ok {
		d.Set("route_table_ids", resp.HAProperties.RouteTables)
	}

	return nil
}

func resourceCVOAWSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting CVO: %#v", d)

	client := meta.(*Client)

	id := d.Id()
	clientID := d.Get("client_id").(string)

	isHA := d.Get("is_ha").(bool)

	deleteErr := client.deleteCVO(id, isHA, clientID)
	if deleteErr != nil {
		log.Print("Error deleting cvo")
		return deleteErr
	}

	return nil
}

func resourceCVOAWSUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating CVO: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)

	// check if svm_password is changed
	if d.HasChange("svm_password") {
		respErr := updateCVOSVMPassword(d, meta, clientID, true, "")
		if respErr != nil {
			return respErr
		}
	}

	//  check if svm_name is changed
	if d.HasChange("svm_name") {
		svmName, svmNewName := d.GetChange("svm_name")
		respErr := client.updateCVOSVMName(d, clientID, svmName.(string), svmNewName.(string), true, "")
		if respErr != nil {
			return respErr
		}
	}

	// check if license_type and instance type are changed
	if d.HasChange("instance_type") || d.HasChange("license_type") {
		respErr := updateCVOLicenseInstanceType(d, meta, clientID, true, "")
		if respErr != nil {
			return respErr
		}
	}

	// check if tier_level is changed
	if d.HasChange("tier_level") && d.Get("capacity_tier").(string) == "S3" {
		respErr := updateCVOTierLevel(d, meta, clientID, true, "")
		if respErr != nil {
			return respErr
		}
	}

	// check if aws_tag has changes
	if d.HasChange("aws_tag") {
		respErr := updateCVOUserTags(d, meta, "aws_tag", clientID, true, "")
		if respErr != nil {
			return respErr
		}
		return resourceCVOAWSRead(d, meta)
	}

	// check if writing_speed_state is changed
	if d.HasChange("writing_speed_state") {
		currentWritingSpeedState, expectWritingSpeedState := d.GetChange("writing_speed_state")
		if currentWritingSpeedState.(string) == "" && strings.ToUpper(expectWritingSpeedState.(string)) == "NORMAL" {
			d.Set("writing_speed_state", expectWritingSpeedState.(string))
			log.Print("writing_speed_state: default value is NORMAL. No change call is needed.")
			return nil
		}
		respErr := updateCVOWritingSpeedState(d, meta, clientID, true, "")
		if respErr != nil {
			return respErr
		}

		return nil
	}

	// upgrade ontap version
	// only when the upgrade_ontap_version is true and the ontap_version is not "latest"
	upgradeErr := client.checkAndDoUpgradeOntapVersion(d, clientID, true, "")
	if upgradeErr != nil {
		return upgradeErr
	}

	return nil
}

func resourceCVOAWSCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	respErr := checkUserTagDiff(diff, "aws_tag", "tag_key")
	if respErr != nil {
		return respErr
	}
	return nil
}

func resourceCVOAWSExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	clientID := d.Get("client_id").(string)

	resID, err := client.getCVOAWS(id, clientID)
	if err != nil {
		log.Print("Error getting cvo")
		return false, err
	}

	if resID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
