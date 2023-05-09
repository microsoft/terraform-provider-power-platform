output "app_user_id" {
  description = "App User ID"
  value       = azuread_service_principal.app_user_principal.application_id
}

output "user_name" {
  description = "AAD user name"
  value       = azuread_user.user.user_principal_name
}

output "user_pass" {
  description = "AAD user password"
  value       = azuread_user.user.password
  sensitive   = true
}

output "environment_name" {
  description = "Power Platform environment name"
  value       = powerplatform_environment.environment.display_name
}

output "environment_id" {
  description = "Power Platform environment ID"
  value       = powerplatform_environment.environment.environment_name
}
