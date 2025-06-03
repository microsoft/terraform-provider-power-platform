# Duplicate Test Config Code Should Be Extracted for DRY Principle

##

/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations_test.go

## Problem

The test configuration for the data source is repeated verbatim in both `TestAccLocationsDataSource_Validate_Read` and `TestUnitLocationsDataSource_Validate_Read`. This breaks the DRY (Don't Repeat Yourself) principle and makes future updates or corrections harder and error-prone.

## Impact

Makes the codebase harder to maintain, as any change to the test config has to be made in several places. Such duplication leads to inconsistencies and increases effort during refactoring or bug correction.

**Severity:** Low

## Location

- Multiple test functions in the same file

## Code Issue

```go
Config: `
	data "powerplatform_locations" "all_locations" {
	}`,
```

## Fix

### Extract configuration as a constant

```go
const locationsDataSourceConfig = `
	data "powerplatform_locations" "all_locations" {
	}`

// Then use in both tests:
Config: locationsDataSourceConfig,
```
