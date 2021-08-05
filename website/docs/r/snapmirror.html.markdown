---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_snapmirror"
sidebar_current: "docs-netapp-cloudmanager-resource-snapmirror"
description: |-
  Provides a netapp-cloudmanager_snapmirror resource. This can be used to create a new snapmirror relationship on Cloud Volumes ONTAP.
---

# netapp-cloudmanager_snapmirror

Provides a netapp-cloudmanager_snapmirror resource. This can be used to create a new snapmirror relationship on Cloud Volumes ONTAP.
Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.

## Example Usages

**Create netapp-cloudmanager_snapmirror:**

```
resource "netapp-cloudmanager_snapmirror" "cl-snapmirror" {
  provider = netapp-cloudmanager
  source_working_environment_id = "xxxxxxxx"
  destination_working_environment_id = "xxxxxxxx"
  source_volume_name = "source"
  source_svm_name = "svm_source"
  destination_volume_name = "source_copy"
  destination_svm_name = "svm_dest"
  policy = "MirrorAllSnapshots"
  schedule = "5min"
  destination_aggregate_name = "aggr1"
  max_transfer_rate = "102400"
  client_id = "xxxxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `source_working_environment_id` - (Optional) The public ID of the source working environment where the snapmirror relationship will be created.
* `destination_working_environment_id` - (Optional) The public ID of the destination working environment where the snapmirror relationship will be created.
* `source_working_environment_name` - (Optional) The source working environment name where the snapmirror relationship will be created. It will be ignored if working_environment_id is provided.
* `destination_working_environment_name` - (Optional) The destination working environment name where the snapmirror relationship will be created. It will be ignored if working_environment_id is provided.
* `source_svm_name` - (Optional) The name of the source SVM. The default SVM name is used, if a name isn't provided.
* `destination_svm_name` - (Optional) The name of the destination SVM. The default SVM name is used, if a name isn't provided.
* `source_volume_name` - (Required) The name of the source volume.
* `destination_volume_name` - (Required) The name of the destination volume to be created for snapmirror relationship.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `policy` - (Optional) The SnapMirror policy name. The default is 'MirrorAllSnapshots'.
* `schedule` - (Optional) Schedule name. The default is '1hour'.
* `max_transfer_rate` - (Required) Maximum transfer rate limit (KB/s). Use 0 for no limit, otherwise use number between 1024 and 2,147,482,624.  The default is 100000.
* `destination_aggregate_name` - (Optional) The aggregate in which the volume will be created. If not provided, Cloud Manager chooses the best aggregate for you.
* `provider_volume_type` - (Optional) The underlying cloud provider volume type. For AWS: ['gp3', 'gp2', 'io1', 'st1', 'sc1']. For Azure: ['Premium_LRS','Standard_LRS','StandardSSD_LRS']. For GCP: ['pd-balanced', 'pd-ssd','pd-standard']
* `capacity_tier` - (Optional) The volume's capacity tier for tiering cold data to object storage: ['S3', 'Blob', 'cloudStorage']. The default values for each cloud provider are as follows: Amazon => 'S3', Azure => 'Blob', GCP => 'cloudStorage'. If none, the capacity tier won't be set on volume creation.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - will be the snapmirror name.

