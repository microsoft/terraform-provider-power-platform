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

data "powerplatform_environments" "all_environments" {}

resource "powerplatform_data_record" "data_record_accounts" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  table_name     = "accounts"
  columns = {
    name = "Sample Account"
    creditonhold = false
    address1_latitude = 47.639583
    description = "This is the description of the sample account"
    revenue = 5000000
    accountcategorycode = 1
  }
}

resource "powerplatform_data_record" "data_record_testones" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  table_name     = "cr4d0_testones"
  columns = {
    cr4d0_multilinefield          = "asdasd\nasd\na\nsdasd"
    cr4d0_wholenumberfield        = 1231
    cr4d0_decimalfield            = 123.12
    cr4d0_dateandtimefield        = "2024-04-10T22:32:12Z"
    cr4d0_dateonlyfield           = "2024-04-10T22:32:12Z"
    cr4d0_currencyfield           = 123
    cr4d0_yesnochoice             = true
    cr4d0_floatfield              = 123.35
    cr4d0_stringfield             = "Testing, 1, 2, 3, 4"
    cr4d0_multiselectchoicefield  = "1, 2"
    cr4d0_name                    = "Test3"
    cr4d0_singleoptionchoicefield = 2
    cr4d0_LookupField = [
      {
        entity_logical_name = "cr4d0_tabletwos"
        data_record_id      = "21715311-9ff6-ee11-a1fd-7c1e5217db96"
      },
      {
        entity_logical_name = "cr4d0_tabletwos"
        data_record_id      = "d3a83be1-1bf9-ee11-a1fd-000d3a4de0ce"
      }
    ]
    cr4d0_TestUserManyToOne = {
      entity_logical_name = "systemusers"
      data_record_id      = "7f054957-2df3-ee11-a1fd-000d3a5389af"
    }
    name_of_the_relation = [
      {
        entity_logical_name = "systemusers"
        data_record_id      = "7f054957-2df3-ee11-a1fd-000d3a5389af"
      },
      {
        entity_logical_name = "accounts"
        data_record_id      = "7f054957-2df3-ee11-a1fd-000d3a5389af"
      },
      {
        entity_logical_name = "contacts"
        data_record_id      = "7f054957-2df3-ee11-a1fd-000d3a5389af"
      }
    ]
  }
}

resource "powerplatform_data_record" "data_record_tabletwos" {
  environment_id = "61ba1e49-21ed-eaba-8192-aaa376d6150d"
  table_name     = "cr4d0_tabletwos"
  record_id      = "21715311-9ff6-ee11-a1fd-7c1e5217db96"
  columns = {
    cr4d0_name = "Updated Set"
  }
}
