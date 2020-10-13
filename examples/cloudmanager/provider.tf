# Specify the provider and access details
provider "netapp-cloudmanager" {
  refresh_token = var.cloudmanager_refresh_token
}

terraform {
  required_version = ">= 0.13"
  required_providers {
    netapp-cloudmanager = {
      source = "hashicorp/netapp-cloudmanager"
      version = "~> 20.10.0"
    }
  }
}
