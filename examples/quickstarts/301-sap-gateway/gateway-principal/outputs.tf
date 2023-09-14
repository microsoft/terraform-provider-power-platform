output "auth_client_id" {
  value       = azuread_service_principal.gateway_principal.id
  description = "output password"
}

output "auth_client_secret" {
  value       = azuread_service_principal_password.gateway_principal_password.value
  description = "output password"
  sensitive   = true
}
