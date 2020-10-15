# Specify CVO resources

provider "azurerm" {
  features {}
  skip_provider_registration = true
}

terraform {
  required_providers {
    netapp-cloudmanager = {
      source = "NetApp/netapp-cloudmanager"
      version = "20.10.0"
    }
  }
}

resource "azurerm_resource_group" "occm-rg" {
  name     = "CMTerraform"
  location = "westus"
}

data "azurerm_subscription" "primary" {
}

resource "netapp-cloudmanager_connector_azure" "cm-azure" {
  depends_on = [azurerm_resource_group.occm-rg]
  provider = netapp-cloudmanager
  name = "CMTerraform"
  location = "westus"
  subscription_id = data.azurerm_subscription.primary.subscription_id
  company = "NetApp"
  resource_group = "CMTerraform"
  vnet_resource_group = "rg_westus"
  subnet_id = "Subnet1"
  vnet_id = "Vnet1"
  account_id = "account-xxxxxxx"
  admin_password = "P@ssword!"
  admin_username = "vmadmin"
}

data "azurerm_virtual_machine" "occm-vm" {
  depends_on = [netapp-cloudmanager_connector_azure.cm-azure]
  name                = "CMTerraform"
  resource_group_name = "CMTerraform"
}

resource "azurerm_role_definition" "occm-role" {
  role_definition_id = "12345678-0000-0000-b9d7-123456789012"
  name        = "Terraform-OCCM-Role"
  scope       = data.azurerm_subscription.primary.id
  description = "This is a custom role created via Terraform"

  permissions {
    actions     = ["Microsoft.Compute/disks/delete",
        "Microsoft.Compute/disks/read",
        "Microsoft.Compute/disks/write",
        "Microsoft.Compute/locations/operations/read",
        "Microsoft.Compute/locations/vmSizes/read",
        "Microsoft.Resources/subscriptions/locations/read",
        "Microsoft.Compute/operations/read",
        "Microsoft.Compute/virtualMachines/instanceView/read",
        "Microsoft.Compute/virtualMachines/powerOff/action",
        "Microsoft.Compute/virtualMachines/read",
        "Microsoft.Compute/virtualMachines/restart/action",
        "Microsoft.Compute/virtualMachines/deallocate/action",
        "Microsoft.Compute/virtualMachines/start/action",
        "Microsoft.Compute/virtualMachines/vmSizes/read",
        "Microsoft.Compute/virtualMachines/write",
        "Microsoft.Compute/images/write",
        "Microsoft.Compute/images/read",
        "Microsoft.Network/locations/operationResults/read",
        "Microsoft.Network/locations/operations/read",
        "Microsoft.Network/networkInterfaces/read",
        "Microsoft.Network/networkInterfaces/write",
        "Microsoft.Network/networkInterfaces/join/action",
        "Microsoft.Network/networkSecurityGroups/read",
        "Microsoft.Network/networkSecurityGroups/write",
        "Microsoft.Network/networkSecurityGroups/join/action",
        "Microsoft.Network/virtualNetworks/read",
        "Microsoft.Network/virtualNetworks/checkIpAddressAvailability/read",
        "Microsoft.Network/virtualNetworks/subnets/read",
        "Microsoft.Network/virtualNetworks/subnets/write",
        "Microsoft.Network/virtualNetworks/subnets/virtualMachines/read",
        "Microsoft.Network/virtualNetworks/virtualMachines/read",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Resources/deployments/operations/read",
        "Microsoft.Resources/deployments/read",
        "Microsoft.Resources/deployments/write",
        "Microsoft.Resources/resources/read",
        "Microsoft.Resources/subscriptions/operationresults/read",
        "Microsoft.Resources/subscriptions/resourceGroups/delete",
        "Microsoft.Resources/subscriptions/resourceGroups/read",
        "Microsoft.Resources/subscriptions/resourcegroups/resources/read",
        "Microsoft.Resources/subscriptions/resourceGroups/write",
        "Microsoft.Storage/checknameavailability/read",
        "Microsoft.Storage/operations/read",
        "Microsoft.Storage/storageAccounts/listkeys/action",
        "Microsoft.Storage/storageAccounts/read",
        "Microsoft.Storage/storageAccounts/delete",
        "Microsoft.Storage/storageAccounts/regeneratekey/action",
        "Microsoft.Storage/storageAccounts/write",
        "Microsoft.Storage/usages/read",
        "Microsoft.Compute/snapshots/write",
        "Microsoft.Compute/snapshots/read",
        "Microsoft.Compute/availabilitySets/write",
        "Microsoft.Compute/availabilitySets/read",
        "Microsoft.Compute/disks/beginGetAccess/action",
        "Microsoft.MarketplaceOrdering/offertypes/publishers/offers/plans/agreements/read",
        "Microsoft.MarketplaceOrdering/offertypes/publishers/offers/plans/agreements/write",
        "Microsoft.Network/loadBalancers/read",
        "Microsoft.Network/loadBalancers/write",
        "Microsoft.Network/loadBalancers/delete",
        "Microsoft.Network/loadBalancers/backendAddressPools/read",
        "Microsoft.Network/loadBalancers/backendAddressPools/join/action",
        "Microsoft.Network/loadBalancers/frontendIPConfigurations/read",
        "Microsoft.Network/loadBalancers/loadBalancingRules/read",
        "Microsoft.Network/loadBalancers/probes/read",
        "Microsoft.Network/loadBalancers/probes/join/action",
        "Microsoft.Authorization/locks/*",
        "Microsoft.Network/routeTables/join/action",
        "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/write",
        "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/read",
        "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/delete",
        "Microsoft.NetApp/netAppAccounts/write",
        "Microsoft.NetApp/netAppAccounts/read",
        "Microsoft.NetApp/netAppAccounts/capacityPools/write",
        "Microsoft.NetApp/netAppAccounts/capacityPools/read",
        "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/delete",
        "Microsoft.Network/privateEndpoints/write",
        "Microsoft.Storage/storageAccounts/PrivateEndpointConnectionsApproval/action",
        "Microsoft.Storage/storageAccounts/privateEndpointConnections/read",
        "Microsoft.Network/privateEndpoints/read",
        "Microsoft.Network/privateDnsZones/write",
        "Microsoft.Network/privateDnsZones/virtualNetworkLinks/write",
        "Microsoft.Network/virtualNetworks/join/action",
        "Microsoft.Network/privateDnsZones/A/write",
        "Microsoft.Network/privateDnsZones/read",
        "Microsoft.Network/privateDnsZones/virtualNetworkLinks/read",
        "Microsoft.Resources/deployments/operationStatuses/read",
        "Microsoft.Insights/Metrics/Read",
        "Microsoft.Compute/virtualMachines/extensions/write",
        "Microsoft.Compute/virtualMachines/extensions/read",
        "Microsoft.Compute/virtualMachines/extensions/delete",
        "Microsoft.Compute/virtualMachines/delete",
        "Microsoft.Network/networkInterfaces/delete",
        "Microsoft.Network/networkSecurityGroups/delete",
        "Microsoft.Resources/deployments/delete",
        "Microsoft.Compute/diskEncryptionSets/read"]
    not_actions = []
  }

  assignable_scopes = [
    data.azurerm_subscription.primary.id, # /subscriptions/00000000-0000-0000-0000-000000000000
  ]
}

