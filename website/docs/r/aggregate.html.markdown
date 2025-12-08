---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_aggregate"
sidebar_current: "docs-netapp-cloudmanager-resource-aggregate"
description: |-
  Provides a netapp-cloudmanager_aggregate resource. This can be used to create a new aggregate on Cloud Volumes ONTAP.
---

# netapp-cloudmanager_aggregate

Provides a netapp-cloudmanager_aggregate resource. This can be used to create a new aggregate on Cloud Volumes ONTAP.
Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.

## Example Usages

**Create netapp-cloudmanager_aggregate:**

```
resource "netapp-cloudmanager_aggregate" "cl-aggregate" {
  provider = netapp-cloudmanager
  name = "aggr2"
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id #
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id #
  number_of_disks = 1
  provider_volume_type = "gp2"
}
```

**Create netapp-cloudmanager_aggregate with EBS Elastic Volumes support (AWS only):**

```
resource "netapp-cloudmanager_aggregate" "cl-aggregate-with-ev" {
  provider = netapp-cloudmanager
  name = "aggr_with_ev_support"
  working_environment_id = netapp-cloudmanager_cvo_aws.cvo-aws.id
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id
  number_of_disks = 3
  provider_volume_type = "gp3"
  disk_size_size = 100
  disk_size_unit = "GB"
  
  # Enable EBS Elastic Volumes support (creation time only)
  initial_ev_aggregate_size = 500
  initial_ev_aggregate_unit = "GB"
}
```

**Update aggregate capacity using EBS Elastic Volumes (AWS only):**

```
resource "netapp-cloudmanager_aggregate" "cl-aggregate-with-capacity" {
  provider = netapp-cloudmanager
  name = "aggr_with_capacity"
  working_environment_id = netapp-cloudmanager_cvo_aws.cvo-aws.id
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id
  number_of_disks = 2
  provider_volume_type = "gp3"
  disk_size_size = 100
  disk_size_unit = "GB"
  
  # EBS Elastic Volumes configuration (required for capacity increase)
  initial_ev_aggregate_size = 500
  initial_ev_aggregate_unit = "GB"
  
  # Increase aggregate capacity (update operation only)
  increase_capacity_size = 200
  increase_capacity_unit = "GB"
}
```

