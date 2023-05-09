output "all_environments_apps" {
  description = "Returns all Power Apps in the tenant"
  value       = data.powerplatform_powerapps.all.apps
}

output "specific_environment_apps" {
  description = "Returns all Power Apps in a specific environment defined in environment_name property of the powerplatform_powerapps data source"
  value       = data.powerplatform_powerapps.all.apps
}

output "specific_app_in_specific_environment" {
  description = "Returns specific Power App from an environment defined in environment_name property of the powerplatform_powerapps data source"
  value       = one(local.only_specific_name).name
}
