terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  host     = var.host
}

resource "powerplatform_environment" "environment" {
  display_name     = "package_import_test"
  location         = "europe"
  language_name    = "1033"
  currency_name    = "USD"
  environment_type = "Sandbox"
}

resource "powerplatform_package" "package" {
  environment_name = powerplatform_environment.environment.name
  package_name     = "package_import_test"
  package_file     = "${path.module}/package_import_test.zip"
  package_settings = ""
}