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
  aws_profile           = var.cloudmanager_aws_profile
  aws_profile_file_path = var.cloudmanager_aws_profile_file_path
  azure_auth_methods    = var.cloudmanager_azure_auth_methods
}
```

## Argument Reference

The following arguments are used to configure the netapp-cloudmanager provider:

* `refresh_token` - (Optional) This is the refresh token for NetApp Cloud Manager API operations. Get the token from [NetApp Cloud Central](https://services.cloud.netapp.com/refresh-token). If sa_client_id and sa_secret_key are provided, the service account will be used and this will be ignored.
* `sa_client_id` - (Optional) This is the service account client ID for NetApp Cloud Manager API operations. The service account can be created on [NetApp Cloud Central](https://services.cloud.netapp.com/). The client id and secret key will be provided on service account creation.
* `sa_secret_key` - (Optional) This is the service account client ID for NetApp Cloud Manager API operations. The service account can be created on [NetApp Cloud Central](https://services.cloud.netapp.com/). The client id and secret key will be provided on service account creation.
* `aws_profile` - (Optional) This is the profile name of the aws credentials file in your home directory, for example,~/.aws/credentials. If not specified, profile named default is used.
* `aws_profile_file_path` - (Optional) Path to the shared credentials file. Shortcuts like $HOME and ~ do not work.
* `azure_auth_methods` - (Optional) List of Azure authentication methods to be used: `env` for environment variables, `cli` for az login.  The methods are tried in sequence.  Defaults to `['cli, 'env']`.   Note that `env` can trigger a 404 BearerAuthorizer error if the credentials provided in the environment variables do not have the expected permissions.

## Configure AWS Credentials
AWS looks for credentials in the following orders:

1. Environment Variables
2. Shared Credentials file
3. Shared Configuration file (if SharedConfig is enabled)
4. EC2 Instance Metadata (credentials only)

If neither aws_profile nor aws_profile_file_path is specified, the provider look for cred in the mentioned order.
If one of aws_profile and aws_profile_file_path is specified, the unspecified option has default value:

If aws_profile_file_path is empty, it will look for "AWS_SHARED_CREDENTIALS_FILE" env variable. If the 
env value is empty will default to current user's home directory.
Linux/OSX: "$HOME/.aws/credentials"
Windows:   "%USERPROFILE%\.aws\credentials"

AWS Profile to extract credentials from the shared credentials file. If empty
will default to environment variable "AWS_PROFILE" or "default" if
environment variable is also not set.


## Configure Azure Credentials
### Option 1: Sign in with Azure CLI
`az login`
### Option 2: Define environment variables
Service principal with client secret
- `AZURE_CLIENT_ID` - ID of an Azure AD application
- `AZURE_TENANT_ID` - ID of the application's Azure AD tenant
- `AZURE_CLIENT_SECRET` - A client secret that was generated for the App Registration
- `AZURE_SUBSCRIPTION_ID` - Subscription identifier
```
export AZURE_TENANT_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
export AZURE_CLIENT_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
export AZURE_CLIENT_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export AZURE_SUBSCRIPTION_ID="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```

By default, the provider will try to authenticate with Azure the CLI (`az login`) and then using environment variables.   This can be configured with `azure_auth_methods`
(az login may set the env variables, so maybe this is redundant.)
## Required Privileges

For additional information on roles and permissions, refer to [NetApp Cloud Manager documentation](https://docs.netapp.com/us-en/occm/).



