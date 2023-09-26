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
}

