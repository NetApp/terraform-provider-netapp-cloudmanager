---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_nss_account"
sidebar_current: "docs-netapp-cloudmanager-resource-nss-account"
description: |-
  Provides a netapp-cloudmanager_nss_account resource. This can be used to create or delete a NetApp Support Site account on the Cloud Manager system.
---

# netapp_cloudmanager_nss_account

Provides a netapp-cloudmanager_nss_account resource. This can be used to create or delete a NetApp Support Site account on the Cloud Manager system.

## Example Usages

**Read netapp-cloudmanager_nss_account:**

```
data "netapp-cloudmanager_nss_account" "nss-account-1" {
		provider = netapp-cloudmanager
		client_id = "Nw4Q2O1kdnLtvhwegGalFnodEHUfPJWh"
		username = "user"
	}
```

**Create netapp-cloudmanager_nss_account:**

```
resource "netapp-cloudmanager_nss_account" "nss-account-2" {
   provider = netapp-cloudmanager
   client_id = "AbCd6kdnLtvhwcgGvlFntdEHUfPJGc"
   username = "user"
   password = "pass"
}
```

## Argument Reference

The following arguments are supported:

* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `username` - (Required) NSS username. Not required in data source.
* `password` - (Required) NSS password. Not required in data source.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier of the account.

