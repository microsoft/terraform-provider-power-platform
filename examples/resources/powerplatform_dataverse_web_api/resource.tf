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

resource "powerplatform_dataverse_web_api" "query" {
  environment_id = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
  create = {
    url    = "/api/data/v9.2/accounts?$select=name,accountid"
    method = "POST"
    body = jsonencode({
      "accountid" : "00000000-0000-0000-0000-000000000033",
      "name" : "powerplatform_dataverse_web_api",
      "creditonhold" : false,
      "address1_latitude" : 47.639583,
      "description" : "This is the description of the sample account",
      "revenue" : 5000000,
      "accountcategorycode" : 1
    })
    headers = [
      {
        name  = "Content-Type"
        value = "application/json; charset=utf-8"
      },
      {
        name  = "OData-MaxVersion"
        value = "4.0"
      },
      {
        name  = "OData-Version"
        value = "4.0"

      },
      {
        name  = "Prefer"
        value = "return=representation"
      }
    ]
  }
}

