---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_aws_fsx"
sidebar_current: "docs-netapp-cloudmanager-resource-aws-fsx"
description: |-
  Provides a netapp-cloudmanager_aws_fsx resource. This can be used to get Cloud ONTAP file system in AWS.
---

# netapp-cloudmanager_aws_fsx

Provides a netapp-cloudmanager_aws_fsx resource. This can be used to get Cloud ONTAP file system in AWS

## Example Usages

**Create netapp-cloudmanager_aws_fsx :**

```
data "netapp-cloudmanager_aws_fsx" "aws-fsx" {
  provider = netapp-cloudmanager
  id = "xxxxxxxxxxxx"
  tenant_id = "account-xxxxxxxx"
}
```


## Argument Reference

The following arguments are supported:

* `id` - (Required) The unique identifier for the working environment.
* `tenant_id` - (Required) The NetApp account ID that the Connector will be associated with.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the AWS FSX.
* `name` - The name of the AWS FSX.
* `region` - The region where the working environment will be created.
* `tenant_id` - The ID of the Cloud Manager workspace/tenant where you want to deploy Cloud Volumes ONTAP.
* `status` - The status of the AWS FSX.
* `lifecycle_status` - The lifecycle of the AWS FSX.

