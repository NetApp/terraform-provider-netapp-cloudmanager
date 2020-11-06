package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceAggregate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAggregateCreate,
		Read:   resourceAggregateRead,
		Delete: resourceAggregateDelete,
		Exists: resourceAggregateExists,
		Update: resourceAggregateUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Required: true,
			},
			"disk_size_size": {
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
		},
	}
}

func resourceAggregateCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating Aggregate: %#v", d)

	client := meta.(*Client)

	client.ClientID = d.Get("client_id").(string)
	aggregate := createAggregateRequest{}
	var workingEnv workingEnvironmentInfo

	if a, ok := d.GetOk("working_environment_id"); ok {
		aggregate.WorkingEnvironmentID = a.(string)
		workingEnvDetail, err := client.getWorkingEnvironmentInfo(aggregate.WorkingEnvironmentID)
		if err != nil {
			return fmt.Errorf("Cannot find working environment by working_environment_id %s", a.(string))
		}
		workingEnv = workingEnvDetail
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		log.Printf("Get environment id %v by %v", workingEnvDetail.PublicID, a.(string))
		aggregate.WorkingEnvironmentID = workingEnvDetail.PublicID
		workingEnv = workingEnvDetail
	} else {
		return fmt.Errorf("Cannot find working environment by working_enviroment_id or working_environment_name")
	}

	aggregate.Name = d.Get("name").(string)
	aggregate.NumberOfDisks = d.Get("number_of_disks").(int)

	if a, ok := d.GetOk("disk_size_size"); ok {
		aggregate.DiskSize.Size, ok = a.(int)
	}
	if a, ok := d.GetOk("disk_size_unit"); ok {
		aggregate.DiskSize.Unit = a.(string)
	}

	if a, ok := d.GetOk("home_node"); ok {
		aggregate.HomeNode = a.(string)
	}
	if a, ok := d.GetOk("provider_yolume_type"); ok {
		aggregate.ProviderVolumeType = a.(string)
		if aggregate.ProviderVolumeType == "io1" {
			if a, ok := d.GetOk("iops"); ok {
				aggregate.Iops = a.(int)
			} else {
				log.Printf("CreateAggregate: provider_volume_type is io1, but iops is not configured.")
			}
		}
	}
	if a, ok := d.GetOk("capacity_tier"); ok {
		aggregate.CapacityTier = a.(string)
		if aggregate.CapacityTier == "NONE" {
			aggregate.CapacityTier = ""
		}
	} else if workingEnv.CloudProviderName == "Amazon" {
		aggregate.CapacityTier = "S3"
	} else if workingEnv.CloudProviderName == "Azure" {
		aggregate.CapacityTier = "Blob"
	} else if workingEnv.CloudProviderName == "GCP" {
		aggregate.CapacityTier = "cloudStorage"
	}

	res, err := client.createAggregate(&aggregate)
	if err != nil {
		log.Print("Error creating aggregate")
		return err
	}

	d.SetId(res.Name)

	log.Printf("Created aggregate: %v", res)

	return resourceAggregateRead(d, meta)
}

// read the specific aggregate with working environemnt Id and aggregate name
func resourceAggregateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading Aggregate: %#v", d)
	client := meta.(*Client)

	client.ClientID = d.Get("client_id").(string)

	aggregate := aggregateRequest{}
	aggregate.WorkingEnvironmentID = d.Get("working_environment_id").(string)
	if a, ok := d.GetOk("working_environment_id"); ok {
		aggregate.WorkingEnvironmentID = a.(string)
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		aggregate.WorkingEnvironmentID = workingEnvDetail.PublicID
	} else {
		return fmt.Errorf("Cannot find working environment by working_enviroment_id or working_environment_name")
	}

	id := d.Id()

	aggr, err := client.getAggregate(aggregate, id)
	if err != nil {
		log.Printf("Error getting aggregate. id = %v", id)
		return err
	}

	if aggr.Name != id {
		return fmt.Errorf("Expected aggregate name %v, Response could not find", aggr.Name)
	}

	return nil
}

func resourceAggregateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting Aggregate: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	request := deleteAggregateRequest{}
	if a, ok := d.GetOk("working_environment_id"); ok {
		request.WorkingEnvironmentID = a.(string)
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		request.WorkingEnvironmentID = workingEnvDetail.PublicID
	} else {
		return fmt.Errorf("Cannot find working environment by working_enviroment_id or working_environment_name")
	}

	request.Name = d.Get("name").(string)

	deleteErr := client.deleteAggregate(request)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceAggregateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating Aggregate: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	request := updateAggregateRequest{}

	if a, ok := d.GetOk("working_environment_id"); ok {
		request.WorkingEnvironmentID = a.(string)
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		request.WorkingEnvironmentID = workingEnvDetail.PublicID
	} else {
		return fmt.Errorf("Cannot find working environment by working_enviroment_id or working_environment_name")
	}

	request.Name = d.Get("name").(string)

	if d.HasChange("number_of_disks") {
		currentNumber, expectNumber := d.GetChange("number_of_disks")
		if expectNumber.(int) > currentNumber.(int) {
			request.NumberOfDisks = expectNumber.(int) - currentNumber.(int)
		} else {
			d.Set("number_of_disks", currentNumber)
			return fmt.Errorf("Aggregate: number_of_disks cannot be reduced")
		}
	}
	updateErr := client.updateAggregate(request)
	if updateErr != nil {
		return updateErr
	}

	log.Printf("Updated aggregate; %v", request.Name)

	return resourceAggregateRead(d, meta)
}

func resourceAggregateExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of Aggregate: %#v", d)
	client := meta.(*Client)

	client.ClientID = d.Get("client_id").(string)
	aggregate := aggregateRequest{}
	if a, ok := d.GetOk("working_environment_id"); ok {
		aggregate.WorkingEnvironmentID = a.(string)
	} else if a, ok = d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return false, fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		aggregate.WorkingEnvironmentID = workingEnvDetail.PublicID
	} else {
		return false, fmt.Errorf("Cannot find working environment by working_enviroment_id or working_environment_name")
	}

	id := d.Id()
	res, err := client.getAggregate(aggregate, id)
	if err != nil {
		log.Print("Error getting aggregate")
		d.SetId("")
		return false, err
	}

	if res.Name != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
