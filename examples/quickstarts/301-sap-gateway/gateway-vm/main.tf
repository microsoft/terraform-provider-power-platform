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

  gallery_application {
    version_id = azurerm_gallery_application_version.example.id
    order      = 1
  }
  /*
  gallery_application {
    version_id = var.sap_nco_version_id
    order      = 2
  }
  */
}

resource "azurerm_shared_image_gallery" "example" {
  name                = "examplegallery"
  resource_group_name = var.resource_group_name
  location            = var.region
}

resource "azurerm_gallery_application" "example" {
  name              = "example-app"
  gallery_id        = azurerm_shared_image_gallery.example.id
  location          = var.region
  supported_os_type = "Windows"
}

resource "azurerm_gallery_application_version" "example" {
  name                   = "0.0.1"
  gallery_application_id = azurerm_gallery_application.example.id
  location               = var.region

  manage_action {
    install = "move .\\PowerShell7 .\\PowerShell-7.3.7-win-x64.msi & start /wait %windir%\\system32\\msiexec.exe /i PowerShell-7.3.7-win-x64.msi /qn /L*V 'C:Install_Test'"
    remove  = "echo"
  }

  source {
    media_link = "https://opdgwsetup.blob.core.windows.net/binaries/PowerShell-7.3.7-win-x64.msi"
  }

  target_region {
    name                   = var.region
    regional_replica_count = 1
  }
}


