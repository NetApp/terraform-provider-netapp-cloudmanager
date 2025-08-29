package cloudmanager

import (
	"fmt"
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
			State: resourceCVOSnapMirrorImport,
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
			"deployment_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Restricted"}, false),
				Default:      "Standard",
			},
			"connector_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delete_destination_volume": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Set to true to delete the destination volume when the snapmirror relationship is destroyed",
			},
		},
	}
}

func resourceCVOSnapMirrorCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating SnapMirror: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	snapMirror := snapMirrorRequest{}

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d, clientID, isSaas, connectorIP)
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
	if s, ok := d.GetOk("iops"); ok {
		snapMirror.ReplicationVolume.Iops = s.(int)
	}
	if s, ok := d.GetOk("throughput"); ok {
		snapMirror.ReplicationVolume.Throughput = s.(int)
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

	res, err := client.buildSnapMirrorCreate(snapMirror, sourceWEInfo.WorkingEnvironmentType, destWEInfo.WorkingEnvironmentType, clientID, isSaas, connectorIP)
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
	clientID := d.Get("client_id").(string)

	snapMirror := snapMirrorRequest{}

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d, clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("Cannot find working environment")
		return err
	}

	snapMirror.ReplicationRequest.SourceWorkingEnvironmentID = sourceWEInfo.PublicID
	snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID = destWEInfo.PublicID
	snapMirror.ReplicationVolume.SourceVolumeName = d.Get("source_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationVolumeName = d.Get("destination_volume_name").(string)
	_, err = client.getSnapMirror(snapMirror, d.Id(), clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("Error getting SnapMirror")
		return err
	}

	return nil
}

func resourceCVOSnapMirrorDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting SnapMirror: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	snapMirror := snapMirrorRequest{}

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return err
	}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d, clientID, isSaas, connectorIP)
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

	err = client.deleteSnapMirror(snapMirror, clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("Error deleting SnapMirror")
		return err
	}

	// If delete_destination_volume is enabled, also delete the destination volume
	if d.Get("delete_destination_volume").(bool) {
		log.Printf("Deleting destination volume: %s", snapMirror.ReplicationVolume.DestinationVolumeName)

		// Prepare volume request for deletion using the same logic as volume resource
		volume := volumeRequest{}
		volume.WorkingEnvironmentID = destWEInfo.PublicID
		volume.WorkingEnvironmentType = destWEInfo.WorkingEnvironmentType
		volume.SvmName = snapMirror.ReplicationVolume.DestinationSvmName
		volume.Name = snapMirror.ReplicationVolume.DestinationVolumeName

		err = client.deleteVolume(volume, clientID, isSaas, connectorIP)
		if err != nil {
			log.Printf("Error deleting destination volume: %s", snapMirror.ReplicationVolume.DestinationVolumeName)
			return err
		}
		log.Printf("Successfully deleted destination volume: %s", snapMirror.ReplicationVolume.DestinationVolumeName)
	}

	return nil
}

func resourceCVOSnapMirrorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of SnapMirror: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)

	snapMirror := snapMirrorRequest{}

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return false, err
	}

	sourceWEInfo, destWEInfo, err := client.getWorkingEnvironmentDetailForSnapMirror(d, clientID, isSaas, connectorIP)
	if err != nil {
		log.Print("Cannot find working environment")
		return false, err
	}

	snapMirror.ReplicationRequest.SourceWorkingEnvironmentID = sourceWEInfo.PublicID
	snapMirror.ReplicationRequest.DestinationWorkingEnvironmentID = destWEInfo.PublicID
	snapMirror.ReplicationVolume.SourceVolumeName = d.Get("source_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationVolumeName = d.Get("destination_volume_name").(string)
	snapMirror.ReplicationVolume.DestinationSvmName = d.Get("destination_svm_name").(string)
	res, err := client.getSnapMirror(snapMirror, d.Id(), clientID, isSaas, connectorIP)
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

func resourceCVOSnapMirrorImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Printf("Importing SnapMirror with ID: %s", d.Id())

	client := meta.(*Client)
	importID := d.Id()

	// Parse the import ID - expect different formats based on deployment mode
	// Standard mode: deployment_mode,client_id,destination_volume_name
	// Restricted mode: deployment_mode,client_id,destination_volume_name,tenant_id,connector_ip
	parts := strings.Split(importID, ",")

	if len(parts) < 3 || (parts[0] != "Standard" && parts[0] != "Restricted") {
		return nil, fmt.Errorf("invalid import ID format. Expected: deployment_mode,client_id,destination_volume_name or deployment_mode,client_id,destination_volume_name,tenant_id,connector_ip for Restricted mode, got: %s", importID)
	}

	deploymentMode := parts[0]
	clientID := parts[1]
	destinationVolumeName := parts[2]

	// Set deployment mode
	d.Set("deployment_mode", deploymentMode)
	d.Set("client_id", clientID)

	// Handle Restricted mode additional parameters
	if deploymentMode == "Restricted" {
		if len(parts) != 5 {
			return nil, fmt.Errorf("invalid import ID format for Restricted mode. Expected: deployment_mode,client_id,destination_volume_name,tenant_id,connector_ip, got: %s", importID)
		}
		d.Set("tenant_id", parts[3])
		d.Set("connector_ip", parts[4])
	} else if deploymentMode == "Standard" {
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid import ID format for Standard mode. Expected: deployment_mode,client_id,destination_volume_name, got: %s", importID)
		}
	}

	// Check deployment mode
	isSaas, connectorIP, err := client.checkDeploymentMode(d, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check deployment mode during import: %v", err)
	}

	// Try to find the snapmirror relationship by destination volume name
	relationship, err := client.findSnapMirrorByDestinationVolume(destinationVolumeName, clientID, isSaas, connectorIP)
	if err != nil {
		return nil, fmt.Errorf("failed to find snapmirror relationship during import: %v", err)
	}

	// Set the ID and populate the required fields
	d.SetId(destinationVolumeName)
	d.Set("client_id", clientID)
	d.Set("source_working_environment_id", relationship.Source.WorkingEnvironmentID)
	d.Set("destination_working_environment_id", relationship.Destination.WorkingEnvironmentID)
	d.Set("source_volume_name", relationship.Source.VolumeName)
	d.Set("destination_volume_name", relationship.Destination.VolumeName)
	d.Set("source_svm_name", relationship.Source.SvmName)
	d.Set("destination_svm_name", relationship.Destination.SvmName)
	d.Set("policy", relationship.Policy)
	d.Set("schedule", relationship.Schedule)
	d.Set("max_transfer_rate", relationship.MaxTransferRate.Size)

	if _, ok := d.GetOk("delete_destination_volume"); !ok {
		d.Set("delete_destination_volume", false)
	}

	// Set optional fields if they exist, otherwise set appropriate defaults
	if relationship.Destination.AggregateName != "" {
		d.Set("destination_aggregate_name", relationship.Destination.AggregateName)
	}
	if relationship.Destination.ProviderVolumeType != "" {
		d.Set("provider_volume_type", relationship.Destination.ProviderVolumeType)
	}
	if relationship.Destination.CapacityTier != "" && relationship.Destination.CapacityTier != "none" {
		d.Set("capacity_tier", relationship.Destination.CapacityTier)
	} else {
		d.Set("capacity_tier", "none")
	}

	// Try to get additional volume information from the destination working environment
	// This is needed because the status API doesn't include all volume details
	if relationship.Destination.WorkingEnvironmentID != "" {
		volumeRequest := volumeRequest{
			WorkingEnvironmentID: relationship.Destination.WorkingEnvironmentID,
			Name:                 relationship.Destination.VolumeName,
		}

		volumes, err := client.getVolume(volumeRequest, clientID, isSaas, connectorIP)
		if err == nil && len(volumes) > 0 {
			volume := volumes[0]
			if volume.AggregateName != "" {
				d.Set("destination_aggregate_name", volume.AggregateName)
			}
			if volume.ProviderVolumeType != "" {
				d.Set("provider_volume_type", volume.ProviderVolumeType)
			}
			if volume.CapacityTier != "" {
				d.Set("capacity_tier", volume.CapacityTier)
			}
		}
	}

	return []*schema.ResourceData{d}, nil
}
