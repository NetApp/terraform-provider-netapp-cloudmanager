package cloudmanager

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceCVOOnPrem() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOOnPremCreate,
		Read:   resourceCVOOnPremRead,
		Delete: resourceCVOOnPremDelete,
		Exists: resourceCVOOnPremExists,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_password": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON_PREM", "AZURE", "AWS", "SOFTLAYER", "GOOGLE", "CLOUD_TIERING"}, false),
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCVOOnPremCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating CVO: %#v", d)

	client := meta.(*Client)

	cvoDetails := createCVOOnPremDetails{}

	cvoDetails.Name = d.Get("name").(string)
	cvoDetails.ClusterAddress = d.Get("cluster_address").(string)
	cvoDetails.ClusterPassword = d.Get("cluster_password").(string)
	cvoDetails.ClusterUserName = d.Get("cluster_user_name").(string)
	cvoDetails.WorkspaceID = d.Get("workspace_id").(string)
	cvoDetails.Location = d.Get("location").(string)
	clientID := d.Get("client_id").(string)

	res, err := client.createCVOOnPrem(cvoDetails, clientID)
	if err != nil {
		log.Print("Error creating instance: ", err)
		return err
	}

	d.SetId(res.PublicID)
	// d.Set("svm_name", res.SvmName)

	log.Printf("Created cvo: %v", res)

	return resourceCVOOnPremRead(d, meta)
}

func resourceCVOOnPremRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	clientID := d.Get("client_id").(string)

	_, err := client.getCVOOnPremByID(id, clientID)
	if err != nil {
		log.Print("Error reading cvo onprem: ", err)
		return err
	}
	return nil
}

func resourceCVOOnPremDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting CVO: %#v", d)

	client := meta.(*Client)

	id := d.Id()
	clientID := d.Get("client_id").(string)

	deleteErr := client.deleteCVOOnPrem(id, clientID)
	if deleteErr != nil {
		log.Print("Error deleting cvo: ", deleteErr)
		return deleteErr
	}

	return nil
}

func resourceCVOOnPremExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of CVO: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	clientID := d.Get("client_id").(string)

	resID, err := client.getCVOOnPrem(id, clientID)
	if err != nil {
		log.Print("Error getting cvo: ", err)
		return false, err
	}

	if resID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
