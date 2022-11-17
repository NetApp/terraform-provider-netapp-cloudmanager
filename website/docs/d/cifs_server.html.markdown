---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_cifs_server"
sidebar_current: "docs-netapp-cloudmanager-datasource-cifs-server"
description: |-
  Provides a netapp-cloudmanager_cifs_server resource. This can be used to read a CIFS server on the Cloud Volume ONTAP system that requires a CIFS volume, based on an Active Directory or Workgroup.
---

# netapp_cloudmanager_cifs_server

Provides a netapp-cloudmanager_cifs_server resource. This can be used to read a CIFS server on the Cloud Volume ONTAP system that requires a CIFS volume, based on an Active Directory or Workgroup.
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

## Argument Reference

The following arguments are supported:

* `client_id` - (Required) The client ID of the Cloud Manager Connector. You can find the ID from a previous create Connector action as shown in the example, or from the Connector tab on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `working_environment_id` - (Optional) The public ID of the working environment where the CIFS server will be created. This argument is optional if working_environment_name is provided. You can find the ID from a previous create Cloud Volumes ONTAP action as shown in the example, or from the information page of the Cloud Volumes ONTAP working environment on [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `working_environment_name` - (Optional) The working environment name where the CIFS server will be created. The argument will be ignored if working_environment_id is provided.
* `svm_name` - (Optional) The name of the SVM. 

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the SVM.
* `domain` - Active Directory domain name. For CIFS AD only.
* `dns_domain` - DNS domain name. For CIFS AD only.
* `ip_addresses` - DNS server IP addresses. For CIFS AD only.
* `netbios` - CIFS server NetBIOS name. For CIFS AD only.
* `organizational_unit` - Organizational Unit in which to register the CIFS server. For CIFS AD only.