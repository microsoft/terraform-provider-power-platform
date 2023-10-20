terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

resource "azurerm_gallery_application" "ps7-igl-app" {
  name              = "ps7-setup.ps1"
  gallery_id        = var.sig_id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "ps7-igl-app-version" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.ps7-igl-app.id
  location               = var.region

  manage_action {
    #install = "msiexec.exe /package PowerShell-7.3.8-win-x64.msi /quiet ADD_EXPLORER_CONTEXT_MENU_OPENPOWERSHELL=1 ADD_FILE_CONTEXT_MENU_RUNPOWERSHELL=1 ENABLE_PSREMOTING=1 REGISTER_MANIFEST=1 USE_MU=1 ENABLE_MU=1 ADD_PATH=1"
    install = "powershell -ExecutionPolicy Unrestricted -File ps7-setup.ps1"
    remove  = "echo" # Uninstall script is not applicable.
  }

  source {
    media_link = var.ps7_setup_link
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}
