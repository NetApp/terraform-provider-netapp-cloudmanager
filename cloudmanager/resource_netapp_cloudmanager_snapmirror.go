package cloudmanager

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCVOSnapMirror() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOSnapMirrorCreate,
		Read:   resourceCVOSnapMirrorRead,
		Delete: resourceCVOSnapMirrorDelete,
		Exists: resourceCVOSnapMirrorExists,
		Update: resourceCVOSnapMirrorUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_working_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_working_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_working_environment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_working_environment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_aggregate_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MirrorAllSnapshots",
			},
			"schedule": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "1hour",
			},
			"max_transfer_rate": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100000,
			},
			"source_svm_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"destination_svm_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_volume_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"destination_volume_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"capacity_tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"S3", "Blob", "cloudStorage", "none"}, false),
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCVOSnapMirrorCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating SnapMirror: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	snapMirror := snapMirrorRequest{}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d)
	if err != nil {
		log.Print("Cannot find working environment")
		return err
	}

	snapMirror.ReplicationRequest.SourceWorkingEnvironmentID = sourceWEInfo.PublicID
	if strings.HasPrefix(destWEInfo.PublicID, "fs-") {
		snapMirror.ReplicationRequest.DestinationFsxID = destWEInfo.PublicID
	} else {
		snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID = destWEInfo.PublicID
	}

	log.Print("PublicIDfsx ", snapMirror.ReplicationRequest.DestinationFsxID)
	log.Print("PublicIDvsa ", snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID)
	snapMirror.ReplicationVolume.SourceVolumeName = d.Get("source_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationVolumeName = d.Get("destination_volume_name").(string)
	snapMirror.ReplicationRequest.PolicyName = d.Get("policy").(string)
	snapMirror.ReplicationRequest.ScheduleName = d.Get("schedule").(string)
	snapMirror.ReplicationRequest.MaxTransferRate = d.Get("max_transfer_rate").(int)

	if s, ok := d.GetOk("destination_aggregate_name"); ok {
		snapMirror.ReplicationVolume.DestinationAggregateName = s.(string)
	}
	if s, ok := d.GetOk("source_svm_name"); ok {
		snapMirror.ReplicationVolume.SourceSvmName = s.(string)
	}
	if s, ok := d.GetOk("destination_svm_name"); ok {
		snapMirror.ReplicationVolume.DestinationSvmName = s.(string)
	}
	if s, ok := d.GetOk("provider_volume_type"); ok {
		snapMirror.ReplicationVolume.DestinationProviderVolumeType = s.(string)
	}
	if s, ok := d.GetOk("capacity_tier"); ok {
		if s.(string) != "none" {
			snapMirror.ReplicationVolume.DestinationCapacityTier = s.(string)
		}
	} else {
		cloudProvider := strings.ToLower(destWEInfo.CloudProviderName)
		if cloudProvider == "aws" {
			snapMirror.ReplicationVolume.DestinationCapacityTier = "S3"
		} else if cloudProvider == "azure" {
			snapMirror.ReplicationVolume.DestinationCapacityTier = "Blob"
		} else if cloudProvider == "gcp" {
			snapMirror.ReplicationVolume.DestinationCapacityTier = "cloudStorage"
		}
	}

	if snapMirror.ReplicationVolume.SourceSvmName == "" {
		snapMirror.ReplicationVolume.SourceSvmName = sourceWEInfo.SvmName
	}
	if snapMirror.ReplicationVolume.DestinationSvmName == "" {
		snapMirror.ReplicationVolume.DestinationSvmName = destWEInfo.SvmName
	}

	res, err := client.buildSnapMirrorCreate(snapMirror, sourceWEInfo.WorkingEnvironmentType, destWEInfo.WorkingEnvironmentType)
	if err != nil {
		log.Print("Error creating SnapMirrorCreate")
		return err
	}

	d.SetId(res.ReplicationVolume.DestinationVolumeName)

	d.Set("source_svm_name", res.ReplicationVolume.SourceSvmName)
	d.Set("destination_svm_name", res.ReplicationVolume.DestinationSvmName)
	d.Set("source_volume_name", res.ReplicationVolume.SourceVolumeName)
	d.Set("destination_volume_name", res.ReplicationVolume.DestinationVolumeName)

	return resourceCVOSnapMirrorRead(d, meta)
}

func resourceCVOSnapMirrorRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching SnapMirror: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)

	snapMirror := snapMirrorRequest{}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d)
	if err != nil {
		log.Print("Cannot find working environment")
		return err
	}

	snapMirror.ReplicationRequest.SourceWorkingEnvironmentID = sourceWEInfo.PublicID
	snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID = destWEInfo.PublicID
	snapMirror.ReplicationVolume.SourceVolumeName = d.Get("source_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationVolumeName = d.Get("destination_volume_name").(string)
	_, err = client.getSnapMirror(snapMirror, d.Id())
	if err != nil {
		log.Print("Error getting SnapMirror")
		return err
	}

	return nil
}

func resourceCVOSnapMirrorDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting SnapMirror: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	snapMirror := snapMirrorRequest{}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d)
	if err != nil {
		log.Print("Cannot find working environment")
		return err
	}

	snapMirror.ReplicationRequest.SourceWorkingEnvironmentID = sourceWEInfo.PublicID
	snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID = destWEInfo.PublicID
	snapMirror.ReplicationVolume.DestinationSvmName = d.Get("destination_svm_name").(string)
	snapMirror.ReplicationVolume.SourceVolumeName = d.Get("source_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationVolumeName = d.Get("destination_volume_name").(string)
	if snapMirror.ReplicationVolume.DestinationVolumeName == "" {
		snapMirror.ReplicationVolume.DestinationVolumeName = snapMirror.ReplicationVolume.SourceVolumeName + "_copy"
	}
	if s, ok := d.GetOk("source_svm_name"); ok {
		snapMirror.ReplicationVolume.SourceSvmName = s.(string)
	}

	err = client.deleteSnapMirror(snapMirror)
	if err != nil {
		log.Print("Error deleting SnapMirror")
		return err
	}
	return nil
}

func resourceCVOSnapMirrorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of SnapMirror: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	snapMirror := snapMirrorRequest{}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d)
	if err != nil {
		log.Print("Cannot find working environment")
		return false, err
	}

	snapMirror.ReplicationRequest.SourceWorkingEnvironmentID = sourceWEInfo.PublicID
	snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID = destWEInfo.PublicID
	snapMirror.ReplicationVolume.SourceVolumeName = d.Get("source_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationVolumeName = d.Get("destination_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationSvmName = d.Get("destination_svm_name").(string)
	res, err := client.getSnapMirror(snapMirror, d.Id())
	if err != nil {
		log.Print("Error getting SnapMirror")
		return false, err
	}

	if res != d.Id() {
		d.SetId("")
		return false, nil
	}
	return true, nil
}

func resourceCVOSnapMirrorUpdate(d *schema.ResourceData, meta interface{}) error {

	return nil
}
