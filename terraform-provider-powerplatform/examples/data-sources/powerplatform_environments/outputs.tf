output "all_environments" {
  description = "All environments existing in the tenant"
  value       = data.powerplatform_environments.all.environments
}