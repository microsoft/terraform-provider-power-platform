provider "powerplatform" {
  username = "${var.username}"
  password = "${var.password}"
  host = "${var.host}"
}


variable "username" {
    default = "user@domain.onmicrosoft.com"
    type = string
}

variable "password" {
    type = string
}

variable "host" {
    default = "http://localhost:8080"
    type = string
}
