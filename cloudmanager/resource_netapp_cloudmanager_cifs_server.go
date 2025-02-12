package cloudmanager

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCVOCIFS() *schema.Resource {
	return &schema.Resource{
		Create: resourceCVOCIFSCreate,
		Read:   resourceCVOCIFSRead,
		Delete: resourceCVOCIFSDelete,
		Exists: resourceCVOCIFSExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
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
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"dns_domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"netbios": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organizational_unit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"svm_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"working_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"working_environment_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.TrimSpace(v) != "" {
						errs = append(errs, fmt.Errorf("using workgroup configuration is deprecated. Create with AD instead"))
					}
					return
				},
			},
			"workgroup_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.TrimSpace(v) != "" {
						errs = append(errs, fmt.Errorf("using workgroup configuration is deprecated. Create with AD instead"))
					}
					return
				},
			},
			"is_workgroup": {
				Type:       schema.TypeBool,
				Optional:   true,
				ForceNew:   true,
				Default:    false,
				Deprecated: "use AD instead",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(bool)
					if v {
						errs = append(errs, fmt.Errorf("using workgroup configuration is deprecated. Create with AD instead"))
					}
					return
				},
			},
		},
	}
}

func resourceCVOCIFSCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating cifs: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	cifs := cifsRequest{}

	cifs.Domain = d.Get("domain").(string)
	cifs.Username = d.Get("username").(string)
	cifs.Password = d.Get("password").(string)
	cifs.DNSDomain = d.Get("dns_domain").(string)
	IPAddresses := d.Get("ip_addresses")

	for _, IPAddress := range IPAddresses.([]interface{}) {
		cifs.IPAddresses = append(cifs.IPAddresses, IPAddress.(string))
	}

	cifs.NetBIOS = d.Get("netbios").(string)
	cifs.OrganizationalUnit = d.Get("organizational_unit").(string)

	workingEnvDetail, err := client.getWorkingEnvironmentDetail(d, clientID, true, "")
	if err != nil {
		return err
	}
	cifs.WorkingEnvironmentID = workingEnvDetail.PublicID

	if v, ok := d.GetOk("svm_name"); ok {
		cifs.SvmName = v.(string)
	} else {
		cifs.SvmName = workingEnvDetail.SvmName
	}

	err = client.createCIFS(cifs, clientID)
	if err != nil {
		log.Print("Error creating cifs")
		return err
	}
	d.SetId(cifs.SvmName)
	return resourceCVOCIFSRead(d, meta)
}

func resourceCVOCIFSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching cifs: %#v", d)

	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	cifs := cifsRequest{}

	workingEnvDetail, err := client.getWorkingEnvironmentDetail(d, clientID, true, "")
	if err != nil {
		return err
	}
	cifs.WorkingEnvironmentID = workingEnvDetail.PublicID

	if v, ok := d.GetOk("svm_name"); ok {
		cifs.SvmName = v.(string)
	} else {
		cifs.SvmName = workingEnvDetail.SvmName
	}
	res, err := client.getCIFS(cifs, clientID)
	if err != nil {
		log.Print("Error reading cifs")
		return err
	}
	for _, cifsConfig := range res {
		log.Printf("cifs config: %#v", cifsConfig)
		// domain, dns_domain, netbios are all converted to low case in the API response
		if cifsConfig.Domain == strings.ToLower(d.Get("domain").(string)) && cifsConfig.DNSDomain == strings.ToLower(d.Get("dns_domain").(string)) &&
			cifsConfig.NetBIOS == strings.ToLower(d.Get("netbios").(string)) && cifsConfig.OrganizationalUnit == d.Get("organizational_unit").(string) {
			d.Set("ip_addresses", cifsConfig.IPAddresses)
			return nil
		}
	}
	return fmt.Errorf("error reading cifs: cifs doesn't exist")
}

func resourceCVOCIFSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting cifs: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	cifs := cifsDeleteRequest{}

	var workingEnvironmentID string
	if v, ok := d.GetOk("working_environment_id"); ok {
		workingEnvironmentID = v.(string)
	} else {
		workingEnvDetail, err := client.getWorkingEnvironmentDetail(d, clientID, true, "")
		if err != nil {
			return err
		}
		workingEnvironmentID = workingEnvDetail.PublicID
	}
	cifs.Username = d.Get("username").(string)
	cifs.Password = d.Get("password").(string)
	if v, ok := d.GetOk("svm_name"); ok {
		cifs.SvmName = v.(string)
	}
	err := client.deleteCIFS(cifs, workingEnvironmentID, clientID)
	if err != nil {
		log.Print("Error deleting cifs")
		return err
	}
	return nil
}

func resourceCVOCIFSExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Fetching existing cifs: %#v", d)
	client := meta.(*Client)
	clientID := d.Get("client_id").(string)
	cifs := cifsRequest{}

	workingEnvDetail, err := client.getWorkingEnvironmentDetail(d, clientID, true, "")
	if err != nil {
		return false, err
	}
	cifs.WorkingEnvironmentID = workingEnvDetail.PublicID
	if v, ok := d.GetOk("svm_name"); ok {
		cifs.SvmName = v.(string)
	} else {
		cifs.SvmName = workingEnvDetail.SvmName
	}
	res, err := client.getCIFS(cifs, clientID)
	if err != nil {
		log.Print("Error reading cifs")
		return false, err
	}
	if len(res) == 0 {
		d.SetId("")
		return false, fmt.Errorf("cifs doesn't exist")
	}

	return true, nil
}

func resourceCVOCIFSUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
