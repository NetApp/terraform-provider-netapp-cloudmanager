---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cvo_azure"
sidebar_current: "docs-netapp-cloudmanager-resource-cvo-azure"
description: |-
  Provides a netapp-cloudmanager_cvo_azure resource. This can be used to create a new Cloud Volume ONTAP system in Azure (single node or HA pair).
---

# netapp-cloudmanager_cvo_azure

Provides a netapp-cloudmanager_cvo_azure resource. This can be used to create a new Cloud Volume ONTAP on Azure (Single or HA).
Requires existence of a Cloud Manager Connector with a role assigned to create Cloud Volumes ONTAP. 'azurerm' provider can be used to create the role and role assignment.

## Example Usages

**Create netapp-cloudmanager_cvo_azure single:**

```
resource "netapp-cloudmanager_cvo_azure" "cl-azure" {
  depends_on = [azurerm_role_assignment.occm-role-assignment]
  provider = netapp-cloudmanager
  name = "TerraformCVOAzure"
  location = "westus"
  availability_zone = 2
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
  svm_password = "P@assword!"
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
  workspace_id = "workspace-fdgsgNse"
  capacity_tier = "Blob"
  writing_speed_state = "NORMAL"
  is_ha = false
}
```

**Create netapp-cloudmanager_cvo_azure HA:**

```
resource "netapp-cloudmanager_cvo_azure" "cl-azure" {
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
  svm_password = "P@assword!"
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
  workspace_id = "workspace-fdgsgNse"
  capacity_tier = "Blob"
  is_ha = true
  license_type = "azure-ha-cot-standard-paygo"
}
```

**Create netapp-cloudmanager_cvo_azure single with WORM:**

