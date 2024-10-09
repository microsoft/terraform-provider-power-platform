variable "applicationid" {
  default     = "00000000-0000-0000-0000-000000000001"
  description = "EntraId clientid of the application"
  type        = string
}

variable "businessunitid" {
  default     = "00000000-0000-0000-0000-000000000001"
  description = "Unique identifier of the business unit"
  type        = string
}

variable "roles" {
  default     = ["System Administrator"]
  type        = set(string)
  description = "The roles that are granted to this application user"
}
