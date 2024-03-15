provider "powerplatform" {
  # Use a service principal to authenticate with the Power Platform service
  client_id     = var.client_id
  client_secret = var.client_secret
  tenant_id     = var.tenant_id

  # Use the Azure CLI to authenticate with the Power Platform service
  # use_cli = true

  # Use OIDC to authenticate with the Power Platform service
  # use_oidc = true
}
