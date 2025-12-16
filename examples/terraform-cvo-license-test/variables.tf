variable "cloudmanager_refresh_token" {
  description = "BlueXP refresh token"
  type        = string
}

variable "cloudmanager_environment" {
  description = "BlueXP environment (prod or stage)"
  type        = string
}

variable "azure_subscription_id" {
  description = "Azure subscription ID"
  type        = string
}

variable "cvo_id" {
  description = "Working environment ID (for import)"
  type        = string
}
variable "cvo_client_id" {
  description = "Connector client ID"
  type        = string
}
variable "cvo_name" {
  description = "CVO name"
  type        = string
}
variable "cvo_location" {
  description = "Azure region"
  type        = string
}
variable "cvo_resource_group" {
  description = "Resource group for CVO"
  type        = string
}
variable "cvo_vnet_id" {
  description = "VNet name"
  type        = string
}
variable "cvo_subnet_id" {
  description = "Subnet name"
  type        = string
}
variable "cvo_vnet_resource_group" {
  description = "Resource group for VNet"
  type        = string
}

variable "cvo_data_encryption_type" {
  description = "Data encryption type (AZURE or NONE)"
  type        = string
}
variable "cvo_storage_type" {
  description = "Storage type (e.g., Premium_LRS)"
  type        = string
}
variable "cvo_disk_size" {
  description = "Disk size"
  type        = number
}
variable "cvo_disk_size_unit" {
  description = "Disk size unit (GB or TB)"
  type        = string
}
variable "cvo_capacity_tier" {
  description = "Capacity tier (Blob or NONE)"
  type        = string
}
variable "cvo_is_ha" {
  description = "Is HA deployment"
  type        = bool
}

variable "cvo_svm_password" {
  description = "SVM admin password"
  type        = string
}
variable "cvo_license_type" {
  description = "License type"
  type        = string
}
# variable "cvo_capacity_package_name" {
#   description = "Capacity package name"
#   type        = string
# }
variable "cvo_instance_type" {
  description = "Instance type"
  type        = string
}
variable "cvo_writing_speed_state" {
  description = "Writing speed state"
  type        = string
}
variable "cvo_ontap_version" {
  description = "ONTAP version"
  type        = string
}
variable "cvo_use_latest_version" {
  description = "Use latest ONTAP version"
  type        = bool
}

variable "cvo_saas_subscription_id" {
  description = "SaaS/NSS subscription ID for capacity-based licensing"
  type        = string
  default     = ""
}

variable "cvo_open_security_group" {
  description = "Use open security group (testing only)"
  type        = bool
  default     = false
}

# variable "azure_tenant_id" {
#   description = "Azure Tenant ID"
#   type        = string
# }

# variable "azure_client_id" {
#   description = "Azure Client ID (App ID)"
#   type        = string
# }

# variable "azure_client_secret" {
#   description = "Azure Client Secret"
#   type        = string
#   sensitive   = true
# }