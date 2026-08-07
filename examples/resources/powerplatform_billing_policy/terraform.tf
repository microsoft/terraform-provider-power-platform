terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "5.0.1"
    }
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}
