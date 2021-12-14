package cloudmanager

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider is the main method for NetApp CloudManager Terraform provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDMANAGER_REFRESH_TOKEN", nil),
				Description: "The refresh_token for OCCM operations.",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDMANAGER_ENVIRONMENT", nil),
				Description: "The environment for OCCM operations.",
				Default:     "prod",
			},
			"sa_secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDMANAGER_SA_SECRET_KEY", nil),
				Description: "The environment for OCCM operations.",
			},
			"sa_client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDMANAGER_SA_CLIENT_ID", nil),
				Description: "The environment for OCCM operations.",
			},
			"simulator": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDMANAGER_SIMULATOR", nil),
				Description: "The environment for OCCM operations.",
				Default:     false,
			},
			"aws_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AWS_PROFILE", nil),
				Default:     "default",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"netapp-cloudmanager_connector_aws":   resourceOCCMAWS(),
			"netapp-cloudmanager_connector_azure": resourceOCCMAzure(),
			"netapp-cloudmanager_connector_gcp":   resourceOCCMGCP(),
			"netapp-cloudmanager_cvo_aws":         resourceCVOAWS(),
			"netapp-cloudmanager_cvo_azure":       resourceCVOAzure(),
			"netapp-cloudmanager_cvo_gcp":         resourceCVOGCP(),
			"netapp-cloudmanager_aggregate":       resourceAggregate(),
			"netapp-cloudmanager_volume":          resourceCVOVolume(),
			"netapp-cloudmanager_cifs_server":     resourceCVOCIFS(),
			"netapp-cloudmanager_snapmirror":      resourceCVOSnapMirror(),
			"netapp-cloudmanager_nss_account":     resourceCVONssAccount(),
			"netapp-cloudmanager_anf_volume":      resourceCVSANFVolume(),
			"netapp-cloudmanager_cvs_gcp_volume":  resourceCVSGCPVolume(),
			"netapp-cloudmanager_aws_fsx":         resourceAWSFSX(),
			"netapp-cloudmanager_aws_fsx_volume":  resourceFsxVolume(),
			"netapp-cloudmanager_cvo_onprem":      resourceCVOOnPrem(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"netapp-cloudmanager_cifs_server": dataSourceCVOCIFS(),
			"netapp-cloudmanager_volume":      dataSourceCVOVolume(),
			"netapp-cloudmanager_nss_account": dataSourceCVONssAccount(),
			"netapp-cloudmanager_aws_fsx":     dataSourceAWSFSX(),
			"netapp-cloudmanager_cvo_aws":     dataSourceCVOAWS(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := configStuct{
		RefreshToken: d.Get("refresh_token").(string),
		Environment:  d.Get("environment").(string),
		SaSecretKey:  d.Get("sa_secret_key").(string),
		SaClientID:   d.Get("sa_client_id").(string),
		Simulator:    d.Get("simulator").(bool),
		AWSProfile:   d.Get("aws_profile").(string),
	}

	return config.clientFun()
}
