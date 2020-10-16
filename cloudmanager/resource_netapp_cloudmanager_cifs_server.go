package cloudmanager

import (
	"fmt"
	"log"

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
				Optional: true,
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
			"dns_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"netbios": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"organizational_unit": {
				Type:     schema.TypeString,
				Optional: true,
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
			},
			"workgroup_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_workgroup": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
		},
	}
}

func resourceCVOCIFSCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Creating cifs: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	cifs := cifsRequest{}
	cifs.isWorkgroup = d.Get("is_workgroup").(bool)
	if cifs.isWorkgroup {
		if v, ok := d.GetOk("server_name"); ok {
			cifs.ServerName = v.(string)
		}
		if v, ok := d.GetOk("workgroup_name"); ok {
			cifs.WorkgroupName = v.(string)
		}
	} else {
		if v, ok := d.GetOk("domain"); ok {
			cifs.Domain = v.(string)
		}
		if v, ok := d.GetOk("username"); ok {
			cifs.Username = v.(string)
		}
		if v, ok := d.GetOk("password"); ok {
			cifs.Password = v.(string)
		}
		if v, ok := d.GetOk("dns_domain"); ok {
			cifs.DNSDomain = v.(string)
		}
		if v, ok := d.GetOk("ip_addresses"); ok {
			ips := make([]string, 0, v.(*schema.Set).Len())
			for _, x := range v.(*schema.Set).List() {
				ips = append(ips, x.(string))
			}
			cifs.IPAddresses = ips
		}
		if v, ok := d.GetOk("netbios"); ok {
			cifs.NetBIOS = v.(string)
		}
		if v, ok := d.GetOk("organizational_unit"); ok {
			cifs.OrganizationalUnit = v.(string)
		}
	}
	if v, ok := d.GetOk("working_environment_id"); ok {
		cifs.WorkingEnvironmentID = v.(string)
		workingEnvDetail, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
		workingEnvDetail, err = client.findWorkingEnvironmentByName(workingEnvDetail.Name)
		if err != nil {
			return err
		}
		cifs.SvmName = workingEnvDetail.SvmName
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentID = workingEnvDetail.PublicID
		cifs.SvmName = workingEnvDetail.SvmName
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	err := client.createCIFS(cifs)
	if err != nil {
		log.Print("Error creating cifs")
		return err
	}
	d.SetId(cifs.WorkingEnvironmentID)
	return resourceCVOCIFSRead(d, meta)
}

func resourceCVOCIFSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching volume: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	cifs := cifsRequest{}
	cifs.isWorkgroup = d.Get("is_workgroup").(bool)
	// TO DO: get cifs-workgroup api missing.
	if cifs.isWorkgroup {
		return nil
	}
	if v, ok := d.GetOk("working_environment_id"); ok {
		cifs.WorkingEnvironmentID = v.(string)
		workingEnvDetail, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
		workingEnvDetail, err = client.findWorkingEnvironmentByName(workingEnvDetail.Name)
		if err != nil {
			return err
		}
		cifs.SvmName = workingEnvDetail.SvmName
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentID = workingEnvDetail.PublicID
		cifs.SvmName = workingEnvDetail.SvmName
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	res, err := client.getCIFS(cifs)
	if err != nil {
		log.Print("Error reading cifs")
		return err
	}
	for _, cifsConfig := range res {
		if _, ok := d.GetOk("domain"); ok {
			d.Set("domain", cifsConfig.Domain)
		}
		if _, ok := d.GetOk("dns_domain"); ok {
			d.Set("dns_domain", cifsConfig.DNSDomain)
		}
		if _, ok := d.GetOk("ip_addresses"); ok {
			d.Set("ip_addresses", cifsConfig.IPAddresses)
		}
		if _, ok := d.GetOk("netbios"); ok {
			d.Set("netbios", cifsConfig.NetBIOS)
		}
		if _, ok := d.GetOk("organizational_unit"); ok {
			d.Set("organizational_unit", cifsConfig.OrganizationalUnit)
		}
		return nil

	}
	return fmt.Errorf("Error reading cifs: cifs doesn't exist")
}

func resourceCVOCIFSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Deleting cifs: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	cifs := cifsRequest{}
	if v, ok := d.GetOk("working_environment_id"); ok {
		cifs.WorkingEnvironmentID = v.(string)
		workingEnvDetail, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
		workingEnvDetail, err = client.findWorkingEnvironmentByName(workingEnvDetail.Name)
		if err != nil {
			return err
		}
		cifs.SvmName = workingEnvDetail.SvmName
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentID = workingEnvDetail.PublicID
		cifs.SvmName = workingEnvDetail.SvmName
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	err := client.deleteCIFS(cifs)
	if err != nil {
		log.Print("Error deleting cifs")
		return err
	}
	return nil
}

func resourceCVOCIFSExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("Fetching cifs: %#v", d)
	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	cifs := cifsRequest{}
	if v, ok := d.GetOk("working_environment_id"); ok {
		cifs.WorkingEnvironmentID = v.(string)
		workingEnvDetail, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return false, nil
		}
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
		workingEnvDetail, err = client.findWorkingEnvironmentByName(workingEnvDetail.Name)
		if err != nil {
			return false, err
		}
		cifs.SvmName = workingEnvDetail.SvmName
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return false, nil
		}
		cifs.WorkingEnvironmentID = workingEnvDetail.PublicID
		cifs.SvmName = workingEnvDetail.SvmName
		cifs.WorkingEnvironmentType = workingEnvDetail.WorkingEnvironmentType
	} else {
		return false, fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	res, err := client.getCIFS(cifs)
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
