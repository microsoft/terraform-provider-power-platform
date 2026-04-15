variable "environment_display_name" {
  default     = "example-git-branch-environment"
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

variable "solution_file" {
  default     = null
  description = "Optional override for the unmanaged solution package path. Leave null to use the bundled sample solution package."
  type        = string
  nullable    = true
}

variable "enable_git_binding" {
  default     = false
  description = "Set to true on the second apply to create the Git integration and solution branch binding after the environment and solution already exist."
  type        = bool
}

variable "git_provider" {
  default     = "AzureDevOps"
  description = "Git provider to bind. Supported value is AzureDevOps."
  type        = string
}

variable "scope" {
  default     = "Solution"
  description = "Source control integration scope. Use Solution for solution-level branch bindings."
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

variable "branch_name" {
  default     = "main"
  description = "Git branch name to bind the solution to."
  type        = string
}

variable "upstream_branch_name" {
  default     = null
  description = "Optional upstream branch name. Leave null to reuse branch_name."
  type        = string
  nullable    = true
}

variable "root_folder_path" {
  default     = "solutions/sample-solution"
  description = "Repository folder path used for the solution binding."
  type        = string
}
