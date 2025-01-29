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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the new aggregate.
* `working_environment_id` - (Optional) The public ID of the working environment where the aggregate will be created. This argument is optional if working_environment_name is provided. You can find the ID from a previous create Cloud Volumes ONTAP action as shown in the example, or from the information page of the Cloud Volumes ONTAP working environment on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Optional) The working environment name where the aggregate will be created. This argument will be ignored if working_environment_id is provided.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `number_of_disks` - (Required) The required number of disks in the new aggregate.
* `disk_size_size` - (Optional) The required size of the disks. The default is '1'. The max number depends on the `provider_volume_type`. Details in this document: AWS: [https://docs.netapp.com/us-en/cloud-volumes-ontap-relnotes/reference-limits-aws.html#aggregate-limits] Azure: [https://docs.netapp.com/us-en/cloud-volumes-ontap-relnotes/reference-limits-azure.html#aggregate-limits] GCP: [https://docs.netapp.com/us-en/cloud-volumes-ontap-relnotes/reference-limits-gcp.html#disk-and-tiering-limits]
* `disk_size_unit` - (Optional) The disk size unit ['GB' or 'TB']. The default is 'TB'
* `home_node` - (Optional) The home node that the new aggregate should belong to. The default is the first node.
* `provider_volume_type` - (Optional) The cloud provider volume type. For AWS: ['gp3', 'gp2', 'io1', 'st1', 'sc1']. For Azure: ['Premium_LRS','Standard_LRS','StandardSSD_LRS']. For GCP: ['pd-balanced', 'pd-ssd','pd-standard']
* `capacity_tier` - (Optional) The aggregate's capacity tier for tiering cold data to object storage: ['S3', 'Blob', 'cloudStorage']. The default values for each cloud provider are as follows: Amazon => 'S3', Azure => 'Blob', GCP => 'cloudStorage'. If NONE, the capacity tier won't be set on aggregate creation.
* `iops` - (Optional) Provisioned IOPS. Needed only when 'providerVolumeType' is 'io1' or 'gp3'
* `throughput` - (Optional) Required only when 'providerVolumeType' is 'gp3'.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - will be the aggregate name.
