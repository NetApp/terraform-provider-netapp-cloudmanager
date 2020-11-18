---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_volume"
sidebar_current: "docs-netapp-cloudmanager-datasource-volume"
description: |-
  Provides a netapp-cloudmanager_volume resource. This can be used to get volumes for Cloud Volumes ONTAP.
---

# netapp-cloudmanager_volume

Provides a netapp-cloudmanager_volume resource. This can be used to get volumes for Cloud Volumes ONTAP.
Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.
NFS, CIFS, and iSCSI volumes are supported.

## Example Usages

**get netapp-cloudmanager_volume:**

```
data "netapp-cloudmanager_volume" "volume-nfs" {
  provider = netapp-cloudmanager
  name = "vol1"
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `working_environment_id` - (Optional) The public ID of the working environment where the volume exists. The ID can be optional if working_environment_name is provided. You can find the ID from the previous create Cloud Volumes ONTAP action as shown in the example, or from the Information page of the Cloud Volumes ONTAP working environment on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `working_environment_name` - (Optional) The working environment name where the volume exists. It will be ignored if working_environment_id is provided.