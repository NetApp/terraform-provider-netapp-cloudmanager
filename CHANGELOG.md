## 22.5.0
NEW FEATURES:

* resource/connector_aws: support returning public_ip_address of the aws connector.

## 22.4.0
BUG FIXES:

* resource/connector_gcp: Support shared vpc.

## 22.2.2
BUG FIXES:

* Support resources operating in parallel
* Allow existing by Node license on CVO creation

ENHANCEMENTS:

* Update the default license_type and capacity_package_name for the CVOs of AWS, AZURE and GCP
* Add snapmirror example

NEW FEATURES:

* resource/connector_gcp: add service_account_key option
* resource/connector_gcp: remove requirement for GCP service account JSON. Support authentication using User Application Default Credentials ("ADCs") as the authentication method. Enable ADCs by running the command gcloud auth application-default login.
* resource/cvo_aws: new option mediator_instance_profile_name.
* resource/aws_fsx: import AWS FSX to CloudManager.
* resource/connector_azure: add storage_account option. Support user defined storage account.


## 22.1.1
BUG FIXES:

* resource/cvo for AWS, AZURE and GCP: add status check on instance_type update
* resource/aws_fsx: validation errors detected on aws_fsx tags

## 22.1.0
ENHANCEMENTS:

* resource/cvo for AWS, AZURE and GCP: add upgrade_ontap_version for ontap_version upgrade
* add Terraform variable aws_profile to support use of profile name in aws credentials file
* resource/snapmirror: add FSX AWS to snapmirror
* resource/volume: add snapshot_policy_name and tiering_policy modification, and check the supporting changeable parameters
* resource/connector_gcp: add network tags option

## 21.12.0
NEW FEATURES:

* resource/aws_fsx_volume: create, update and delete FSx volume.
* resource/cvo_onnprem: This can be used to register an onprem ONTAP system into CloudManager.

## 21.11.1
ENHANCEMENTS:

* resource/cvo for AWS, AZURE and GCP: display svm_name in state file on the CVO creation
* resource/aws_fsx: add name tag as the basic tag on aws_fsx creation

BUG FIXES:

* resource/cvo_aws route_table_ids parameter force recreation ([#77](https://github.com/NetApp/terraform-provider-netapp-cloudmanager/issues/77))
* resource/cvo_azure import function is disabled. Error out if terraform import is used
