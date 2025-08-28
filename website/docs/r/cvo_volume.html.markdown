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

**Create netapp-cloudmanager_volume for restricted mode:**

```
 "netapp-cloudmanager_volume" "cvo-volume-restricted" {
  provider = netapp-cloudmanager
  name = "test_vol"
  size = 10
  unit = "GB"
  snapshot_policy_name = "default"
  working_environment_name = "tfgcprestricted"
  provider_volume_type = "pd-standard"
  export_policy_type = "custom"
  export_policy_ip = ["0.0.0.0/0"]
  export_policy_nfs_version = ["nfs3"]
  export_policy_rule_access_control = "readwrite"
  export_policy_rule_super_user = true
  comment = "test"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
  deployment_mode = "Restricted"
  connector_ip = "10.10.10.10"
  tenant_id = "account-22Vu41zs"
}
```

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
  export_policy_rule_access_control = "readwrite"
  export_policy_rule_super_user = true
  snapshot_policy_name = "sp1"
  snapshot_policy {
     schedule {
       schedule_type = "5min"
       retention = 10
    }
    schedule {
       schedule_type = "hourly"
       retention = 5
    }
  }
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

**Create netapp-cloudmanager_volume on OnPrem:**

```
resource "netapp-cloudmanager_volume" "cvo-volume-onprem" {
  provider = netapp-cloudmanager
  name = "onprem_test_vol"
  volume_protocol = "nfs"
  provider_volume_type = "onprem"
  size = 10
  unit = "GB"
  export_policy_type = "custom"
  export_policy_ip = ["0.0.0.0/0"]
   svm_name = "test_onprem"
  working_environment_name = "cvo-name"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
}
```

## Argument Reference

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `name` - (Required) The name of the volume.
* `svm_name` - (Optional) The name of the SVM. The default SVM name is used, if a name isn't provided.
* `size` - (Required) The volume size, supported with decimal numbers.
* `size_unit` - (Required) ['Byte' or 'KB' or 'MB' or 'GB' or 'TB'].
* `provider_volume_type` - (Required) The underlying cloud provider volume type. For AWS: ['gp3', 'gp2', 'io1', 'st1', 'sc1'] (ebs_volume_type on AWS CVO). For Azure: ['Premium_LRS','Standard_LRS','StandardSSD_LRS', 'Premium_ZRS'] (storage_type on Azure CVO). For GCP: ['pd-balanced', 'pd-ssd','pd-standard'] (gcp_volume_type on GCP CVO). For onPrem: 'onprem'.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `connector_ip` - (Optional) The private IP of the connector, this is only required for Restricted mode.
* `tenant_id` - (Optional) The NetApp tenant ID that the Connector will be associated with.  You can find the tenant ID in the Identity & Access Management in Settings, Organization tab of BlueXP at [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `deployment_mode` - (Optional) The mode of deployment to use for the working environment: ['Standard', 'Restricted']. The default is 'Standard'. To know more on deployment modes [https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/](https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/).
* `enable_thin_provisioning` - (Optional) Enable thin provisioning.
* `enable_compression` - (Optional) Enable compression.
* `enable_deduplication` - (Optional) Enable deduplication.
* `aggregate_name ` - (Optional, Computed) The aggregate in which the volume will be created. If not provided, Cloud Manager chooses the best aggregate for you. For OnPrem, aggregate input is required.
* `volume_protocol` - (Optional) The protocol for the volume: ['nfs', 'cifs', 'iscsi']. This affects the provided parameters. The default is 'nfs'.
* `working_environment_id` - (Optional) The public ID of the working environment where the volume will be created. The ID can be optional if working_environment_name is provided. You can find the ID from the previous create Cloud Volumes ONTAP action as shown in the example, or from the Information page of the Cloud Volumes ONTAP working environment on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Optional) The working environment name where the aggregate will be created. It will be ignored if working_environment_id is provided.
* `capacity_tier` - (Optional) The volume's capacity tier for tiering cold data to object storage: ['S3', 'Blob', 'cloudStorage']. The default values for each cloud provider are as follows: Amazon => 'S3', Azure => 'Blob', GCP => 'cloudStorage'. If none, the capacity tier won't be set on volume creation.
* `export_policy_name` - (Optional) The export policy name. (NFS protocol parameters)
* `export__policy_type` - (Optional) The export policy type. (NFS protocol parameters)
* `export_policy_ip` - (Optional) Custom export policy list of IPs. Order matters. (NFS protocol parameters)
* `export_policy_nfs_version` - (Optional) Export policy protocol. (NFS protocol parameters)
* `export_policy_rule_access_control` (Optional) Choice of 'readonly', 'readwrite', 'none'. (NFS protocol parameters) 
* `export_policy_rule_super_user` - (Optional) Boolean option to specify super user or not. (NFS protocol parameters)
  `export__policy_type`, `export_policy_ip`, `export_policy_nfs_version`, `export_policy_nfs_version` and  `export_policy_rule_super_user` are required together for export policy.
