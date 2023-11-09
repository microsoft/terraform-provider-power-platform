output "all_environments" {
  description = "All environments"
  value       = data.powerplatform_environments.all_environments
}

output "all_applications" {
  description = "All applications"
  value       = data.powerplatform_applications.all_applications
}

/*
output "name" {
  description = "Display name of the environment"
  value       = powerplatform_application.development.display_name
}

output "id" {
  description = "Unique identifier of the environment"
  value       = powerplatform_application.development.environment_name
}

output "url" {
  description = "URL of the environment"
  value       = powerplatform_application.development.url
}
*/
