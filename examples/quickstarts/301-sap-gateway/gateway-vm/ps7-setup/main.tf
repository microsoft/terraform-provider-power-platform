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
