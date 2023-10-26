output "all_policies" {
  description = "All Power Platform Data Loss Prevention Policies"
  value       = data.powerplatform_data_loss_prevention_policies.all_policies
}
