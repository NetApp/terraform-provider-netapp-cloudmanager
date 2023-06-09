---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_aws_fsx_volume"
sidebar_current: "docs-netapp-cloudmanager-resource-aws-fsx-volume"
description: |-
  Provides a netapp-cloudmanager_aws-fsx-volume resource. This can be used to create, update, and delete volumes for Amazon FSx ONTAP.
---

# netapp-cloudmanager_fsx_aws_volume

Requires a Amazon FSx ONTAP system.
NFS and CIFS volumes are supported.

## Example Usages

**Create netapp-cloudmanager_aws_fsx_volume of type NFS:**

```
resource "netapp-cloudmanager_aws_fsx_volume" "fsx-volume-nfs" {
  provider = netapp-cloudmanager
  volume_protocol = "nfs"
  name = "vol1"
  size = 10
  unit = "GB"
  export_policy_type = "custom"
  export_policy_ip = ["0.0.0.0/0"]
  export_policy_nfs_version = ["nfs4"]
  file_system_id = "your-file-system-id"
  tags = {
    "abc" = "xyz"
  }
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}
```

**Create netapp-cloudmanager_aws_fsx_volume of type CIFS:**

```
resource "netapp-cloudmanager_aws_fsx_volume" "fsx-volume-cifs" {
  provider = netapp-cloudmanager
  name = "cifs_vol"
  volume_protocol = "cifs"
  size = 10
  unit = "GB"
  share_name = "share_cifs"
  permission = "full_control"
  users = ["Everyone"]
  file_system_id = "your-file-system-id"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `svm_name` - (Optional) The name of the SVM. The default SVM name is used, if a name isn't provided.
* `size` - (Required) The volume size, supported with decimal numbers.
* `size_unit` - (Required) ['Byte' or 'KB' or 'MB' or 'GB' or 'TB'].
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous created Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `enable_storage_efficiency` - (Optional) Enable storage efficiency.
* `export__policy_type` - (Optional) The export policy type. (NFS protocol parameters)
* `export_policy_ip` - (Optional) Custom export policy list of IPs. (NFS protocol parameters)
* `export_policy_nfs_version` - (Optional) Export policy protocol. (NFS protocol parameters)
* `snapshot_policy_name` - (Optional) Snapshot policy name. The default is 'default'. (NFS protocol parameters)
* `share_name` (Optional) Share name. (CIFS protocol parameters)
* `permission` (Optional) CIFS share permission type. (CIFS protocol parameters)
* `users` (Optional) List of users with the permission. (CIFS protocol parameters)
* `volume_protocol` - (Required) The protocol for the volume: ['nfs', 'cifs']. This affects the provided parameters.
*  `tags` - (Optional) Set tags for the volume during creation. The API doesn't contain any information about tags so the provider doesn't guarantee tags will be added successfully and detect any drift after create.
* `tenant_id` - (Required) The workspace id.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The uuid of the volume.

