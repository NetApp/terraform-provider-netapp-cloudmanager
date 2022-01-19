package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCVOAzure() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOAzureCreate,
		Read:   resourceCVOAzureRead,
		Delete: resourceCVOAzureDelete,
		Update: resourceCVOAzureUpdate,
		Exists: resourceCVOAzureExists,
		Importer: &schema.ResourceImporter{
			State: resourceCVOAzureImport,
		},
		CustomizeDiff: resourceCVOAzureCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscription_id": {
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
				Default:      "AZURE",
				ValidateFunc: validation.StringInSlice([]string{"AZURE", "NONE"}, false),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Premium_LRS",
				ValidateFunc: validation.StringInSlice([]string{"Premium_LRS", "Standard_LRS", "StandardSSD_LRS"}, false),
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"disk_size_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "TB",
				ValidateFunc: validation.StringInSlice([]string{"GB", "TB"}, false),
			},
			"azure_encryption_parameters": {
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
				Default:      "azure-cot-standard-paygo",
				ValidateFunc: validation.StringInSlice(AzureLicenseTypes, false),
			},
			"capacity_package_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Essential", "Professional", "Freemium"}, false),
			},
			"provided_license": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Standard_DS4_v2",
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnet_resource_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr": {
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
				ValidateFunc: validation.StringInSlice([]string{"normal", "cool"}, false),
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
				Default:      "Blob",
				ValidateFunc: validation.StringInSlice([]string{"Blob", "NONE"}, false),
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
			"allow_deploy_in_existing_rg": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"azure_tag": {
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
			"ha_enable_https": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
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

func resourceCVOAzureCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating CVO Azure: %#v", d)

	client := meta.(*Client)

	cvoDetails := createCVOAzureDetails{}

	cvoDetails.Name = d.Get("name").(string)
	cvoDetails.Region = d.Get("location").(string)
	cvoDetails.SubscriptionID = d.Get("subscription_id").(string)
	cvoDetails.DataEncryptionType = d.Get("data_encryption_type").(string)
	cvoDetails.WorkspaceID = d.Get("workspace_id").(string)
	cvoDetails.StorageType = d.Get("storage_type").(string)
	cvoDetails.SvmPassword = d.Get("svm_password").(string)
	capacityTier := d.Get("capacity_tier").(string)
	if capacityTier == "Blob" {
		cvoDetails.CapacityTier = capacityTier
		cvoDetails.TierLevel = d.Get("tier_level").(string)
	}
	cvoDetails.OptimizedNetworkUtilization = true
	cvoDetails.BackupVolumesToCbs = d.Get("backup_volumes_to_cbs").(bool)
	cvoDetails.EnableCompliance = d.Get("enable_compliance").(bool)
	cvoDetails.EnableMonitoring = d.Get("enable_monitoring").(bool)
	if c, ok := d.GetOk("azure_tag"); ok {
		tags := c.(*schema.Set)
		if tags.Len() > 0 {
			cvoDetails.AzureTags = expandUserTags(tags)
		}
	}
	cvoDetails.DiskSize.Size = d.Get("disk_size").(int)
	cvoDetails.DiskSize.Unit = d.Get("disk_size_unit").(string)
	cvoDetails.VsaMetadata.OntapVersion = d.Get("ontap_version").(string)
	cvoDetails.VsaMetadata.UseLatestVersion = d.Get("use_latest_version").(bool)
	cvoDetails.VsaMetadata.LicenseType = d.Get("license_type").(string)
	cvoDetails.VsaMetadata.InstanceType = d.Get("instance_type").(string)

	if cvoDetails.DataEncryptionType == "AZURE" {
		if c, ok := d.GetOk("azure_encryption_parameters"); ok {
			cvoDetails.AzureEncryptionParameters.Key = c.(string)
		}
	}
	if c, ok := d.GetOk("cidr"); ok {
		cvoDetails.Cidr = c.(string)
	}

	if c, ok := d.GetOk("writing_speed_state"); ok {
		cvoDetails.WritingSpeedState = c.(string)
	}

	if c, ok := d.GetOk("nss_account"); ok {
		cvoDetails.NssAccount = c.(string)
	}

	if c, ok := d.GetOk("security_group_id"); ok {
		cvoDetails.SecurityGroupID = c.(string)
	}

	if c, ok := d.GetOk("cloud_provider_account"); ok {
		cvoDetails.CloudProviderAccount = c.(string)
	}

	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	if c, ok := d.GetOk("capacity_package_name"); ok {
		cvoDetails.VsaMetadata.CapacityPackageName = c.(string)
		log.Print("Get capacity_package_name")
	}

	if c, ok := d.GetOk("provided_license"); ok {
		cvoDetails.VsaMetadata.ProvidedLicense = c.(string)
	}

	if c, ok := d.GetOk("resource_group"); ok {
		cvoDetails.ResourceGroup = c.(string)
		cvoDetails.AllowDeployInExistingRg = d.Get("allow_deploy_in_existing_rg").(bool)
	} else {
		cvoDetails.ResourceGroup = cvoDetails.Name + "-rg"
	}

	if c, ok := d.GetOk("serial_number"); ok {
		cvoDetails.SerialNumber = c.(string)
	}

	cvoDetails.IsHA = d.Get("is_ha").(bool)
	if cvoDetails.IsHA == true {
		if c, ok := d.GetOk("platform_serial_number_node1"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode1 = c.(string)
		}

		if c, ok := d.GetOk("platform_serial_number_node2"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode2 = c.(string)
		}
		cvoDetails.HAParams.EnableHTTPS = d.Get("ha_enable_https").(bool)
	}

	err := validateCVOAzureParams(cvoDetails)
	if err != nil {
		log.Print("Error validating parameters")
		return err
	}

	cvoDetails.VnetForInternal = d.Get("vnet_id").(string)

	resourceGroup := cvoDetails.ResourceGroup
	if c, ok := d.GetOk("vnet_resource_group"); ok {
		cvoDetails.VnetResourceGroup = c.(string)
		resourceGroup = cvoDetails.VnetResourceGroup
	}

	resourceGroupPath := fmt.Sprintf("subscriptions/%s/resourceGroups/%s", cvoDetails.SubscriptionID, resourceGroup)
	vnetFormat := "/%s/providers/Microsoft.Network/virtualNetworks/%s"
	if client.GetSimulator() {
		log.Print("In simulator env...")
		vnetFormat = "%s/%s"
	}

	vnet := fmt.Sprintf(vnetFormat, resourceGroupPath, cvoDetails.VnetForInternal)
	cvoDetails.VnetID = vnet
	cvoDetails.SubnetID = fmt.Sprintf("%s/subnets/%s", vnet, d.Get("subnet_id").(string))

	res, err := client.createCVOAzure(cvoDetails)
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(res.PublicID)
	d.Set("svm_name", res.SvmName)
	log.Printf("Created cvo: %v", res)

	return resourceCVOAzureRead(d, meta)
}

