package cloudmanager

import (
	"log"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCVOAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOAWSCreate,
		Read:   resourceCVOAWSRead,
		Delete: resourceCVOAWSDelete,
		Exists: resourceCVOAWSExists,
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
				ValidateFunc: validation.StringInSlice([]string{"gp2", "io1", "sc1", "st1"}, false),
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
			"ontap_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "latest",
			},
			"use_latest_version": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "cot-standard-paygo",
				ValidateFunc: validation.StringInSlice(AWSLicenseTypes, false),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ForceNew:  true,
				Sensitive: true,
			},
			"tier_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "normal",
				ValidateFunc: validation.StringInSlice([]string{"normal", "ia", "ia-single", "intelligent"}, false),
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
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "S3",
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
				Default:  false,
			},
			"enable_compliance": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
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
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"aws_tag": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCVOAWSCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating CVO: %#v", d)

	client := meta.(*Client)

	cvoDetails := createCVOAWSDetails{}

	cvoDetails.Name = d.Get("name").(string)
	cvoDetails.Region = d.Get("region").(string)
	cvoDetails.DataEncryptionType = d.Get("data_encryption_type").(string)
	cvoDetails.WorkspaceID = d.Get("workspace_id").(string)
	cvoDetails.EbsVolumeType = d.Get("ebs_volume_type").(string)
	cvoDetails.SvmPassword = d.Get("svm_password").(string)
	capacityTier := d.Get("capacity_tier").(string)
	if capacityTier == "S3" {
		cvoDetails.CapacityTier = capacityTier
		cvoDetails.TierLevel = d.Get("tier_level").(string)
	}
	cvoDetails.OptimizedNetworkUtilization = d.Get("optimized_network_utilization").(bool)
	cvoDetails.InstanceTenancy = d.Get("instance_tenancy").(string)
	cvoDetails.BackupVolumesToCbs = d.Get("backup_volumes_to_cbs").(bool)
	cvoDetails.EnableCompliance = d.Get("enable_compliance").(bool)
	cvoDetails.EnableMonitoring = d.Get("enable_monitoring").(bool)
	if c, ok := d.GetOk("aws_tag"); ok {
		tags := c.(*schema.Set)
		if tags.Len() > 0 {
			cvoDetails.AwsTags = expandAWSTags(tags)
		}
	}
	cvoDetails.EbsVolumeSize.Size = d.Get("ebs_volume_size").(int)
	cvoDetails.EbsVolumeSize.Unit = d.Get("ebs_volume_size_unit").(string)
	cvoDetails.VsaMetadata.OntapVersion = d.Get("ontap_version").(string)
	cvoDetails.VsaMetadata.UseLatestVersion = d.Get("use_latest_version").(bool)
	cvoDetails.VsaMetadata.LicenseType = d.Get("license_type").(string)
	cvoDetails.VsaMetadata.InstanceType = d.Get("instance_type").(string)

	if c, ok := d.GetOk("vpc_id"); ok {
		cvoDetails.VpcID = c.(string)
	}

	if c, ok := d.GetOk("writing_speed_state"); ok {
		cvoDetails.WritingSpeedState = c.(string)
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

	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	cvoDetails.IsHA = d.Get("is_ha").(bool)

	if cvoDetails.IsHA == true {
		cvoDetails.HAParams.FailoverMode = d.Get("failover_mode").(string)
		cvoDetails.HAParams.Node1SubnetID = d.Get("node1_subnet_id").(string)
		cvoDetails.HAParams.Node2SubnetID = d.Get("node2_subnet_id").(string)
		cvoDetails.HAParams.MediatorSubnetID = d.Get("mediator_subnet_id").(string)
		cvoDetails.HAParams.MediatorKeyPairName = d.Get("mediator_key_pair_name").(string)

		if o, ok := d.GetOkExists("mediator_assign_public_ip"); ok {
			mediatorAssignPublicIP := o.(bool)
			cvoDetails.HAParams.MediatorAssignPublicIP = &mediatorAssignPublicIP
		}

		cvoDetails.HAParams.ClusterFloatingIP = d.Get("cluster_floating_ip").(string)
		cvoDetails.HAParams.DataFloatingIP = d.Get("data_floating_ip").(string)
		cvoDetails.HAParams.DataFloatingIP2 = d.Get("data_floating_ip2").(string)
		cvoDetails.HAParams.SvmFloatingIP = d.Get("svm_floating_ip").(string)
		routeTableIds := d.Get("route_table_ids")
		for _, routeTableID := range routeTableIds.([]interface{}) {
			cvoDetails.HAParams.RouteTableIds = append(cvoDetails.HAParams.RouteTableIds, routeTableID.(string))
		}
		if c, ok := d.GetOk("platform_serial_number_node1"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode1 = c.(string)
		}

		if c, ok := d.GetOk("platform_serial_number_node2"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode2 = c.(string)
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

	res, err := client.createCVOAWS(cvoDetails)
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(res.PublicID)

	log.Printf("Created cvo: %v", res)

	return resourceCVOAWSRead(d, meta)
}

func resourceCVOAWSRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceCVOAWSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting CVO: %#v", d)

	client := meta.(*Client)

	id := d.Id()
	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	isHA := d.Get("is_ha").(bool)

	deleteErr := client.deleteCVO(id, isHA)
	if deleteErr != nil {
		log.Print("Error deleting cvo")
		return deleteErr
	}

	return nil
}

func resourceCVOAWSExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	resID, err := client.getCVOAWS(id)
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
