---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cvo_aws"
sidebar_current: "docs-netapp-cloudmanager-datasource-cvo-aws"
description: |-
  Provides a netapp-cloudmanager_cvo_aws resource. This can be used to get AWS Cloud Volumes ONTAP.
---

# netapp-cloudmanager_cvo_aws

Provides a netapp-cloudmanager_cvo_aws resource. This can be used to get AWS Cloud Volumes ONTAP.

## Example Usages

**get netapp-cloudmanager_cvo_aws:**

```
data "netapp-cloudmanager_cvo_aws" "aws-cvo-1" {
  provider = netapp-cloudmanager
  name = "awsha"
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `name` - (Required) The name of the cvo aws.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `svm_name` - The name of the SVM.
* `id` - The id of this working environment.
