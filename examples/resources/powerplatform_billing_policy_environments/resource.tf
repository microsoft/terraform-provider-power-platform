terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {

}

resource "powerplatform_billing_policy_environments" "pay_as_you_go_policy_environm" {
  billing_policy_id = "123"
  environments      = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"]
}


