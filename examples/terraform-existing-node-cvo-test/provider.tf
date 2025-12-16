terraform {
  required_providers {
    netapp-cloudmanager = {
      source  = "netapp/netapp-cloudmanager"
      version = "26.0.0"  
    }
  }
}

provider "netapp-cloudmanager" {
  refresh_token = var.cloudmanager_refresh_token
  environment   = var.cloudmanager_environment
}