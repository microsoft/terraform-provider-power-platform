
output "azure_openai_connection" {
  description = "New Azure Open AI connection"
  value       = powerplatform_connection.azure_openai_connection
  sensitive   = true
}