func resourceCVOAzureRead(d *schema.ResourceData, meta interface{}) error {
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

func resourceCVOAzureDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting CVO: %#v", d)

	client := meta.(*Client)

	id := d.Id()
	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	isHA := d.Get("is_ha").(bool)

	deleteErr := client.deleteCVOAzure(id, isHA)
	if deleteErr != nil {
		log.Print("Error deleting cvo")
		return deleteErr
	}

	return nil
}

func resourceCVOAzureUpdate(d *schema.ResourceData, meta interface{}) error {
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
	if d.HasChange("tier_level") && d.Get("capacity_tier").(string) == "Blob" {
		respErr := updateCVOTierLevel(d, meta)
		if respErr != nil {
			return respErr
		}
	}

	// check if aws_tag has changes
	if d.HasChange("azure_tag") {
		respErr := updateCVOUserTags(d, meta, "azure_tag")
		if respErr != nil {
			return respErr
		}
		return resourceCVOAzureRead(d, meta)
	}
	// upgrade ontap version
	// only when the upgrade_ontap_version is true and the ontap_version is not "latest"
	upgradeErr := client.checkAndDoUpgradeOntapVersion(d)
	if upgradeErr != nil {
		return upgradeErr
	}

	return nil
}

func resourceCVOAzureCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	respErr := checkUserTagDiff(diff, "azure_tag", "tag_key")
	if respErr != nil {
		return respErr
	}
	return nil
}

func resourceCVOAzureExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}

	resID, err := client.getCVOAzure(id)
	if err != nil {
		log.Print("Error getting cvo")
		return false, err
	}

	log.Print(resID)
	log.Print(id)

	if resID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}

func resourceCVOAzureImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, fmt.Errorf("CVO Azure resource's import function is disabled")
}
