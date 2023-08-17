output "all_environments_apps" {
  description = "Returns all Solutions in an environment"
  value       = data.powerplatform_solutions.all
}