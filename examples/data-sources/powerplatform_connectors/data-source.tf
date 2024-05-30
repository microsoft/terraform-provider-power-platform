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

data "powerplatform_connectors" "all_connectors" {}
