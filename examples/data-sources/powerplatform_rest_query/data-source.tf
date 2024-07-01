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

resource "powerplatform_environment" "development" {
  display_name     = "example_environment"
  location         = "europe"
  azure_region     = "northeurope"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

data "powerplatform_rest_query" "webapi_query" {
  scope                = "${powerplatform_environment.development.dataverse.url}/.default"
  url                  = "${powerplatform_environment.development.dataverse.url}/api/data/v9.2/RetrieveCurrentOrganization(AccessType=@p1)?@p1=Microsoft.Dynamics.CRM.EndpointAccessType'Default'"
  method               = "GET"
  expected_http_status = [200]
}
