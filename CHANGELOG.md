## 21.11.1 (Unreleased)
ENHANCEMENTS:

* resource/cvo for AWS, AZURE and GCP: display svm_name in state file on the CVO creation

BUG FIXES:

* resource/cvo_aws route_table_ids parameter force recreation ([#77](https://github.com/NetApp/terraform-provider-netapp-cloudmanager/issues/77))
* resource/cvo_azure import function is disabled. Error out if terraform import is used