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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceCVONssAccountRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Getting nss account: %s", d.Get("name").(string))
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	res, err := client.getNssAccount(d.Get("name").(string))
	if err != nil {
		log.Printf("Error getting nss account: %s", d.Get("name").(string))
		return err
	}
	if res == nil {
		return fmt.Errorf("account name doesn't exist")
	}
	d.Set("name", res["accountName"])
	d.Set("username", res["nssUserName"])
	d.SetId(res["publicId"].(string))
	return nil
}
