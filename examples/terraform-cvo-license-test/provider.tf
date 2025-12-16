terraform {
  required_providers {
    netapp-cloudmanager = {
      source  = "netapp/netapp-cloudmanager"
      version = ">= 24.12.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.0"
    }
  }
}

provider "netapp-cloudmanager" {
  refresh_token = var.cloudmanager_refresh_token
  environment   = var.cloudmanager_environment
}

provider "azurerm" {
  features {}
  
  subscription_id = var.azure_subscription_id
  # tenant_id       = var.azure_tenant_id
  # client_id       = var.azure_client_id
  # client_secret   = var.azure_client_secret
  
  # Skip provider registration since connector handles this
  skip_provider_registration = true
}