terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

# Provider configuration including partner ID for Customer Usage Attribution
provider "powerplatform" {
  use_cli    = true
  partner_id = "00000000-0000-0000-0000-000000000000"
}
