output "all_connectors" {
  description = "All Billing Policies"
  value       = data.powerplatform_billing_policies.all_policies
}
