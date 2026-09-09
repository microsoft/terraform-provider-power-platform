output "environment_id" {
  description = "Unique identifier of the example environment."
  value       = powerplatform_environment.example.id
}

output "git_integration_id" {
  description = "Unique identifier of the environment Git integration binding."
  value       = powerplatform_environment_git_integration.example.id
}
