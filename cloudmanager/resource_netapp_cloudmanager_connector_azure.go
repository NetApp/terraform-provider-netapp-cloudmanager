package cloudmanager

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOCCMAzure() *schema.Resource {
	return &schema.Resource{
		Create: resourceOCCMAzureCreate,
		Read:   resourceOCCMAzureRead,
		Delete: resourceOCCMAzureDelete,
		Exists: resourceOCCMAzureExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vnet_resource_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_security_resource_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"virtual_machine_size": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Standard_DS3_v2",
			},
			"network_security_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"company": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"proxy_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proxy_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proxy_password": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"proxy_certificates": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"associate_public_ip_address": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"admin_username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"admin_password": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceOCCMAzureCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Location = d.Get("location").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.VnetID = d.Get("vnet_id").(string)
	occmDetails.SubscriptionID = d.Get("subscription_id").(string)
	occmDetails.Company = d.Get("company").(string)
	occmDetails.AdminUsername = d.Get("admin_username").(string)
	occmDetails.AdminPassword = d.Get("admin_password").(string)
	occmDetails.VirtualMachineSize = d.Get("virtual_machine_size").(string)
	occmDetails.NetworkSecurityGroupName = d.Get("network_security_group_name").(string)
	if o, ok := d.GetOk("vnet_resource_group"); ok {
		occmDetails.VnetResourceGroup = o.(string)
	}

	if o, ok := d.GetOk("network_security_resource_group"); ok {
		occmDetails.NetworkSecurityResourceGroup = o.(string)
	}

	if o, ok := d.GetOk("proxy_url"); ok {
		occmDetails.ProxyURL = o.(string)
	}

	if o, ok := d.GetOk("proxy_user_name"); ok {
		if occmDetails.ProxyURL != "" {
			occmDetails.ProxyUserName = o.(string)
		} else {
			return fmt.Errorf("Mission proxy_url")
		}
	}

	if o, ok := d.GetOk("proxy_password"); ok {
		if occmDetails.ProxyURL != "" {
			occmDetails.ProxyPassword = o.(string)
		} else {
			return fmt.Errorf("Mission proxy_url")
		}
	}

	var proxyCertificates []string
	if certificateFiles, ok := d.GetOk("proxy_certificates"); ok {
		for _, cFile := range certificateFiles.([]interface{}) {
			// read file
			b, err := ioutil.ReadFile(cFile.(string))
			if err != nil {
				return fmt.Errorf("Cannot read certificate file: %s", err)
			}
			// endcode certificate
			encodedCertificate := base64.StdEncoding.EncodeToString(b)
			log.Printf("CFile: %s, Org cert: %s, encoded cert: %s", cFile.(string), string(b), string(encodedCertificate))
			proxyCertificates = append(proxyCertificates, encodedCertificate)
		}
	}

	if o, ok := d.GetOk("resource_group"); ok {
		occmDetails.ResourceGroup = o.(string)
	}

	if o, ok := d.GetOk("account_id"); ok {
		client.AccountID = o.(string)
	}

	if o, ok := d.GetOkExists("associate_public_ip_address"); ok {
		associatePublicIPAddress := o.(bool)
		occmDetails.AssociatePublicIPAddress = &associatePublicIPAddress
	}

	res, err := client.createOCCMAzure(occmDetails, proxyCertificates, "")
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(occmDetails.Name)
	log.Print("Set ID: ", occmDetails.Name)
	if err := d.Set("client_id", res.ClientID); err != nil {
		return fmt.Errorf("Error reading occm client_id: %s", err)
	}

	if err := d.Set("account_id", res.AccountID); err != nil {
		return fmt.Errorf("Error reading occm account_id: %s", err)
	}

	log.Printf("Created occm: %v", res)

	return resourceOCCMAzureRead(d, meta)
}

func resourceOCCMAzureRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading OCCM: %#v", d)
	client := meta.(*Client)
	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Location = d.Get("location").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.VnetID = d.Get("vnet_id").(string)
	occmDetails.SubscriptionID = d.Get("subscription_id").(string)
	occmDetails.Company = d.Get("company").(string)

	if o, ok := d.GetOk("vnet_resource_group"); ok {
		occmDetails.VnetResourceGroup = o.(string)
	}

	if o, ok := d.GetOk("resource_group"); ok {
		occmDetails.ResourceGroup = o.(string)
	}

	id := d.Id()

	resID, err := client.getdeployAzureVM(occmDetails, id)
	if err != nil {
		log.Print("Error getting occm")
		return err
	}

	if resID != id {
		return fmt.Errorf("Expected occm ID %v, Response could not find", id)
	}

	return nil
}

func resourceOCCMAzureDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := deleteOCCMDetails{}

	id := d.Id()
	occmDetails.InstanceID = id
	occmDetails.Name = d.Get("name").(string)
	occmDetails.SubscriptionID = d.Get("subscription_id").(string)
	occmDetails.Location = d.Get("location").(string)
	if o, ok := d.GetOk("resource_group"); ok {
		occmDetails.ResourceGroup = o.(string)
	}
	clientID := d.Get("client_id").(string)
	client.AccountID = d.Get("account_id").(string)

	deleteErr := client.deleteOCCMAzure(occmDetails, clientID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceOCCMAzureExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of OCCM: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Location = d.Get("location").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.VnetID = d.Get("vnet_id").(string)
	occmDetails.SubscriptionID = d.Get("subscription_id").(string)
	occmDetails.Company = d.Get("company").(string)

	if o, ok := d.GetOk("vnet_resource_group"); ok {
		occmDetails.VnetResourceGroup = o.(string)
	}

	if o, ok := d.GetOk("resource_group"); ok {
		occmDetails.ResourceGroup = o.(string)
	}

	resID, err := client.getdeployAzureVM(occmDetails, id)
	if err != nil {
		log.Print("Error getting occm")
		return false, err
	}

	if resID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
