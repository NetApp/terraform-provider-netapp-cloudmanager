package cloudmanager

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOCCMGCP() *schema.Resource {
	return &schema.Resource{
		Create: resourceOCCMGCPCreate,
		Read:   resourceOCCMGCPRead,
		Delete: resourceOCCMGCPDelete,
		Exists: resourceOCCMGCPExists,
		Update: resourceOCCMGCPUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_account_email": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_account_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"service_account_key"},
				Default:       "",
			},
			"service_account_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"service_account_path"},
				Default:       "",
				Sensitive:     true,
			},
			"machine_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "n2-standard-4",
				ForceNew: true,
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "default",
				ForceNew:     true,
				ValidateFunc: validateSubnet(),
			},
			"network_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"firewall_tags": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
			"associate_public_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceOCCMGCPCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.GCPCommonSuffixName = "-vm-boot-deployment"
	occmDetails.Name = d.Get("name").(string)
	occmDetails.GCPProject = d.Get("project_id").(string)
	occmDetails.Zone = d.Get("zone").(string)
	occmDetails.Region = string(occmDetails.Zone[0 : len(occmDetails.Zone)-2])
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	occmDetails.MachineType = d.Get("machine_type").(string)
	occmDetails.ServiceAccountEmail = d.Get("service_account_email").(string)
	var err error
	client.GCPServiceAccountKey, err = getGCPServiceAccountKey(d)
	if err != nil {
		return err
	}

	occmDetails.FirewallTags = d.Get("firewall_tags").(bool)
	occmDetails.AssociatePublicIP = d.Get("associate_public_ip").(bool)
	occmDetails.Company = d.Get("company").(string)

	if o, ok := d.GetOk("network_project_id"); ok {
		occmDetails.NetworkProjectID = o.(string)
	}

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

	if o, ok := d.GetOk("account_id"); ok {
		client.AccountID = o.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		o := v.(*schema.Set)
		tags := make([]string, 0, o.Len())
		for _, v := range o.List() {
			tags = append(tags, v.(string))
		}
		occmDetails.Tags = tags
	}

	res, err := client.deployGCPVM(occmDetails, proxyCertificates, "")
	if err != nil {
		log.Print("Error creating instance")
		return err
	}

	d.SetId(occmDetails.Name)
	if err := d.Set("client_id", res.ClientID); err != nil {
		return fmt.Errorf("Error reading occm client_id: %s", err)
	}

	if err := d.Set("account_id", res.AccountID); err != nil {
		return fmt.Errorf("Error reading occm account_id: %s", err)
	}

	log.Printf("Created occm: %v", res)

	return resourceOCCMGCPRead(d, meta)
}

func resourceOCCMGCPRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading OCCM: %#v", d)
	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.GCPCommonSuffixName = "-vm-boot-deployment"
	occmDetails.Name = d.Get("name").(string)
	occmDetails.GCPProject = d.Get("project_id").(string)
	occmDetails.Region = d.Get("zone").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	var err error
	client.GCPServiceAccountKey, err = getGCPServiceAccountKey(d)
	if err != nil {
		return err
	}
	occmDetails.Company = d.Get("company").(string)
	clientID := d.Get("client_id").(string)

	id := d.Id() + "-vm-boot-deployment"

	resID, err := client.getdeployGCPVM(occmDetails, id, clientID)
	if err != nil {
		log.Print("Error getting occm")
		return err
	}

	if resID != id {
		return fmt.Errorf("Expected occm ID %v, Response could not find", id)
	}

	if _, ok := d.GetOk("tags"); ok {
		instance, err := client.getVMInstance(occmDetails, clientID)
		tagItems := instance["tags"].(map[string]interface{})
		tags := tagItems["items"].([]interface{})
		if err != nil {
			return err
		}
		var current []string
		if d.Get("firewall_tags").(bool) {
			for _, v := range tags {
				if v.(string) != "firewall-tag-bvsu" && v.(string) != "http-server" && v.(string) != "https-server" {
					current = append(current, v.(string))
				}
			}
			d.Set("tags", current)
		} else {
			d.Set("tags", tags)
		}
	}

	return nil
}

func resourceOCCMGCPDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting OCCM: %#v", d)

	client := meta.(*Client)

	occmDetails := deleteOCCMDetails{}

	occmDetails.GCPCommonSuffixName = "-vm-boot-deployment"
	id := d.Id() + occmDetails.GCPCommonSuffixName
	occmDetails.InstanceID = id
	occmDetails.Name = d.Get("name").(string)
	occmDetails.Project = d.Get("project_id").(string)
	var err error
	client.GCPServiceAccountKey, err = getGCPServiceAccountKey(d)
	if err != nil {
		return err
	}
	occmDetails.Region = d.Get("zone").(string)
	clientID := d.Get("client_id").(string)
	client.AccountID = d.Get("account_id").(string)

	deleteErr := client.deleteOCCMGCP(occmDetails, clientID)
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func resourceOCCMGCPExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Checking existence of OCCM: %#v", d)
	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.GCPCommonSuffixName = "-vm-boot-deployment"
	occmDetails.Name = d.Get("name").(string)
	occmDetails.GCPProject = d.Get("project_id").(string)
	occmDetails.Region = d.Get("zone").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	var err error
	client.GCPServiceAccountKey, err = getGCPServiceAccountKey(d)
	if err != nil {
		return false, err
	}
	occmDetails.Company = d.Get("company").(string)
	clientID := d.Get("client_id").(string)

	id := d.Id() + occmDetails.GCPCommonSuffixName

	resID, err := client.getdeployGCPVM(occmDetails, id, clientID)
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

func resourceOCCMGCPUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Updating OCCM: %#v", d)
	client := meta.(*Client)

	occmDetails := createOCCMDetails{}

	occmDetails.GCPCommonSuffixName = "-vm-boot-deployment"
	occmDetails.Name = d.Get("name").(string)
	occmDetails.GCPProject = d.Get("project_id").(string)
	occmDetails.Region = d.Get("zone").(string)
	occmDetails.SubnetID = d.Get("subnet_id").(string)
	var err error
	client.GCPServiceAccountKey, err = getGCPServiceAccountKey(d)
	if err != nil {
		return err
	}
	occmDetails.Company = d.Get("company").(string)
	clientID := d.Get("client_id").(string)

	if d.HasChange("tags") {
		instance, err := client.getVMInstance(occmDetails, clientID)
		tagItems := instance["tags"].(map[string]interface{})
		fingerprint := tagItems["fingerprint"].(string)
		o := d.Get("tags").(*schema.Set)
		tags := make([]string, 0, o.Len())
		for _, v := range o.List() {
			tags = append(tags, v.(string))
		}
		if d.Get("firewall_tags").(bool) {
			tags = append(tags, []string{"firewall-tag-bvsu", "http-server", "https-server"}...)
		}
		occmDetails.Tags = tags
		err = client.setVMInstaceTags(occmDetails, fingerprint, clientID)
		if err != nil {
			return err
		}
	}
	return resourceOCCMGCPRead(d, meta)
}

func getGCPServiceAccountKey(d *schema.ResourceData) (string, error) {
	serviceAccountPath := d.Get("service_account_path").(string)
	serviceAccountKey := d.Get("service_account_key").(string)
	if serviceAccountPath != "" {
		serviceAccountKey, err := ioutil.ReadFile(serviceAccountPath)
		if err != nil {
			return "", fmt.Errorf("Cannot read service account file: %s", err)
		}
		return string(serviceAccountKey), nil
	} else if serviceAccountKey != "" {
		return serviceAccountKey, nil
	} else {
		// Check if default application credential file exists
		_, ok := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS")
		if ok {
			return "", nil
		}
		if _, err := os.Stat(wellKnownFile()); err == nil {
			return "", nil
		}
	}
	return "", fmt.Errorf("Neither service_account_path nor service_account_key is set, unable to proceed")
}

func guessUnixHomeDir() string {
	// Prefer $HOME over user.Current
	if v := os.Getenv("HOME"); v != "" {
		return v
	}
	// Else, fall back to user.Current:
	if u, err := user.Current(); err == nil {
		return u.HomeDir
	}
	return ""
}

func wellKnownFile() string {
	const f = "application_default_credentials.json"
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), "gcloud", f)
	}
	return filepath.Join(guessUnixHomeDir(), ".config", "gcloud", f)
}

func validateSubnet() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}
		v = strings.ToLower(v)
		slices := strings.Split(v, "/")
		errorMessage := "invalid format of subnet, the correct format is either <subnetID> or projects/<projectID>/regions/<region>/subnetworks/<subnetID>"
		if len(slices) != 1 && len(slices) != 6 {
			es = append(es, fmt.Errorf(errorMessage))
			return
		}

		if len(slices) == 6 && (slices[0] != "projects" || slices[2] != "regions" || slices[4] != "subnetworks") {
			es = append(es, fmt.Errorf(errorMessage))
			return
		}
		return
	}
}
