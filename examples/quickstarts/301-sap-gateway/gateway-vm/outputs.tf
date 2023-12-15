output "vm_opgw_principal_id" {
  value = azurerm_windows_virtual_machine.vm-opgw.identity.0.principal_id
}
