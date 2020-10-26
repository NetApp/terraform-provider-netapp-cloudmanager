package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCVOGCP() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOGCPCreate,
		Read:   resourceCVOGCPRead,
		Delete: resourceCVOGCPDelete,
		Exists: resourceCVOGCPExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"gcp_volume_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "pd-ssd",
				ValidateFunc: validation.StringInSlice([]string{"pd-standard", "pd-ssd"}, false),
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
				Default:      "gcp-cot-standard-paygo",
				ValidateFunc: validation.StringInSlice(GCPLicenseTypes, false),
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
				ForceNew:  true,
				Sensitive: true,
			},
			"tier_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
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
				Default:  "NORMAL",
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
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"label_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
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
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
	cvoDetails.CapacityTier = d.Get("capacity_tier").(string)
	cvoDetails.TierLevel = d.Get("tier_level").(string)
	cvoDetails.GCPVolumeSize.Size = d.Get("gcp_volume_size").(int)
	cvoDetails.GCPVolumeSize.Unit = d.Get("gcp_volume_size_unit").(string)
	cvoDetails.VsaMetadata.OntapVersion = d.Get("ontap_version").(string)
	cvoDetails.VsaMetadata.UseLatestVersion = d.Get("use_latest_version").(bool)
	cvoDetails.VsaMetadata.LicenseType = d.Get("license_type").(string)
	cvoDetails.VpcID = d.Get("vpc_id").(string)
	cvoDetails.Project = d.Get("project_id").(string)
	cvoDetails.WritingSpeedState = d.Get("writing_speed_state").(string)
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
	}

	if networkProjectID != "" {
		cvoDetails.SubnetID = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", networkProjectID, cvoDetails.Region[0:len(cvoDetails.Region)-2], subnetID)
	} else {
		cvoDetails.SubnetID = fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", cvoDetails.Project, cvoDetails.Region[0:len(cvoDetails.Region)-2], subnetID)
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

	if c, ok := d.GetOk("serial_number"); ok {
		cvoDetails.SerialNumber = c.(string)
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

	deleteErr := client.deleteCVOGCP(id)
	if deleteErr != nil {
		log.Print("Error deleting cvo")
		return deleteErr
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
