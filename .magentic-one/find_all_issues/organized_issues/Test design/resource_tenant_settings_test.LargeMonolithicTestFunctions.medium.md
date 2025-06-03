# Title

Excessively Large and Monolithic Test Functions Reduce Maintainability

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings_test.go

## Problem

The test functions are very large, with many repeated calls to `resource.TestCheckResourceAttr`, covering dozens of settings in one test function. These large tests can make it difficult for contributors to identify which individual behavior is being tested or which attribute fails if something is broken. They also hinder pinpointing test purposes and increase merge conflict risks as the codebase evolves.

## Impact

Reduces test granularity and maintainability. Small changes in the resource may break a large test, making it hard to understand or resolve. It also makes code review more burdensome and discourages adding narrowly scoped, meaningful, easier-to-debug tests. Severity: medium.

## Location

Inside test structure, e.g., `Check: resource.ComposeAggregateTestCheckFunc(`, where there are 30+ attribute checks in a single test.

## Code Issue

```go
Check: resource.ComposeAggregateTestCheckFunc(
    resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "disable_capacity_allocation_by_environment_admins", "false"),
    ...
    resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", "false"),
),
```
(Repeated in multiple test steps/functions, with many attribute checks.)

## Fix

Break apart the huge list of `resource.TestCheckResourceAttr` calls into helper functions grouped by logical functionality (e.g., basic booleans, licensing, governance, or other substructures). Then, compose the whole test from these helpers for readability and future maintainability. For instance:

```go
func checkBasicBooleans(expected string) resource.TestCheckFunc {
    return resource.ComposeTestCheckFunc(
        resource.TestCheckResourceAttr("powerplatform_tenant_settings.settings", "walk_me_opt_out", expected),
        ...
    )
}
```
Then,

```go
Check: resource.ComposeAggregateTestCheckFunc(
    checkBasicBooleans("false"),
    checkLicensingAttrs("85"),
    ...
),
```

This allows maintainers to more easily locate tests for a particular section of the resource and makes the file more readable and maintainable.
