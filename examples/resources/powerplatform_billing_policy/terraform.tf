terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.80.0"
    }
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}
