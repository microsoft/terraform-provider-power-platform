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
  username = "${var.username}"
  password = "${var.password}"
  host = "http://localhost:8080"
}


provider "azuread" {
  client_id     = "${var.aad_client_id}"
  client_secret = "${var.aad_client_secret}"
  tenant_id     = "11111111-2222-3333-4444-555555555555"
}


data "azuread_domains" "aad_domains" {
  only_initial = true
}

variable "users" {
  type = map(object({
    name = string
    roles = set(string)
  }))
  default = {
    "lara" = {
      name = "user1"
      roles = ["System Administrator"]
    }
    "kanaka" = {
      name = "user2"
      roles = ["System Administrator", "Environment Maker"]
    }
    "matt" = {
      name = "user3"
      roles = ["System Customizer"]
    }
    "hannes" = {
      name = "user4"
      roles = ["Environment Maker"]
    }
  }
}

resource "azuread_user" "user" {
  for_each = var.users

  user_principal_name = "${each.value.name}@${data.azuread_domains.aad_domains.domains[0].domain_name}"
  display_name        = each.value.name
  mail_nickname       = each.value.name
  given_name          = each.value.name
  surname             = each.value.name
  password            = "${var.new_user_password}"
  usage_location      = "CH"
}

resource "azuread_group_member" "power_platform_licenses_group" {
  for_each = var.users

  group_object_id  = "11111111-2222-3333-4444-555555555555"
  member_object_id = azuread_user.user[each.key].id

  depends_on = [
    azuread_user.user
  ]
}

resource "powerplatform_user" "ppuser" {
  for_each = var.users

  is_app_user = false
  aad_id = azuread_user.user[each.key].id
  user_principal_name = azuread_user.user[each.key].user_principal_name
  first_name = azuread_user.user[each.key].given_name
  last_name = azuread_user.user[each.key].surname
  environment_name = powerplatform_environment.environment.environment_name
  security_roles =each.value.roles

  depends_on = [
    azuread_group_member.power_platform_licenses_group,
    powerplatform_solution.solution
  ]
}

resource "powerplatform_environment" "environment" {
  display_name = "my-environment"
  location = "europe"
  language_name = "1033"
  currency_name = "USD"
  environment_type = "Sandbox"
}

resource "powerplatform_solution" "solution" {
  solution_file = "${path.module}/solution_1_0_0_0.zip"
  settings_file = ""
  solution_name = "solution"
  environment_name = powerplatform_environment.environment.environment_name
}
