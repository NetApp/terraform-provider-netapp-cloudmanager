output "cvo_id" {
  description = "Working environment ID of the CVO"
  value       = netapp-cloudmanager_cvo_azure.AzureLicenseTest72.id
}

output "cvo_name" {
  description = "Name of the CVO"
  value       = netapp-cloudmanager_cvo_azure.AzureLicenseTest72.name
}

output "license_type" {
  description = "Current license type"
  value       = netapp-cloudmanager_cvo_azure.AzureLicenseTest72.license_type
}

output "capacity_package" {
  description = "Current capacity package"
  value       = netapp-cloudmanager_cvo_azure.AzureLicenseTest72.capacity_package_name
}

output "instance_type" {
  description = "Current instance type"
  value       = netapp-cloudmanager_cvo_azure.AzureLicenseTest72.instance_type
}