## Argument Reference

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the new aggregate.
* `working_environment_id` - (Optional, Forces new resource) The public ID of the working environment where the aggregate will be created. This argument is optional if working_environment_name is provided. You can find the ID from a previous create Cloud Volumes ONTAP action as shown in the example, or from the information page of the Cloud Volumes ONTAP working environment on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Optional, Forces new resource) The working environment name where the aggregate will be created. This argument will be ignored if working_environment_id is provided.
* `client_id` - (Required, Forces new resource) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `connector_ip` - (Optional) The IP of the connector, this is only required for 'Restricted' mode account.
* `tenant_id` - (Optional) The NetApp tenant ID that the Connector will be associated with. This is required for the Restricted deployment mode. You can find the tenant ID in the Identity & Access Management in Settings, Organization tab of BlueXP at [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `deployment_mode` - (Optional) The mode of deployment to use for the working environment: ['Standard', 'Restricted']. The default is 'Standard'. To know more on deployment modes [https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/](https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/)
* `number_of_disks` - (Optional) The required number of disks in the new aggregate.
* `disk_size_size` - (Optional, Forces new resource) The required size of the disks. The max number depends on the `provider_volume_type`. Details in this document: AWS: [https://docs.netapp.com/us-en/cloud-volumes-ontap-relnotes/reference-limits-aws.html#aggregate-limits] Azure: [https://docs.netapp.com/us-en/cloud-volumes-ontap-relnotes/reference-limits-azure.html#aggregate-limits] GCP: [https://docs.netapp.com/us-en/cloud-volumes-ontap-relnotes/reference-limits-gcp.html#disk-and-tiering-limits] **Note: Must be provided together with `disk_size_unit`**
* `disk_size_unit` - (Optional, Forces new resource) The disk size unit ['GB' or 'TB']. **Note: Must be provided together with `disk_size_size`**
* `home_node` - (Optional, Forces new resource) The home node that the new aggregate should belong to. The default is the first node.
* `provider_volume_type` - (Optional, Forces new resource) The cloud provider volume type. For AWS: ['gp3', 'gp2', 'io1', 'st1', 'sc1']. For Azure: ['Premium_LRS','Standard_LRS','StandardSSD_LRS']. For GCP: ['pd-balanced', 'pd-ssd','pd-standard']
* `capacity_tier` - (Optional, Forces new resource) The aggregate's capacity tier for tiering cold data to object storage: ['S3', 'Blob', 'cloudStorage']. The default values for each cloud provider are as follows: Amazon => 'S3', Azure => 'Blob', GCP => 'cloudStorage'. If NONE, the capacity tier won't be set on aggregate creation.
* `iops` - (Optional, Forces new resource) Provisioned IOPS. Needed only when 'providerVolumeType' is 'io1' or 'gp3'
* `throughput` - (Optional, Forces new resource) Required only when 'providerVolumeType' is 'gp3'.
* `initial_ev_aggregate_size` - (Optional, Forces new resource) Initial size for EBS Elastic Volumes aggregate (AWS only). This enables the aggregate to support capacity expansion using Amazon EBS Elastic Volumes. **Creation time only** - cannot be modified after aggregate creation. **Note: Must be provided together with `initial_ev_aggregate_unit`**
* `initial_ev_aggregate_unit` - (Optional, Forces new resource) Unit for initial EBS Elastic Volumes aggregate size (GB, TB, GiB, or TiB). Only used with `initial_ev_aggregate_size`. Defaults to 'GB' if not specified. **Creation time only** - cannot be modified after aggregate creation. **Note: Must be provided together with `initial_ev_aggregate_size`**
* `increase_capacity_size` - (Optional, Computed) Additional capacity to add to the aggregate using Amazon EBS Elastic Volumes. **Only supported for AWS aggregates with EBS Elastic Volumes enabled**. **Update operation only** - cannot be used during aggregate creation. The aggregate must be created with `initial_ev_aggregate_size` to support capacity increases. **Important:** After a successful capacity increase operation, remove the parameter from your configuration to prevent unnecessary state changes and achieve idempotency in subsequent Terraform runs. **Note: Must be provided together with `increase_capacity_unit`**
* `increase_capacity_unit` - (Optional, Computed) Unit for the additional capacity (Byte, KB, MB, GB, or TB). Only used with `increase_capacity_size`. **Update operation only** - cannot be used during aggregate creation. **Important:** After a successful capacity increase operation, remove the parameter from your configuration to prevent unnecessary state changes and achieve idempotency in subsequent Terraform runs. **Note: Must be provided together with `increase_capacity_size`**

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - will be the aggregate name.
* `total_capacity_size` - The total capacity of the aggregate.
* `total_capacity_unit` - The unit of the total capacity.
* `available_capacity_size` - The available capacity of the aggregate.
* `available_capacity_unit` - The unit of the available capacity.

## Import

This resource supports import, which allows you to import existing aggregates into the state of this resource.

#### Standard Mode
Import requires deployment_mode,client_id,working_environment_name and aggregate name, separated by a comma.

id = `deployment_mode`,`client_id`,`working_environment_name`,`name`

#### Restricted Mode
Import requires deployment_mode,client_id,working_environment_name,aggregate name,tenant_id and connector_ip separated by a comma.

id = `deployment_mode`,`client_id`,`working_environment_name`,`name`,`tenant_id`,`connector_ip`

### Terraform Import

For example

```shell
 terraform import netapp-cloudmanager_aggregate.example Standard,xxxxxx,cvo,aggr1
```

!> The terraform import CLI command can only import resources into the state. Importing via the CLI does not generate configuration. If you want to generate the accompanying configuration for imported resources, use the import block instead.

### Terraform Import Block

This requires Terraform 1.5 or higher, and will auto create the configuration for you

First create the block

```terraform
import {
  to = netapp-cloudmanager_aggregate.aggregate_import
  id = "Standard,xxxxxx,cvo,aggr1"
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

# __generated__ by Terraform from "Standard,xxxxxx,cvo,aggr1"
resource "netapp-cloudmanager_aggregate" "aggregate_import" {
  available_capacity_size = 100
  available_capacity_unit = "GB"
  capacity_tier           = "S3"
  client_id               = "xxxxxxx"
  deployment_mode         = "Standard"
  disk_size_size          = 100
  disk_size_unit          = "GB"
  home_node               = "node1"
  id                      = "aggr1"
  name                    = "aggr1"
  number_of_disks         = 6
  provider_volume_type    = "gp2"
  total_capacity_size     = 600
  total_capacity_unit     = "GB"
  working_environment_id  = "xxxxxx"
  working_environment_name = "cvo"
}
```
