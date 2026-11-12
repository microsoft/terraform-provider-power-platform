---
page_title: "powerplatform_data_record Resource - powerplatform"
description: |-
  The Power Platform Data Record Resource allows the management of configuration records that are stored in Dataverse as records. This resource is not recommended for managing business data or other data that may be changed by Dataverse users in the context of normal business activities.
---

# powerplatform_data_record (Resource)

The Power Platform Data Record Resource allows the management of configuration records that are stored in Dataverse as records. This resource is not recommended for managing business data or other data that may be changed by Dataverse users in the context of normal business activities.

~> **Use this resource for configuration data only.** `powerplatform_data_record` is intended to manage configuration that is stored as Dataverse rows, such as business units, security roles, teams, application users or environment variable values. Do not use it to manage business data: rows that people or business processes create and change as part of normal work, such as accounts, contacts, leads, opportunities or cases. Terraform treats your configuration as the source of truth for every row it manages, so changes made to those rows outside Terraform can be overwritten or deleted on the next apply.

Data Record is a special type of a resources, that allows creation of any type Dataverese table record. The syntax for working with `data_record` resource is simmilar to raw WebAPI HTTP requests that this record uses:

- [WebAPI overview - Power Platform | Microsoft Learn](https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/overview)

## Lookup (relationship) columns

A lookup column is set by nesting a `table_logical_name` / `data_record_id` object under the column name. The provider turns that into an OData `@odata.bind` binding, using **the attribute name exactly as you wrote it**.

That name must be the relationship's `ReferencingEntityNavigationPropertyName`, which is not always the lowercase lookup column name. Most system relationships do match (`primarycontactid`, `businessunitid`), but some expose a PascalCase navigation property instead — for example `environmentvariablevalue` -> `environmentvariabledefinition` binds as `EnvironmentVariableDefinitionId`, not `environmentvariabledefinitionid`. Using the column name where the navigation property differs fails with `0x80048d19` ("undeclared property").

Look the correct name up against your own environment:

```text
GET https://<environment>.crm.dynamics.com/api/data/v9.2/EntityDefinitions(LogicalName='<table>')/ManyToOneRelationships?$select=ReferencingEntityNavigationPropertyName,ReferencedEntity,ReferencingAttribute
```

## Relationship (collection) columns

A collection-valued relationship (one-to-many or many-to-many, for example `teamroles_association`) is set by assigning a set of `table_logical_name` / `data_record_id` objects to the relationship name. On every apply the provider associates the listed rows with this record.

When you remove an item from that set, the optional `relationship_removal_behavior` attribute decides what happens to the related row:

- `delete_related_record` (the default when the attribute is not set): the related row itself is deleted from Dataverse. This is the behavior of earlier provider versions.
- `disassociate`: only the link between the two rows is removed, with a `DELETE <record>/<relationship>(<related id>)/$ref` request, and the related row is kept.

Use `disassociate` when the related rows are managed elsewhere (for example security roles assigned to a team or an application user), so that removing an assignment does not delete the role. This setting only applies to items removed from a relationship during an update; it does not change what happens when the `powerplatform_data_record` resource itself is destroyed.

## Example Usage

The following examples show how to use the `data_record` resource to configure some of the most common Dataverse settings.  These are minimal examples just to show the syntax, and do not include all possible configuration options.  Use these as a starting point if you need to set additional fields.

### Business Units

