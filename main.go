package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/netapp/terraform-provider-netapp-cloudmanager/cloudmanager"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cloudmanager.Provider,
	})
}
