terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "4.27.0"
    }
    azurecaf = {
      source = "aztfmod/azurecaf"
    }
  }
}
