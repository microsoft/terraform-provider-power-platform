output "policy_system_id" {
  description = "The system ID of the enterprise policy created in Azure"
  value       = azapi_resource.powerplatform_policy.output.properties.systemId
}

output "policy_id" {
  description = "The ID of the enterprise policy resource in Azure"
  value       = azapi_resource.powerplatform_policy.output.id
}

output "policy_resource" {
  description = "The enterprise encryption policy resource"
  value       = powerplatform_enterprise_policy.encryption
}