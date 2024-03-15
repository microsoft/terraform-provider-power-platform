terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id
}

resource "powerplatform_billing_policy" "pay_as_you_go" {
  name     = "payAsYouGoBillingPolicyExample"
  location = "europe"
  status   = "Enabled"
  billing_instrument = {
    resource_group  = "resource_group_name"
    subscription_id = "00000000-0000-0000-0000-000000000000"
  }
}


