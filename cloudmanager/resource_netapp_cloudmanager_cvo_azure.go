package cloudmanager

import (
	"fmt"
	"log"
	"strings"

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
			"saas_subscription_id": {
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
				Default:      "AZURE",
				ValidateFunc: validation.StringInSlice([]string{"AZURE", "NONE"}, false),
			},
			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "Premium_LRS",
				ValidateFunc: validation.StringInSlice([]string{"Premium_LRS", "Standard_LRS", "StandardSSD_LRS", "Premium_ZRS"}, false),
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
				Type:     schema.TypeSet,
				MaxItems: 1,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"vault_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"user_assigned_identity": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
				ValidateFunc: validation.StringInSlice(AzureLicenseTypes, false),
			},
			"capacity_package_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Essential", "Professional", "Freemium", "Edge", "Optimized"}, false),
			},
			"provided_license": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true, // specify Standard_E8ds_v5 for default value
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
			"availability_zone": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone_node1": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"availability_zone_node2": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"nss_account": {
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
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
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
			"open_security_group": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
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
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == ""
				},
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
						"root_volume_aggregate": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								// Detailed validation performed in CustomizeDiff with full context
								return
							},
						},
					},
				},
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
			"storage_account_network_access": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Disabled", "Enabled", "SecuredByPerimeter"}, false),
				Description:  "Controls network access for the storage account. Accepted values: Enabled, Disabled, SecuredByPerimeter.",
			},
		},
	}
}

func resourceCVOAzureCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating CVO Azure: %#v", d)

	client := meta.(*Client)
	client.Retries = d.Get("retries").(int)

	cvoDetails := createCVOAzureDetails{}

	cvoDetails.Name = d.Get("name").(string)
	clientID := d.Get("client_id").(string)
	cvoDetails.Region = d.Get("location").(string)
	cvoDetails.SubscriptionID = d.Get("subscription_id").(string)
	cvoDetails.DataEncryptionType = d.Get("data_encryption_type").(string)
	cvoDetails.WorkspaceID = d.Get("workspace_id").(string)
	cvoDetails.StorageType = d.Get("storage_type").(string)
	cvoDetails.SvmPassword = d.Get("svm_password").(string)

	if c, ok := d.GetOk("saas_subscription_id"); ok {
		cvoDetails.SaasSubscriptionID = c.(string)
	}
	if c, ok := d.GetOk("svm_name"); ok {
		cvoDetails.SvmName = c.(string)
	}
	capacityTier := d.Get("capacity_tier").(string)
	if capacityTier == "Blob" {
		cvoDetails.CapacityTier = capacityTier
		cvoDetails.TierLevel = d.Get("tier_level").(string)
	}
	if c, ok := d.GetOk("availability_zone"); ok {
		cvoDetails.AvailabilityZone = c.(int)
	}
	cvoDetails.OptimizedNetworkUtilization = true
	if c, ok := d.GetOk("backup_volumes_to_cbs"); ok {
		cvoDetails.BackupVolumesToCbs = c.(bool)
	}

	if c, ok := d.GetOk("enable_compliance"); ok {
		cvoDetails.EnableCompliance = c.(bool)
	}
	cvoDetails.EnableMonitoring = d.Get("enable_monitoring").(bool)
	if c, ok := d.GetOk("azure_tag"); ok {
		tags := c.(*schema.Set)
		if tags.Len() > 0 {
			cvoDetails.AzureTags = expandUserTags(tags)
		}
	}
	if c, ok := d.GetOk("open_security_group"); ok {
		cvoDetails.OpenSecurityGroup = c.(bool)
	}
	cvoDetails.DiskSize.Size = d.Get("disk_size").(int)
	cvoDetails.DiskSize.Unit = d.Get("disk_size_unit").(string)
	cvoDetails.VsaMetadata.OntapVersion = d.Get("ontap_version").(string)
	cvoDetails.VsaMetadata.UseLatestVersion = d.Get("use_latest_version").(bool)
	cvoDetails.VsaMetadata.LicenseType = d.Get("license_type").(string)
	licenseType := d.Get("license_type").(string)
	if licenseType != "capacity-paygo" && licenseType != "ha-capacity-paygo" {
		return fmt.Errorf(
			"node-based licenses are no longer supported for new CVOs. "+
				"Please use 'capacity-paygo' for single node or 'ha-capacity-paygo' for HA. "+
				"Management of existing node-based CVOs created with last provider version is supported. Migration to capacity based licenses from NetApp console is recommended."+
				"Current license_type: '%s'",
			licenseType,
		)
	}
	if c, ok := d.GetOk("capacity_package_name"); ok {
		cvoDetails.VsaMetadata.CapacityPackageName = c.(string)
	} else {
		// by Capacity - set default capacity package name
		if strings.HasSuffix(cvoDetails.VsaMetadata.LicenseType, "capacity-paygo") {
			cvoDetails.VsaMetadata.CapacityPackageName = "Essential"
		}
	}
	cvoDetails.VsaMetadata.InstanceType = d.Get("instance_type").(string)

	if c, ok := d.GetOk("cidr"); ok {
		cvoDetails.Cidr = c.(string)
	}

	if c, ok := d.GetOk("writing_speed_state"); ok {
		cvoDetails.WritingSpeedState = strings.ToUpper(c.(string))
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

	if cvoDetails.DataEncryptionType == "AZURE" {
		if c, ok := d.GetOk("azure_encryption_parameters"); ok {
			parametersSet := c.(*schema.Set)
			cvoDetails.AzureEncryptionParameters = expendEncryptionParameters(parametersSet)
		}
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
		if c, ok := d.GetOk("platform_serial_number_node1"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode1 = c.(string)
		}

		if c, ok := d.GetOk("platform_serial_number_node2"); ok {
			cvoDetails.HAParams.PlatformSerialNumberNode2 = c.(string)
		}
		cvoDetails.HAParams.EnableHTTPS = d.Get("ha_enable_https").(bool)
		if c, ok := d.GetOk("availability_zone_node1"); ok {
			cvoDetails.HAParams.AvailabilityZoneNode1 = c.(int)
		}
		if c, ok := d.GetOk("availability_zone_node2"); ok {
			cvoDetails.HAParams.AvailabilityZoneNode2 = c.(int)
		}
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

	if _, ok := d.GetOk("azure_encryption_parameters"); ok {
		if !strings.HasPrefix(cvoDetails.AzureEncryptionParameters.UserAssignedIdentity, "/subscriptions") {
			cvoDetails.AzureEncryptionParameters.UserAssignedIdentity = fmt.Sprintf("/%s/providers/Microsoft.ManagedIdentity/userAssignedIdentities/%s",
				resourceGroupPath, cvoDetails.AzureEncryptionParameters.UserAssignedIdentity)
		}
	}

	if c, ok := d.GetOk("storage_account_network_access"); ok {
		cvoDetails.StorageAccountNetworkAccess = c.(string)
	}

	if client.GetSimulator() {
		log.Print("In simulator env...")
		vnetFormat = "%s/%s"
	}

	vnet := fmt.Sprintf(vnetFormat, resourceGroupPath, cvoDetails.VnetForInternal)
	cvoDetails.VnetID = vnet
	cvoDetails.SubnetID = fmt.Sprintf("%s/subnets/%s", vnet, d.Get("subnet_id").(string))

	res, err := client.createCVOAzure(cvoDetails, clientID)
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(res.PublicID)
	d.Set("svm_name", res.SvmName)
	log.Printf("Created cvo: %v", res)

	// Add SVMs on Azure CVO (supports SN and HA)
	if c, ok := d.GetOk("svm"); ok {
		svms := c.(*schema.Set)
		svmList := expandAzureSVMs(svms)
		for _, svm := range svmList {
			if err := client.addSVMtoCVOAzure(res.PublicID, clientID, svm.SvmName, d.Get("is_ha").(bool), svm.RootVolumeAggregate); err != nil {
				log.Printf("Error adding SVM %v: %v", svm.SvmName, err)
				return err
			}
		}
	}

	return resourceCVOAzureRead(d, meta)
}

func resourceCVOAzureRead(d *schema.ResourceData, meta interface{}) error {
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
	return nil
}

func resourceCVOAzureDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting CVO: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	id := d.Id()
	isHA := d.Get("is_ha").(bool)

	deleteErr := client.deleteCVOAzure(id, isHA, clientID)
	if deleteErr != nil {
		log.Print("Error deleting cvo")
		return deleteErr
	}

	return nil
}

