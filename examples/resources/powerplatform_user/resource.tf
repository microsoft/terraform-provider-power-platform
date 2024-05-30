terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
    random = {
      source = "hashicorp/random"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

provider "azuread" {
  use_cli = true
}

data "azuread_domains" "aad_domains" {
  only_initial = true
}

locals {
  domain_name = data.azuread_domains.aad_domains.domains[0].domain_name
}

resource "random_password" "passwords" {
  length           = 16
  special          = true
  override_special = "_%@"
}

resource "azuread_user" "test_user" {
  user_principal_name = "user_example@${local.domain_name}"
  display_name        = "user_example"
  mail_nickname       = "user_example"
  password            = random_password.passwords.result
  usage_location      = "US"
}

resource "powerplatform_environment" "dataverse_user_example" {
  display_name     = "user_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_user" "new_user" {
  environment_id = powerplatform_environment.dataverse_user_example.id
  security_roles = [
    "e0d2794e-82f3-e811-a951-000d3a1bcf17", // bot author
  ]
  aad_id         = azuread_user.test_user.id
  disable_delete = false
}
