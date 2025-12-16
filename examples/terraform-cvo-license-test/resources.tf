resource "netapp-cloudmanager_cvo_azure" "AzureLicenseTest72" {
  name                 = var.cvo_name
  location             = var.cvo_location
  client_id            = var.cvo_client_id
  subscription_id      = var.azure_subscription_id
  resource_group       = var.cvo_resource_group
  allow_deploy_in_existing_rg = true
  
  subnet_id            = var.cvo_subnet_id
  vnet_id              = var.cvo_vnet_id
  vnet_resource_group  = var.cvo_vnet_resource_group
  
  data_encryption_type = var.cvo_data_encryption_type
  storage_type         = var.cvo_storage_type
  disk_size            = var.cvo_disk_size
  disk_size_unit       = var.cvo_disk_size_unit
  capacity_tier        = var.cvo_capacity_tier
  is_ha                = var.cvo_is_ha

  svm_password          = var.cvo_svm_password
  license_type          = var.cvo_license_type
  # capacity_package_name = var.cvo_capacity_package_name
  instance_type         = var.cvo_instance_type
  writing_speed_state   = var.cvo_writing_speed_state
  ontap_version         = var.cvo_ontap_version
  use_latest_version    = var.cvo_use_latest_version

  saas_subscription_id  = var.cvo_saas_subscription_id
  open_security_group   = var.cvo_open_security_group

  azure_tag {
    tag_key   = "Creator"
    tag_value = "sr13947"
  }
  
  azure_tag {
    tag_key   = "KeepMe"
    tag_value = "9"
  }
  
  azure_tag {
    tag_key   = "KeepMeOn"
    tag_value = "9"
  }



}