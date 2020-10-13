---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cifs_server"
sidebar_current: "docs-netapp-cloudmanager-resource-cifs-server"
description: |-
  Provides a netapp-cloudmanager_cifs_server resource. This can be used to create or delete a CIFS server on the Cloud Volume ONTAP system that requires a CIFS volume, based on an Active Directory or Workgroup.
---

# netapp_cloudmanager_cifs_server

Provides a netapp-cloudmanager_cifs_server resource. This can be used to create or delete a CIFS server on the Cloud Volume ONTAP system that requires a CIFS volume, based on an Active Directory or Workgroup.
Requires existence of a Cloud Manager Connector and a Cloud Volumes ONTAP system.

## Example Usages

**Read netapp-cloudmanager_cifs_server:**

```
data "netapp-cloudmanager_cifs_server" "cvo-cifs" {
	provider = netapp-cloudmanager
   client_id = "AbCd6kdnLtvhwcgGvlFntdEHUfPJGc"
   working_environment_name = "CvoName"
}
```

**Create netapp-cloudmanager_cifs_server with AD:**

```
resource "netapp-cloudmanager_cifs_server" "cl-cifs" {
   provider = netapp-cloudmanager
   domain = "test.com"
   username = "admin"
   password = "abcde"
   dns_domain = "test.com"
   ip_addresses = ["1.0.0.1"]
   netbios = "cvoname"
   organizational_unit = "CN=Computers"
   client_id = "AbCd6kdnLtvhwcgGvlFntdEHUfPJGc"
   working_environment_name = "CvoName"
}
```

**Create netapp-cloudmanager_cifs_server with workgroup:**

```
resource "netapp-cloudmanager_cifs_server" "cl-cifs-wg" {
   provider = netapp-cloudmanager
   server_name = "server"
   workgroup_name  = "workgroup"
   client_id = netapp-cloudmanager_connector_gcp.cm-gcp.client_id
   working_environment_name = "CvoName"
   is_workgroup = true
}
```

## Argument Reference

The following arguments are supported:

* `working_environment_id` - (Optional) The public ID of the working environment where the CIFS server will be created. This argument is optional if working_environment_name is provided. You can find the ID from a previous create Cloud Volumes ONTAP action as shown in the example, or from the information page of the Cloud Volumes ONTAP working environment on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `working_environment_name` - (Optional) The working environment name where the CIFS server will be created. The argument will be ignored if working_environment_id is provided.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `domain` - (Optional) Active Directory domain name. For CIFS AD only.
* `username` - (Optional) Active Directory admin user name. For CIFS AD only.
* `password` - (Optional) Active Directory admin password. For CIFS AD only.
* `dns_domain` - (Optional) DNS domain name. For CIFS AD only.
* `ip_addresses` - (Optional) DNS server IP addresses. For CIFS AD only.
* `netbios` - (Optional) CIFS server NetBIOS name. For CIFS AD only.
* `organizational_unit` - (Optional) Organizational Unit in which to register the CIFS server. For CIFS AD only.
* `is_workgroup` - (Optional) For CIFS workgroup operations, set to true.
* `server_name` - (Optional) Server name. For CIFS workgroup only.
* `workgroup_name` - (Optional) Workgroup name. For CIFS workgroup only.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier of the working environment.

