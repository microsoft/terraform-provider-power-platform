output "all_environments_apps" {
  description = "Returns all Power Apps in the tenant"
  value       = data.powerplatform_environment_powerapps.all.powerapps
}
