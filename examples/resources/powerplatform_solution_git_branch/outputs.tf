output "environment_id" {
  description = "Unique identifier of the example environment."
  value       = powerplatform_environment.example.id
}

output "solution_id" {
  description = "Provider-formatted ID of the example unmanaged solution."
  value       = powerplatform_solution.example.id
}

output "git_integration_id" {
  description = "Unique identifier of the environment Git integration binding, if enabled."
  value       = try(powerplatform_environment_git_integration.example[0].id, null)
}

output "solution_git_branch_id" {
  description = "Unique identifier of the solution Git branch binding, if enabled."
  value       = try(powerplatform_solution_git_branch.example[0].id, null)
}
