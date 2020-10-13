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
* Create volumes all types of volumes: NFS, CIFS, and iSCSI

## Example Usage


# Configure the netapp-cloudmanager Provider
```
provider "netapp-cloudmanager" {
  refresh_token         = var.cloudmanager_refresh_token
}
```

## Argument Reference

The following arguments are used to configure the netapp-cloudmanager provider:

* `refresh_token` - (Required) This is the refresh token for NetApp Cloud Manager API operations. Get the token from [NetApp Cloud Central](https://services.cloud.netapp.com/refresh-token)

## Required Privileges

For additional information on roles and permissions, refer to [NetApp Cloud Manager documentation](https://docs.netapp.com/us-en/occm/).



