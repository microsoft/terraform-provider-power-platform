---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "powerplatform_data_record Resource - powerplatform"
subcategory: ""
description: |-
  Resource for managing PowerPlatform Data Record
---

# powerplatform_data_record (Resource)

Resource for managing PowerPlatform Data Record

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
  tenant_id     = var.tenant_id
  client_id     = var.client_id
  client_secret = var.client_secret
}

data "powerplatform_environments" "all_environments" {}

resource "powerplatform_data_record" "data_record_by_environment_id" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  table_name     = "accounts"
  columns = {
    name = "Sample Account"
    creditonhold = false
    address1_latitude = 47.639583
    description = "This is the description of the sample account"
    revenue = 5000000
    accountcategorycode = 1
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `columns` (Dynamic) Columns of the data record table
- `environment_id` (String) Id of the Dynamics 365 environment
- `table_name` (String) Name of the data record table

### Read-Only

- `id` (String) Unique id (guid)