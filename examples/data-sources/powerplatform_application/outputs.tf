output "all_environments" {
  description = "All environments"
  value       = data.powerplatform_environments.all_environments
}

output "all_applications" {
  description = "All applications"
  value       = data.powerplatform_applications.all_applications
}
