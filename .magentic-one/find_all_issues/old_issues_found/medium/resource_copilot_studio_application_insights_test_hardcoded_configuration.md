# Title

Static Configuration Inside Resource Blocks

##

`/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights_test.go`

## Problem

The configuration blocks in the acceptance tests contain hardcoded values such as `location` and `schema_validation_enabled`. This reduces flexibility and limits the extent to which these tests can be reused across regions or schemas.

## Impact

- **Severity:** Medium
- Hardcoded configurations limit portability and adaptability of tests.
- Test cases may fail or not run correctly if environments differ.

## Location

Lines 26-34 and Lines 91-96.

## Code Issue

```go
	type     = "Microsoft.Resources/resourceGroups@2021-04-01"
	location = "East US"
	name     = "power-platform-app-insights-rg-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
```

```go
	schema_validation_enabled = false
	type = "Microsoft.Insights/components@2020-02-02"
```

## Fix

Use dynamically parameterized values or environment-specific configurations instead of hardcoding.

```go
	resource "azapi_resource" "rg_example" {
		type     = "Microsoft.Resources/resourceGroups@2021-04-01"
		location = var.region // Use variable instead of hardcoded value
		name     = var.resource_group_name_prefix + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) // Make `resource_group_name_prefix` dynamic
	}

	resource "azapi_resource" "app_insights" {
		schema_validation_enabled = var.schema_validation_enabled // Replace hardcoded false with variable
		type = var.resource_type // Parameterize resource type
		location = var.region // Use parameterized region
	}
```
