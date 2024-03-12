output "all_Security_" {
  description = "Returns all Power Apps in the tenant"
  value       = data.powerplatform_securityroles.all.security_roles
}
