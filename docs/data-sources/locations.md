---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_locations Data Source - Power Platform"
subcategory: ""
description: |-
  Fetches the list of available Dynamics 365 locations. For more information see Power Platform Geos https://learn.microsoft.com/power-platform/admin/regions-overview
---

# powerplatform_locations (Data Source)

Fetches the list of available Dynamics 365 locations. For more information see [Power Platform Geos](https://learn.microsoft.com/power-platform/admin/regions-overview)

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

data "powerplatform_locations" "all_locations" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))

### Read-Only

- `locations` (Attributes List) List of available locations (see [below for nested schema](#nestedatt--locations))

<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.


<a id="nestedatt--locations"></a>
### Nested Schema for `locations`

Read-Only:

- `azure_regions` (List of String) List of Azure regions
- `can_provision_customer_engagement_database` (Boolean) Can the location provision a customer engagement database
- `can_provision_database` (Boolean) Can the location provision a database
- `code` (String) Code of the location
- `display_name` (String) Display name of the location
- `id` (String) Unique identifier of the location
- `is_default` (Boolean) Is the location default
- `is_disabled` (Boolean) Is the location disabled
- `name` (String) Name of the location
