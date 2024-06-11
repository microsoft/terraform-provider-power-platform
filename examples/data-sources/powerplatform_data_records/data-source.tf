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

resource "powerplatform_environment" "env" {
  display_name     = "powerplatform_data_record_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

data "powerplatform_data_records" "data_query" {
  environment_id    = powerplatform_environment.env.id
  entity_collection = "systemusers"
  filter            = "isdisabled eq false"
  select            = ["firstname", "lastname", "domainname"]
  top               = 2
  order_by          = "lastname asc"

  expand = [
    {
      navigation_property = "systemuserroles_association",
      select              = ["name"],
    },
    {
      navigation_property = "teammembership_association",
      select              = ["name"],
    }
  ]
}
