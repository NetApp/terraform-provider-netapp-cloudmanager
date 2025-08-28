---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cbs"
sidebar_current: "docs-netapp-cloudmanager-resource-cbs"
description: |-
  Provides a netapp-cloudmanager_cbs resource. This can be used to enable or disable cloud backup on the volume and snapshot in the Cloud Volume ONTAP system.
---

# netapp-cloudmanager_cbs

Provides a netapp-cloudmanager_cbs resource. This can be used to enable cloud backup on a specific working environment Cloud Volumes ONTAP on AWS and Azure.
Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.

## Example Usages

**Create netapp-cloudmanager_aggregate:**

```
resource "netapp-cloudmanager_cbs" "aws-cbs" {
  provider = netapp-cloudmanager
  cloud_provider = "AWS"
  region = netapp-cloudmanager_cvo_aws.cvo-aws.region
  account_id = "account-xxxxxx"
  aws_cbs_parameters {
    aws_account_id = "123452054321"
    archive_storage_class = "GLACIER"
  }
  backup_policy {
    name = "abc"
    object_lock = "GOVERNANCE"
    policy_rules {
      rule {
        label = "Daily"
        retention = "30"
      }
      rule {
        label = "Weekly"
        retention = "4"
      }
    }
  }
  volumes {
    volume_name = "test"
    mode = "SCHEDULED"
    backup_policy {
      name = "xxxxxxx"
    }
  }
  volumes {
    volume_name = "test2"
    mode = "SCHEDULED"
    backup_policy {
      name = "xxxxxxx"
    }
  }
  working_environment_id = netapp-cloudmanager_cvo_aws.cvo-aws.id
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id
}

resource "netapp-cloudmanager_cbs" "azure-cbs" {
  provider = netapp-cloudmanager
  cloud_provider = "AZURE"
  region = netapp-cloudmanager_cvo_azure.cvo-azure.location
  account_id = "account-xxxxxxxx"
  azure_cbs_parameters {
    resource_group = "xxxxxxxxx"
    subscription = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    key_vault_id = "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/xxxxxxxxx/providers/Microsoft.KeyVault/vaults/xxxxxx"
    key_name: "xxxxxx"
  }
  backup_policy {
    name = "def"
    object_lock = "NONE"
    policy_rules {
      rule {
        label = "Daily"
        retention = "30"
      }
      rule {
        label = "Hourly"
        retention = "24"
      }
    }
  }
  volumes {
    volume_name = "test"
    mode = "SCHEDULED"
    backup_policy {
      name = "xxxxxxx"
    }
  }
  volumes {
    volume_name = "test2"
    mode = "SCHEDULED"
    backup_policy {
      name = "xxxxxxx"
    }
  }
  working_environment_id = netapp-cloudmanager_cvo_azure.cvo-azure.id
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
}
```

## Argument Reference

Arguments marked with “Forces new resource” will cause the resource to be recreated if their value is changed after creation.

The following arguments are supported:

* `working_environment_id` - (Optional, Forces new resource) The public ID of the working environment where the aggregate will be created. This argument is optional if working_environment_name is provided. You can find the ID from a previous create Cloud Volumes ONTAP action as shown in the example, or from the information page of the Cloud Volumes ONTAP working environment on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Optional, Forces new resource) The working environment name where the aggregate will be created. This argument will be ignored if working_environment_id is provided.
* `client_id` - (Required, Forces new resource) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `account_id` - (Required, Forces new resource) The NetApp account ID that the backup cloud will be associated with. You can find the account ID in the account tab of Cloud Manager at [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `cloud_provider` - (Required, Forces new resource) Need to be one of ['AWS', 'AZURE', 'GCP']
* `region` - (Required, Forces new resource) The region where the working environment created.
* `bucket`- (Optional, Forces new resource)
* `ip_space` - (Optional, Forces new resource)
* `backup_policy` - (Optional)
* `auto_backup_enabled` - (Optional, Forces new resource) Auto backup all volumes in working environments.
* `max_transfer_rate` - (Optional, Forces new resource) Modifies node level throttling of an ONTAP cluster. Value to be specified in kilo bytes per second(kbps). A value of 0 implies Unlimited throttling.
* `export_existing_snapshots` - (Optional, Forces new resource) Export pre-existing Snapshot copies to object storage
* `aws_cbs_parameters` - (Optional, Forces new resource)
* `azure_cbs_parameters` - (Optional, Forces new resource)

The `aws_cbs_parameters` block supports the following:
* `aws_account_id` - (Optional) Required when the provider is AWS.
* `archive_storage_class` - (Optional) Required for AWS to specify which storage class to use for archiving.
* `access_key` - (Optional)
* `secret_password` - (Optional)
* `kms_key_id` - (Optional) Input field for a customer-managed key use case
* `private_endpoint_id` - (Optional)

The `azure_cbs_parameters` block supports the following:
* `resource_group` - (Required) The resource group name.
* `storage_account` - (Optional) The storage account.
* `subscription` - (Required) The subscription ID.
* `private_endpoint_id` - (Optional) The id can be found with private endpoints with JSON view in Azure.
* `key_vault_id` - (Optional) The id can be found with key vault JSON View in Azure. e.g. "/subscriptions/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/resourceGroups/xxxxxxxxx/providers/Microsoft.KeyVault/vaults/xxxxxx"
* `key_name` - (Optional) Key vault name.

The `backup_policy` block supports the followings:
* `name` - (Required)
* `policy_rules` - (Optional)
* `archive_after_days` - (Optional)
* `object_lock` - (Optional) For AWS, DataLock and Ransomware Protection can be enabled in the "GOVERNANCE" mode or "COMPLIANCE" mode. For Azure, DataLock and Ransomware Protection can be enabled in the "UNLOCKED" mode or "LOCKED" mode.

The `policy_rules` block supports the followings:
* `rule` - (Optional)

The `rule` blocks support the followings:
* `label` - (Optional, Forces new resource) ['Hourly', 'Daily', 'Weekly', 'Monthly', 'Yearly']
* `retention` - (Optional, Forces new resource) The number value goes with the `label`

The `volumes` block supports the followings:
* `volume_name` - (Required) Name of the volume to enable backup.
* `mode` - (Optional) type of mode to create snapshot copies.
* `backup_policy` - (Optional)

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the cloud backup service.