* `snapshot_policy_name` - (Optional) Snapshot policy name. The default is 'default'. (NFS protocol parameters)
* `iops` - (Optional) Provisioned IOPS. Needed only when 'provider_volume_type' is 'io1' or 'gp3'.
* `throughput` - (Optional) Required only when 'provider_volume_type' is 'gp3'.
* `share_name` (Optional) Share name. (CIFS protocol parameters)
* `permission` (Optional) CIFS share permission type. (CIFS protocol parameters)
* `users` (Optional) List of users with the permission. (CIFS protocol parameters)
* `igroups` (Optional) List of igroups. (iSCSI protocol parameters)
* `os_name` (Optional) Operating system. (iSCSI protocol parameters)
* `comment` - (Optional) Sets a comment associated with the volume. 
* `initiator` (Optional) Set of attributes of Initiator. (iSCSI protocol parameters)
*  `tags` - (Optional) Set tags for the volume during creation. The API doesn't contain any information about tags so the provider doesn't guarantee tags will be added successfully and detect any drift after create.

The `initiator` block supports:
* `alias` (Required) Initiator alias. (iSCSI protocol parameters)
*  `iqn` (Required) Initiator IQN. (iSCSI protocol parameters)

The `snapshot_policy` block supports:
* `schedule` - (Required) The schedule configuration for creating snapshot policy. When `snapshot_policy_name` does not exist, the snapshot policy will be created with `schedule`(s) and named as `snapshot_policy_name`. It supports the volume creation based on the AWS, AZURE and GCP CVO.

The `schedule` block supports:
* `schedule_type` - (Required, Forces new resource) snapshot policy schedule type. Must be one of '5min', '8hour', 'hourly', 'daily', 'weekly', 'monthly'.
* `retention` - (Required, Forces new resource) snapshot policy retention.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the volume.

## Import

This resource supports import, which allows you to import existing volumes into the state of this resource.

#### Standard Mode
Import requires deployment_mode,client_id,working_environment_name and volume name, separated by a comma.

id = `deployment_mode`,`client_id`,`working_environment_name`,`name`

#### Restricted Mode
Import requires deployment_mode,client_id,working_environment_name,volume name,tenant_id and connector_ip separated by a comma.

id = `deployment_mode`,`client_id`,`working_environment_name`,`name`,`tenant_id`,`connector_ip`

### Terraform Import

For example

```shell
 terraform import netapp-cloudmanager_volume.example Standard,xxxxxx,cvo,xxxxx,xxxxx,10.10.10.10
```

!> The terraform import CLI command can only import resources into the state. Importing via the CLI does not generate configuration. If you want to generate the accompanying configuration for imported resources, use the import block instead.

### Terraform Import Block

This requires Terraform 1.5 or higher, and will auto create the configuration for you

First create the block

```terraform
import {
  to = netapp-cloudmanager_volume.volume_import
  id = "Standard,xxxxxx,cvo,xxxxx,xxxxx,10.10.10.10"
}
```

Next run, this will auto create the configuration for you

```shell
terraform plan -generate-config-out=generated.tf
```

This will generate a file called generated.tf, which will contain the configuration for the imported resource

```terraform
# __generated__ by Terraform
# Please review these resources and move them into your main configuration files.

# __generated__ by Terraform from "Standard,xxxxxx,cvo,xxxxx,xxxxx,10.10.10.10"
resource "netapp-cloudmanager_volume" "volume_import" {
  space_guarantee = "volume"
  state           = "online"
  svm_name        = "svm1"
  tiering = {
    minimum_cooling_days = 0
    policy_name          = "none"
  }
  type = "rw"

  aggregate_name = "aggr"
  capacity_tier = "test"
  client_id = "xxxxxxx"
  comment = "test"
  deployment_mode = "Standard"
  id = "xxxxxxxxx"
  name = "test"
  size = 10
  svm_name = "svm"
  unit = "mb"
  working_environment_id = "xxxxxx"
  working_environment_name = "test"
}
```
