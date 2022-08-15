---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cvo_aws"
sidebar_current: "docs-netapp-cloudmanager-resource-cvo-aws"
description: |-
  Provides a netapp-cloudmanager_cvo_aws resource. This can be used to create a new Cloud Volume ONTAP system in AWS (single node or an HA pair).
---

# netapp-cloudmanager_cvo_aws

Provides a netapp-cloudmanager_cvo_aws resource. This can be used to create a new Cloud Volume ONTAP system in AWS (single node or an HA pair). The environment needs to be configured with the proper credentials before it can be used (run this command: aws configure).

## Example Usages

**Create netapp-cloudmanager_cvo_aws single:**

```
resource "netapp-cloudmanager_cvo_aws" "cvo-aws" {
  provider = netapp-cloudmanager
  name = "TerraformCVO"
  region = "us-west-2"
  subnet_id = "subnet-xxxxxxx"
  vpc_id = "vpc-xxxxxxxx"
  aws_tag {
              tag_key = "abcd"
              tag_value = "ABCD"
            }
  aws_tag {
              tag_key = "xxx"
              tag_value = "YYY"
            }
  svm_password = "P@assword!"
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id
  writing_speed_state = "NORMAL"
}
```


**Create netapp-cloudmanager_cvo_aws HA:**

```
resource "netapp-cloudmanager_cvo_aws" "cvo-aws" {
  provider = netapp-cloudmanager
  name = "TerraformCVO"
  region = "us-west-2"
  subnet_id = "subnet-xxxxxxx"
  vpc_id = "vpc-xxxxxxxx"
  aws_tag {
              tag_key = "abcd"
              tag_value = "ABCD"
            }
  aws_tag {
              tag_key = "xxx"
              tag_value = "YYY"
            }
  svm_password = "P@assword!"
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id
  is_ha = true
  failover_mode = "FloatingIP"
  node1_subnet_id = "subnet-1"
  node2_subnet_id = "subnet-1"
  mediator_subnet_id = "subnet-xxxxxx"
  mediator_key_pair_name = "key1"
  cluster_floating_ip = "2.1.1.1"
  data_floating_ip = "2.1.1.2"
  data_floating_ip2 = "2.1.1.3"
  svm_floating_ip = "2.1.1.4"
  route_table_ids = ["rt-1","rt-2"]
  license_type = "ha-cot-standard-paygo"
}
```




## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Volumes ONTAP working environment.
* `region` - (Required) The region where the working environment will be created.
* `subnet_id` - (Optional) The subnet id where the working environment will be created. Required when single mode only.
* `svm_password` - (Required) The admin password for Cloud Volumes ONTAP.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `vpc_id` - (Optional) The VPC ID where the working environment will be created. If this argument isn't provided, the VPC will be calculated by using the provided subnet ID.
* `workspace_id` - (Optional) The ID of the Cloud Manager workspace where you want to deploy Cloud Volumes ONTAP. If not provided, Cloud Manager uses the first workspace. You can find the ID from the Workspace tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `data_encryption_type` - (Optional) The type of encryption to use for the working environment: ['AWS', 'NONE']. The default is 'AWS'.
* `aws_encryption_kms_key_id` - (Optional) AWS encryption parameters. It is required if using aws encryption. Only one of KMS key id or KMS arn should be specified
* `aws_encryption_kms_key_arn` - (Optional) AWS encryption parameters. It is requried if using aws encryption. Only one of KMS key id or KMS arn should be specified
* `ebs_volume_size` - (Optional) EBS volume size for the first data aggregate. For GB, the unit can be: [100 or 500]. For TB, the unit can be: [1,2,4,8,16]. The default is '1' .
* `ebs_volume_size_unit` - (Optional) ['GB' or 'TB']. The default is 'TB'.
* `ebs_volume_type` - (Optional) The EBS volume type for the first data aggregate ['gp3', 'gp2','io1','st1','sc1']. The default is 'gp2'.
* `iops` - (Optional) Provisioned IOPS. Required only when 'ebs_volume_type' is 'io1' or 'gp3'.
* `throughput` - (Optional) Required only when 'ebs_volume_type' is 'gp3'.
* `ontap_version` - (Optional) The required ONTAP version. Ignored if 'use_latest_version' is set to true. The default is to use the latest version.
* `use_latest_version` - (Optional) Indicates whether to use the latest available ONTAP version. The default is 'true'.
* `license_type` - (Optional) The type of license to use. For single node: (by Capacity): ['capacity-paygo'], (by Node paygo): ['cot-explore-paygo','cot-standard-paygo', 'cot-premium-paygo'], (by Node byol): ['cot-premium-byol']. For HA: (by Capacity): ['ha-capacity-paygo'], (by Node paygo): ['ha-cot-explore-paygo','ha-cot-standard-paygo','ha-cot-premium-paygo'], (by Node byol): 'ha-cot-premium-byol']. The default is 'capacity-paygo' for single node, and 'ha-capacity-paygo' for HA.
* `capacity_package_name` - (Optional) The capacity package name: ['Essential', 'Professional', 'Freemium']. Default is 'Essential'.
* `instance_type` - (Optional) The instance type to use, which depends on the license type: Explore:['m5.xlarge'], Standard:['m5.2xlarge','r5.xlarge'], Premium:['m5.4xlarge','r5.2xlarge','c4.8xlarge'], BYOL: all instance types defined for PayGo. For more supported instance types, refer to Cloud Volumes ONTAP Release Notes. The default is 'm5.2xlarge'.
* `platform_serial_number` - (Optional) The serial number for the cluster. This is required when 'license_type' is set 'cot-premium-byol'.
* `capacity_tier` - (Optional) Whether to enable data tiering for the first data aggregate: ['S3','NONE']. The default is 'S3'.
* `tier_level` - (Optional) The tiering level when 'capacity_tier' is set to 'S3' ['normal','ia','ia-single','intelligent']. The default is 'normal'.
* `nss_account` - (Optional) The NetApp Support Site account ID to use with this Cloud Volumes ONTAP system. If the license type is BYOL and an NSS account isn't provided, Cloud Manager tries to use the first existing NSS account.
* `writing_speed_state` - (Optional) The write speed setting for Cloud Volumes ONTAP: ['NORMAL','HIGH']. The default is 'NORMAL'. This argument is not relevant for HA pairs.
* `instance_tenancy` - (Optional) The EC2 instance tenancy: ['default','dedicated']. The default is 'default'.
* `instance_profile_name` - (Optional) The instance profile name for the working environment. If not provided, Cloud Manager creates the instance profile.
* `security_group_id` - (Optional) The ID of the security group for the working environment. If not provided, Cloud Manager creates the security group.
* `cloud_provider_account` - (Optional) The cloud provider credentials id to use when deploying the Cloud Volumes ONTAP system. You can find the ID in Cloud Manager from the Settings > Credentials page. If not specified, Cloud Manager uses the instance profile of the Connector.
* `backup_volumes_to_cbs` - (Optional) Automatically enable back up of all volumes to S3 [true, false].
* `enable_compliance` - (Optional) Enable the Cloud Compliance service on the working environment [true, false].
* `enable_monitoring` - (Optional) Enable the Monitoring service on the working environment [true, false]. The default is false.
* `optimized_network_utilization` - (Optional) Use optimized network utilization [true, false]. The default is true.
* `is_ha` - (Optional) Indicate whether the working environment is an HA pair or not [true, false]. The default is false.
* `failover_mode` - (Optional) For HA, the failover mode for the HA pair: ['PrivateIP', 'FloatingIP']. 'PrivateIP' is for a single availability zone and 'FloatingIP' is for multiple availability zones.
* `mediator_assign_public_ip` - (Optional) bool option to assign public IP. The default is 'true'.
* `mediator_instance_profile_name` - (Optional) name of the mediator instance profile.
* `platform_serial_number_node1` - (Optional) For HA BYOL, the serial number for the first node. This is required when using 'ha-cot-premium-byol'.
* `platform_serial_number_node2` - (Optional) For HA BYOL, the serial number for the second node. This is required when using 'ha-cot-premium-byol'.
* `node1_subnet_id` - (Optional) For HA, the subnet ID of the first node.
* `node2_subnet_id` - (Optional) For HA, the subnet ID of the second node.
* `mediator_subnet_id` - (Optional) For HA, the subnet ID of the mediator.
* `mediator_key_pair_name` - (Optional) For HA, the key pair name for the mediator instance.
* `cluster_floating_ip` - (Optional) For HA FloatingIP, the cluster management floating IP address.
* `data_floating_ip` - (Optional) For HA FloatingIP, the data floating IP address.
* `svm_floating_ip` - (Optional) For HA FloatingIP, the SVM management floating IP address.
* `route_table_ids` - (Optional) For HA FloatingIP, the list of route table IDs that will be updated with the floating IPs.
* `upgrade_ontap_version` - (Optional) Indicates whether to upgrade ontap image with `ontap_version`. To upgrade ontap image, `ontap_version` cannot be 'latest' and `use_latest_version` needs to be false.

The `aws_tag` block supports the following:
* `tag_key` - (Required) The key of the tag.
* `tag_value` - (Required) The tag value.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the working environment.
* `svm_name` - The name of the SVM.


## Terraform Variables

* `aws_profile` - (Optional) This is the profile name of the aws credentials file in your home directory, for example,~/.aws/credentials. If not specified, profile named default is used.
