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

* `working_environment_id` - (Optional) The public ID of the working environment where the CIFS server will be created. This argument is optional if working_environment_name is provided. You can find the ID from a previous create Cloud Volumes ONTAP action as shown in the example, or from the information page of the Cloud Volumes ONTAP working environment on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `working_environment_name` - (Optional) The working environment name where the CIFS server will be created. The argument will be ignored if working_environment_id is provided.
* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://console.bluexp.netapp.com/](https://console.bluexp.netapp.com/).
* `domain` - (Required) Active Directory domain name. For CIFS AD only.
* `username` - (Required) Active Directory admin user name. For CIFS AD only.
* `password` - (Required) Active Directory admin password. For CIFS AD only.
* `dns_domain` - (Required) DNS domain name. For CIFS AD only.
* `ip_addresses` - (Required) DNS server IP addresses. For CIFS AD only.
* `netbios` - (Required) CIFS server NetBIOS name. For CIFS AD only.
* `organizational_unit` - (Required) Organizational Unit in which to register the CIFS server. For CIFS AD only.
* `svm_name` - (Optional) The name of the SVM. API will use the svmName from the CVO if it is not provided here.
* `is_workgroup` - (Deprecated) For CIFS workgroup operations, set to true. Creating cifs server with workgroup is deprecated.
* `server_name` - (Deprecated) Server name. For CIFS workgroup only. Creating cifs server with workgroup is deprecated.
* `workgroup_name` - (Deprecated) Workgroup name. For CIFS workgroup only. Creating cifs server with workgroup is deprecated.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the SVM.

