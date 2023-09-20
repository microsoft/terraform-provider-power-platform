variable "client_id" {
  description = "The username of the Power Platform API in user@domain format"
  type        = string


}
variable "secret" {
  description = "The password of the Power Platform API user"
  sensitive   = true
  type        = string
  default     = null

}
variable "tenant_id" {
  description = "The AAD tenant id of service principal or user"
  type        = string

}
