---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_snapmirror"
sidebar_current: "docs-netapp-cloudmanager-resource-snapmirror"
description: |-
  Provides a netapp-cloudmanager_snapmirror resource. This can be used to create a new snapmirror relationship from any CVO to any CVO, any CVO to ONPREM, ONPREM to any CVO, CVO to FSX. Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.
---

# netapp-cloudmanager_snapmirror

Provides a netapp-cloudmanager_snapmirror resource. This can be used to create a new snapmirror relationship from any CVO to any CVO, any CVO to ONPREM, ONPREM to any CVO, CVO to FSX. Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.

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

**Create netapp-cloudmanager_snapmirror with automatic destination volume deletion:**

```
resource "netapp-cloudmanager_snapmirror" "cl-snapmirror-with-volume-deletion" {
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
  delete_destination_volume = true
  client_id = "xxxxxxxxxxx"
}
```

## Argument Reference

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `source_working_environment_id` - (Optional) The public ID of the source working environment where the snapmirror relationship will be created.
* `destination_working_environment_id` - (Optional) The public ID of the destination working environment where the snapmirror relationship will be created.
* `source_working_environment_name` - (Optional) The source working environment name where the snapmirror relationship will be created. It will be ignored if working_environment_id is provided.
* `destination_working_environment_name` - (Optional) The destination working environment name where the snapmirror relationship will be created. It will be ignored if working_environment_id is provided.
* `source_svm_name` - (Optional, Computed) The name of the source SVM. The default SVM name is used, if a name isn't provided.
* `destination_svm_name` - (Optional, Computed) The name of the destination SVM. The default SVM name is used, if a name isn't provided.
* `source_volume_name` - (Required) The name of the source volume.
* `destination_volume_name` - (Required) The name of the destination volume to be created for snapmirror relationship.
* `connector_ip` - (Optional) The private IP of the connector, this is only required for Restricted mode account.
* `tenant_id` - (Optional, Forces new resource) The NetApp tenant ID that the Connector will be associated with. To be used in FSX or when `deployment_mode` is `Restricted`.  You can find the tenant ID in the Identity & Access Management in Settings, Organization tab of BlueXP at [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `deployment_mode` - (Optional) The mode of deployment to use for the working environment: ['Standard', 'Restricted']. The default is 'Standard'. To know more on deployment modes [https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/](https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/).
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `policy` - (Optional) The SnapMirror policy name. The default is 'MirrorAllSnapshots'.
* `schedule` - (Optional) Schedule name. The default is '1hour'.
* `max_transfer_rate` - (Required) Maximum transfer rate limit (KB/s). Use 0 for no limit, otherwise use number between 1024 and 2,147,482,624.  The default is 100000.
* `destination_aggregate_name` - (Optional) The aggregate in which the volume will be created. If not provided, Cloud Manager chooses the best aggregate for you.
* `provider_volume_type` - (Optional) The underlying cloud provider volume type. For AWS: ['gp3', 'gp2', 'io1', 'st1', 'sc1']. For Azure: ['Premium_LRS','Standard_LRS','StandardSSD_LRS']. For GCP: ['pd-balanced', 'pd-ssd','pd-standard']
* `capacity_tier` - (Optional) The volume's capacity tier for tiering cold data to object storage: ['S3', 'Blob', 'cloudStorage']. The default values for each cloud provider are as follows: Amazon => 'S3', Azure => 'Blob', GCP => 'cloudStorage'. If none, the capacity tier won't be set on volume creation.
* `delete_destination_volume` - (Optional) Set to true to delete the destination volume when the snapmirror relationship is destroyed. The default is false.
* `iops` - (Optional, Forces new resource) The number of IOPS to provision for the volume.
* `throughput` - (Optional, Forces new resource) The throughput to provision for the volume.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - will be the snapmirror name.


## Import

This resource supports import, which allows you to import existing snapmirror relationships into the state of this resource.

### Standard Mode

Import requires `deployment_mode`, `client_id` and `destination_volume_name` separated by commas.

id = `deployment_mode,client_id,destination_volume_name`

### Restricted Mode

Import requires `deployment_mode`, `client_id`, `destination_volume_name`, `tenant_id` and `connector_ip` separated by commas.

id = `deployment_mode,client_id,destination_volume_name,tenant_id,connector_ip`

#### Terraform Import

For example:

```shell
terraform import netapp-cloudmanager_snapmirror.example Standard,xxxxxxx,dest_volume_copy
```

For Restricted mode:

```shell
terraform import netapp-cloudmanager_snapmirror.example Restricted,xxxxxxx,dest_volume_copy,account-xxxxx,10.10.10.10
```

> The terraform import CLI command can only import resources into the state. Importing via the CLI does not generate configuration. If you want to generate the accompanying configuration for imported resources, use the import block instead.

#### Terraform Import Block

This requires Terraform 1.5 or higher, and will auto create the configuration for you.

First create the block:

```terraform
import {
  to = netapp-cloudmanager_snapmirror.snapmirror_import
  id = "Standard,xxxxxxx,dest_volume_copy"
}
```

Next run, this will auto create the configuration for you:

```shell
terraform plan -generate-config-out=generated.tf
```

This will generate a file called `generated.tf`, which will contain the configuration for the imported resource:

```terraform
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "Standard,xxxxxxx,dest_volume_copy"
resource "netapp-cloudmanager_snapmirror" "snapmirror_import" {
  client_id                          = "xxxxxxx"
  destination_volume_name            = "dest_volume_copy"
  destination_working_environment_id = "VsaWorkingEnvironment-xxxxxxx"
  max_transfer_rate                  = 100000
  policy                            = "MirrorAllSnapshots"
  schedule                          = "1hour"
  source_volume_name                = "source_volume"
  source_working_environment_id     = "VsaWorkingEnvironment-xxxxxxx"
  deployment_mode                   = "Standard"
}
```

