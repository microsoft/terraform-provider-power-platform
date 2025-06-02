# Title

Unit Tests and Acceptance Tests Share Too Much Copied Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings_test.go

## Problem

Both the unit tests and acceptance tests share the exact same, very long `Config` resource HCL code (sometimes with only small changes, e.g., swapping `true`/`false`). This makes the test file much larger than necessary and increases the likelihood that updates to the resource definition or the keys must be duplicated in multiple places. This repetitive pattern is prone to error and difficult to maintain over time.

## Impact

High risk to maintainability, code duplication, and potential drift between test cases. If a change is needed, it must be made in multiple locations. This also makes it harder to read and reason about the test content. Severity: medium.

## Location

Occurs in every test step where

```go
Config: `
resource "powerplatform_tenant_settings" "settings" {
...
}
`,
```

is copied into multiple test step definitions, especially across both unit and acceptance test blocks.

## Code Issue

```go
Config: `
resource "powerplatform_tenant_settings" "settings" {
    walk_me_opt_out = false
    ...
    user_management_settings = {
        enable_delete_disabled_user_in_all_environments = false
    }
}
`,
```

Repeated in multiple tests, with large bodies of identical HCL.

## Fix

Refactor the repeated configuration string out to helper functions or constants. For example:

```go
const baseTenantSettingsConfig = `
resource "powerplatform_tenant_settings" "settings" {
    walk_me_opt_out = %v
    ...
    user_management_settings = {
        enable_delete_disabled_user_in_all_environments = %v
    }
}
`

// Then, in the test:
Config: fmt.Sprintf(baseTenantSettingsConfig, "false", "false"),
```

Or, if the settings blocks are mostly the same, use test utility helpers to construct the strings or parts of the HCL dynamically.

This will make the tests more readable and much easier to maintain, while reducing the risk of errors due to missed updates in one location.
