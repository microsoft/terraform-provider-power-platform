# Title

Magic Strings for Resource Attribute Names and Expected Values

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings_test.go

## Problem

Throughout the tests, resource attribute names and expected values (e.g., `"disable_capacity_allocation_by_environment_admins"`, `"false"`, etc.) are hardcoded in every test step of every test function as string literals. This introduces risk: typos or copy-paste errors could cause silent failures, and any refactor or rename of attributes must be hunted down through the file as string matches.

## Impact

Maintainability and type-safety are reduced. Detection of typos or refactoring is much harder, and any change to the resource schema will require repetitive edits with risk of missing occurrences. This is a common cause of hard-to-detect bugs. Severity: medium.

## Location

Every call to `resource.TestCheckResourceAttr`:

```go
resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "false"),
...
resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "85"),
...
```

## Code Issue

```go
resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "power_platform.licensing.storage_capacity_consumption_warning_threshold", "85")
```

## Fix

Refactor the attribute names (and, if useful, frequently-repeated expected values) into typed constants at the top of the file. For example:

```go
const attrDisableCapacity = "disable_capacity_allocation_by_environment_admins"
const attrLicensingThreshold = "power_platform.licensing.storage_capacity_consumption_warning_threshold"
const expectedFalse = "false"
const expectedTrue  = "true"
```

Then use:

```go
resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", attrDisableCapacity, expectedFalse)
```

This greatly improves readability, enables automatic refactoring tools, and prevents subtle string mismatches.
