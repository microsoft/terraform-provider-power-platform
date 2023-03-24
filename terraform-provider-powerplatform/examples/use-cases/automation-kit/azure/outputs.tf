output "key_vault_name" {
  value = "${azurerm_key_vault.kv_dev.id}"
}

output "key_vault_secret_client_id_name" {
  value = "${azurerm_key_vault_secret.client_id.name}"
}

output "key_vault_secret_client_password_name" {
  value = "${azurerm_key_vault_secret.app_user_secret.name}"
}

output "key_vault_client_secret_tenant_id_name" {
  value = "${azurerm_key_vault_secret.tenant_id.name}"
}

output "automation_kit_application_id" {
  value = "${azuread_application.automation_kit_app_user.application_id}"
}