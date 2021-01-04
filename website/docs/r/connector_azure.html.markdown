---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_connector_azure"
sidebar_current: "docs-netapp-cloudmanager-resource-connector-azure"
description: |-
  Provides a netapp-cloudmanager_connector_azure resource. This can be used to create a new Cloud Manager Connector in Azure.
---

# netapp-cloudmanager_connector_azure

Provides a netapp-cloudmanager_connector_azure resource. This can be used to create a new Cloud Manager Connector in AZURE.
The environment needs to be configured with the proper credentials before it can be used (az login).
The minimum required policy can be found at [Connector deployment policy for Azure](https://s3.amazonaws.com/occm-sample-policies/Policy_for_Setup_As_Service_Azure.json)

In order for the Connector to create a Cloud Volumes ONTAP system, it requires a role assignment. This can be done with azurerm provider. The following role is required: [Cloud Manager policy for Azure](https://occm-sample-policies.s3.amazonaws.com/Policy_for_cloud_Manager_Azure_3.8.7.json)


<!---
i think we need to create section for terraform and point to there
-->

## Example Usages

**Create netapp-cloudmanager_connector_azure:**

```
resource "netapp-cloudmanager_connector_azure" "cl-occm-azure" {
  provider = netapp-cloudmanager
  name = "TF-ConnectorAzure"
  location = "westus"
  subscription_id = "xxxxxxxxxxxxxxxx"
  company = "NetApp"
  resource_group = "rg_westus"
  subnet_id = "Subnet1"
  vnet_id = "Vnet1"
  network_security_group_name = "OCCM_SG"
  associate_public_ip_address = true
  account_id = "account-ABCNJGB0X"
  admin_password = "P@ssword123456"
  admin_username = "vmadmin"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Manager Connector.
* `location` - (Required) The location where the Cloud Manager Connector will be created.
* `company` - (Required) The name of the company of the user.
* `resource_group` - (Required) The resource group in Azure where the resources will be created.
* `subnet_id` - (Required) The name of the subnet for the virtual machine.
* `subscription_id` - (Required) The ID of the Azure subscription.
* `vnet_id` - (Required) The name of the virtual network.
* `network_security_group_name` - (Required) The name of the security group for the instance.
* `admin_username` - (Required) The user name for the Connector.
* `admin_password` - (Required) The password for the Connector.
* `vnet_resource_group` - (Optional) The resource group in Azure associated with the virtual network. If not provided, it’s assumed that the VNet is within the previously specified resource group.
* `network_security_resource_group` - (Optional) The resource group in Azure associated with the security group. If not provided, it’s assumed that the security group is within the previously specified resource group.
* `virtual_machine_size` - (Optional) The virtual machine type. (for example, Standard_D2s_v3). At least 4 CPU and 16 GB of memory are required.
* `proxy_url` - (Optional) The proxy URL, if using a proxy to connect to the internet.
* `proxy_user_name` - (Optional) The proxy user name, if using a proxy to connect to the internet.
* `proxy_password` - (Optional) The proxy password, if using a proxy to connect to the internet.
* `proxy_certificates` - (Optional) The proxy certificates. A list of certificate file names.
* `associate_public_ip_address` - (Optional) Indicates whether to associate the public IP address to the virtual machine.
* `account_id` - (Optional) The NetApp account ID that the Connector will be associated with. If not provided, Cloud Manager uses the first account. If no account exists, Cloud Manager creates a new account. You can find the account ID in the account tab of Cloud Manager at [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the virtual machine.
* `client_id` - The unique client ID of the connector, can be used in other resources.
* `account_id` - The NetApp tenancy account ID.

