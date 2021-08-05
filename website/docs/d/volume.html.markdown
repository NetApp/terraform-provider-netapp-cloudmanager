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

* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `name` - (Required) The name of the volume.
* `working_environment_id` - (Optional) The public ID of the working environment where the volume exists. The ID can be optional if working_environment_name is provided. You can find the ID from the previous create Cloud Volumes ONTAP action as shown in the example, or from the Information page of the Cloud Volumes ONTAP working environment on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `working_environment_name` - (Optional) The working environment name where the volume exists. It will be ignored if working_environment_id is provided.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `svm_name` - The name of the SVM.
* `size` - The volume size, supported with decimal numbers.
* `size_unit` - ['Byte' or 'KB' or 'MB' or 'GB' or 'TB'].
* `provider_volume_type` - The underlying cloud provider volume type. For AWS: ['gp3', 'gp2', 'io1', 'st1', 'sc1']. For Azure: ['Premium_LRS','Standard_LRS','StandardSSD_LRS']. For GCP: ['pd-balanced', 'pd-ssd','pd-standard']
* `enable_thin_provisioning` - Enable thin provisioning. The default is 'true'.
* `enable_compression` - Enable compression. The default is 'true'.
* `enable_deduplication` - Enable deduplication. The default is 'true'.
* `aggregate_name ` - The aggregate in which the volume will be created.
* `volume_protocol` - The protocol for the volume: ["nfs", "cifs", "iscsi"]. The default is 'nfs'
* `capacity_tier` - The volume's capacity tier for tiering cold data to object storage: ['S3', 'Blob', 'cloudStorage']. The default values for each cloud provider are as follows: Amazon => 'S3', Azure => 'Blob', GCP => 'cloudStorage'.
* `mount_point` The mount point.
* `export_policy_name` - The export policy name. (NFS protocol parameters)
* `export_policy_type` - The export policy type. (NFS protocol parameters)
* `export_policy_ip` - Custom export policy list of IPs. (NFS protocol parameters)
* `export_policy_nfs_version` - Export policy protocol. (NFS protocol parameters)
* `snapshot_policy_name` - Snapshot policy name. The default is 'default'. (NFS protocol parameters)
* `share_name` Share name. (CIFS protocol parameters)
* `permission` CIFS share permission type. (CIFS protocol parameters)
* `users` List of users with the permission. (CIFS protocol parameters)
