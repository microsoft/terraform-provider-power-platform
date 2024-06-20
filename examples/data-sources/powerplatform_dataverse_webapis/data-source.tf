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


data "powerplatform_dataverse_web_apis" "webapi_query" {
  environment_id = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
  url            = "api/data/v9.2/WhoAmI"
  method         = "GET"

  lifecycle {
    postcondition {
      condition     = self.output.status == 200
      error_message = "expected a 200 result code from API"
    }
  }
}
