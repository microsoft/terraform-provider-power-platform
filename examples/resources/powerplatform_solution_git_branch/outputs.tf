output "git_integration_id" {
  description = "Unique identifier of the environment Git integration binding."
  value       = powerplatform_environment_git_integration.example.id
}

output "solution_git_branch_id" {
  description = "Unique identifier of the solution Git branch binding."
  value       = powerplatform_solution_git_branch.example.id
}
