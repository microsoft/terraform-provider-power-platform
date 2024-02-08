variable "client_id" {
  description = "The username of the Power Platform API in user@domain format"
  type        = string
}
variable "client_secret" {
  description = "The client secret of the Power Platform API user"
  sensitive   = true
  type        = string
}
variable "tenant_id" {
  description = "The tenant id of the AAD tenant"
  type        = string
}
