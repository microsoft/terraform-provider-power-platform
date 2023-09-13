variable "username" {
  description = "The username of the Power Platform API in user@domain format"
  type        = string
  required    = false
  default     = null
  validation {
    condition     = (var.username != null && var.password != null) || (var.username == null && var.password == null)
    error_message = "Both username and password must be provided, or both should be null."
  }
}
variable "password" {
  description = "The password of the Power Platform API user"
  sensitive   = true
  type        = string
  default     = null
  required    = false
  validation {
    condition     = (var.username != null && var.password != null) || (var.username == null && var.password == null)
    error_message = "Both username and password must be provided, or both should be null."
  }
}
variable "tenant_id" {
  description = "The AAD tenant id of service principal or user"
  type        = string
  required    = true
}
