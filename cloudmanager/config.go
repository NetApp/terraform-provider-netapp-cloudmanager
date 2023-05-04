package cloudmanager

import (
	"fmt"
	"log"
	"strings"
)

// Config is a struct for user input
type configStruct struct {
	RefreshToken       string
	SaSecretKey        string
	SaClientID         string
	Environment        string
	CVOHostName        string
	Simulator          bool
	AWSProfile         string
	AWSProfileFilePath string
	AzureAuthMethods   []string
	ConnectorHost      string
}

// Client is the main function to connect to the APi
func (c *configStruct) clientFun() (*Client, error) {
	var client *Client
	if c.Environment == "prod" {
		log.Print("Prod Environment")
		client = &Client{
			CloudManagerHost:     "https://cloudmanager.cloud.netapp.com",
			AuthHost:             "https://netapp-cloud-account.auth0.com/oauth/token",
			SaAuthHost:           "https://cloudmanager.cloud.netapp.com/auth/oauth/token",
			Audience:             "https://api.cloud.netapp.com",
			Auth0Client:          "Mu0V1ywgYteI6w1MbD15fKfVIUrNXGWC",
			AMIFilter:            "Setup-As-Service-AMI-Prod*",
			AWSAccount:           "952013314444",
			GCPDeploymentManager: "https://www.googleapis.com",
			GCPCompute:           "https://compute.googleapis.com",
			GCPImageProject:      "netapp-cloudmanager",
			GCPImageFamily:       "cloudmanager",
			CVSHostName:          "https://api.services.cloud.netapp.com",
		}
	} else if c.Environment == "stage" {
		log.Print("Stage Environment")
		client = &Client{
			CloudManagerHost:        "https://staging.cloudmanager.cloud.netapp.com",
			AuthHost:                "https://staging-netapp-cloud-account.auth0.com/oauth/token",
			SaAuthHost:              "https://staging.cloudmanager.cloud.netapp.com/auth/oauth/token",
			Audience:                "https://api.cloud.netapp.com",
			Auth0Client:             "O6AHa7kedZfzHaxN80dnrIcuPBGEUvEv",
			AMIFilter:               "Setup-As-Service-AMI-*",
			AWSAccount:              "282316784512",
			GCPDeploymentManager:    "https://www.googleapis.com",
			GCPCompute:              "https://compute.googleapis.com",
			GCPImageProject:         "tlv-automation",
			GCPImageFamily:          "occm-automation",
			AzureEnvironmentForOCCM: "stage",
			CVSHostName:             "https://staging.api.services.cloud.netapp.com",
		}
	} else {
		return &Client{}, fmt.Errorf("expected environment to be one of [prod stage], %s", c.Environment)
	}

	if c.SaSecretKey != "" && c.SaClientID != "" {
		client.SetServiceCredential(c.SaSecretKey, c.SaClientID)
	} else if c.RefreshToken != "" {
		client.SetRefreshToken(c.RefreshToken)
	} else {
		return &Client{}, fmt.Errorf("expected refresh_token or sa_secret_key and sa_client_id")
	}

	if c.ConnectorHost != "" {
		if strings.HasPrefix(c.ConnectorHost, "https://") {
			client.CloudManagerHost = c.ConnectorHost
		} else if strings.HasPrefix(c.ConnectorHost, "http://") {
			client.CloudManagerHost = c.ConnectorHost
		} else {
			client.CloudManagerHost = fmt.Sprintf("https://%s", c.ConnectorHost)
		}
	}

	if c.Simulator {
		client.SetSimulator(c.Simulator)
	}
	client.AWSProfile = c.AWSProfile
	client.AWSProfileFilePath = c.AWSProfileFilePath
	client.AzureAuthMethods = c.AzureAuthMethods

	return client, nil
}
