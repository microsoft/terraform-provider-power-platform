terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.15.0"
    }
  }
}


provider "powerplatform" {
  username = "${var.username}"
  password = "${var.password}"
  host = "http://localhost:8080"
}


provider "azuread" {
  client_id     = "${var.aad_client_id}"
  client_secret = "${var.aad_client_secret}"
  tenant_id     = "11111111-2222-3333-4444-555555555555"
}


data "azuread_domains" "aad_domains" {
  only_initial = true
}

#Add user to AAD
resource "azuread_user" "user" {
  user_principal_name = "test1@${data.azuread_domains.aad_domains.domains[0].domain_name}"
  display_name        = "Test_1"
  mail_nickname       = "test1"
  given_name          = "Tester"
  surname             = "Testowski"
  password            = "${var.new_user_password}"
  usage_location      = "CH"
}

#AAD Group with licences assigned
resource "azuread_group_member" "power_platform_licenses_group" {
  group_object_id  = "11111111-2222-3333-4444-555555555555"
  member_object_id = azuread_user.user.id
}

#Add service principal to AAD
data "azuread_application_published_app_ids" "well_known" {}

resource "azuread_application" "app_user" {
  display_name = "PowerApps-AppService"

  required_resource_access {
      resource_app_id = data.azuread_application_published_app_ids.well_known.result.DynamicsCrm
      resource_access {
          id = "78ce3f0f-a1ce-49c2-8cde-64b5c0896db4" #"user_impersonation"
          type = "Scope"
      }
  }

  web {
    redirect_uris = ["https://localhost/"]
  }

}

resource "azuread_service_principal" "app_user_principal" {
  application_id = azuread_application.app_user.application_id
}

resource "azuread_service_principal_password" "example" {
  service_principal_id = azuread_service_principal.app_user_principal.object_id
}

#Add new environment
resource "powerplatform_environment" "environment" {
  display_name = "my-environment"
  location = "europe"
  language_name = "1033"
  currency_name = "USD"
  environment_type = "Sandbox"
}

#Add environment to DLP policy
resource "powerplatform_data_loss_prevention_policy" "my_policy" {
    display_name = "My Policy"
    
    environment_type = "ExceptEnvironments"
    environment {
          name = powerplatform_environment.environment.environment_name
    }
    default_connectors_classification = "Blocked"

    connector_group {
      classification = "Confidential"
      connector {
        id = "/providers/Microsoft.PowerApps/apis/shared_sql"
        name = "SQL Server"
      }

      connector {
        id = "/providers/Microsoft.PowerApps/apis/shared_sql"
        name = "SQL Server"
      }
    }
}

resource "powerplatform_user" "user" {
  is_app_user = false
  aad_id = azuread_user.user.id
  user_principal_name = azuread_user.user.user_principal_name
  firstname = azuread_user.user.given_name
  lastname = azuread_user.user.surname
  environments = [ powerplatform_environment.environment.environment_name ]
  security_roles = [ "System Administrator", "Basic User" ]
}


#Outputs
output "app_user_id" {
  value = azuread_service_principal.app_user_principal.application_id
}

output "user_name" {
  value = azuread_user.user.user_principal_name
}

output "user_pass" {
  value = azuread_user.user.password
  sensitive = true
}

output "environment_name" {
  value = powerplatform_environment.environment.display_name
}

output "environment_id" {
  value = powerplatform_environment.environment.environment_name
}