Example of how to create a [Business Unit](https://learn.microsoft.com/power-platform/admin/create-edit-business-units)

```terraform
terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

variable "environment_id" {
  type = string
}

variable "name" {
  type = string
}

variable "costcenter" {
  type = string
}

variable "parent_business_unit_id" {
  type = string
}

resource "powerplatform_data_record" "business_unit" {
  environment_id     = var.environment_id
  table_logical_name = "businessunit"
  disable_on_destroy = true
  columns = {
    name       = var.name
    costcenter = var.costcenter
    parentbusinessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.parent_business_unit_id
    }
  }
}

output "resource_id" {
  value = powerplatform_data_record.business_unit.id
}

output "resource" {
  value = powerplatform_data_record.business_unit
}
```

### Application User

Example of how to create an [Application User](https://learn.microsoft.com/power-platform/admin/manage-application-users)

```terraform
terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

variable "environment_id" {
  description = "The unique identifier of the environment"
  type        = string
  validation {
    condition     = length(var.environment_id) > 0
    error_message = "The environment id must not be empty"
  }

}

variable "application_id" {
  description = "EntraId clientid of the application"
  type        = string
  validation {
    condition     = length(var.application_id) > 0
    error_message = "The application id must not be empty"
  }
}

variable "business_unit_id" {
  description = "Unique identifier of the business unit"
  type        = string
  validation {
    condition     = length(var.business_unit_id) > 0
    error_message = "The business unit id must not be empty"
  }
}

variable "role_ids" {
  type        = set(string)
  description = "The role ids that are granted to this application user"
}


resource "powerplatform_data_record" "app_user" {
  table_logical_name = "systemuser"
  environment_id     = var.environment_id
  disable_on_destroy = true # Application Users cannot be deleted without being disabled first
  # Removing a role from role_ids only unassigns it; the default would delete the role itself.
  relationship_removal_behavior = "disassociate"
  columns = {
    applicationid = var.application_id
    businessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.business_unit_id
    }
    systemuserroles_association =toset([for rid in var.role_ids : { table_logical_name = "role", data_record_id = tostring(rid) }])
  }
}

output "application_user_id" {
  value = powerplatform_data_record.app_user.id
}
```

### Role

Example of how to create a [Role](https://learn.microsoft.com/power-platform/admin/create-edit-security-role#create-a-security-role)

```terraform
terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

variable "environment_id" {
  description = "The unique identifier of the environment"
  type        = string
  validation {
    condition     = length(var.environment_id) > 0
    error_message = "The environment id must not be empty"
  }
}

variable "role_name" {
  description = "The name of the role"
  type        = string
  validation {
    condition     = length(var.role_name) > 0
    error_message = "The role name must not be empty"
  }
}

variable "business_unit_id" {
  description = "The unique identifier of the business unit"
  type        = string
  validation {
    condition     = length(var.business_unit_id) > 0
    error_message = "The business unit id must not be empty"
  }
}

resource "powerplatform_data_record" "role" {
  environment_id     = var.environment_id
  table_logical_name = "role"

  columns = {
    name = var.role_name

    businessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.business_unit_id
    }
  }
}

output "role_id" {
  value = powerplatform_data_record.role.id
}
```

### Team

Example of how to create a [Team](https://learn.microsoft.com/power-platform/admin/manage-teams)

```terraform
terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

variable "environment_id" {
  description = "The unique identifier of the environment"
  type        = string
  validation {
    condition     = length(var.environment_id) > 0
    error_message = "The environment id must not be empty"
  }

}

variable "team_name" {
  description = "The name of the team"
  type        = string
  validation {
    condition     = length(var.team_name) > 0
    error_message = "The team name must not be empty"
  }
}

variable "team_description" {
  description = "The description of the team"
  type        = string
}

variable "role_ids" {
  type        = set(string)
  description = "The role ids that are granted to this team"
  
}

resource "powerplatform_data_record" "team" {
  environment_id     = var.environment_id
  table_logical_name = "team"
  # Removing a role from role_ids only unassigns it; the default would delete the role itself.
  relationship_removal_behavior = "disassociate"
  columns = {
    name        = var.team_name
    description = var.team_description

    teamroles_association = tolist([for rid in var.role_ids : { table_logical_name = "role", data_record_id = tostring(rid) }])
  }
}

output "team_id" {
  value = powerplatform_data_record.team.id
}
```

## End to End Example

```terraform
terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

provider "azuread" {
  use_cli = true
}

resource "powerplatform_environment" "data_record_example_env" {
  display_name     = "powerplatform_data_record_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

# get the root business unit by querying for the business unit without a parent
data "powerplatform_data_records" "root_business_unit" {
  environment_id    = powerplatform_environment.data_record_example_env.id
  entity_collection = "businessunits"
  filter            = "parentbusinessunitid eq null"
  select            = ["name"]
}

# Create a new business unit with the root business unit as parent
module "business_unit" {
  source                  = "./res_business_unit"
  environment_id          = powerplatform_environment.data_record_example_env.id
  name                    = "Sales"
  costcenter              = "123"
  parent_business_unit_id = one(data.powerplatform_data_records.root_business_unit.rows).businessunitid
}

# Create a new role
module "custom_role" {
  source           = "./res_role"
  environment_id   = powerplatform_environment.data_record_example_env.id
  role_name        = "my custom role"
  business_unit_id = one(data.powerplatform_data_records.root_business_unit.rows).businessunitid
}

module "team" {
  source           = "./res_team"
  environment_id   = powerplatform_environment.data_record_example_env.id
  team_name        = "main team"
  team_description = "main team description"
  role_ids         = [module.custom_role.role_id]

}

resource "azuread_application_registration" "data_record_app_user" {
  display_name = "powerplatform_data_record_example"
}

resource "azuread_service_principal" "data_record_app_user" {
  client_id = azuread_application_registration.data_record_app_user.client_id
}

module "application_user" {
  source           = "./res_application_user"
  environment_id   = powerplatform_environment.data_record_example_env.id
  application_id   = azuread_application_registration.data_record_app_user.client_id
  business_unit_id = one(data.powerplatform_data_records.root_business_unit.rows).businessunitid
  role_ids         = [module.custom_role.role_id]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `columns` (Dynamic) Columns of the data record table
- `environment_id` (String) Id of the Dynamics 365 environment
- `table_logical_name` (String) Logical name of the data record table

### Optional

- `disable_on_destroy` (Boolean) If true, the resource will either set isdisabled to true or statecode to 1 with a PATCH request, before attempting to delete the record.
- `relationship_removal_behavior` (String) How to handle an item that is removed from a collection-valued relationship in `columns`. `delete_related_record` (the default when not set) deletes the related row itself. `disassociate` only removes the link between the two rows and keeps the related row.
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `id` (String) Unique id (guid)

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
# Data record resource can be imported using the environment id, the table logical name and the data record id
# in the format: <environment_id>/<table_logical_name>(<data_record_id>)
terraform import powerplatform_data_record.example "00000000-0000-0000-0000-000000000000/account(00000000-0000-0000-0000-000000000001)"
```
