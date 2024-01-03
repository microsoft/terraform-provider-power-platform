output "storage_blob_ps7_setup_link" {
  value = azurerm_storage_blob.storage_blob_ps7_setup.url
}

output "storage_blob_java_runtime_link" {
  value = azurerm_storage_blob.storage_blob_java_runtime.url
}

output "storage_blob_sapnco_install_link" {
  value = azurerm_storage_blob.storage_blob_sapnco_install.url
}

output "storage_blob_runtime_setup_link" {
  value = azurerm_storage_blob.storage_blob_runtime_setup.url
}
