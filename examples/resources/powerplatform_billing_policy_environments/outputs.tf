output "policy_name" {
  description = "Unique Id of the Power Platform Billing Policy"
  value       = powerplatform_billing_policy_environments.pay_as_you_go_policy_envs.billing_policy_id
}
