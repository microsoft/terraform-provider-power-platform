terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.15.0"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  host     = "http://localhost:8080"
}

provider "azuread" {
  tenant_id     = var.aad_tenant_id
  client_id     = var.aad_client_id
  client_secret = var.aad_client_secret
}

data "azuread_domains" "aad_domains" {
  only_initial = true
}

data "azuread_application_published_app_ids" "well_known" {}

#Add user to AAD
resource "azuread_user" "user" {
  user_principal_name = "test_1@${data.azuread_domains.aad_domains.domains[0].domain_name}"
  display_name        = "Test_1"
  mail_nickname       = "test1"
  given_name          = "Tester"
  surname             = "Testowski"
  password            = var.new_user_password
  usage_location      = "CH"
}

#AAD Group with licences assigned
resource "azuread_group_member" "power_platform_licenses_group" {
  group_object_id  = var.aad_licensing_security_group
  member_object_id = azuread_user.user.id
}

resource "azuread_application" "app_user" {
  display_name = "PowerApps-AppService-Test-Terraform"

  required_resource_access {
    resource_app_id = data.azuread_application_published_app_ids.well_known.result.DynamicsCrm
    resource_access {
      id   = "78ce3f0f-a1ce-49c2-8cde-64b5c0896db4" #"user_impersonation"
      type = "Scope"
    }
  }

  web {
    redirect_uris = ["https://localhost/"]
  }

}

resource "azuread_service_principal" "app_user_principal" {
  application_id = azuread_application.app_user.application_id
}

resource "powerplatform_environment" "environment" {
  display_name     = "User Test Environment"
  location         = "europe"
  language_name    = "1033"
  currency_name    = "USD"
  environment_type = "Sandbox"
}

resource "powerplatform_user" "user" {
  environment_name    = powerplatform_environment.environment.environment_name
  is_app_user         = false
  aad_id              = azuread_user.user.id
  user_principal_name = azuread_user.user.user_principal_name
  first_name          = azuread_user.user.given_name
  last_name           = azuread_user.user.surname
  security_roles      = ["Basic User", "Environment Maker"]
}

resource "powerplatform_user" "app_user" {
  environment_name = powerplatform_environment.environment.environment_name
  is_app_user      = true
  application_id   = azuread_application.app_user.application_id
  first_name       = "Terraform Test"
  last_name        = "App User"
  security_roles   = ["System Administrator", "Basic User"]
}