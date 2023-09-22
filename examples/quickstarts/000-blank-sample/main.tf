terraform {
  required_version = ">= 1.5"
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "powerplatform" {
  client_id = var.client_id
  secret    = var.secret
  tenant_id = var.tenant_id
}


resource "powerplatform_environment" "development" {
  display_name      = "example_environment"
  location          = "europe"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  domain            = "mydomain"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

data "powerplatform_connectors" "all_connectors" {}

# data "azurerm_resource_group" "example" {
#   name     = "example"
#   location = "West Europe"
# }
