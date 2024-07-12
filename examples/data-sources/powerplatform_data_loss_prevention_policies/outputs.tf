output "all_applications" {
  description = "All policies"
  value       = data.powerplatform_data_loss_prevention_policies.tenant_data_loss_prevention_policies.policies
}
