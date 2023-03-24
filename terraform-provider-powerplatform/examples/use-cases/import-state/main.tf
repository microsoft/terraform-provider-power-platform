terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
  }
}


provider "powerplatform" {
  username = "${var.username}"
  password = "${var.password}"
}

#terraform import powerplatform_environment.test 11111111-2222-3333-4444-555555555555
#terraform show
#terraform plan

resource "powerplatform_environment" "test" {
    currency_name                     = "USD"
    display_name                      = "Test123"
    environment_type                  = "Sandbox"
    language_name                     = 1033
    location                          = "unitedstates"
}