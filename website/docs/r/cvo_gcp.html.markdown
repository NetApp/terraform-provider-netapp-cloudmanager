---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cvo_gcp"
sidebar_current: "docs-netapp-cloudmanager-resource-cvo-gcp"
description: |-
  Provides a netapp-cloudmanager_cvo_gcp resource. This can be used to create a new Cloud Volume ONTAP system in GCP.
---

# netapp-cloudmanager_cvo_gcp

Provides a netapp-cloudmanager_cvo_gcp resource. This can be used to create a new Cloud Volume ONTAP system in GCP.

## Example Usages

**Create netapp-cloudmanager_cvo_gcp:**

```
resource "netapp-cloudmanager_cvo_gcp" "cl-cvo-gcp" {
  provider = netapp-cloudmanager
  name = "terraformcvogcp"
  project_id = "occm-project"
  zone = "us-east1-b"
  gcp_service_account = "fabric-pool@occm-project.iam.gserviceaccount.com"
  svm_password = "netapp1!"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
  workspace_id = "workspace-IDz6Nnwl"
  gcp_label {
        label_key = "abcd"
        label_value = "ABCD"
      }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Volumes ONTAP working environment.
* `project_id` - (Required) The ID of the GCP project.
* `zone` - (Required) The zone of the region where the working environment will be created.
* `gcp_service_account` - (Required) The gcp_service_account email in order to enable tiering of cold data to Google Cloud Storage.
* `svm_password` - (Required) The admin password for Cloud Volumes ONTAP.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `workspace_id` - (Optional) The ID of the Cloud Manager workspace where you want to deploy Cloud Volumes ONTAP. If not provided, Cloud Manager uses the first workspace. You can find the ID from the Workspace tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `data_encryption_type` - (Optional) The type of data encryption to use for the working environment: ['GCP', 'NONE']. The default is 'GCP'.
* `gcp_encryption_parameters` - (Optional) Required if using gcp encryption with custom key. Key format is 'projects/default-project/locations/global/keyRings/test/cryptoKeys/key1'.
* `gcp_volume_type` - (Optional) The type of the storage for the first data aggregate: ['pd-balanced', 'pd-standard', 'pd-ssd']. The default is 'pd-ssd'
* `subnet_id` - (Optional) The name of the subnet for Cloud Volumes ONTAP. The default is: 'default'.
* `network_project_id` - (Optional) The project id in GCP associated with the Subnet. If not provided, it’s assumed that the Subnet is within the previously specified project id.
* `vpc_id` - (Optional) The name of the VPC.
* `gcp_volume_size` - (Optional) The GCP volume size for the first data aggregate. For GB, the unit can be: [100 or 500]. For TB, the unit can be: [1,2,4,8]. The default is '1' .
* `gcp_volume_size_unit` - (Optional) ['GB' or 'TB']. The default is 'TB'.
* `ontap_version` - (Optional) The required ONTAP version. Ignored if 'use_latest_version' is set to true. The default is to use the latest version.
* `use_latest_version` - (Optional) Indicates whether to use the latest available ONTAP version. The default is 'true'.
* `license_type` - (Optional) The type of license to use. For single node: (by Capacity): ['capacity-paygo'], (by Node paygo): ['gcp-cot-explore-paygo', 'gcp-cot-standard-paygo', 'gcp-cot-premium-paygo'], (by Node byol): ['gcp-cot-premium-byol'], For HA: (by Capacity): ['ha-capacity-paygo'], (by Node paygo): ['gcp-ha-cot-explore-paygo', 'gcp-ha-cot-standard-paygo', 'gcp-ha-cot-premium-paygo'], (by Node byol): ['gcp-ha-cot-premium-byol']. The default is 'capacity-paygo' for single node, and 'ha-capacity-paygo'for HA.
* `capacity_package_name` - (Optional) The capacity package name: ['Essential', 'Professional', 'Freemium']. Default is 'Essential'.
* `instance_type` - (Optional) The type of instance to use, which depends on the license type you choose: Explore:['custom-4-16384'], Standard:['n1-standard-8'], Premium:['n1-standard-32'], BYOL: all instance types defined for PayGo. For more supported instance types, refer to Cloud Volumes ONTAP Release Notes. default is 'n1-standard-8' .
* `serial_number` - (Optional) The serial number for the system. Required when using 'gcp-cot-premium-byol'.
* `capacity_tier` - (Optional) Indicates the type of data tiering to use: ['cloudStorage', 'NONE']. The default is 'cloudStorage'.
* `tier_level` - (Optional) In case capacity_tier is cloudStorage, this argument indicates the tiering level: ['standard', 'nearline', 'coldline']. The default is: 'standard'.
* `nss_account` - (Optional) The NetApp Support Site account ID to use with this Cloud Volumes ONTAP system. If the license type is BYOL and an NSS account isn't provided, Cloud Manager tries to use the first existing NSS account.
* `writing_speed_state` - (Optional) The write speed setting for Cloud Volumes ONTAP: ['NORMAL','HIGH']. The default is 'NORMAL'. This argument is not relevant for HA pairs.
* `firewall_rule` - (Optional) The name of the firewall rule for Cloud Volumes ONTAP. If not provided, Cloud Manager generates the rule.
* `backup_volumes_to_cbs` - (Optional) Automatically enable back up of all volumes to Google Cloud buckets [true, false].
* `enable_compliance` - (Optional) Enable the Cloud Compliance service on the working environment [true, false].
* `is_ha` - (Optional) Indicate whether the working environment is an HA pair or not [true, false]. The default is false.
* `platform_serial_number_node1` - (Optional) For HA BYOL, the serial number for the first node.
* `platform_serial_number_node2` - (Optional) For HA BYOL, the serial number for the second node.
* `node1_zone` - (Optional)  Zone for node 1.
* `node2_zone` - (Optional) Zone for node 2.
* `mediator_zone` - (Optional) Zone for mediator.
* `vpc0_node_and_data_connectivity` - (Optional) VPC path for nic1, required for node and data connectivity. If using shared VPC, netwrok_project_id must be provided.
* `vpc1_cluster_connectivity` - (Optional) VPC path for nic2, required for cluster connectivity.
* `vpc2_ha_connectivity` - (Optional) VPC path for nic3, required for HA connectivity.
* `vpc3_data_replication` - (Optional) VPC path for nic4, required for data replication.
* `subnet0_node_and_data_connectivity` - (Optional) Subnet path for nic1, requered for node and data connectivity. If using shared VPC, netwrok_project_id must be provided.
* `subnet1_cluster_connectivity` - (Optional) Subnet path for nic2, required for cluster connectivity.
* `subnet2_ha_connectivity` - (Optional) Subnet path for nic3, required for HA connectivity.
* `subnet3_data_replication` - (Optional) Subnet path for nic4, required for data replication.
* `vpc0_firewall_rule_name` - (Optional) Firewall rule name for vpc1.
* `vpc1_firewall_rule_name` - (Optional) Firewall rule name for vpc2.
* `vpc2_firewall_rule_name` - (Optional) Firewall rule name for vpc3.
* `vpc3_firewall_rule_name` - (Optional) Firewall rule name for vpc4.
* `upgrade_ontap_version` - (Optional) Indicates whether to upgrade ontap image with `ontap_version`. To upgrade ontap image, `ontap_version` cannot be 'latest' and `use_latest_version` needs to be false.

The `gcp_label` block supports:
* `label_key` - (Required) The key of the tag.
* `label_value` - (Required) The tag value.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the working environment.
* `svm_name` - The name of the SVM.
