---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cvo_onprem"
sidebar_current: "docs-netapp-cloudmanager-resource-cvo-onprem"
description: |-
  Provides a netapp-cloudmanager_cvo_onprem resource. This can be used to register an onprem ONTAP system into CloudManager.
---

# netapp-cloudmanager_cvo_onprem

Provides a netapp-cloudmanager_cvo_onprem resource. This can be used to register an onprem ONTAP system into CloudManager.

## Example Usages

**Create netapp-cloudmanager_cvo_onprem:**

```
resource "netapp-cloudmanager_cvo_onprem" "cvo-onprem" {
  provider = netapp-cloudmanager
  name = "onprem"
  cluster_address = "10.10.10.10"
  cluster_user_name = "admin"
  cluster_password = "netapp1!"
  client_id = "xxxxxxxxx"
  location = "ON_PREM"
}
```


## Argument Reference

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the Cloud Volumes ONTAP working environment.
* `cluster_address` - (Required, Forces new resource) The ip address of the cluster management interface.
* `cluster_user_name` - (Required, Forces new resource) The admin user name for the onprem ONTAP system.
* `cluster_password` - (Required, Forces new resource) The admin password for the onprem ONTAP system.
* `client_id` - (Required, Forces new resource) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `workspace_id` - (Optional, Forces new resource) The ID of the Cloud Manager workspace where you want to deploy Cloud Volumes ONTAP. If not provided, Cloud Manager uses the first workspace. You can find the ID from the Workspace tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `location` - (Required, Forces new resource) The type of location to use for the working environment: ['ON_PREM', 'AZURE', 'AWS', 'SOFTLAYER', 'GOOGLE', 'CLOUD_TIERING'].

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the working environment.