# Configure the Power Platform Provider using username/password
provider "powerplatform" {
  username  = var.username
  password  = var.password
  tenant_id = var.tenant_id
}
