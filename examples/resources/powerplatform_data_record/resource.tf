terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

provider "azuread" {
  use_cli = true
}

resource "powerplatform_environment" "data_record_example_env" {
  display_name     = "powerplatform_data_record_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

# get the root business unit by querying for the business unit without a parent
data "powerplatform_data_records" "root_business_unit" {
  environment_id    = powerplatform_environment.data_record_example_env.id
  entity_collection = "businessunits"
  filter            = "parentbusinessunitid eq null"
  select            = ["name"]
}

# Create a new business unit with the root business unit as parent
module "business_unit" {
  source                  = "./res_business_unit"
  environment_id          = powerplatform_environment.data_record_example_env.id
  name                    = "Sales"
  costcenter              = "123"
  parent_business_unit_id = one(data.powerplatform_data_records.root_business_unit.rows).businessunitid
}

# Create a new role
module "custom_role" {
  source           = "./res_role"
  environment_id   = powerplatform_environment.data_record_example_env.id
  role_name        = "my custom role"
  business_unit_id = one(data.powerplatform_data_records.root_business_unit.rows).businessunitid
}

module "team" {
  source           = "./res_team"
  environment_id   = powerplatform_environment.data_record_example_env.id
  team_name        = "main team"
  team_description = "main team description"
  role_ids         = [module.custom_role.role_id]

}

resource "azuread_application_registration" "data_record_app_user" {
  display_name = "powerplatform_data_record_example"
}

resource "azuread_service_principal" "data_record_app_user" {
  client_id = azuread_application_registration.data_record_app_user.client_id
}

module "application_user" {
  source           = "./res_application_user"
  environment_id   = powerplatform_environment.data_record_example_env.id
  application_id   = azuread_application_registration.data_record_app_user.client_id
  business_unit_id = one(data.powerplatform_data_records.root_business_unit.rows).businessunitid
  role_ids         = [module.custom_role.role_id]
}