```
resource "netapp-cloudmanager_cvo_azure" "cl-azure" {
  depends_on = [azurerm_role_assignment.occm-role-assignment]
  provider = netapp-cloudmanager
  name = "TerraformCVOAzure"
  location = "westus"
  availability_zone = 2
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
  svm_password = "P@assword!"
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
  workspace_id = "workspace-fdgsgNse"
  writing_speed_state = "NORMAL"
  is_ha = false
  worm_retention_period_length = 2
  worm_retention_period_unit = "days"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Volumes ONTAP working environment.
* `location` - (Required) The location where the working environment will be created.
* `availability_zone` - (Optional) The availability zone on the location configuration.
* `subscription_id` - (Required) The ID of the Azure subscription.
* `subnet_id` - (Required) The name of the subnet for the Cloud Volumes ONTAP system.
* `vnet_id` - (Required) The name of the virtual network.
* `vnet_resource_group` - (Required) The resource group in Azure associated to the virtual network.
* `workspace_id` - (Optional) The ID of the Cloud Manager workspace where you want to deploy Cloud Volumes ONTAP. If not provided, Cloud Manager uses the first workspace. You can find the ID from the Workspace tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `data_encryption_type` - (Optional) The type of encryption to use for the working environment: ['AZURE', 'NONE']. The default is 'AZURE'.
* `azure_encryption_parameters` - (Optional) AZURE encryption parameters. It is required if using AZURE encryption.
* `storage_type` - (Optional) The type of storage for the first data aggregate: ['Premium_LRS', 'Standard_LRS', 'StandardSSD_LRS', 'Premium_ZRS']. The default is 'Premium_LRS'
* `svm_password` - (Required) The admin password for Cloud Volumes ONTAP.
* `svm_name` - (Optional) The name of the SVM.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `resource_group` - (Optional) The resource_group where Cloud Volumes ONTAP will be created. If not provided, Cloud Manager creates the resource group (name of the working environment with suffix '-rg').
* `allow_deploy_in_existing_rg` - (Optional) Indicates if to allow creation in existing resource group, Default is false.
* `cidr` - (Optional) The CIDR of the VNET. If not provided, resource needs az login to authorize and fetch the cidr details from Azure.
* `disk_size` - (Optional) Azure volume size for the first data aggregate. For GB, the unit can be: [100 or 500]. For TB, the unit can be: [1,2,4,8,16]. The default is '1' .
* `disk_size_unit` - (Optional) ['GB' or 'TB']. The default is 'TB'.
* `ontap_version` - (Optional) The required ONTAP version. Ignored if 'use_latest_version' is set to true. The default is to use the latest version.
* `use_latest_version` - (Optional) Indicates whether to use the latest available ONTAP version. The default is 'true'.
* `license_type` - (Optional) The type of license to be use. For single node: (by Capacity): ['capacity-paygo'], (by Node paygo): ['azure-cot-explore-paygo', 'azure-cot-standard-paygo', 'azure-cot-premium-paygo'], (by Node byol): ['azure-cot-premium-byol']. For HA: (by Capacity): ['ha-capacity-paygo'], (by Node paygo): ['azure-ha-cot-standard-paygo', 'azure-ha-cot-premium-paygo'], (by Node byol): ['azure-ha-cot-premium-byol']. The default is 'capacity-paygo' for single node, and 'ha-capacity-paygo'for HA.
* `capacity_package_name` - (Optional) The capacity package name: ['Essential', 'Professional', 'Freemium']. Default is 'Essential'.
* `instance_type` - (Optional) The type of instance to use, which depends on the license type you chose: Explore:['Standard_DS3_v2'], Standard:['Standard_DS4_v2,Standard_DS13_v2,Standard_L8s_v2'], Premium:['Standard_DS5_v2','Standard_DS14_v2'], BYOL: all instance types defined for PayGo. For more supported instance types, refer to Cloud Volumes ONTAP Release Notes. The default is 'Standard_DS4_v2' .
* `serial_number` - (Optional) The serial number for the cluster. Required when using one of these: ['azure-cot-premium-byol' or 'azure-ha-cot-premium-byol'].
* `capacity_tier` - (Optional) Whether to enable data tiering for the first data aggregate: ['Blob', 'NONE']. The default is 'BLOB'.
* `tier_level` - (Optional) If capacity_tier is Blob, this argument indicates the tiering level: ['normal', 'cool']. The default is: 'normal'.
* `nss_account` - (Optional) The NetApp Support Site account ID to use with this Cloud Volumes ONTAP system. If the license type is BYOL and an NSS account isn't provided, Cloud Manager tries to use the first existing NSS account.
* `writing_speed_state` - (Optional) The write speed setting for Cloud Volumes ONTAP: ['NORMAL','HIGH']. The default is 'NORMAL'. This argument is not relevant for HA pairs.
* `security_group_id` - (Optional) The name of the security group (full identifier: /subscriptions/xxxxxx/resourceGroups/rg_westus/providers/Microsoft.Network/networkSecurityGroups/CVO-SG). If not provided, Cloud Manager creates the security group.
* `cloud_provider_account` - (Optional) The cloud provider credentials id to use when deploying the Cloud Volumes ONTAP system. You can find the ID in Cloud Manager from the Settings > Credentials page. If not specified, Cloud Manager uses the managed service identity of the Connector virtual machine.
* `backup_volumes_to_cbs` - (Optional) Automatically enable back up of all volumes to Azure Blob [true, false].
* `enable_compliance` - (Optional) Enable the Cloud Compliance service on the working environment [true, false].
* `enable_monitoring` - (Optional) Enable the Monitoring service on the working environment [true, false]. The default is false.
* `is_ha` - (Optional) Indicate whether the working environment is an HA pair or not [true, false]. The default is false.
* `platform_serial_number_node1` - (Optional) For HA BYOL, the serial number for the first node.
* `platform_serial_number_node2` - (Optional) For HA BYOL, the serial number for the second node.
* `availability_zone_node1` - (Optional) For HA, the availability zone for the first node.
* `availability_zone_node2` - (Optional) For HA, the availability zone for the second node.
* `ha_enable_https` - (Optional) For HA, enable the HTTPS connection from CVO to storage accounts. This can impact write performance. The default is false.
* `upgrade_ontap_version` - (Optional) Indicates whether to upgrade ontap image with `ontap_version`. To upgrade ontap image, `ontap_version` cannot be 'latest' and `use_latest_version` needs to be false.
* `retries` - (Optional) The number of attempts to wait for the completion of creating the CVO with 60 seconds apart for each attempt. For HA, this value is incremented by 30. The default is '60'.
* `worm_retention_period_length` - (Optional) WORM retention period length. Once specified retention period, the WORM is enabled. When WORM storage is activated, data tiering to object storage can’t be enabled.
* `worm_retention_period_unit` - (Optional) WORM retention period unit: ['years','months','days','hours','minutes','seconds'].

The `azure_tag` block supports:
* `tag_key` - (Required) The key of the tag.
* `tag_value` - (Required) The tag value.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the working environment.
* `svm_name` - The default name of the SVM will be exported if it is not provided in the resource.
