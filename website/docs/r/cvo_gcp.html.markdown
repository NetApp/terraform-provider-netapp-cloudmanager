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
  workspace_id = "workspace-******"
  license_type = "capacity-paygo"
  gcp_label {
        label_key = "abcd"
        label_value = "ABCD"
  }
  svm {
    svm_name = "svm01"
  }
  svm {
    svm_name = "svm03"
  }
}
```

**Create netapp-cloudmanager_cvo_gcp for restricted mode:**

```
resource "netapp-cloudmanager_cvo_gcp" "cl-cvo-gcp" {
  provider = netapp-cloudmanager
  name = "terraformcvogcp"
  project_id = "occm-project"
  zone = "us-east1-b"
  capacity_package_name = "Freemium"
  subnet_id = "cvs-terraform-abc"
  vpc_id = "cvs-terraform-abc"
  gcp_volume_type = "pd-ssd"
  data_encryption_type = "GCP"
  svm_password = "netapp1!"
  ontap_version = "latest"
  use_latest_version = true
  instance_type = "n2-standard-4"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
  workspace_id = "workspace-******"
  writing_speed_state = "NORMAL"
  license_type = "capacity-paygo"
  enable_compliance = true
  gcp_volume_size = 500
  gcp_volume_size_unit = "GB"
  deployment_mode = "Restricted"
  connector_ip = "10.10.10.10"
  tenant_id = "account-******"
}
```

**Create netapp-cloudmanager_cvo_gcp with WORM:**

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
  license_type = "capacity-paygo"
  gcp_label {
        label_key = "abcd"
        label_value = "ABCD"
  }
  svm {
    svm_name = "svm01"
  }
  svm {
    svm_name = "svm03"
  }
  worm_retention_period_length = 2
  worm_retention_period_unit = "hours"
}
```

**Create netapp-cloudmanager_cvo_gcp HA with HIGH writing speed:**

```
resource "netapp-cloudmanager_cvo_gcp" "cl-cvo-gcp-ha" {
  provider = netapp-cloudmanager
  name = "tfcvohahigh"
  project_id = "occm-project"
  zone = "us-east4-a"
  subnet_id = "default"
  vpc_id = "default"
  gcp_service_account = "abcdefg@tlv-support.iam.gserviceaccount.com"
  is_ha = true
  svm_password = "netapp11!"
  use_latest_version = false
  ontap_version = "ONTAP-9.13.0.T1.gcpha"
  gcp_volume_type = "pd-ssd"
  capacity_package_name = "Professional"
  instance_type = "n2-standard-16"
  license_type = "ha-capacity-paygo"
  mediator_zone = "us-east4-c"
  node1_zone = "us-east4-a"
  node2_zone =  "us-east4-b"
  subnet0_node_and_data_connectivity = "projects/tlv-support/regions/us-east4/subnetworks/default"
  subnet1_cluster_connectivity = "projects/tlv-support/regions/us-east4/subnetworks/rn-cluster-subnet"
  subnet2_ha_connectivity = "projects/tlv-support/regions/us-east4/subnetworks/rn-rdma-subnet"
  subnet3_data_replication = "projects/tlv-support/regions/us-east4/subnetworks/rn-replication-subnet"
  vpc0_node_and_data_connectivity = "projects/tlv-support/global/networks/default"
  vpc1_cluster_connectivity = "projects/tlv-support/global/networks/rnicholl-vpc1-cluster-internal"
  vpc2_ha_connectivity = "projects/tlv-support/global/networks/rnicholl-vpc2-rdma-internal"
  vpc3_data_replication = "projects/tlv-support/global/networks/rnicholl-vpc3-replication-internal"
  nss_account = "cd12x234-f876-4567-8f1f-12345678xxx"
  client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
  writing_speed_state = "HIGH"
  flash_cache = true
}
```

