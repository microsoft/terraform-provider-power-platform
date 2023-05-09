terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  host     = var.host
}

resource "powerplatform_environment" "development" {
  display_name                              = "Example Environment"
  location                                  = "europe"
  language_name                             = "1033"
  currency_name                             = "USD"
  environment_type                          = "Sandbox"
  is_custom_controls_in_canvas_apps_enabled = true
}

