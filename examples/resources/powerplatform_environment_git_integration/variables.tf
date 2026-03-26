variable "environment_display_name" {
  default     = "example-git-integration-environment"
  description = "Display name of the example environment."
  type        = string
}

variable "location" {
  default     = "europe"
  description = "Power Platform geography for the example environment."
  type        = string
}

variable "azure_region" {
  default     = "northeurope"
  description = "Azure region for the Dataverse-backed example environment."
  type        = string
}

variable "security_group_id" {
  default     = "00000000-0000-0000-0000-000000000000"
  description = "Security group ID for Dataverse provisioning. Use the zero GUID to disable."
  type        = string
}

variable "git_provider" {
  default     = "AzureDevOps"
  description = "Git provider to bind. Supported value is AzureDevOps."
  type        = string
}

variable "scope" {
  default     = "Environment"
  description = "Source control integration scope. Use Environment for environment-scoped bindings or Solution when pairing with powerplatform_solution_git_branch."
  type        = string
}

variable "organization_name" {
  default     = "example-org"
  description = "Git organization or owner name."
  type        = string
}

variable "project_name" {
  default     = "example-project"
  description = "Git project name used for Azure DevOps bindings."
  type        = string
}

variable "repository_name" {
  default     = "example-repo"
  description = "Git repository name to bind to the environment."
  type        = string
}
