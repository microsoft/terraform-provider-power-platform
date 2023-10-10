terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">=1.2.26"
    }
  }
}

resource "azurerm_gallery_application" "igl-app" {
  name              = "PowerShell7"
  gallery_id        = var.sig_id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "igl-app-version" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.igl-app.id
  location               = var.region

  manage_action {
    install = "move .\\PowerShell7 .\\installps7.ps1 & powershell -ExecutionPolicy Unrestricted -File installps7.ps1"
    remove  = "echo"
  }

  source {
    media_link = var.installps7_link
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}
