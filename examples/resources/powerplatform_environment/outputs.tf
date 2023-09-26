output "name" {
  description = "Display name of the environment"
  value       = powerplatform_environment.development.display_name
}

output "id" {
  description = "Unique identifier of the environment"
  value       = powerplatform_environment.development.environment_id
}

output "url" {
  description = "URL of the environment"
  value       = powerplatform_environment.development.url
}
