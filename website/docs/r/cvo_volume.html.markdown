---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_volume"
sidebar_current: "docs-netapp-cloudmanager-resource-volume"
description: |-
  Provides a netapp-cloudmanager_volume resource. This can be used to create, update, and delete volumes for Cloud Volumes ONTAP.
---

# netapp-cloudmanager_volume

Provides a netapp-cloudmanager_volume resource. This can be used to create, update, and delete volumes for Cloud Volumes ONTAP.
Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.
NFS, CIFS, and iSCSI volumes are supported.

## Example Usages

**Create netapp-cloudmanager_volume of type NFS:**

```
resource "netapp-cloudmanager_volume" "cvo-volume-nfs" {
  depends_on = [netapp-cloudmanager_volume.cifs-volume-1]
  provider = netapp-cloudmanager
  volume_protocol = "nfs"
  name = "vol1"
  size = 10
  unit = "GB"
  provider_volume_type = "pd-standard"
  export_policy_type = "custom"
  export_policy_ip = ["0.0.0.0/0"]
  export_policy_nfs_version = ["nfs4"]
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}
```

**Create netapp-cloudmanager_volume of type CIFS:**

```
resource "netapp-cloudmanager_volume" "cvo-volume-cifs" {
  depends_on = [netapp-cloudmanager_cifs_server.cvo-cifs-workgroup]
  provider = netapp-cloudmanager
  name = "cifs_vol2"
  volume_protocol = "cifs"
  provider_volume_type = "pd-ssd"
  size = 10
  unit = "GB"
  share_name = "share_cifs"
  permission = "full_control"
  users = ["Everyone"]
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}
```

**Create netapp-cloudmanager_volume of type ISCSI:**

```
resource "netapp-cloudmanager_volume" "cvo-volume-iscsi" {
  provider = netapp-cloudmanager
  name = "iscsi_test_vol"
  volume_protocol = "iscsi"
  size = 10
  unit = "GB"
  igroups = ["test_igroup"]
  initiator {
    alias = "test_alias"
    iqn = "test_iqn"
  }
  os_name = "linux"
  working_environment_name = "cvo-name"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `svm_name` - (Optional) The name of the SVM. The default SVM name is used, if a name isn't provided.
* `size` - (Required) The volume size, supported with decimal numbers.
* `size_unit` - (Required) ['Byte' or 'KB' or 'MB' or 'GB' or 'TB'].
* `provider_volume_type` - (Required) The underlying cloud provider volume type. For AWS: ['gp3', 'gp2', 'io1', 'st1', 'sc1']. For Azure: ['Premium_LRS','Standard_LRS','StandardSSD_LRS']. For GCP: ['pd-balanced', 'pd-ssd','pd-standard']. For onPrem: 'onprem'.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `enable_thin_provisioning` - (Optional) Enable thin provisioning.
* `enable_compression` - (Optional) Enable compression.
* `enable_deduplication` - (Optional) Enable deduplication.
* `aggregate_name ` - (Optional) The aggregate in which the volume will be created. If not provided, Cloud Manager chooses the best aggregate for you.
* `volume_protocol` - (Optional) The protocol for the volume: ['nfs', 'cifs', 'iscsi']. This affects the provided parameters. The default is 'nfs'
* `working_environment_id` - (Optional) The public ID of the working environment where the volume will be created. The ID can be optional if working_environment_name is provided. You can find the ID from the previous create Cloud Volumes ONTAP action as shown in the example, or from the Information page of the Cloud Volumes ONTAP working environment on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `working_environment_name` - (Optional) The working environment name where the aggregate will be created. It will be ignored if working_environment_id is provided.
* `capacity_tier` - (Optional) The volume's capacity tier for tiering cold data to object storage: ['S3', 'Blob', 'cloudStorage']. The default values for each cloud provider are as follows: Amazon => 'S3', Azure => 'Blob', GCP => 'cloudStorage'. If none, the capacity tier won't be set on volume creation.
* `export_policy_name` - (Optional) The export policy name. (NFS protocol parameters)
* `export__policy_type` - (Optional) The export policy type. (NFS protocol parameters)
* `export_policy_ip` - (Optional) Custom export policy list of IPs. (NFS protocol parameters)
* `export_policy_nfs_version` - (Optional) Export policy protocol. (NFS protocol parameters)
* `snapshot_policy_name` - (Optional) Snapshot policy name. The default is 'default'. (NFS protocol parameters)
* `iops` - (Optional) Provisioned IOPS. Needed only when 'provider_volume_type' is 'io1' or 'gp3'
* `throughput` - (Optional) Required only when 'provider_volume_type' is 'gp3'.
* `share_name` (Optional) Share name. (CIFS protocol parameters)
* `permission` (Optional) CIFS share permission type. (CIFS protocol parameters)
* `users` (Optional) List of users with the permission. (CIFS protocol parameters)
* `igroups` (Optional) List of igroups. (iSCSI protocol parameters)
* `os_name` (Optional) Operating system. (iSCSI protocol parameters)
* `initiator` (Optional) Set of attributes of Initiator. (iSCSI protocol parameters)

The `initiator` block supports:
* `alias` (Required) Initiator alias. (iSCSI protocol parameters)
*  `iqn` (Required) Initiator IQN. (iSCSI protocol parameters)

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the volume.

