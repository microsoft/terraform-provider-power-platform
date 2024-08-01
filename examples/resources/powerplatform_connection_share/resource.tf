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


resource "powerplatform_connection" "new_sql_connection" {
  environment_id = var.environment_id
  name           = "shared_sql"
  display_name   = "My SQL Connection"
  connection_parameters_set = jsonencode({
    "name" : "oauthSP",
    "values" : {
      "token" : {
        "value" : "https://global.consent.azure-apim.net/redirect/sql"
      },
      "token:TenantId" : {
        "value" : "${var.tenant_id}"
      },
      "token:clientId" : {
        "value" : "${var.client_id}"
      },
      "token:clientSecret" : {
        "value" : "${var.client_secret}"
      }
    }
  })

  lifecycle {
    ignore_changes = [
      connection_parameters_set
    ]
  }
}

resource "powerplatform_connection_share" "share_with_admin" {
  environment_id = var.environment_id
  connector_name = powerplatform_connection.new_sql_connection.name
  connection_id  = powerplatform_connection.new_sql_connection.id
  role_name      = "CanEdit"
  principal = {
    entra_object_id = var.user_object_id
  }
}

