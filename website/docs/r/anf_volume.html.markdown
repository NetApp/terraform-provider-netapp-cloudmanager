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

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the volume.
* `size` - (Required, Forces new resource) The volume size, supported with decimal numbers.
* `size_unit` - (Required, Forces new resource) [ 'GB' ].
* `volume_path` - (Required, Forces new resource) The volume path.
* `protocol_types` (Required, Forces new resource) [ 'NFSv3' ].
* `location` - (Required, Forces new resource) The location of the account.
* `service_level` - (Required, Forces new resource) ['Premium' or 'Standard' or 'Ultra'].
* `subnet` - (Required, Forces new resource) The name of the subnet.
* `virtual_network`  - (Required, Forces new resource) The name of the virtual network.
* `account` - (Required, Forces new resource) The name of the account.
* `netapp_account` - (Required, Forces new resource) The name of the netapp account.
* `subscription`  - (Required, Forces new resource) The name of the subscription.
* `resource_groups` - (Required, Forces new resource) The name of the resource group in Azure where the volume will be created.
* `capacity_pool` - (Required, Forces new resource) The name of the capacity pool.
* `client_id` - (Required, Forces new resource) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Required, Forces new resource) The working environment name.
* `export_policy` - (Optional, Forces new resource) The rules of the export policy.

The `export_policy` block supports:

* `rule` - (Optional, Forces new resource) The rule of the export policy.

The `rule` block supports:

* `allowed_clients` - (Optional, Forces new resource) allowed clients.
* `rule_index` - (Optional) rule index.
* `nfsv3` - (Optional, Forces new resource) Boolean.
* `unix_read_only` - (Optional, Forces new resource) Boolean.
* `unix_read_write` - (Optional) Boolean.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the volume.

