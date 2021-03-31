---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_nss_account"
sidebar_current: "docs-netapp-cloudmanager-resource-nss-account"
description: |-
  Provides a netapp-cloudmanager_nss_account resource. This can be used to read a NetApp Support Site account on the Cloud Manager system.
---

# netapp_cloudmanager_nss_account

Provides a netapp-cloudmanager_nss_account resource. This can be used to read a NetApp Support Site account on the Cloud Manager system.

## Example Usages

**Read netapp-cloudmanager_nss_account:**

```
data "netapp-cloudmanager_nss_account" "nss-account-1" {
		provider = netapp-cloudmanager
		client_id = "Rw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		username = "user"
	}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `username` - (Required) The user name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier of the account.
