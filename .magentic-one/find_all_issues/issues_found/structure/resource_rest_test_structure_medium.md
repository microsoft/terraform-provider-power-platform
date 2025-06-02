# Title

Redundant configuration blocks in test steps reduce clarity and maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest_test.go

## Problem

The configuration code for resources in the test steps is heavily duplicated between test steps in the acceptance test (i.e., the strings defined in the `Config` field of each `resource.TestStep`). This duplication makes the test harder to refactor and maintain, especially when changes to resource definition or structure are required—they must be updated in every location, increasing the risk of inconsistencies and copy-paste errors.

## Impact

This structural issue impacts maintainability and readability of tests. Updates to the test configuration require effort to update in several places, and changes might be missed. This is a medium severity issue, not impacting runtime directly but significantly affecting code hygiene and long-term maintenance.

## Location

- Each `Config:` string field in the `resource.TestStep` slices inside `TestAccTestRest_Validate_Create`:

```go
  {
    Config: `
      ...
      resource "powerplatform_environment" "env" {
        ...
      }
      ...
      resource "powerplatform_rest" "query" {
        ...
      }
      ...
    `,
    ...
  },
  {
    Config: `
      ...
      resource "powerplatform_environment" "env" {
        ...
      }
      ...
      resource "powerplatform_rest" "query" {
        ...
      }
      ...
    `,
    ...
  },
```

## Fix

Refactor common configuration fragments into constants or variables, or use Go's multiline string variables to house common resource configuration blocks. Interpolate variables where values change, for example:

```go
const envConfig = `
resource "powerplatform_environment" "env" {
  display_name     = "%s"
  location         = "unitedstates"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}
`

const restConfig = `
locals {
   ...
}
resource "powerplatform_rest" "query" {
   ...
}
`

testConfigStep1 := fmt.Sprintf(envConfig, mocks.TestName()) + restConfig
testConfigStep2 := fmt.Sprintf(envConfig, mocks.TestName()) + restConfigStep2Diff // customize as needed

Steps: []resource.TestStep{
  {
    Config: testConfigStep1,
    ...
  },
  {
    Config: testConfigStep2,
    ...
  },
}
```

Or, if configuration difference is minimal, parameterize using Go string formatting and compose as needed. This will reduce duplication, making tests easier to update and understand.
