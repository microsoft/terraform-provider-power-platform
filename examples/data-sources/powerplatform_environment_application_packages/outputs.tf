output "all_applications" {
  description = "All applications"
  value       = data.powerplatform_environment_application_packages.all_applications
}

output "all_applications_from_publisher" {
  description = "All applications from publisher"
  value       = data.powerplatform_environment_application_packages.all_applications_from_publisher
}

output "specific_application" {
  description = "Specific application"
  value       = data.powerplatform_environment_application_packages.specific_application
}
