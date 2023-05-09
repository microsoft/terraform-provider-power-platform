output "all_dlp_policies" {
  description = "All DLP policies existing in the tenant"
  value       = data.powerplatform_data_loss_prevention_policies.all
}
