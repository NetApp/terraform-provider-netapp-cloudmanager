---
layout: "netapp_cloudmanager"
page_title: "NetApp_CloudManager: netapp_cloudmanager_connector_gcp"
sidebar_current: "docs-netapp-cloudmanager-resource-connector-gcp"
description: |-
  Provides a NetApp_CloudManager connector GCP resource. This can be used to create a new Cloud Manager Connector in GCP.
---

# netapp-cloudmanager_connector_gcp

Provides a NetApp_CloudManager connector GCP resource. This can be used to create a new Cloud Manager Connector in GCP.
In order to use that resource, you should have a service account key file. The minimum required policy can be found here: [Connector deployment policy for GCP](https://occm-sample-policies.s3.amazonaws.com/Setup_As_Service_3.7.3_GCP.yaml).
The service account for the Connector VM instance must have the permissions defined in [Cloud Manager policy for GCP](https://occm-sample-policies.s3.amazonaws.com/Policy_for_Cloud_Manager_3.8.0_GCP.yaml)

<!---
i think we need to create section for terraform and point to there
-->

## Example Usages

**Create netapp-cloudmanager_connector_gcp:**

```
resource "netapp-cloudmanager_connector_gcp" "cl-occm-gcp" {
  provider = netapp-cloudmanager
  name = "occm-gcp"
  project_id = "xxxxxxx"
  zone = "us-east4-b"
  company = "NetApp"
  service_account_email = "xxxxxxxxxxxxxxxx"
  service_account_path = "gcp_creds.json"
  account_id = "account-moKEW1b5"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, non-modifiable) The name of the Cloud Manager Connector.
* `project_id` - (Required, non-modifiable) The GCP project_id where the connector will be created.
* `zone` - (Required, non-modifiable) The GCP zone where the Connector will be created.
* `company` - (Required, non-modifiable) The name of the company of the user.
* `service_account_email` - (Required, non-modifiable) The email of the service_account for the connector instance. This service account is used to allow the Connector to create Cloud Volume ONTAP.
* `service_account_path` - (Optional, non-modifiable) The local path of the service_account JSON file for GCP authorization purposes. This service account is used to create the Connector in GCP. Conflicts with `service_account_key`. It does not need to be set when using the user application default credential authentication by running the command gcloud auth application-default login..
* `service_account_key` - (Optional, non-modifiable) A JSON string for the service account key for GCP authorization purposes. This service account is used to create the Connector in GCP. Conflicts with `service_account_path`. It does not need to be set when using the user application default credential authentication by running the command gcloud auth application-default login.
* `subnet_id` - (Optional, non-modifiable) The name of the subnet for the virtual machine. The default value is "Default". Two formats are supported: either <subnetID> or projects/<projectID>/regions/<region>/subnetworks/<subnetID>. 
* `network_project_id` - (Optional, non-modifiable) The project id in GCP associated with the Subnet. If not provided, it’s assumed that the Subnet is within the previously specified project id.
* `machine_type` - (Optional, non-modifiable) The machine_type for the Connector VM. The default value is "n2-standard-4"
* `firewall_tags` - (Optional, non-modifiable) Indicates whether to add firewall_tags to the connector VM (HTTP and HTTP). The default is "true". If true, "firewall-tag-bvsu", "http-server", "https-server" are set as tags, and they can't be managed using tag options.
* `associate_public_ip` - (Optional, non-modifiable) Indicates whether to associate a public IP address to the virtual machine. The default is "true"
* `proxy_url` - (Optional, non-modifiable) The proxy URL, if using a proxy to connect to the internet.
* `proxy_user_name` - (Optional, non-modifiable) The proxy user name, if using a proxy to connect to the internet.
* `proxy_password` - (Optional, non-modifiable) The proxy password, if using a proxy to connect to the internet.
* `proxy_certificates` - (Optional, non-modifiable) The proxy certificates. A list of certificate file names.
* `account_id` - (Optional, non-modifiable) The NetApp account ID that the Connector will be associated with. If not provided, Cloud Manager uses the first account. If no account exists, Cloud Manager creates a new account. You can find the account ID in the account tab of Cloud Manager at [https://cloudmanager.netapp.com](https://cloudmanager.netapp.com).
* `tags` - (Optional) The list of network tags.
* `labels` - (Optional) The map of labels.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The name of the virtual machine.
* `client_id` - The unique client ID of the Connector. Can be used in other resources.
* `account_id` - The NetApp tenancy account ID.

