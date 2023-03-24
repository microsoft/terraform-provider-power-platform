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
  host = "http://localhost:8080"
}

resource "powerplatform_environment" "development" {
  display_name = "DevelopmentEnvironment"
  location = "europe"
  language_name = "1033"
  currency_name = "USD"
  environment_type = "Sandbox"

  is_custom_controls_in_canvas_apps_enabled = true
}

output "name" {
  value = powerplatform_environment.development.display_name
}

output "id" {
  value = powerplatform_environment.development.environment_name
}

output "url" {
  value = powerplatform_environment.development.url
}

