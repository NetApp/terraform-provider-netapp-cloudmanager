package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAWSFSX() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAWSFSXRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAWSFSXRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching aws fsx: %#v", d)

	client := meta.(*Client)

	id := d.Get("id").(string)
	tenantID := d.Get("tenant_id").(string)

	res, err := client.getAWSFSXByID(id, tenantID)
	if err != nil {
		log.Print("Error getting AWS FSX")
		return err
	}

	d.SetId(res.ID)
	err = d.Set("name", res.Name)
	if err != nil {
		return fmt.Errorf("Error setting fsx name: %s", err.Error())
	}
	err = d.Set("status", res.ProviderDetails.Status.Status)
	if err != nil {
		return fmt.Errorf("Error setting fsx status: %s", err.Error())
	}
	err = d.Set("lifecycle_status", res.ProviderDetails.Status.Lifecycle)
	if err != nil {
		return fmt.Errorf("Error setting fsx lifecycle: %s", err.Error())
	}
	err = d.Set("region", res.Region)
	if err != nil {
		return fmt.Errorf("Error setting fsx region: %s", err.Error())
	}
	return nil
}
