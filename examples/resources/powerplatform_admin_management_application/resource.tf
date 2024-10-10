terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~>3.0"
    }
    azurecaf = {
      source = "aztfmod/azurecaf"
      version = "~>1.2"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

provider "azuread" {
  use_cli = true
}

resource "azuread_application_registration" "example_app" {
  display_name = "Power Platform Example Admin Management Application"
}

resource "azuread_service_principal" "example_sp" {
  client_id = azuread_application_registration.example_app.client_id
}

resource "powerplatform_admin_management_application" "example_registration" {
  id = azuread_application_registration.example_app.client_id
}
