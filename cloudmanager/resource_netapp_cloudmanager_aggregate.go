package cloudmanager

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAggregate() *schema.Resource {
	return &schema.Resource{
		Create:        resourceAggregateCreate,
		Read:          resourceAggregateRead,
		Delete:        resourceAggregateDelete,
		Exists:        resourceAggregateExists,
		Update:        resourceAggregateUpdate,
		CustomizeDiff: resourceAggregateCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceAggregateImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"working_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"working_environment_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"number_of_disks": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disk_size_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"disk_size_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"GB", "TB"}, true),
			},
			"home_node": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"provider_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "gp2",
				ForceNew: true,
			},
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "S3", "Blob", "cloudStorage"}, true),
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
			"increase_capacity_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Additional capacity to add to the aggregate (only available during updates)",
				ForceNew:    false,
			},
			"increase_capacity_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Byte", "KB", "MB", "GB", "TB"}, true),
				Description:  "Unit for the additional capacity (Byte, KB, MB, GB, or TB) (only available during updates)",
				ForceNew:     false,
			},
			"initial_ev_aggregate_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Initial size for EBS Elastic Volumes aggregate. This enables the aggregate to support capacity expansion.",
			},
			"initial_ev_aggregate_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Byte", "KB", "MB", "GB", "TB"}, true),
				Description:  "Unit for initial EBS Elastic Volumes aggregate size",
			},
			"total_capacity_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Total capacity of the aggregate",
			},
			"total_capacity_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unit of the total capacity",
			},
			"available_capacity_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Available capacity of the aggregate",
			},
			"available_capacity_unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unit of the available capacity",
			},
		},
	}
}

func resourceAggregateCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating Aggregate: %#v", d)

	client := meta.(*Client)

	clientID := d.Get("client_id").(string)
	aggregate := createAggregateRequest{}

	// Check deployment mode
	isSaaS, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	workingEnv, err := client.getWorkingEnvironmentDetail(d, clientID, isSaaS, connectorIP)

	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	aggregate.WorkingEnvironmentID = workingEnv.PublicID

	// Validate that capacity increase fields are not used during creation
	if capacitySize, ok := d.GetOk("increase_capacity_size"); ok && capacitySize.(int) > 0 {
		return fmt.Errorf("increase_capacity_size can only be used during aggregate updates, not during creation")
	}
	if capacityUnit, ok := d.GetOk("increase_capacity_unit"); ok && capacityUnit.(string) != "" {
		return fmt.Errorf("increase_capacity_unit can only be used during aggregate updates, not during creation")
	}

	aggregate.Name = d.Get("name").(string)

	if a, ok := d.GetOk("number_of_disks"); ok {
		aggregate.NumberOfDisks, _ = a.(int)
	}
	if a, ok := d.GetOk("disk_size_size"); ok {
		aggregate.DiskSize.Size, _ = a.(int)
	}
	if a, ok := d.GetOk("disk_size_unit"); ok {
		aggregate.DiskSize.Unit = a.(string)
	}

	if a, ok := d.GetOk("home_node"); ok {
		aggregate.HomeNode = a.(string)
	}
	if a, ok := d.GetOk("provider_volume_type"); ok {
		aggregate.ProviderVolumeType = a.(string)
		if aggregate.ProviderVolumeType == "io1" {
			if a, ok := d.GetOk("iops"); ok {
				aggregate.Iops = a.(int)
			} else {
				log.Printf("CreateAggregate: provider_volume_type is io1, but iops is not configured.")
			}
		}
		if aggregate.ProviderVolumeType == "gp3" {
			if a, ok := d.GetOk("iops"); ok {
				aggregate.Iops = a.(int)
			} else {
				log.Printf("CreateAggregate: provider_volume_type is gp3, but iops is not configured.")
			}
			if a, ok := d.GetOk("throughput"); ok {
				aggregate.Throughput = a.(int)
			} else {
				log.Printf("CreateAggregate: provider_volume_type is gp3, but throughput is not configured.")
			}
		}
	}
	if a, ok := d.GetOk("capacity_tier"); ok {
		if a.(string) != "NONE" {
			aggregate.CapacityTier = a.(string)
		}
	} else if workingEnv.CloudProviderName == "Amazon" {
		aggregate.CapacityTier = "S3"
	} else if workingEnv.CloudProviderName == "Azure" {
		aggregate.CapacityTier = "Blob"
	} else if workingEnv.CloudProviderName == "GCP" {
		aggregate.CapacityTier = "cloudStorage"
	}

	// Handle initial EV aggregate size for AWS EBS Elastic Volumes
	if initialSize, ok := d.GetOk("initial_ev_aggregate_size"); ok {
		if workingEnv.CloudProviderName != "Amazon" {
			return fmt.Errorf("initial_ev_aggregate_size is only supported for Amazon Web Services (AWS) environments")
		}

		aggregate.InitialEvAggregateSize.Size = initialSize.(int)

		if initialUnit, ok := d.GetOk("initial_ev_aggregate_unit"); ok {
			aggregate.InitialEvAggregateSize.Unit = initialUnit.(string)
		}

		log.Printf("Setting initial EV aggregate size: %d %s", aggregate.InitialEvAggregateSize.Size, aggregate.InitialEvAggregateSize.Unit)
	}

	res, err := client.createAggregate(&aggregate, clientID, isSaaS, connectorIP)
	if err != nil {
		log.Print("Error creating aggregate")
		return err
	}

	d.SetId(res.Name)

	log.Printf("Created aggregate: %v", res)

	return resourceAggregateRead(d, meta)
}

