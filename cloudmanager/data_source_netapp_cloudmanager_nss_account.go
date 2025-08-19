package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCVONssAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCVONssAccountRead,

		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceCVONssAccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Getting nss account: %s", d.Get("username").(string))
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	res, err := client.getNssAccount(d.Get("username").(string), clientID)
	if err != nil {
		log.Printf("Error getting nss account: %s", d.Get("username").(string))
		return err
	}
	if res == nil {
		return fmt.Errorf("Failed to find account: %s", d.Get("username"))
	}
	d.Set("username", res["nssUserName"])
	d.SetId(res["publicId"].(string))
	return nil
}
