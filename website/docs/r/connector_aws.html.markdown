---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_connector_aws"
sidebar_current: "docs-netapp-cloudmanager-resource-connector-aws"
description: |-
  Provides a NetApp_CloudManager Connector AWS resource. This can be used to create a new Cloud Manager Connector in AWS.
---

# netapp-cloudmanager_connector_aws

Provides a NetApp_CloudManager Connector AWS resource. This can be used to create a new Cloud Manager Connector in AWS.
The environment needs to be configured with the proper credentials before it can be used (run this command: aws configure).
The minimum required policy can be found at [Connector deployment policy for AWS](https://s3.amazonaws.com/occm-sample-policies/Policy_for_Setup_As_Service.json)

<!---
i think we need to create section for terraform and point to there
-->

## Example Usages

**Create NetApp_CloudManager aws:**

```
resource "netapp-cloudmanager_connector_aws" "cl-occm-aws" {
  provider = netapp-cloudmanager
  name = "TF-ConnectorAWS"
  region = "us-west-1"
  key_name = "automation_key"
  company = "NetApp"
  instance_type = "t3.xlarge"
  aws_tag {
              tag_key = "abcd"
              tag_value = "ABCD"
            }
  aws_tag {
              tag_key = "xxx"
              tag_value = "YYY"
            }
  instance_metadata{
    http_put_response_hop_limit = 2
    http_tokens = "required"
    http_endpoint = "enabled"
  }
  subnet_id = "subnet-xxxxx"
  security_group_id = "sg-xxxxxxxxx"
  iam_instance_profile_name = "OCCM_AUTOMATION"
  account_id = "account-ABCNJGB0X"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Manager Connector.
* `region` - (Required) The region where the Cloud Manager Connector will be created.
* `company` - (Required) The name of the company of the user.
* `key_name` - (Required) The name of the key pair to use for the Connector instance.
* `instance_type` - (Required) The type of instance (for example, t3.xlarge). At least 4 CPU and 16 GB of memory are required.
* `subnet_id` - (Required) The ID of the subnet for the instance.
* `security_group_id` - (Required) The ID of the security group for the instance, multiple security groups can be provided separated by ','.
* `iam_instance_profile_name` - (Required) The name of the instance profile for the Connector.
* `proxy_url` - (Optional) The proxy URL, if using a proxy to connect to the internet.
* `proxy_user_name` - (Optional) The proxy user name, if using a proxy to connect to the internet.
* `proxy_password` - (Optional) The proxy password, if using a proxy to connect to the internet.
* `proxy_certificates` - (Optional) The proxy certificates. A list of certificate file names.
* `associate_public_ip_address` - (Optional) Indicates whether to associate a public IP address to the instance. If not provided, the association will be done based on the subnet's configuration.
* `enable_termination_protection` - (Optional) Indicates whether to enable termination protection on the instance, default is false.
* `account_id` - (Optional) The NetApp account ID that the Connector will be associated with. If not provided, Cloud Manager uses the first account. If no account exists, Cloud Manager creates a new account. You can find the account ID in the account tab of Cloud Manager at [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `instance_metadata` - (Optional,Computed) The block of AWS EC2 instance metadata.

The `aws_tag` block supports the following:
* `tag_key` - (Required) The key of the tag.
* `tag_value` - (Required) The tag value.

The `instance_metadata` block supports the following:
* `http_endpoint` - (Optional, Computed) If the value is disabled, you cannot access your instance metadata. Choices: ["enabled", "disabled"]
* `http_tokens` - (Optional, Computed) Indicates whether IMDSv2 is required. Choices: ["optional", "required"]
* `http_put_response_hop_limit` - (Optional, Computed) The desired HTTP PUT response hop limit for instance metadata requests. The larger the number, the further instance metadata requests can travel. Possible values: Integers from 1 to 64.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The EC2 instance ID.
* `client_id` - The unique client ID of the Connector. Can be used in other resources.
* `account_id` - The NetApp tenancy account ID.
* `public_ip_address` - The public IP of the connector.

## Unique id versus name

With netapp-cloudmanager_connector_aws, every resource has a unique ID, but names are not necessarily unique.

## Connector Import
The id used to import is constructed with two attributes: client id and connector id. The format is CLIENT_ID:CONNECTOR_ID

