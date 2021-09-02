---
layout: "netapp-cloudmanager"
page_title: "Provider: NetApp_CloudManager"
sidebar_current: "docs-netapp-cloudmanager-index"
description: |-
  The netapp-cloudmanager provider is used to interact with NetApp Cloud Manager in order to create and manage Cloud Volumes ONTAP in AWS, Azure, and GCP. The provider needs to be configured with the proper credentials before it can be used.
---

# netapp-cloudmanager Provider

The netapp-cloudmanager provider is used to interact with NetApp Cloud Manager in order to create and manage Cloud Volumes ONTAP in AWS, Azure, and GCP.
The provider needs to be configured with the proper credentials before it can be used.


Use the navigation to the left to read about the available resources.

~> **NOTE:** The netapp-cloudmanager provider currently represents _initial support_
and therefore may undergo significant changes as the community improves it.

### The following actions are supported in all cloud providers (AWS, Azure, and GCP):
* Create a Cloud Manager Connector
* Create a Cloud Volumes ONTAP system (single node or HA pair)
* Create aggregates
* Create a CIFS server to enable CIFS volume creation
* Create volumes of any type: NFS, CIFS, or iSCSI
* Create snapmirror relationship
* Create Netapp Support Site account
* Create a AWS working environment for FSX

## Example Usage


# Configure the netapp-cloudmanager Provider
```
provider "netapp-cloudmanager" {
  refresh_token         = var.cloudmanager_refresh_token
  sa_secret_key         = var.cloudmanager_sa_secret_key
  sa_client_id          = var.cloudmanager_sa_client_id
}
```

## Argument Reference

The following arguments are used to configure the netapp-cloudmanager provider:

* `refresh_token` - (Optional) This is the refresh token for NetApp Cloud Manager API operations. Get the token from [NetApp Cloud Central](https://services.cloud.netapp.com/refresh-token). If sa_client_id and sa_secret_key are provided, the service account will be used and this will be ignored.
* `sa_client_id` - (Optional) This is the service account client ID for NetApp Cloud Manager API operations. The service account can be created on [NetApp Cloud Central](https://services.cloud.netapp.com/). The client id and secret key will be provided on service account creation.
* `sa_secret_key` - (Optional) This is the service account client ID for NetApp Cloud Manager API operations. The service account can be created on [NetApp Cloud Central](https://services.cloud.netapp.com/). The client id and secret key will be provided on service account creation.

## Required Privileges

For additional information on roles and permissions, refer to [NetApp Cloud Manager documentation](https://docs.netapp.com/us-en/occm/).



