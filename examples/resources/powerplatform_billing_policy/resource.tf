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

resource "powerplatform_billing_policy" "pay_as_you_go" {
  name     = "pay_as_you_go"
  location = "europe"
  status   = "Enabled" //
  billing_instrument = {
    resource_group  = "sddssd"
    subscription_id = "sdsdsd"
  }

}


