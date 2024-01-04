output "id" {
  description = "Unique identifier of the environment"
  value       = powerplatform_environment.xpp-dev1.id
}
output "name" {
  description = "Display name of the environment"
  value       = powerplatform_environment.xpp-dev1.display_name
}

output "url" {
  description = "URL of the environment"
  value       = powerplatform_environment.xpp-dev1.url
}

output "linked_app_type" {
  description = "Type of the linked D365 app"
  value       = powerplatform_environment.xpp-dev1.linked_app_type
}

output "linked_app_id" {
  description = "Unique identifier of the linked D365 Finance app"
  value       = powerplatform_environment.xpp-dev1.linked_app_id
}

output "linked_app_url" {
  description = "URL of the linked D365 Finance app"
  value       = powerplatform_environment.xpp-dev1.linked_app_url
}
