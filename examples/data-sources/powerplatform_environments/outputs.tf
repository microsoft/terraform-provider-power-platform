output "all_environments" {
  description = "All environments in the tenant"
  value       = data.powerplatform_environments.all_environments
}