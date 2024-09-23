terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.3.0"
    }
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

provider "azurerm" {
  features {}
  use_cli = true
}

provider "azurecaf" {
}

data "azurerm_client_config" "current" {
}

resource "azurecaf_name" "rg_example_name" {
  name          = "power-platform-billing"
  resource_type = "azurerm_resource_group"
  random_length = 5
  clean_input   = true
}

resource "azurerm_resource_group" "rg_example" {
  name     = azurecaf_name.rg_example_name.result
  location = "westeurope"
}

resource "powerplatform_billing_policy" "pay_as_you_go" {
  name     = "payAsYouGoBillingPolicyExample"
  location = "europe"
  status   = "Enabled"
  billing_instrument = {
    resource_group  = azurerm_resource_group.rg_example.name
    subscription_id = data.azurerm_client_config.current.subscription_id
  }
}
