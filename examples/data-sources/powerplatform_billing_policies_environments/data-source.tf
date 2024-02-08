terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}

data "powerplatform_billing_policies_environments" "all_pay_as_you_go_policy_envs" {
  billing_policy_id = "00000000-0000-0000-0000-000000000000"
}
