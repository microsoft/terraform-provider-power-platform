terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_environment" "user_example" {
  display_name      = "user_example"
  location          = "europe"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

resource "powerplatform_user" "new_user" {
  environment_id = powerplatform_environment.user_example.id
  security_roles = [] //["a1801436-efd6-e811-a96e-000d3a3ab886"]
  aad_id         = "ad7b0121-6fca-440b-99ae-0d54d89a3ac7"
}


