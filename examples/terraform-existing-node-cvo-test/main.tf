# Test 2.1: Create node-based CVO with OLD provider
# resource "netapp-cloudmanager_cvo_azure" "existing_node_cvo" {
#   name                  = "ExistingNodeCVO01"
#   location              = "eastus2"
#   subscription_id       = "54b91999-b3e6-4599-908e-416e0b8516b3"
#   client_id             = "Mek5haDZqsYmnNvhcazcpNzgBXk1vXi8clients"
  
#   # Network
#   subnet_id             = "cvo_qasubnet4"
#   vnet_id               = "cvo_vnet4qa1"
#   vnet_resource_group   = "AdServer-rg"
#   resource_group        = "ExistingNodeCVO01-rg"
  
#   # Credentials
#   svm_password          = "NetApp123!"
  
#   # ✅ NODE-BASED LICENSE (allowed in old provider)
#   license_type          = "azure-cot-standard-paygo"
  
#   # Storage
#   data_encryption_type  = "AZURE"
#   storage_type          = "Premium_LRS"
#   disk_size             = 1
#   disk_size_unit        = "TB"
#   capacity_tier         = "NONE"
  
#   # Instance
#   instance_type         = "Standard_DS4_v2"
  
#   # ONTAP
#   ontap_version         = "latest"
#   use_latest_version    = true
#   writing_speed_state   = "NORMAL"
  
#   # HA
#   is_ha                 = false
  
#   # Resource Group
#   allow_deploy_in_existing_rg = true
  
#   # Tags
#   azure_tag {
#     tag_key   = "Creator"
#     tag_value = "sr13947"
#   }
  
#   azure_tag {
#     tag_key   = "KeepMe"
#     tag_value = "9"
#   }
  
#   azure_tag {
#     tag_key   = "KeepMeOn"
#     tag_value = "9"
#   }
# }

# # Output for tracking
# output "cvo_id" {
#   value = netapp-cloudmanager_cvo_azure.existing_node_cvo.id
# }

# output "license_type" {
#   value = netapp-cloudmanager_cvo_azure.existing_node_cvo.license_type
# }

# output "svm_name" {
#   value = netapp-cloudmanager_cvo_azure.existing_node_cvo.svm_name
# }

resource "netapp-cloudmanager_cvo_azure" "new_node_cvo" {
  name                  = "NewNodeCVO02"
  location              = "eastus2"
  subscription_id       = "54b91999-b3e6-4599-908e-416e0b8516b3"
  client_id             = "Mek5haDZqsYmnNvhcazcpNzgBXk1vXi8clients"  
  
  subnet_id             = "cvo_qasubnet4"
  vnet_id               = "cvo_vnet4qa1"
  vnet_resource_group   = "AdServer-rg"
  resource_group        = "NewNodeCVO02-rg"  
  
  svm_password          = "NetApp123!"
  
  license_type          = "capacity-paygo"
  
  data_encryption_type  = "AZURE"
  storage_type          = "Premium_LRS"
  disk_size             = 1
  disk_size_unit        = "TB"
  capacity_tier         = "NONE"
  
  instance_type         = "Standard_DS4_v2"
  
  ontap_version         = "latest"
  use_latest_version    = true
  writing_speed_state   = "NORMAL"
  
  is_ha                 = false
  
  allow_deploy_in_existing_rg = true
  
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
  
  azure_tag {
    tag_key   = "TestCase"
    tag_value = "NewNodeCVO"
  } 

  # lifecycle {
  #   ignore_changes = [license_type]
  # }

  
}

# # Outputs
# output "existing_cvo_id" {
#   value = netapp-cloudmanager_cvo_azure.existing_node_cvo.id
# }

# output "existing_license_type" {
#   value = netapp-cloudmanager_cvo_azure.existing_node_cvo.license_type
# }

# output "existing_svm_name" {
#   value = netapp-cloudmanager_cvo_azure.existing_node_cvo.svm_name
# }

# output "new_cvo_id" {
#   value = netapp-cloudmanager_cvo_azure.new_node_cvo.id
# }

# output "new_license_type" {
#   value = netapp-cloudmanager_cvo_azure.new_node_cvo.license_type
# }

# output "new_svm_name" {
#   value = netapp-cloudmanager_cvo_azure.new_node_cvo.svm_name
# }