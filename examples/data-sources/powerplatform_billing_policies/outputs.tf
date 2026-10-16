output "all_billing_policies" {
  description = "All Billing Policies"
  value       = data.powerplatform_billing_policies.all_policies
}

output "copilot_studio_billing_policies" {
  description = "Names of the billing policies with the Copilot Studio pay-as-you-go meter turned on"
  value       = [for policy in data.powerplatform_billing_policies.all_policies.billing_policies : policy.name if policy.meters.copilot_studio]
}
