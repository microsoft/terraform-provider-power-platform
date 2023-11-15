output "all_settings" {
  description = "Tenant settings for the current tenant"
  value       = data.powerplatform_tenant_settings.settings
}
