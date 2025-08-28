---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_aws_fsx"
sidebar_current: "docs-netapp-cloudmanager-resource-aws-fsx"
description: |-
  Provides a netapp-cloudmanager_aws_fsx resource. This can be used to create a new Cloud ONTAP file system in AWS.
---

# netapp-cloudmanager_aws_fsx

Provides a netapp-cloudmanager_aws_fsx resource. This can be used to create a new Cloud ONTAP file system in AWS

## Example Usages

**Create netapp-cloudmanager_aws_fsx :**

```
resource "netapp-cloudmanager_aws_fsx" "aws-fsx" {
  provider = netapp-cloudmanager
  name = "TerraformAWSFSX"
  region = "us-west-2"
  primary_subnet_id = "subnet-xxxxxxx"
  secondary_subnet_id = "subnet-xxxxxxx"
  tenant_id = "account-xxxxxxxx"
  workspace_id = "workspace-xxxxxxx"
  tags {
            tag_key = "abcd"
            tag_value = "ABCD"
        }
  tags {
            tag_key = "xxx"
            tag_value = "YYY"
        }
  fsx_admin_password = "P@assword!"
  throughput_capacity = 512
  storage_capacity_size = 1024
  storage_capacity_size_unit = "GiB"
  aws_credentials_name = "abcd"
}
```


## Argument Reference

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the working environment.
* `aws_credentials_name` - (Required, Forces new resource) The name of the AWS Credentials account name.
* `region` - (Required, Forces new resource) The region where the working environment will be created.
* `primary_subnet_id` - (Required, Forces new resource) For HA, the subnet ID of the first node.
* `secondary_subnet_id` - (Required, Forces new resource) For HA, the subnet ID of the second node.
* `fsx_admin_password` - (Required, Forces new resource) The admin password for working environment.
* `tenant_id` - (Required, Forces new resource) The NetApp account ID that the working environment will be associated with.
* `workspace_id` - (Required, Forces new resource) The ID of the Cloud Manager workspace of working environment.
* `kms_key_id` - (Optional, Forces new resource) AWS encryption parameters. It is required if using aws encryption.
* `minimum_ssd_iops` - (Optional, Forces new resource) Provisioned SSD IOPS.
* `storage_capacity_size` - (Optional, Forces new resource) EBS volume size for the first data aggregate. For GB, the unit can be: [100 or 500]. For TB, the unit can be: [1,2,4,8,16]. The default is '1'.
* `storage_capacity_size_unit` - (Optional, Forces new resource) ['GB' or 'TB']. The default is 'TB'.
* `throughput_capacity` - (Optional, Forces new resource) capacity of the throughput.
* `security_group_ids` - (Optional, Forces new resource) The ID of the security group for the working environment.
* `endpoint_ip_address_range` - (Optional, Forces new resource) The endpoint IP address range.
* `route_table_ids` - (Optional, Forces new resource) The list of route table IDs that will be updated with the floating IPs.
* `import_file_system` - (Optional, Forces new resource) bool option to existing import AWS file system to CloudManager. The default is 'false'.
* `file_system_id` - (Optional, Forces new resource) The AWS file system ID to import to CloudManager. Required when import_file_system is 'true'.

The `tags` block supports the following (Forces new resource):
* `tag_key` - (Required) The key of the tag.
* `tag_value` - (Required) The tag value.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the working environment.

