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

resource "azurerm_gallery_application" "java-igl-app" {
  name              = "JavaRE"
  gallery_id        = var.sig_id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "java-igl-app-version" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.java-igl-app.id
  location               = var.region

  manage_action {
    install = "move .\\JavaRE .\\java-setup.ps1 & powershell -ExecutionPolicy Unrestricted -File java-setup.ps1"
    remove  = "echo"
  }

  source {
    media_link = var.java_setup_link
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}
