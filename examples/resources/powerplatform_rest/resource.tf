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

locals {
  body = jsonencode({
    "accountid" : "00000000-0000-0000-0000-000000000002",
    "name" : "Sample Account1",
    "creditonhold" : true,
    "address1_latitude" : 47.6396,
    "description" : "This is the updated description of the sample account",
    "revenue" : 6000000,
    "accountcategorycode" : 2
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

# resource "powerplatform_rest" "query" {
#   environment_id = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
#   create = {
#     url                  = "/api/data/v9.2/accounts?$select=name,accountid"
#     method               = "POST"
#     body                 = local.body
#     headers              = local.headers
#     expected_http_status = [201]
#   }
#   read = {
#     url                  = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000002)?$select=name,accountid"
#     method               = "GET"
#     expected_http_status = [200]
#   }
#   update = {
#     url                  = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000002)?$select=name,accountid"
#     method               = "PATCH"
#     body                 = local.body
#     headers              = local.headers
#     expected_http_status = [200]
#   }
#   destroy = {
#     url                  = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000002)"
#     method               = "DELETE"
#     expected_http_status = [204]
#   }
# }



resource "powerplatform_rest" "create_multiple" {
  environment_id = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
  create = {
    url                  = "api/data/v9.2/accounts/Microsoft.Dynamics.CRM.CreateMultiple"
    method               = "POST"
    expected_http_status = [200]
    body = jsonencode({
      "Targets" : [
        {
          "name" : "company 1"
          "@odata.type" : "Microsoft.Dynamics.CRM.account"
        },
        {
          "name" : "company 2"
          "@odata.type" : "Microsoft.Dynamics.CRM.account"
        },
        {
          "name" : "company 3"
          "@odata.type" : "Microsoft.Dynamics.CRM.account"
        }
      ]
    })
  }
}
