terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_environment" "data_record_app_user_example_env" {
  display_name     = "powerplatform_data_record_app_user_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

data "powerplatform_security_roles" "all_roles" {
  business_unit_id = var.businessunitid
  environment_id   = powerplatform_environment.data_record_app_user_example_env
}

locals {
  selected_roles = [for role in coalesce(data.powerplatform_security_roles.all_roles.security_roles, []) : role.role_id if contains(var.roles, role.name)]
}

resource "powerplatform_data_record" "app_user" {
  table_logical_name = "systemuser"
  environment_id     = powerplatform_environment.data_record_app_user_example_env.id
  disable_on_destroy = true
  columns = {
    applicationid = var.applicationid
    businessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.businessunitid
    }
    systemuserroles_association = tolist([for rid in local.selected_roles : { table_logical_name = "role", data_record_id = tostring(rid) }])
    isdisabled                  = false
  }
}
