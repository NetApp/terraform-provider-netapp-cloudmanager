---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_anf_volume"
sidebar_current: "docs-netapp-cloudmanager-resource-anf-volume"
description: |-
  Provides a netapp-cloudmanager_anf_volume resource. This can be used to create, and delete volumes for Azure NetApp Files.
---

# netapp-cloudmanager_anf_volume

Provides a netapp-cloudmanager_anf_volume resource. This can be used to create, and delete volumes for Azure NetApp Files.
Requires existence of a Cloud Manager Connector.

## Example Usages

**Create netapp-cloudmanager_volume:**

```
resource "netapp-cloudmanager_anf_volume" "test-1" {
  provider = netapp-cloudmanager
  name = "test_vol"
  size = 105
  size_unit = "gb"
  volume_path = "volume-path"
  protocol_types = ["NFSv3"]
  location = "eastus"
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
  service_level = "Standard"
  subnet = "default"
  virtual_network = "mynetwork"
  working_environment_name = "ANF_environment"
  account = "Demo_SIM"
  netapp_account = "test"
  subscription = "My Subscription"
  resource_groups = "myRG-eastus"
  capacity_pool = "ANFPool"
  rules { 
        rule {
            allowed_clients = "1.0.0.1"
            rule_index = 1
            nfsv3 = true
            unix_read_only = true
            }
        rule {
            allowed_clients = "1.0.0.2"
            rule_index = 2
            nfsv3 = true
            unix_read_only = true
            unix_read_write = false
        }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `size` - (Required) The volume size, supported with decimal numbers.
* `size_unit` - (Required) [ 'GB' ].
* `volume_path` - (Required) The volume path.
* `protocol_types` (Required) [ 'NFSv3' ].
* `location` - (Required) The location of the account.
* `service_level` - (Required) ['Premium' or 'Standard' or 'Ultra'].
* `subnet` - (Required) The name of the subnet.
* `virtual_network`  - (Required) The name of the virtual network.
* `account` - (Required) The name of the account.
* `netapp_account` - (Required) The name of the netapp account.
* `subscription`  - (Required) The name of the subscription.
* `resource_groups` - (Required) The name of the resource group in Azure where the volume will be created.
* `capacity_pool` - (Required) The name of the capacity pool.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Required) The working environment name.
* `export_policy` - (Optional) The rules of the export policy.


The `export_policy` block supports:
* `rule` - (Optional) The rule of the export policy.

The `rule` block supports:
* `allowed_clients` - (Optional) allowed clients.
* `rule_index` - (Optional) rule index.
* `nfsv3` - (Optional) Boolean.
* `unix_read_only` - (Optional) Boolean.
* `unix_read_write` - (Optional) Boolean.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the volume.

