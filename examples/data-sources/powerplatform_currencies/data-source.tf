terraform {
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_client = true
}

data "powerplatform_locations" "all_locations" {}

data "powerplatform_currencies" "all_currencies_by_location" {
  location = data.powerplatform_locations.all_locations.locations[0].name
}
