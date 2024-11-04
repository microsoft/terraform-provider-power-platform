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

resource "powerplatform_enterprise_policy" "network_injection" {
  environment_id = "00000000-0000-0000-0000-000000000000"
  system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000000"
  policy_type    = "NetworkInjection"
}

# resource "powerplatform_enterprise_policy" "encryption" {
#   environment_id = "00000000-0000-0000-0000-000000000000"
#   system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/00000000-0000-0000-0000-000000000000"
#   policy_type    = "Encryption"
# }
