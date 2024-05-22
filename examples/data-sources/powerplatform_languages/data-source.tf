terraform {
  required_providers {
    power-platform = {
      source  = "microsoft/power-platform"
    }
  }
}

provider "power-platform" {
  use_cli = true
}

data "powerplatform_locations" "all_locations" {}

data "powerplatform_languages" "all_languages_by_location" {
  location = data.powerplatform_locations.all_locations.locations[0].name
}
