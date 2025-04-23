terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

# Provider configuration with Continuous Access Evaluation (CAE) enabled
# CAE allows authentication tokens to be invalidated in near real-time when security policies change
provider "powerplatform" {
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret

  # Enable Continuous Access Evaluation for real-time security policy enforcement
  enable_continuous_access_evaluation = true
}

# Configuration variables
variable "tenant_id" {
  type        = string
  description = "The Azure AD tenant ID"
}

variable "client_id" {
  type        = string
  description = "The client/application ID"
}

variable "client_secret" {
  type        = string
  description = "The client secret"
  sensitive   = true
}
