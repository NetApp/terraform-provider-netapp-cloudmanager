package cloudmanager

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCVOAWS() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCVOAWSRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"svm_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceCVOAWSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("Reading CVO: %#v", d)
	client := meta.(*Client)

	if c, ok := d.GetOk("client_id"); ok {
		client.ClientID = c.(string)
	}
	if a, ok := d.GetOk("id"); ok {
		WorkingEnvironmentID := a.(string)
		workingEnvDetail, err := client.findWorkingEnvironmentByID(WorkingEnvironmentID)
		if err != nil {
			return fmt.Errorf("Cannot find working environment by working_environment_id %s", WorkingEnvironmentID)
		}
		d.SetId(workingEnvDetail.PublicID)
		d.Set("name", workingEnvDetail.Name)
		d.Set("svm_name", workingEnvDetail.SvmName)
	} else if a, ok = d.GetOk("name"); ok {
		workingEnvDetail, err := client.findWorkingEnvironmentByName(a.(string))
		if err != nil {
			return fmt.Errorf("Cannot find working environment by working_environment_name %s", a.(string))
		}
		d.SetId(workingEnvDetail.PublicID)
		d.Set("name", workingEnvDetail.Name)
		d.Set("svm_name", workingEnvDetail.SvmName)
	}

	return nil
}
