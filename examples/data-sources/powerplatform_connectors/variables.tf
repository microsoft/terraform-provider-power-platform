variable "client_id" {
  description = "The client ID of the of the service principal"
  type        = string

}
variable "secret" {
  description = "The client secret of the service principal"
  sensitive   = true
  type        = string
}
variable "tenant_id" {
  description = "The tenant id of the AAD tenant"
  type        = string
}
variable "username" {
  description = "The username of the user"
  type        = string
  sensitive   = false
}
variable "password" {
  description = "The password of the user"
  type        = string
  sensitive   = true
}
