output "all_policy_environments" {
  description = "All Billing Policies for an Environment"
  value       = data.powerplatform_billing_policies_environments.all_pay_as_you_go_policy_envs
}
