terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

data "powerplatform_billing_policies_environments" "all_pay_as_you_go_policy_envs" {
  billing_policy_id = "00000000-0000-0000-0000-000000000000"
}
