terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

resource "azurerm_gallery_application" "sapnco-install-igl-app" {
  name              = "sapnco.msi"
  gallery_id        = var.sig_id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "sapnco-install-igl-app-version" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.sapnco-install-igl-app.id
  location               = var.region

  manage_action {
    install = "msiexec /i sapnco.msi /qn /norestart FolderForm_AllUsers=ALL GAC_WMI=2 /log install.log"
    remove  = "echo" # Uninstall script is not applicable.
  }

  source {
    media_link = var.sapnco_install_link
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}
