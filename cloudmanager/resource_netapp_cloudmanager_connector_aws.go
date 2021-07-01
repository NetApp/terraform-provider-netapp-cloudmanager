package cloudmanager

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOCCMAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceOCCMAWSCreate,
		Read:   resourceOCCMAWSRead,
		Delete: resourceOCCMAWSDelete,
		Exists: resourceOCCMAWSExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ami": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "t3.xlarge",
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"iam_instance_profile_name": {
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
			"enable_termination_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"associate_public_ip_address": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"aws_tag": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tag_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"tag_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceOCCMAWSCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)
	if o, ok := d.GetOk("proxy_url"); ok {
		occmDetails.ProxyURL = o.(string)
	}

	if o, ok := d.GetOk("proxy_user_name"); ok {
		if occmDetails.ProxyURL != "" {
			occmDetails.ProxyUserName = o.(string)
		} else {
			return fmt.Errorf("Missing proxy_url")
		}
	}

	if o, ok := d.GetOk("proxy_password"); ok {
		if occmDetails.ProxyURL != "" {
			occmDetails.ProxyPassword = o.(string)
		} else {
			return fmt.Errorf("Missing proxy_url")
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

	if o, ok := d.GetOk("ami"); ok {
		occmDetails.AMI = o.(string)
	}

	if o, ok := d.GetOk("account_id"); ok {
		client.AccountID = o.(string)
	}

	if o, ok := d.GetOkExists("associate_public_ip_address"); ok {
		associatePublicIPAddress := o.(bool)
		occmDetails.AssociatePublicIPAddress = &associatePublicIPAddress
	}

	if o, ok := d.GetOkExists("enable_termination_protection"); ok {
		enableTerminationProtection := o.(bool)
		occmDetails.EnableTerminationProtection = &enableTerminationProtection
	}

	if o, ok := d.GetOk("aws_tag"); ok {
		tags := o.(*schema.Set)
		if tags.Len() > 0 {
			occmDetails.AwsTags = expandUserTags(tags)
		}
	}

	res, err := client.createOCCM(occmDetails, proxyCertificates)

	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(res.InstanceID)
	if err := d.Set("client_id", res.ClientID); err != nil {
		return fmt.Errorf("Error reading occm client_id: %s", err)
	}

	if err := d.Set("account_id", res.AccountID); err != nil {
		return fmt.Errorf("Error reading occm account_id: %s", err)
	}

	log.Printf("Created occm: %v", res)

	return resourceOCCMAWSRead(d, meta)
}

func resourceOCCMAWSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading OCCM: %#v", d)
	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)

	if o, ok := d.GetOk("ami"); ok {
		occmDetails.AMI = o.(string)
	}

	id := d.Id()

	res, err := client.getAWSInstance(occmDetails, id)
	if err != nil {
		log.Print("Error getting occm")
		return err
	}

	if res.InstanceID != id {
		return fmt.Errorf("Expected occm ID %v, Response could not find", id)
	}

	return nil
}

func resourceOCCMAWSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := deleteOCCMDetails{}

	id := d.Id()
	occmDetails.InstanceID = id
	occmDetails.Region = d.Get("region").(string)
	client.ClientID = d.Get("client_id").(string)
	client.AccountID = d.Get("account_id").(string)

	deleteErr := client.deleteOCCM(occmDetails)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceOCCMAWSExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of OCCM: %#v", d)
	client := meta.(*Client)

	id := d.Id()
	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)

	if o, ok := d.GetOk("ami"); ok {
		occmDetails.AMI = o.(string)
	}

	res, err := client.getAWSInstance(occmDetails, id)
	if err != nil {
		log.Print("Error getting occm")
		return false, err
	}

	if res.InstanceID != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}
