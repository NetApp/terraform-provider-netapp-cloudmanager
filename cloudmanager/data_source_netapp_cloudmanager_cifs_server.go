package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCVOCIFS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCVOCIFSRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"netbios": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organizational_unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"svm_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"working_environment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"working_environment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workgroup_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_workgroup": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceCVOCIFSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Fetching cifs: %#v", d)

	client := meta.(*Client)
	client.ClientID = d.Get("client_id").(string)
	cifs := cifsRequest{}
	if v, ok := d.GetOk("working_environment_id"); ok {
		cifs.WorkingEnvironmentID = v.(string)
		weInfo, err := client.getWorkingEnvironmentInfo(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
		weInfo, err = client.findWorkingEnvironmentByName(weInfo.Name)
		if err != nil {
			return err
		}
		cifs.SvmName = weInfo.SvmName
	} else if v, ok := d.GetOk("working_environment_name"); ok {
		weInfo, err := client.findWorkingEnvironmentByName(v.(string))
		if err != nil {
			return nil
		}
		cifs.WorkingEnvironmentID = weInfo.PublicID
		cifs.SvmName = weInfo.SvmName
		cifs.WorkingEnvironmentType = weInfo.WorkingEnvironmentType
	} else {
		return fmt.Errorf("either working_environment_id or working_environment_name is required")
	}
	res, err := client.getCIFS(cifs)
	if err != nil {
		log.Print("Error reading cifs")
		return err
	}
	for _, cifsConfig := range res {
		d.SetId(cifs.WorkingEnvironmentID)
		d.Set("domain", "cifsConfig.Domain")
		d.Set("dns_domain", cifsConfig.DNSDomain)
		d.Set("ip_addresses", cifsConfig.IPAddresses)
		d.Set("netbios", cifsConfig.NetBIOS)
		d.Set("organizational_unit", cifsConfig.OrganizationalUnit)
		return nil
	}
	return fmt.Errorf("Error reading cifs: cifs doesn't exist")
}
