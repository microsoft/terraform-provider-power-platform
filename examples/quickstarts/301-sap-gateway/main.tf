module "gateway_principal" {
  source    = "./gateway-principal"
  client_id = var.client_id
  secret    = var.secret
  tenant_id = var.tenant_id
}
