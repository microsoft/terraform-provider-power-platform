---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_connections Data Source - Power Platform"
subcategory: ""
description: |-
  Fetches a list of Connection https://learn.microsoft.com/en-us/power-apps/maker/canvas-apps/add-manage-connections for a given environment. Each connection represents an connection instance to an external data source or service.
---

# powerplatform_connections (Data Source)

Fetches a list of [Connection](https://learn.microsoft.com/en-us/power-apps/maker/canvas-apps/add-manage-connections) for a given environment. Each connection represents an connection instance to an external data source or service.

## Example Usage

```terraform
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

data "powerplatform_environments" "all_environments" {}

data "powerplatform_connections" "all_connections" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) Environment Id. The unique identifier of the environment that the connection are associated with.

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `connections` (Attributes List) List of Connections (see [below for nested schema](#nestedatt--connections))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.


<a id="nestedatt--connections"></a>
### Nested Schema for `connections`

Read-Only:

- `connection_parameters` (String) Connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.
- `connection_parameters_set` (String) Set of connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.
- `display_name` (String) Display name of the connection.
- `id` (String) Unique connection id
- `name` (String) Name of the connection.
- `status` (Set of String) List of connection statuses
