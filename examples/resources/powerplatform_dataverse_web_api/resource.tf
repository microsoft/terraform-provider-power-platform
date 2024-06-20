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
    "accountid" : "00000000-0000-0000-0000-000000000001",
    "name" : "powerplatform_dataverse_web_api",
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

//we should skip output.status as it changes every time
//create should not be a nested item because we have to repeat content for update nested item
//how to get the id of the created item?
# resource "powerplatform_dataverse_web_api" "query" {
#   environment_id = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
#   create = {
#     url     = "/api/data/v9.2/accounts?$select=name,accountid"
#     method  = "POST"
#     body    = local.body
#     headers = local.headers
#   }
#   read = {
#     url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
#     method = "GET"
#   }
#   update = {
#     url     = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)?$select=name,accountid"
#     method  = "PATCH"
#     body    = local.body
#     headers = local.headers
#   }
#   delete = {
#     url    = "api/data/v9.2/accounts(00000000-0000-0000-0000-000000000001)"
#     method = "DELETE"
#   }
# }

# output "query_output_body" {
#   description = "Query Output Body"
#   value       = jsondecode(resource.powerplatform_dataverse_web_api.query.output.body).accountid
# }

# output "query_output_status" {
#   description = "Query Output Status"
#   value       = resource.powerplatform_dataverse_web_api.query.output.status
# }



# resource "powerplatform_dataverse_web_api" "whoami" {
#   environment_id = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
#   create = {
#     url     = "/api/data/v9.2/WhoAmI"
#     method  = "GET"
#     headers = local.headers
#   }
#   read = {
#     url    = "/api/data/v9.2/WhoAmI"
#     method = "GET"
#   }
# }


#is lack of read ok?
#what about delete?
resource "powerplatform_dataverse_web_api" "create_multiple" {
  environment_id = "a1e605fb-80ad-e1b2-bae0-f046efc0e641"
  create = {
    url    = "api/data/v9.2/accounts/Microsoft.Dynamics.CRM.CreateMultiple"
    method = "POST"
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
