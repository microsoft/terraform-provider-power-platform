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


data "powerplatform_rest_query" "webapi_query" {
  environment_id       = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
  url                  = "api/data/v9.2/WhoAmI"
  method               = "GET"
  expected_http_status = [200]
}
