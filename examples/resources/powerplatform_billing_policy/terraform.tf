terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.62.1"
    }
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}
