terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  username  = var.username
  password  = var.password
  tenant_id = var.tenant_id
}

resource "powerplatform_billing_policy_environments" "pay_as_you_go_policy_envs" {
  billing_policy_id = "00000000-0000-0000-0000-000000000000"
  environments      = ["00000000-0000-0000-0000-000000000001", "00000000-0000-0000-0000-000000000002"]
}


