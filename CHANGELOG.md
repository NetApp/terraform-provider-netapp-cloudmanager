## 23.3.4
NEW FEATURES:
* resource/cbs: new resource cloud backup service.

BUG FIXES:
* resource/connector for GCP: Make service_account_key optional.

## 23.3.3
BUG FIXES:
* resource/connector: revert the retries change
* resource/cvo: revert the update values on the retries parameter

## 23.3.2
BUG FIXES:
* Documenation formate fix

NEW ENHANCEMENTS:
* resource/cvo for AZURE, GCP and AWS: change `retries` default to 75 and 100 for HA.
* resource/connector for AZURE, GCP and AWS: add `retries` option.

## 23.3.0
NEW ENHANCEMENTS:
* resource/cvo for AZURE and GCP: support two new capacity based package 'Edge' and 'Optimized' in `capacity_package_name`

BUG FIXES:
* cifs_server: fix bug on reading cifs based on the proper svm

## 23.1.1
NEW ENHANCEMENTS:
* BlueXP update domains adjustment
* Update AZURE document

BUG FIXES:
* resource/cvo_gcp:`zone` is not required in HA case. `node1_zone` will be used when `zone` is not provided in HA.
* resource/cvo_volume: update the volume with the proper `svm_name`

## 23.1.0
NEW FEATURES:
* resource/cvo_volume: add `tags` option.
* resource/cvo for AWS, AZURE and GCP: add `worm_retention_period_length` and `worm_retention_period_unit` to support WORM for creating a CVO.

BUG FIXES:
* resource/azure_cvo: fix bug on the `azureEncryptionParameters` with proper type.

## 22.12.0
NEW FEATURES:
* resource/cvo_volume: support create and delete onPrem volume.
* resource/cvo_volume: support create snapshot policy for AWS, AZURE and GCP if the snapshot policy is not available.

## 22.11.0
NEW FEATURES:
* resource/cvo_gcp HA: support add, rename and delete SVMs.
* resource/connector_gcp: add `labels` option.

BUG FIXES:
* resource/cvo_gcp: both capacity_tier and tierl_level should be optional.
* cifs_server on resource and data source: CIFS server with workgroup is depreciated. Since creating CIFS server with AD is the only way, updated the param attributes accordingly.

## 22.10.0
NEW FEATURES:
* resource/connector_snapmirror: support fsx as a source for snapmirror relationship with cvo.
* resource/cvo_aws: add `retries` parameter to increase wait time when creating CVO.
* resource/cvo_azure: add `retries` parameter to increase wait time when creating CVO.
* resource/cvo_gcp: add `retries` parameter to increase wait time when creating CVO.
* resoruce/cvs for AWS, AZURE and GCP: add `svm_name` an optional parameter. The modification is supported.

NEW ENHANCEMENTS:
* resource/connector_azure: display the deployed virtual machine principal_id in state file on the connector azure creation.
* resource/cvo_azure: add availability_zone_node1 and availability_zone_node2 to support HA deployment.
* resoruce/cvo_azure: add new support value "Premium_ZRS" in paramter storage_type.

## 22.9.1
NEW FEATURES:
* resource/connector_snapmirror: support fsx as a source for snapmirror relationship with fsx/onprem.

NEW ENHANCEMENTS:
* resource/cvo_azure: add availability_zone parameter for single node deployment.
* Use sensitive flag on the password of each resource.

BUG FIXES:
* azure: change default authentication to `['cli', 'env']` to give priority to `az login`.

## 22.9.0
NEW FEATURES:
* add retries on the task status check for handling status 504 cases
* azure: support ENV variables in addition to `az login`.  Added `azure_auth_methods` to define which methods to use.

BUG FIXES:
* resource/connector_gcp: update machine_type default value.
* resource/connector for GCP: validate `zone` is not empty or a empty space.

## 22.8.3
NEW FEATURES:
* resource/cvo_aws: add mediator_security_group_id option.

NEW ENHANCEMENTS:
* resource/cvo for AWS: add cluster_key_pair_name parameter for key pair on SSH authentication.

## 22.8.2
NEW ENHANCEMENTS:
* resource/cvo for AWS, AZURE and GCP: support backup_volumes_to_cbs and enable_compliance.
* resource/connector for AWS, AZURE and GCP: wait for creation to complete increase to 15 minutes.

BUG FIXES:
* resource/volume: remove default values of enable_thin_provisioning, enable_compression and enable_deduplication.

## 22.8.1
BUG FIXES:
* resource/connector_aws: fix bug whe get instance returns error, but error is not returned to upstream.

## 22.8.0
NEW FEATURES:

* resource/connector_aws: support returning public_ip_address of the aws connector.
* add Terraform variable aws_profile_file_path to specify aws credentials file location. ([#90](https://github.com/NetApp/terraform-provider-netapp-cloudmanager/issues/90))

ENHANCEMENTS:
* resource/connector_azure: support full subnet_id and vnet_id
* resource/aggregate: better handle multi creation with count
* resource/cvo for AWS, AZURE and GCP: add writing_speed_state update

BUG FIXES:
* resource/cvo for AWS, AZURE and GCP: fix force recreation on writing_speed_state  ([#104](https://github.com/NetApp/terraform-provider-netapp-cloudmanager/issues/104))

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
