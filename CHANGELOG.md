## 25.4.0
IMPROVEMENTS:
* resource/connector_gcp: replaced GCP Deployment Manager API with individual GCP Compute Engine APIs for VM and disk management.

ENHANCEMENTS:
* Update all the resources documenation by adding `Forces new resource` if the modification is not supported.


## 25.3.0
NEW FEATURES:
* resource/volume supports import.

BUG FIXES:
* resource/volume: Fixed an issue where `export_policy_rule_super_user` was incorrectly updated to `none` regardless of the specified value.
* resource/volume: Fixed incorrect references to parameters in documentation.

## 25.2.0
NEW FEATURES:
* resource/cvo supports Restricted mode for GCP.
* resource/volume supports Restricted mode for GCP.
* resource/snapmirror supports Restricted mode for GCP.
* resource/aggregate supports Restricted mode for GCP.

BUG FIXES:
* add exponential backoff retries to task status check for handling status 504 case.

## 25.1.0
BUG FIXES:
* Update document with proper version terraform 1.1.0 and Go 1.21.
* resource/aggregate: fix `number_of_disks` update failure.
* Fix document typo and aggregate disk size description.
* resource/connector_gcp: Fix error in creation of restricted mode when `associate_public_ip` is false.

## 24.11.3
ENHANCEMENTS:
* Update GCP storage package to support GCP identity federation. This version requires terraform 1.1 and the Go 1.21.
* Update document: indicate the minimum required terraform version.

BUG FIXES:
* resource/connector_gcp: Fix schema structure while creating Restricted mode.

## 24.11.2
BUG FIXES:
* add `azure_tag` option in documentation.

## 24.11.1
BREAKING CHANGE:
* resource/connector_aws: update `instance_type` default value from `t3.xlarge` to `t3.2xlarge`
* resource/connector_azure: update `virtual_machines_size` default value from `Standard_DS3_v2` to `Standard_D8s_v3`
* resource/connector_gcp: update `machine_type` default value from `n2-standard-4` to  `n2-standard-8`

## 24.11.0
ENHANCEMENTS:
* resource/connector_azure: adding `azure_tag` option, now supports tags.

NEW FEATURES:
* Azure and GCP connectors now support Restricted mode.

BUG FIXES:
* auth user accesToken: Fix 403 issue with authorizer API token.

## 24.5.1
ENHANCEMENTS:
* remove duplicated volume documentation page.

## 24.5.0
BUG FIXES:
* resoruce/volume: support `comment` update with adding 3 minutes wait time.

## 24.2.0
NEW FEATURES:
* resource/connector_aws: support `instance_metadata` block.

## 24.1.0
ENHANCEMENTS:
* resource/cvo_gcp: fix typo on `vpc3_firewall_rule_tag_name`.
* add logging to API calls.


## 23.11.0
*ENHANCEMENTS: add retries when http returns >200 status code in getWorkingEnvironmentProperties. 

## 23.10.0
BUG FIXES:
* resource/aws_fsx: handling for a situation in which the status does not exist yet.
* resource/cifs_server: fix the read function on domain, dns_domain and netbios checking with case insensitive.
* resource/volume: add `export_policy_rule_super_user` and `export_policy_rule_access_control` options. Fix export policy update error.

## 23.8.2
BUG FIXES:
* resource/volume: fix schema structure for export policy response in volume.

ENHANCEMENTS:
* support dev environment.
* resources/cvo: update the documentation on the `ontap_version` naming convention and the way to find the available upgrade versions on `upgrade_ontap_version`.

## 23.8.1
NEW FEATURES:
* resource/cvo_gcp: support LDM/flashCache on both single and HA.

BUG FIXES:
* resource/connector_gcp: fix gcp config flags backward compatible issue.

## 23.8.0
BUG FIXES:
* resource/volume: fix documentation name for volume and add an example for creating on_prem volume.
* ressource/cvo_aws, cvo_azure, cvo_gcp: remove force new from `retries`.

NEW ENHANCEMENTS:
* resource/cvo_aws and cvo_gcp: add `saas_subscription_id`.

NEW FEATURES:
* resource/cvo_gcp: support adding `firewall_tag_name_rule` and `firewall_ip_ranges`.

## 23.7.0
NEW FEATURES:
* resource/connector_gcp: support adding gcp keys `gcp_block_project_ssh_keys`, `gcp_serial_port_enable`, `gcp_enable_os_login` and `gcp_enable_os_login_sk` to the config.

## 23.6.1
NEW FEATURES:
* resource/cvo_gcp: support HIGH `writing_speed_state` in GCP HA. Make `gcp_service_account` optional.

## 23.6.0
NEW FEATURES:
* resource/aws_fsx_volume: add new option`tags`.

## 23.5.2
NEW ENHANCEMENTS:
* resource/cvo_azure: add `saas_subscription_id`.
* provider: add `connector_host` for restricted mode.

## 23.5.1
NEW ENHANCEMENTS:
* resoruce/cvo_aws: add `assume_role_arn` for AWS CVO HA.

## 23.5.0
BUG FIXES:
* resource/cbs: fix cbs to work without volume.

## 23.5.0
BUG FIXES:
* resource/connector_gcp: fix label key and value do not acccpt only numeric value.
* resource/volume: disable `enable_deduplication` and `enable_compression` in read function.

NEW FEATURES:
* resource/cbs: new resource cloud backup service which supports AWS a nd Azure.
* resource/volume: add new option`comment`.

BUG FIXES:
* resource/connector for GCP: Make `service_account_key` optional.

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
* resource/cvo_azure: add `availability_zone_node1` and `availability_zone_node2` to support HA deployment.
* resoruce/cvo_azure: add new support value "Premium_ZRS" in parameter `storage_type`.

## 22.9.1
NEW FEATURES:
* resource/connector_snapmirror: support fsx as a source for snapmirror relationship with fsx/onprem.

NEW ENHANCEMENTS:
* resource/cvo_azure: add `availability_zone` parameter for single node deployment.
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