// read the specific aggregate with working environment Id and aggregate name
func resourceAggregateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading Aggregate: %#v", d)
	client := meta.(*Client)

	clientID := d.Get("client_id").(string)
	aggregate := aggregateRequest{}

	// Check deployment mode
	isSaaS, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	workingEnv, err := client.getWorkingEnvironmentDetail(d, clientID, isSaaS, connectorIP)

	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	aggregate.WorkingEnvironmentID = workingEnv.PublicID

	importing := false
	if strings.Contains(d.Id(), ",") {
		importing = true
	}

	id := d.Id()
	if importing {
		// During import, use the name from the import ID
		id = d.Get("name").(string)
	}

	aggr, err := client.getAggregate(aggregate, id, workingEnv.WorkingEnvironmentType, clientID, isSaaS, connectorIP)
	if err != nil {
		log.Printf("Error getting aggregate. id = %v", id)
		return err
	}

	if importing {
		// During import, set the ID to the aggregate name
		d.SetId(aggr.Name)
		d.Set("number_of_disks", len(aggr.Disks))
		d.Set("working_environment_name", workingEnv.Name)
		d.Set("disk_size_size", aggr.ProviderVolumes[0].Size.Size)
		d.Set("disk_size_unit", aggr.ProviderVolumes[0].Size.Unit)
		d.Set("provider_volume_type", aggr.ProviderVolumes[0].DiskType)
	}

	if aggr.Name != d.Get("name").(string) {
		return fmt.Errorf("expected aggregate name %v, Response could not find", aggr.Name)
	}

	// Set computed capacity values
	d.Set("total_capacity_size", aggr.TotalCapacity.Size)
	d.Set("total_capacity_unit", aggr.TotalCapacity.Unit)
	d.Set("available_capacity_size", aggr.AvailableCapacity.Size)
	d.Set("available_capacity_unit", aggr.AvailableCapacity.Unit)

	return nil
}

func resourceAggregateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting Aggregate: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	request := deleteAggregateRequest{}

	// Check deployment mode
	isSaaS, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	workingEnvDetail, err := client.getWorkingEnvironmentDetail(d, clientID, isSaaS, connectorIP)

	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	request.WorkingEnvironmentID = workingEnvDetail.PublicID

	request.Name = d.Get("name").(string)

	deleteErr := client.deleteAggregate(request, clientID, isSaaS, connectorIP)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceAggregateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating Aggregate: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	request := updateAggregateRequest{}

	// Check deployment mode
	isSaaS, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	workingEnvDetail, err := client.getWorkingEnvironmentDetail(d, clientID, isSaaS, connectorIP)

	if err != nil {
		return fmt.Errorf("cannot find working environment")
	}
	request.WorkingEnvironmentID = workingEnvDetail.PublicID

	aggregate := aggregateRequest{}
	aggregate.WorkingEnvironmentID = workingEnvDetail.PublicID
	id := d.Id()
	aggr, err := client.getAggregate(aggregate, id, workingEnvDetail.WorkingEnvironmentType, clientID, isSaaS, connectorIP)
	if err != nil {
		log.Printf("Error getting aggregate. id = %v", id)
		return err
	}
	currentNumber := len(aggr.Disks)
	request.Name = d.Get("name").(string)

	log.Printf("Current number of disks: %v", currentNumber)

	// Handle capacity increase
	if d.HasChange("increase_capacity_size") || d.HasChange("increase_capacity_unit") {
		capacitySize := d.Get("increase_capacity_size").(int)
		capacityUnit := d.Get("increase_capacity_unit").(string)

		if capacitySize > 0 {
			// Check if this is an Amazon (AWS) environment first
			if workingEnvDetail.CloudProviderName != "Amazon" {
				return fmt.Errorf("aggregate capacity increase is only supported for Amazon Web Services (AWS) environments, current environment is %s", workingEnvDetail.CloudProviderName)
			}

			increaseRequest := increaseAggregateCapacityRequest{
				WorkingEnvironmentID: workingEnvDetail.PublicID,
				AggregateName:        request.Name,
				CapacityToAdd: diskSize{
					Size: capacitySize,
					Unit: capacityUnit,
				},
			}

			err := client.increaseAggregateCapacity(increaseRequest, clientID, isSaaS, connectorIP)
			if err != nil {
				return fmt.Errorf("failed to increase aggregate capacity: %v", err)
			}

			log.Printf("Successfully increased aggregate capacity by %d %s", capacitySize, capacityUnit)

			// As long as config has increase fields, every run acts as an update by resetting the increase fields in state file.
			// only if increase fields are not used in the configuration, the resource can achieve idempotency
			d.Set("increase_capacity_size", 0)
			d.Set("increase_capacity_unit", "")
		}
	}

	// Handle disk count update
	if d.HasChange("number_of_disks") {
		expectNumber := d.Get("number_of_disks")
		log.Printf("Expect number of disks: %v", expectNumber.(int))
		if expectNumber.(int) > currentNumber {
			request.NumberOfDisks = expectNumber.(int) - currentNumber
		} else {
			d.Set("number_of_disks", currentNumber)
			return fmt.Errorf("aggregate: number_of_disks cannot be reduced")
		}

		updateErr := client.updateAggregate(request, clientID, isSaaS, connectorIP)
		if updateErr != nil {
			d.Set("number_of_disks", currentNumber)
			return updateErr
		}
	}

	log.Printf("Updated aggregate: %v", request.Name)

	return resourceAggregateRead(d, meta)
}

func resourceAggregateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of Aggregate: %#v", d)
	client := meta.(*Client)

	clientID := d.Get("client_id").(string)
	aggregate := aggregateRequest{}

	isSaaS, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return false, err
	}

	workingEnv, err := client.getWorkingEnvironmentDetail(d, clientID, isSaaS, connectorIP)

	if err != nil {
		return false, fmt.Errorf("cannot find working environment")
	}
	aggregate.WorkingEnvironmentID = workingEnv.PublicID

	importing := false
	if strings.Contains(d.Id(), ",") {
		importing = true
	}

	id := d.Id()
	if importing {
		// During import, use the name from the import ID
		id = d.Get("name").(string)
	}

	res, err := client.getAggregate(aggregate, id, workingEnv.WorkingEnvironmentType, clientID, isSaaS, connectorIP)
	if err != nil {
		log.Print("Error getting aggregate")
		d.SetId("")
		return false, err
	}

	if importing {
		// During import, compare with name
		if res.Name != d.Get("name").(string) {
			d.SetId("")
			return false, nil
		}
	} else {
		// Normal operation, compare with ID
		if res.Name != id {
			d.SetId("")
			return false, nil
		}
	}

	return true, nil
}

func resourceAggregateCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {
	// Validate disk_size_size and disk_size_unit are provided together
	diskSizeSize := diff.Get("disk_size_size")
	diskSizeUnit := diff.Get("disk_size_unit")

	hasDiskSizeSize := diskSizeSize != nil && diskSizeSize.(int) > 0
	hasDiskSizeUnit := diskSizeUnit != nil && diskSizeUnit.(string) != ""

	if hasDiskSizeSize && !hasDiskSizeUnit {
		return fmt.Errorf("disk_size_unit is required when disk_size_size is specified")
	}

	if hasDiskSizeUnit && !hasDiskSizeSize {
		return fmt.Errorf("disk_size_size is required when disk_size_unit is specified")
	}

	// Validate initial_ev_aggregate_size and initial_ev_aggregate_unit are provided together
	initialEvSize := diff.Get("initial_ev_aggregate_size")
	initialEvUnit := diff.Get("initial_ev_aggregate_unit")

	hasInitialEvSize := initialEvSize != nil && initialEvSize.(int) > 0
	hasInitialEvUnit := initialEvUnit != nil && initialEvUnit.(string) != ""

	if hasInitialEvSize && !hasInitialEvUnit {
		return fmt.Errorf("initial_ev_aggregate_unit is required when initial_ev_aggregate_size is specified")
	}

	if hasInitialEvUnit && !hasInitialEvSize {
		return fmt.Errorf("initial_ev_aggregate_size is required when initial_ev_aggregate_unit is specified")
	}

	// Validate increase_capacity_size and increase_capacity_unit are provided together
	increaseSize := diff.Get("increase_capacity_size")
	increaseUnit := diff.Get("increase_capacity_unit")

	hasIncreaseSize := increaseSize != nil && increaseSize.(int) > 0
	hasIncreaseUnit := increaseUnit != nil && increaseUnit.(string) != ""

	if hasIncreaseSize && !hasIncreaseUnit {
		return fmt.Errorf("increase_capacity_unit is required when increase_capacity_size is specified")
	}

	if hasIncreaseUnit && !hasIncreaseSize {
		return fmt.Errorf("increase_capacity_size is required when increase_capacity_unit is specified")
	}

	return nil
}

func resourceAggregateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if parts[0] != "Standard" && parts[0] != "Restricted" {
		return []*schema.ResourceData{}, fmt.Errorf("wrong option for deployment_mode: %s, options for deployment_mode are 'Standard' and 'Restricted'", parts[0])
	}

	if parts[0] == "Standard" && len(parts) != 4 {
		return []*schema.ResourceData{}, fmt.Errorf("wrong format of resource: %s. Please input in the format 'deployment_mode,client_id,working_environment_name,name'", d.Id())
	}

	if parts[0] == "Restricted" && len(parts) != 6 {
		return []*schema.ResourceData{}, fmt.Errorf("wrong format of resource: %s. Please input in the format 'deployment_mode,client_id,working_environment_name,name,tenant_id,connector_ip'", d.Id())
	}

	d.Set("deployment_mode", parts[0])
	d.Set("client_id", parts[1])
	d.Set("working_environment_name", parts[2])
	d.Set("name", parts[3])
	if parts[0] == "Restricted" {
		d.Set("tenant_id", parts[4])
		d.Set("connector_ip", parts[5])
	}

	return []*schema.ResourceData{d}, nil

}