func resourceCVOAzureUpdate(d *schema.ResourceData, meta interface{}) error {
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
		if svmNewName.(string) != "" {
			respErr := client.updateCVOSVMName(d, clientID, svmName.(string), svmNewName.(string), true, "")
			if respErr != nil {
				return respErr
			}
		} else {
			d.Set("svm_name", svmName.(string))
			log.Print("svm_name is set empty. Keep current svm_name. No change.")
		}
	}

	// check if license_type and instance type are changed
	if d.HasChange("instance_type") || d.HasChange("license_type") {
		respErr := updateCVOLicenseInstanceType(d, meta, clientID, true, "")
		if respErr != nil {
			return respErr
		}
	}

	// check if svm list changes (supports SN and HA)
	if d.HasChange("svm") {
		respErr := updateCVOSVMAzure(d, client, clientID)
		if respErr != nil {
			return respErr
		}
	}

	// check if tier_level is changed
	if d.HasChange("tier_level") && d.Get("capacity_tier").(string) == "Blob" {
		respErr := updateCVOTierLevel(d, meta, clientID, true, "")
		if respErr != nil {
			return respErr
		}
	}

	// check if aws_tag has changes
	if d.HasChange("azure_tag") {
		respErr := updateCVOUserTags(d, meta, "azure_tag", clientID, true, "")
		if respErr != nil {
			return respErr
		}
		return resourceCVOAzureRead(d, meta)
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

func resourceCVOAzureCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	respErr := checkUserTagDiff(diff, "azure_tag", "tag_key")
	if respErr != nil {
		return respErr
	}
	// Disallow root_volume_aggregate on creation (aggregates not available yet)
	if diff.Id() == "" {
		if svmSet, ok := diff.GetOk("svm"); ok {
			svms := svmSet.(*schema.Set)
			for _, v := range svms.List() {
				svm := v.(map[string]interface{})
				if rootVol, ok := svm["root_volume_aggregate"].(string); ok && rootVol != "" {
					svmName := svm["svm_name"].(string)
					return fmt.Errorf("root_volume_aggregate cannot be specified for SVM '%s' during CVO creation. Aggregates do not exist at creation time. You can only specify this parameter when adding SVMs to an existing CVO", svmName)
				}
			}
		}
	}

	// Prevent modifying root_volume_aggregate for existing SVMs
	if diff.HasChange("svm") {
		oldRaw, newRaw := diff.GetChange("svm")
		oldSet := oldRaw.(*schema.Set)
		newSet := newRaw.(*schema.Set)

		oldMap := make(map[string]string)
		for _, v := range oldSet.List() {
			svm := v.(map[string]interface{})
			name := svm["svm_name"].(string)
			val := ""
			if x, ok := svm["root_volume_aggregate"].(string); ok {
				val = x
			}
			oldMap[name] = val
		}
		newMap := make(map[string]string)
		for _, v := range newSet.List() {
			svm := v.(map[string]interface{})
			name := svm["svm_name"].(string)
			val := ""
			if x, ok := svm["root_volume_aggregate"].(string); ok {
				val = x
			}
			newMap[name] = val
		}

		for name, oldVal := range oldMap {
			if newVal, exists := newMap[name]; exists {
				if oldVal != newVal {
					return fmt.Errorf("cannot modify root_volume_aggregate for existing SVM '%s' (from '%s' to '%s'). The root_volume_aggregate is immutable after SVM creation. To change it, you must delete and recreate the entire CVO resource", name, oldVal, newVal)
				}
			}
		}
	}
	return nil
}

func resourceCVOAzureExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	clientID := d.Get("client_id").(string)

	resID, err := client.getCVOAzure(id, clientID)
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

func expendEncryptionParameters(azEncryptParameterList *schema.Set) azureEncryptionParameters {
	var params azureEncryptionParameters
	for _, v := range azEncryptParameterList.List() {
		paramSet := v.(map[string]interface{})
		params.Key = paramSet["key"].(string)
		params.VaultName = paramSet["vault_name"].(string)
		params.UserAssignedIdentity = paramSet["user_assigned_identity"].(string)
	}
	return params
}
