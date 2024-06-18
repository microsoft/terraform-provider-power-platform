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


data "powerplatform_data_records" "root_business_unit" {
  environment_id    = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
  entity_collection = "businessunits"
  filter            = "parentbusinessunitid eq null"
  select            = ["name"]
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

resource "powerplatform_data_record" "role" {
  environment_id     = powerplatform_environment.data_record_example_env.id
  table_logical_name = "role"

  columns = {
    name = "my custom role"

    businessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = data.powerplatform_data_records.root_business_unit.rows[0].businessunitid
    }
  }
}

resource "powerplatform_data_record" "team" {
  environment_id     = powerplatform_environment.data_record_example_env.id
  table_logical_name = "team"
  columns = {
    name        = "main team"
    description = "main team description"

    teamroles_association = [
      {
        table_logical_name = "role"
        data_record_id     = powerplatform_data_record.role.id
      }
    ]
  }
}
