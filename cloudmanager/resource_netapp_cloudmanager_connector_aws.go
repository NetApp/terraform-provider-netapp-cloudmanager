package cloudmanager

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOCCMAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceOCCMAWSCreate,
		Read:   resourceOCCMAWSRead,
		Delete: resourceOCCMAWSDelete,
		Exists: resourceOCCMAWSExists,
		Update: resourceOCCMAWSUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceOCCMAWSImport,
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.Contains(v, " ") {
						errs = append(errs, fmt.Errorf("%q must not contain space, got: %s", key, v))
					}
					return
				},
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
			"public_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

	res, err := client.createOCCM(occmDetails, proxyCertificates, "")

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
	var clientID string
	occmDetails := createOCCMDetails{}

	occmDetails.Name = d.Get("name").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)

	if v, ok := d.GetOk("client_id"); ok {
		clientID = v.(string)
	}

	if o, ok := d.GetOk("ami"); ok {
		occmDetails.AMI = o.(string)
	}

	id := d.Id()

	res, err := client.getAWSInstance(occmDetails, id)
	if err != nil {
		log.Print("Error getting occm")
		return err
	}

	if *res.InstanceId != id {
		return fmt.Errorf("Expected occm ID %v, Response could not find", id)
	}

	if occmDetails.Region == "" {
		occmDetails.Region = *res.Placement.AvailabilityZone
		occmDetails.Region = occmDetails.Region[:len(occmDetails.Region)-1]
		d.Set("region", occmDetails.Region)
	}
	occmDetails.InstanceID = *res.InstanceId
	disableAPITermination, err := client.CallAWSDescribeInstanceAttribute(occmDetails)
	if err != nil {
		return err
	}
	d.Set("enable_termination_protection", disableAPITermination)
	d.Set("instance_type", res.InstanceType)
	d.Set("subnet_id", res.SubnetId)
	var sgIDs string
	var remoteSgIDsArray []string
	for _, sgID := range res.SecurityGroups {
		remoteSgIDsArray = append(remoteSgIDsArray, *sgID.GroupId)
	}
	// preserve order of the ids in state file
	localSgIDs := d.Get("security_group_id").(string)
	localSgIDsArray := strings.Split(localSgIDs, ",")
	m := make(map[string]bool)

	for _, item := range localSgIDsArray {
		m[item] = true
	}

	for _, item := range remoteSgIDsArray {
		if _, ok := m[item]; !ok {
			localSgIDsArray = append(localSgIDsArray, item)
		}
	}
	sgIDs = strings.Join(localSgIDsArray, ",")
	d.Set("security_group_id", sgIDs)
	d.Set("key_name", res.KeyName)
	iamInstanceProfile := *res.IamInstanceProfile.Arn
	slashIndex := strings.Index(iamInstanceProfile, "/")
	iamInstanceProfile = iamInstanceProfile[slashIndex+1:]
	d.Set("iam_instance_profile_name", iamInstanceProfile)
	if _, ok := d.GetOk("ami"); ok {
		d.Set("ami", res.ImageId)
	}
	// The following tags are ignored.
	excludedTags := [...]string{"Name", "OCCMInstance", "Owner", "PrincipalId"}
	tags := make([]map[string]string, 0)
	for _, tag := range res.Tags {
		var exclusion bool
		for _, exTag := range excludedTags {
			if *tag.Key == exTag {
				exclusion = true
			}
		}
		if exclusion == false {
			tagMap := make(map[string]string)
			tagMap["tag_key"] = *tag.Key
			tagMap["tag_value"] = *tag.Value
			tags = append(tags, tagMap)
		}
		if *tag.Key == "Name" {
			d.Set("name", *tag.Value)
		}
	}
	d.Set("aws_tag", tags)

	if res.PublicIpAddress == nil {
		d.Set("associate_public_ip_address", false)
	} else {
		d.Set("associate_public_ip_address", true)
		d.Set("public_ip_address", *res.PublicIpAddress)
	}

	if _, ok := d.GetOk("company"); !ok {
		company, err := client.getCompany(clientID)
		if err != nil {
			log.Printf("Error when reading system info from cloudmanager.")
			return err
		}
		d.Set("company", company)
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
	clientID := d.Get("client_id").(string)
	client.AccountID = d.Get("account_id").(string)

	deleteErr := client.deleteOCCM(occmDetails, clientID)
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

	if res.InstanceId == nil {
		d.SetId("")
		return false, nil
	}

	if *res.InstanceId != id {
		d.SetId("")
		return false, nil
	}

	return true, nil
}

// resourceOCCMAWSUpdate updates occm. Currently only tags can be updated.
func resourceOCCMAWSUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)

	occmDetails := createOCCMDetails{}
	deleteAwsTags := []userTags{}
	addModifyAwsTags := []userTags{}
	if d.HasChange("aws_tag") {
		old, new := d.GetChange("aws_tag")
		oldTags := old.(*schema.Set)
		newTags := new.(*schema.Set)
		addModifyTags := newTags.Difference(oldTags)
		// the firt difference function gives the set that only in old tags but not in new tags, which includes tags to be deleted or updated.
		// the second difference function excludes the tags to be modified.
		deleteTags := oldTags.Difference(newTags)
		deleteTags = deleteTags.Difference(addModifyTags)
		if deleteTags.Len() > 0 {
			deleteAwsTags = expandUserTags(deleteTags)
		}
		if addModifyTags.Len() > 0 {
			addModifyAwsTags = expandUserTags(addModifyTags)
		}
	}
	clientID := d.Get("client_id").(string)
	occmDetails.Region = d.Get("region").(string)
	occmDetails.Name = d.Get("name").(string)
	occmDetails.InstanceType = d.Get("instance_type").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.SecurityGroupID = d.Get("security_group_id").(string)
	occmDetails.KeyName = d.Get("key_name").(string)
	occmDetails.IamInstanceProfileName = d.Get("iam_instance_profile_name").(string)
	occmDetails.Company = d.Get("company").(string)
	occmDetails.InstanceID = d.Id()

	err := client.updateOCCM(occmDetails, nil, deleteAwsTags, addModifyAwsTags, clientID)

	if err != nil {
		log.Print("Error updating instance")
		return err
	}
	return nil
}

func resourceOCCMAWSImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Wrong format of resource: %s. Please follow 'client_id/connector_id'", d.Id())
	}

	d.SetId(parts[1])
	d.Set("client_id", parts[0])

	return []*schema.ResourceData{d}, nil

}
