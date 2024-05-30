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

resource "powerplatform_billing_policy_environment" "pay_as_you_go_policy_envs" {
  billing_policy_id = "00000000-0000-0000-0000-000000000000"
  environments      = ["00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"]
}
