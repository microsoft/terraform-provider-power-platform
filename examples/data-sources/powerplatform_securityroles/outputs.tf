output "all_Security_" {
  description = "Returns all Power Apps in the tenant"
  value       = data.powerplatform_security_roles.all.security_roles
}
