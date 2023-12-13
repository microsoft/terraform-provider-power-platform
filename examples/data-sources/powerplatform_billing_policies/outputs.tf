output "all_billing_policies" {
  description = "All Billing Policies"
  value       = data.powerplatform_billing_policies.all_policies
}
