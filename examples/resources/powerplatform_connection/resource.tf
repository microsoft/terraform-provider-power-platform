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


resource "powerplatform_connection" "flow_management" {
  environment_id = "0f555a0d-488a-ecd5-995c-47a85a167255"
  name           = "flow_management"
  display_name   = "flow management conn 1"
}
