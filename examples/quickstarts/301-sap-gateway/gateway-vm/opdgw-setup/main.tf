terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

resource "azurerm_gallery_application" "opdgw-igl-app" {
  name              = "opdgw-setup.ps1"
  gallery_id        = var.sig_id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "opdgw-igl-app-version" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.opdgw-igl-app.id
  location               = var.region

  manage_action {
    #install = "C:\\powershell7\\7\\pwsh.exe -ExecutionPolicy Unrestricted -File .\\opdgw-setup.ps1 -keyVaultUri ${var.keyVaultUri} -secretPPName ${var.secretPPName} -userAdmin ${var.userIdAdmin_pp}"
    install = "mkdir C:\\sapint & copy opdgw-setup.ps1 C:\\sapint"
    remove  = "echo" # Uninstall script is not applicable.
  }

  source {
    media_link = var.opdgw_setup_link
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}
