package cloudmanager

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAWSFSX() *schema.Resource {
	return &schema.Resource{
		Create: resourceAWSFSXCreate,
		Read:   resourceAWSFSXRead,
		Delete: resourceAWSFSXDelete,
		Exists: resourceAWSFSXExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		CustomizeDiff: resourceAWSFSXCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aws_credentials_name": {
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
				Required: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_capacity_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1,
			},
			"storage_capacity_size_unit": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "TB",
				ValidateFunc: validation.StringInSlice([]string{"GiB", "TiB"}, false),
			},
			"fsx_admin_password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"throughput_capacity": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{512, 1024, 2048}),
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
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
			"primary_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secondary_subnet_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"minimum_ssd_iops": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(401),
			},
			"endpoint_ip_address_range": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"import_file_system": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAWSFSXCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating AWS FSX: %#v", d)

	client := meta.(*Client)

	fsxDetails := createAWSFSXDetails{}

	fsxDetails.Name = d.Get("name").(string)
	fsxDetails.AWSCredentials = d.Get("aws_credentials_name").(string)
	fsxDetails.Region = d.Get("region").(string)
	fsxDetails.WorkspaceID = d.Get("workspace_id").(string)
	fsxDetails.ThroughputCapacity = d.Get("throughput_capacity").(int)
	fsxDetails.FSXAdminPassword = d.Get("fsx_admin_password").(string)
	fsxDetails.TenantID = d.Get("tenant_id").(string)

	if d.Get("import_file_system").(bool) == true {
		if d.Get("file_system_id").(string) == "" {
			return fmt.Errorf("need file_system_id when importing file system")
		}
		fileSystemID := d.Get("file_system_id").(string)
		fsxID, err := client.importAWSFSX(fsxDetails, fileSystemID)
		if err != nil {
			log.Print("Error importing AWS FSX")
			return err
		}

		d.SetId(fsxID)

		log.Printf("Created AWS FSX: %v", fsxID)

		return resourceAWSFSXRead(d, meta)
	}

	addNameTag := true
	if c, ok := d.GetOk("tags"); ok {
		tags := c.(*schema.Set)
		if tags.Len() > 0 {
			fsxDetails.AwsFSXTags = expandfsxTags(tags)
			if hasNameTag(fsxDetails.AwsFSXTags) {
				addNameTag = false
			}
		}
	}
	if addNameTag {
		// add name tag
		fsxTag := fsxTags{}
		fsxTag.TagKey = "name"
		fsxTag.TagValue = fsxDetails.Name
		fsxDetails.AwsFSXTags = append(fsxDetails.AwsFSXTags, fsxTag)
	}
	log.Print("fsxdetails: ", fsxDetails)
	fsxDetails.StorageCapacity.Size = d.Get("storage_capacity_size").(int)
	fsxDetails.StorageCapacity.Unit = d.Get("storage_capacity_size_unit").(string)

	securityGroupIds := d.Get("security_group_ids")
	for _, securityGroupID := range securityGroupIds.([]interface{}) {
		fsxDetails.SecurityGroupIds = append(fsxDetails.SecurityGroupIds, securityGroupID.(string))
	}

	if c, ok := d.GetOk("kms_key_id"); ok {
		fsxDetails.KmsKeyID = c.(string)
	}

	if c, ok := d.GetOk("minimum_ssd_iops"); ok {
		fsxDetails.MinimumSsdIops = c.(int)
	}

	if c, ok := d.GetOk("endpoint_ip_address_range"); ok {
		fsxDetails.EndpointIPAddressRange = c.(string)
	}

	fsxDetails.PrimarySubnetID = d.Get("primary_subnet_id").(string)
	fsxDetails.SecondarySubnetID = d.Get("secondary_subnet_id").(string)

	routeTableIds := d.Get("route_table_ids")
	for _, routeTableID := range routeTableIds.([]interface{}) {
		fsxDetails.RouteTableIds = append(fsxDetails.RouteTableIds, routeTableID.(string))
	}

	res, err := client.createAWSFSX(fsxDetails)
	if err != nil {
		log.Print("Error creating AWS FSX")
		return err
	}

	d.SetId(res.ID)

	log.Printf("Created AWS FSX: %v", res)

	return resourceAWSFSXRead(d, meta)
}

func resourceAWSFSXRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading AWS FSX: %#v", d)
	client := meta.(*Client)

	id := d.Id()

	tenantID := d.Get("tenant_id").(string)

	_, err := client.getAWSFSX(id, tenantID, true, "")
	if err != nil {
		log.Print("Error getting AWS FSX")
		return err
	}

	return nil
}

func resourceAWSFSXDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting AWS FSX: %#v", d)

	client := meta.(*Client)

	id := d.Id()

	tenantID := d.Get("tenant_id").(string)

	deleteErr := client.deleteAWSFSX(id, tenantID)
	if deleteErr != nil {
		log.Print("Error deleting AWS FSX")
		return deleteErr
	}

	return nil
}

func resourceAWSFSXCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	_ = ctx
	respErr := checkUserTagDiff(diff, "tags", "tag_key")
	if respErr != nil {
		return respErr
	}
	return nil
}

func resourceAWSFSXExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of AWS FSX: %#v", d)
	client := meta.(*Client)

	id := d.Id()

	tenantID := d.Get("tenant_id").(string)

	resID, err := client.getAWSFSX(id, tenantID, true, "")
	if err != nil {
		log.Print("Error getting AWS FSX")
		return false, err
	}

	if resID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