resource "azurerm_role_assignment" "occm-role-assignment" {
  depends_on = [azurerm_role_definition.occm-role]
  scope              = data.azurerm_subscription.primary.id
  role_definition_id = azurerm_role_definition.occm-role.role_definition_resource_id 
  principal_id       = data.azurerm_virtual_machine.occm-vm.identity.0.principal_id
}

resource "netapp-cloudmanager_cvo_azure" "cvo-azure" {
  depends_on = [azurerm_role_assignment.occm-role-assignment]
  provider = netapp-cloudmanager
  name = "TerraformCVOAzure"
  location = "westus"
  subscription_id = data.azurerm_subscription.primary.subscription_id
  subnet_id = "Subnet1"
  vnet_id = "Vnet1"
  vnet_resource_group = "rg_westus"
  data_encryption_type = "AZURE"
  azure_tag {
              tag_key = "abcd"
              tag_value = "ABCD"
            }
  azure_tag {
              tag_key = "xxx"
              tag_value = "YYY"
            }
  storage_type = "Premium_LRS"
  svm_password = "********"
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
  workspace_id = "workspace-xxxxxx"
  capacity_tier = "Blob"
  writing_speed_state = "NORMAL"
  is_ha = false
}

resource "netapp-cloudmanager_cvo_azure" "cvo-azure-ha" {
  depends_on = [azurerm_role_assignment.occm-role-assignment]
  provider = netapp-cloudmanager
  name = "TerraformCVOAzureHA"
  location = "westus"
  subscription_id = data.azurerm_subscription.primary.subscription_id
  subnet_id = "Subnet1"
  vnet_id = "Vnet1"
  vnet_resource_group = "rg_westus"
  data_encryption_type = "AZURE"
  azure_tag {
              tag_key = "abcd"
              tag_value = "ABCD"
            }
  azure_tag {
              tag_key = "xxx"
              tag_value = "YYY"
            }
  storage_type = "Premium_LRS"
  svm_password = "********"
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
  workspace_id = "workspace-xxxxxx"
  capacity_tier = "Blob"
  writing_speed_state = "NORMAL"
  is_ha = true
  license_type = "azure-ha-cot-standard-paygo"
}