## Argument Reference

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the Cloud Volumes ONTAP working environment.
* `project_id` - (Required, Forces new resource) The ID of the GCP project.
* `zone` - (Optional, Forces new resource) The zone of the region where the working environment will be created. It is required in single.
* `gcp_service_account` - (Optional, Forces new resource) The gcp_service_account email in order to enable tiering of cold data to Google Cloud Storage.
* `svm_password` - (Required) The admin password for Cloud Volumes ONTAP.
* `svm_name` - (Optional) The name of the SVM.
* `connector_ip` - (Optional) The private IP of the connector, this is only required for Restricted mode.
* `tenant_id` - (Optional) The NetApp tenant ID that the Connector will be associated with.  You can find the tenant ID in the Identity & Access Management in Settings, Organization tab of BlueXP at [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `deployment_mode` - (Optional) The mode of deployment to use for the working environment: ['Standard', 'Restricted']. The default is 'Standard'. To know more on deployment modes [https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/](https://docs.netapp.com/us-en/bluexp-setup-admin/concept-modes.html/).
* `client_id` - (Required, Forces new resource) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `workspace_id` - (Optional, Forces new resource) The ID of the Cloud Manager workspace where you want to deploy Cloud Volumes ONTAP. If not provided, Cloud Manager uses the first workspace. You can find the ID from the Workspace tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `data_encryption_type` - (Optional, Forces new resource) The type of data encryption to use for the working environment: ['GCP', 'NONE']. The default is 'GCP'.
* `gcp_encryption_parameters` - (Optional, Forces new resource) Required if using gcp encryption with custom key. Key format is 'projects/default-project/locations/global/keyRings/test/cryptoKeys/key1'.
* `gcp_volume_type` - (Optional, Forces new resource) The type of the storage for the first data aggregate: ['pd-balanced', 'pd-standard', 'pd-ssd']. The default is 'pd-ssd'
* `subnet_id` - (Optional, Forces new resource) The name of the subnet for Cloud Volumes ONTAP. The default is: 'default'.
* `network_project_id` - (Optional, Forces new resource) The project id in GCP associated with the Subnet. If not provided, it’s assumed that the Subnet is within the previously specified project id.
* `vpc_id` - (Optional, Forces new resource) The name of the VPC.
* `gcp_volume_size` - (Optional, Forces new resource) The GCP volume size for the first data aggregate. For GB, the unit can be: [100 or 500]. For TB, the unit can be: [1,2,4,8]. The default is '1' .
* `gcp_volume_size_unit` - (Optional, Forces new resource) ['GB' or 'TB']. The default is 'TB'.
* `ontap_version` - (Optional) The required ONTAP version. Ignored if `use_latest_version` is set to true. The default is to use the latest version. The naming convention: 

|Release|Naming convention|Example|
|-------|-----------------|-------|
|Patch Single | `ONTAP-${version}.gcp` | ONTAP-9.13.1P1.gcp|
|Patch HA | `ONTAP-${version}.gcpha` | ONTAP-9.13.1P1.gcpha|
|Regular Single | `ONTAP-${version}.T1.gcp` | ONTAP-9.10.1RC1.T1.gcp|
|Regular HA | `ONTAP-${version}.T1.gcpha` | ONTAP-9.13.0.T1.gcpha|

* `use_latest_version` - (Optional) Indicates whether to use the latest available ONTAP version. The default is 'true'.
* `license_type` - (Optional) The type of license to use. For single node: (by Capacity): ['capacity-paygo'], (by Node paygo): ['gcp-cot-explore-paygo', 'gcp-cot-standard-paygo', 'gcp-cot-premium-paygo'], (by Node byol): ['gcp-cot-premium-byol'], For HA: (by Capacity): ['ha-capacity-paygo'], (by Node paygo): ['gcp-ha-cot-explore-paygo', 'gcp-ha-cot-standard-paygo', 'gcp-ha-cot-premium-paygo'], (by Node byol): ['gcp-ha-cot-premium-byol']. The default is 'capacity-paygo' for single node, and 'ha-capacity-paygo'for HA.
* `capacity_package_name` - (Optional) The capacity package name: ['Essential', 'Professional', 'Freemium', 'Edge', 'Optimized']. Default is 'Essential'. 'Edge' and 'Optimized' need ontap version 9.11.1 or above.
* `instance_type` - (Required) The type of instance to use, which depends on the license type you choose: Explore:['custom-4-16384'], Standard:['n1-standard-8'], Premium:['n1-standard-32'], BYOL: all instance types defined for PayGo. For more supported instance types, refer to Cloud Volumes ONTAP Release Notes. The default is 'n2-standard-8’ but the users will have to specify the default value explicitly during CVO creation. 
* `serial_number` - (Optional, Forces new resource) The serial number for the system. Required when using 'gcp-cot-premium-byol'.
* `capacity_tier` - (Optional, Forces new resource) Indicates the type of data tiering to use: ['cloudStorage']. The default is 'cloudStorage'.
* `tier_level` - (Optional) In case capacity_tier is cloudStorage, this argument indicates the tiering level: ['standard', 'nearline', 'coldline']. The default is: 'standard'.
* `saas_subscription_id` - (Optional, Forces new resource) SaaS Subscription ID. It is needed if the subscription is not paygo type.
* `nss_account` - (Optional, Forces new resource) The NetApp Support Site account ID to use with this Cloud Volumes ONTAP system. If the license type is BYOL and an NSS account isn't provided, Cloud Manager tries to use the first existing NSS account.
* `writing_speed_state` - (Optional) The write speed setting for Cloud Volumes ONTAP: ['NORMAL','HIGH']. The default is 'NORMAL'. For single node system, HIGH write speed is supported with all machine types. For HA, Flash Cache, high write speed, and a higher maximum transmission unit (MTU) of 8,896 bytes are available through the High write speed option with the n2-standard-16, n2-standard-32, n2-standard-48, and n2-standard-64 instance types.
* `flash_cache` - (Optional, Forces new resource) Enable Flash Cache. In GCP HA (version 9.13.0), HIGH write speed and FlashCache are coupled together both needs to be activated, one cannot be activated without the other. For GCP single (version 9.13.1) is supported. Only the instance_type is one of the followings: n2-standard-16,32,48,64
* `firewall_rule` - (Optional, Forces new resource) The name of the firewall rule for a single node cluster. If not provided, the rule will be generated automatically.
* `firewall_tag_name_rule` - (Optional, Forces new resource) Target tag of the firewall when creating a CVO with an existing firewall. It is used for a single node cluster.
* `firewall_ip_ranges` - (Optional, Forces new resource) Define the allowed inbound traffic for the generated policy. It is used when selecting create a new firewall. Recommend set false: Allow traffic within the selected VPC only. Allow inbound traffic only from the cluster node VPCs.
* `backup_volumes_to_cbs` - (Optional, Forces new resource) Automatically enable back up of all volumes to Google Cloud buckets [true, false].
* `enable_compliance` - (Optional, Forces new resource) Enable the Cloud Compliance service on the working environment [true, false].
* `is_ha` - (Optional, Forces new resource) Indicate whether the working environment is an HA pair or not [true, false]. The default is false.
* `platform_serial_number_node1` - (Optional, Forces new resource) For HA BYOL, the serial number for the first node.
* `platform_serial_number_node2` - (Optional, Forces new resource) For HA BYOL, the serial number for the second node.
* `node1_zone` - (Optional, Forces new resource)  Zone for node 1. It will also be used in the 'zone' if it is not provided in HA.
* `node2_zone` - (Optional, Forces new resource) Zone for node 2.
* `mediator_zone` - (Optional, Forces new resource) Zone for mediator.
* `vpc0_node_and_data_connectivity` - (Optional, Forces new resource) VPC path for nic1, required for node and data connectivity. If using shared VPC, `network_project_id` must be provided.
* `vpc1_cluster_connectivity` - (Optional, Forces new resource) VPC path for nic2, required for cluster connectivity.
* `vpc2_ha_connectivity` - (Optional, Forces new resource) VPC path for nic3, required for HA connectivity.
* `vpc3_data_replication` - (Optional, Forces new resource) VPC path for nic4, required for data replication.
* `subnet0_node_and_data_connectivity` - (Optional, Forces new resource) Subnet path for nic1, required for node and data connectivity. If using shared VPC, `network_project_id` must be provided.
* `subnet1_cluster_connectivity` - (Optional, Forces new resource) Subnet path for nic2, required for cluster connectivity.
* `subnet2_ha_connectivity` - (Optional, Forces new resource) Subnet path for nic3, required for HA connectivity.
* `subnet3_data_replication` - (Optional, Forces new resource) Subnet path for nic4, required for data replication.
* `vpc0_firewall_rule_name` - (Optional, Forces new resource) Firewall rule name for vpc1.
* `vpc1_firewall_rule_name` - (Optional, Forces new resource) Firewall rule name for vpc2.
* `vpc2_firewall_rule_name` - (Optional, Forces new resource) Firewall rule name for vpc3.
* `vpc3_firewall_rule_name` - (Optional, Forces new resource) Firewall rule name for vpc4.
* `vpc0_firewall_rule_tag_name` - (Optional, Forces new resource) Firewall rule tag name for vpc1.
* `vpc1_firewall_rule_tag_name` - (Optional, Forces new resource) Firewall rule tag name for vpc2.
* `vpc2_firewall_rule_tag_name` - (Optional, Forces new resource) Firewall rule tag name for vpc3.
* `vpc3_firewall_rule_tag_name` - (Optional, Forces new resource) Firewall rule tag name for vpc4.
* `upgrade_ontap_version` - (Optional) Indicates whether to upgrade ontap image with `ontap_version`. To upgrade ontap image, `ontap_version` cannot be 'latest' and `use_latest_version` needs to be false. The available versions can be found in BlueXP UI. Click the CVO -> click **New Version Available** under **Notifications** -> the latest available version will be shown. The list of available versions can be found in **Select older versions**. Update the `ontap_version` by follow the naming conversion.
* `retries` - (Optional) The number of attempts to wait for the completion of creating the CVO with 60 seconds apart for each attempt. For HA, this value is incremented by 30. The default is '60'.
* `worm_retention_period_length` - (Optional, Forces new resource) WORM retention period length. Once specified retention period, the WORM is enabled. When WORM storage is activated, data tiering to object storage can’t be enabled.
* `worm_retention_period_unit` - (Optional, Forces new resource) WORM retention period unit: ['years','months','days','hours','minutes','seconds'].

The `gcp_label` block supports:
* `label_key` - (Required) The key of the tag.
* `label_value` - (Required) The tag value.

The `svm` block supports:
* `svm_name` - (Required) The extra SVM name for CVO HA.
* `root_volume_aggregate` - (Optional) Specifies the aggregate where the root volume of the SVM will be created. This attribute could only be used after CVO creation to add SVM to an existing CVO. 

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the working environment.
* `svm_name` - The default name of the SVM will be exported if it is not provided in the resource.
