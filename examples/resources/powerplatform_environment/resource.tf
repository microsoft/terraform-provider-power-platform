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

// when env type == dev then owner_id is required
// when ower_id is provider then env type should be developer
// when owner_id is provider then dataverse should be provided
// when owner_id is provided then security_group_id should not be 
// when environemt_group_id is provided then dataverse should be provided
// when security_group_id is provided then owner_id should not be provided
// validate that ownerid and security_group_id are valid guids

resource "powerplatform_environment" "development" {
  display_name         = "example_environment"
  description          = "example environment description"
  location             = "europe"
  azure_region         = "northeurope"
  environment_type     = "Sandbox"
  cadence              = "Moderate"
  environment_group_id = ""
  owner_id             = "00000000-0000-0000-0000-000000000000"
  dataverse = {
    language_code = "1033"
    currency_code = "USD"
    //security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}
