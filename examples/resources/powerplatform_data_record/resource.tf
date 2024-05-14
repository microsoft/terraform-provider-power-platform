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

resource "random_string" "random_suffix" {
  length  = 5
  special = false
  upper   = false
}

data "powerplatform_environments" "all_environments" {}

resource "powerplatform_data_record" "data_record_sample_contact1" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "contact"
  columns = {
    firstname          = "John"
    lastname           = "Doe ${random_string.random_suffix.result}"
    telephone1         = "555-555-5555"
    emailaddress1      = "johndoe@contoso.com"
    address1_composite = "123 Main St\nRedmond\nWA\n98052\nUS"
    anniversary        = "2024-04-10"
    annualincome       = 1234.56
    birthdate          = "2024-04-10"
    description        = "This is the description of the the terraform \n\nsample contact"
  }
}

resource "powerplatform_data_record" "data_record_sample_contact2" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "contact"
  columns = {
    firstname          = "Jane"
    lastname           = "Doe ${random_string.random_suffix.result}"
    telephone1         = "555-555-5555"
    emailaddress1      = "janedoe@contoso.com"
    address1_composite = "123 Main St\nRedmond\nWA\n98052\nUS"
    anniversary        = "2024-04-11"
    annualincome       = 1234.56
    birthdate          = "2024-04-11"
    description        = "This is the description of the the terraform \n\nsample contact"
  }
}

resource "powerplatform_data_record" "data_record_accounts" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "account"
  columns = {
    name                = "Sample Account ${random_string.random_suffix.result}"
    creditonhold        = false
    address1_latitude   = 47.639583
    description         = "This is the description of the sample account"
    revenue             = 5000000
    accountcategorycode = 1
  }
}

resource "powerplatform_data_record" "data_record_tabletwos" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "cr4d0_tabletwo"
  columns = {
    cr4d0_tabletwoid = "21715311-9ff6-ee11-a1fd-7c1e5217db96"
    cr4d0_name       = "Set ${random_string.random_suffix.result}"
  }
}
