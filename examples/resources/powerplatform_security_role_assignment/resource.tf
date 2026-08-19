terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_environment" "env" {
  display_name     = "Example Environment"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_application_user" "application_user" {
  environment_id = powerplatform_environment.env.id
  application_id = var.application_id
}

resource "powerplatform_security_role_assignment" "example" {
  environment_id     = powerplatform_environment.env.id
  principal_id       = powerplatform_application_user.application_user.id
  security_role_name = "Basic User"
}

# A security role can also be assigned to a team. Teams live in their own Dataverse
# table, so use team_id instead of system_user_id. Exactly one of the two is required.
resource "powerplatform_security_role_assignment" "team_admin" {
  environment_id     = powerplatform_environment.env.id
  team_id            = var.team_id
  security_role_name = "System Administrator"
}

variable "team_id" {
  description = "Dataverse teamid of the team to assign the role to"
  type        = string
}
