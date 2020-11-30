package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCVONssAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVONssAccountCreate,
		Read:   resourceCVONssAccountRead,
		Delete: resourceCVONssAccountDelete,
		Exists: resourceCVONssAccountExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceCVONssAccountCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating nss account: %s", d.Get("name").(string))
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	nssAcc := nssAccountRequest{}
	nssAcc.Name = d.Get("name").(string)
	nssAcc.VsaList = make([]string, 0, 0)
	if v, ok := d.GetOk("username"); ok {
		nssAcc.AccountCredentials.Username = v.(string)
	} else {
		return fmt.Errorf("username is required to create nss account")
	}
	if v, ok := d.GetOk("password"); ok {
		nssAcc.AccountCredentials.Password = v.(string)
	} else {
		return fmt.Errorf("password is required to create nss account")
	}
	res, err := client.createNssAccount(nssAcc)
	if err != nil {
		log.Printf("Error creating nss account: %s", d.Get("name").(string))
		return err
	}
	d.SetId(res["publicId"].(string))
	return resourceCVONssAccountRead(d, meta)
}

func resourceCVONssAccountRead(d *schema.ResourceData, meta interface{}) error {
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
	if _, ok := d.GetOk("name"); ok {
		d.Set("name", res["accountName"])
	}
	if _, ok := d.GetOk("username"); ok {
		d.Set("username", res["nssUserName"])
	}
	return nil
}

func resourceCVONssAccountDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting nss account: %s", d.Get("name").(string))
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	err := client.deleteNssAccount(d.Id())
	if err != nil {
		log.Printf("Error deleting nss account: %s", d.Get("name").(string))
		return err
	}
	d.SetId("")
	return nil
}

func resourceCVONssAccountExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of nss account: %s", d.Get("name").(string))
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	res, err := client.getNssAccount(d.Get("name").(string))
	if err != nil {
		log.Printf("Error checking existence of nss account: %s", d.Get("name").(string))
		return false, err
	}
	if res == nil {
		return false, nil
	}
	return true, nil
}
