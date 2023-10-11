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

resource "azurecaf_name" "sig" {
  name          = var.base_name
  resource_type = "azurerm_shared_image_gallery"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_shared_image_gallery" "sig" {
  name                = azurecaf_name.sig.result
  resource_group_name = var.resource_group_name
  location            = var.region
}

# Create PowerShell 7 version in Shared Image Gallery
module "powershell-setup" {
  source              = "./powershell-setup"
  prefix              = var.prefix
  base_name           = var.base_name
  resource_group_name = var.resource_group_name
  region              = var.region
  sig_id              = azurerm_shared_image_gallery.sig.id
  installps7_link     = var.installps7_link
}

# Create Java Runtime version in Shared Image Gallery
module "java-runtime-setup" {
  source              = "./java-runtime-setup"
  prefix              = var.prefix
  base_name           = var.base_name
  resource_group_name = var.resource_group_name
  region              = var.region
  sig_id              = azurerm_shared_image_gallery.sig.id
  java_setup_link     = var.java_setup_link
}

resource "azurecaf_name" "vm-opgw" {
  name          = var.base_name
  resource_type = "azurerm_windows_virtual_machine"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_windows_virtual_machine" "vm-opgw" {
  name                  = azurecaf_name.vm-opgw.result
  location              = var.region
  resource_group_name   = var.resource_group_name
  network_interface_ids = [var.nic_id]

  size                                                   = "Standard_D4s_v5"
  admin_username                                         = "sapadmin"
  admin_password                                         = var.vm_pwd
  computer_name                                          = "vmopgw"
  enable_automatic_updates                               = true
  bypass_platform_safety_checks_on_user_schedule_enabled = false
  patch_assessment_mode                                  = "ImageDefault"
  patch_mode                                             = "AutomaticByOS"

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
    disk_size_gb         = 128
    name                 = "myosdisk1"
  }

  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-smalldisk"
    version   = "latest"
  }

  # Setup PowerShell 7
  gallery_application {
    version_id = module.powershell-setup.powershell_version_id
    order      = 1
  }

  # Setup Java Runtime
  gallery_application {
    version_id = module.java-runtime-setup.java_runtime_version_id
    order      = 2
  }

}
