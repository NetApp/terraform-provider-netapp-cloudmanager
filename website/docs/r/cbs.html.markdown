---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cbs"
sidebar_current: "docs-netapp-cloudmanager-resource-cbs"
description: |-
  Provides a netapp-cloudmanager_cbs resource. This can be used to enable or disable cloud backup on the volume and snapshot in the Cloud Volume ONTAP system.
---

# netapp-cloudmanager_cbs

Provides a netapp-cloudmanager_cbs resource. This can be used to enable cloud backup on a specific working environment Cloud Volumes ONTAP on AWS, Aure, and GCP.
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
  working_environment_id = netapp-cloudmanager_cvo_aws.cvo-aws.id
  client_id = netapp-cloudmanager_connector_aws.cm-aws.client_id
}

resource "netapp-cloudmanager_cbs" "zure-cbs" {
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
    object_lock = "NONE
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
  working_environment_id = netapp-cloudmanager_cvo_azure.cvo-azure.id
  client_id = netapp-cloudmanager_connector_azure.cm-azure.client_id
}

resource "netapp-cloudmanager_cbs" "gcp-cbs" {
  provider = netapp-cloudmanager
  cloud_provider = "GCP"
  region = "us-east4"
  account_id = netapp-cloudmanager_connector_gcp.cl-occm-gcp.account_id
  gcp_cbs_parameters {
    project_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.project_id
  }
  backup_policy {
    name = "xyz"
    object_lock = "NONE"
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
  client_id = netapp-cloudmanager_connector_gcp.cl-occm-gcp.client_id
  working_environment_id = netapp-cloudmanager_cvo_gcp.cvo-gcp.id
}
```

## Argument Reference

The following arguments are supported:

* `working_environment_id` - (Optional) The public ID of the working environment where the aggregate will be created. This argument is optional if working_environment_name is provided. You can find the ID from a previous create Cloud Volumes ONTAP action as shown in the example, or from the information page of the Cloud Volumes ONTAP working environment on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Optional) The working environment name where the aggregate will be created. This argument will be ignored if working_environment_id is provided.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `account_id` - (Required) The NetApp account ID that the backup cloud will be associated with. You can find the account ID in the account tab of Cloud Manager at [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `cloud_provider` - (Required) Need to be one of ['AWS', 'AZURE', 'GCP']
* `region` - (Required) The region where the working environment created.
* `bucket`- (Optional)
* `ip_space` - (Optional)
* `backup_policy` - (Optional)
* `auto_backup_enabled` - (Optional) Auto backup all volumes in working environments.
* `max_transfer_rate` - (Optional) Modifies node level throttling of an ONTAP cluster. Value to be specified in kilo bytes per second(kbps). A value of 0 implies Unlimited throttling.
* `export_existing_snapshots` - (Optional) Export pre-existing Snapshot copies to object storage
* `aws_cbs_parameters` - (Optional)
* `azure_cbs_parameters` - (Optional)
* `gcp_cbs_parameters` - (Optional)

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

The `gcp_cbs_parameters` block supports the following:
* `project_id` - (Required) The ID of the GCP project.
* `access_key` - (Optional) The GCP access key.
* `secret_password` - (Optional) The GCP secret password.
* `kms_key_ring_id` - (Optional) GCP KMS key ring ID.
* `kms_crypto_key_id` - (Optional) GCP KMS crypto key ID.

The `backup_policy` block suports the followings:
* `name` - (Required)
* `policy_rules` - (Optional)
* `archive_after_days` - (Optional)
* `object_lock` - (Optional) For AWS, DataLock and Ransomware Protection can be enabled in the "GOVERNANCE" mode or "COMPLIANCE" mode. For Azure, DataLock and Ransomware Protection can be enabled in the "UNLOCKED" mode or "LOCKED" mode.

The `policy_rules` block suports the folowings:
* `rule` - (Optional)

The `rule` blocks support the followings:
* `label` - (Optional) ['Hourly', 'Daily', 'Weekly', 'Monthly', 'Yearly']
* `retention` - (Optional) The number value goes with the `label`