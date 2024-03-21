terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

data "powerplatform_locations" "all_locations" {}

data "powerplatform_languages" "all_languages_by_location" {
  location_id = data.powerplatform_locations.all_locations.locations[0].name
}
