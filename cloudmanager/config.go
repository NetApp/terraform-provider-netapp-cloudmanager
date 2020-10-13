package cloudmanager

import (
	"fmt"
	"log"
)

// Config is a struct for user input
type configStuct struct {
	RefreshToken string
	Environment  string
	CVOHostName  string
}

// Client is the main function to connect to the APi
func (c *configStuct) clientFun() (*Client, error) {
	var client *Client
	if c.Environment == "prod" {
		log.Print("Prod Environment")
		client = &Client{
			CloudManagerHost:     "https://cloudmanager.cloud.netapp.com",
			AuthHost:             "https://netapp-cloud-account.auth0.com/oauth/token",
			Audience:             "https://api.cloud.netapp.com/",
			Auth0Client:          "Mu0V1ywgYteI6w1MbD15fKfVIUrNXGWC",
			AMIFilter:            "Setup-As-Service-AMI-Prod*",
			AWSAccount:           "952013314444",
			GCPDeploymentManager: "https://www.googleapis.com",
			GCPImageProject:      "netapp-cloudmanager",
			GCPImageFamily:       "cloudmanager",
		}
	} else if c.Environment == "stage" {
		log.Print("Stage Environment")
		client = &Client{
			CloudManagerHost:        "https://staging.cloudmanager.cloud.netapp.com",
			AuthHost:                "https://staging-netapp-cloud-account.auth0.com/oauth/token",
			Audience:                "https://api.cloud.netapp.com/",
			Auth0Client:             "O6AHa7kedZfzHaxN80dnrIcuPBGEUvEv",
			AMIFilter:               "Setup-As-Service-AMI-*",
			AWSAccount:              "282316784512",
			GCPDeploymentManager:    "https://www.googleapis.com",
			GCPImageProject:         "tlv-automation",
			GCPImageFamily:          "occm-automation",
			AzureEnvironmentForOCCM: "stage",
		}
	} else {
		return &Client{}, fmt.Errorf("expected environment to be one of [prod stage], %s", c.Environment)
	}

	client.SetRefreshToken(c.RefreshToken)

	return client, nil
}
