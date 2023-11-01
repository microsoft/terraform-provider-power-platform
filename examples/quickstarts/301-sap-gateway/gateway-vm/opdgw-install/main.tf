terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

resource "azurerm_gallery_application" "opdgw-install-igl-app" {
  name              = "opdgw-install.ps1"
  gallery_id        = var.sig_id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "opdgw-install-igl-app-version" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.opdgw-install-igl-app.id
  location               = var.region

  manage_action {
    install = "C:\\powershell7\\7\\pwsh.exe -ExecutionPolicy Unrestricted -File .\\opdgw-install.ps1"
        install = "C:\\powershell7\\7\\pwsh.exe -ExecutionPolicy Unrestricted -File .\\opdgw-install.ps1"
    remove  = "echo" # Uninstall script is not applicable.
  }

  source {
    media_link = var.opdgw_install_link
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}
