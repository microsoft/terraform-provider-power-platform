terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

resource "azurerm_gallery_application" "runtime-igl-app" {
  name              = "runtime-setup.ps1"
  gallery_id        = var.sig_id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "runtime-igl-app-version" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.runtime-igl-app.id
  location               = var.region

  manage_action {
    install = "mkdir C:\\sapint & copy runtime-setup.ps1 C:\\sapint & C:\\powershell7\\7\\pwsh.exe -ExecutionPolicy Unrestricted -command \"&{Install-Module -Name DataGateway -Force}\""
    remove  = "echo" # Uninstall script is not applicable.
  }

  source {
    media_link = var.runtime_setup_link
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}
