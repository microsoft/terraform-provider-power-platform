# Issue 3

## Title

Hardcoded Configuration in Test Step Makes Test Non-Reusable

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings_test.go`

## Problem

The `Config` in the test step for `TestAccTestEnvironmentSettingsDataSource_Validate_Read` uses hardcoded values for `location`, `environment_type`, `language_code`, `currency_code`, and other attributes. Hardcoded configurations make tests difficult to reuse and maintain.

## Impact

Hardcoded values reduce test flexibility and increase maintenance costs, as changes to the tested configuration require editing the test code directly. Severity: **Low**

## Location

```go
Config: `
resource "powerplatform_environment" "example_environment_settings" {
	display_name      = "` + mocks.TestName() + `"
	location          = "unitedstates"
	environment_type  = "Sandbox"
	dataverse = {
		language_code     = "1033"
		currency_code     = "USD"
		security_group_id = "00000000-0000-0000-0000-000000000000"
	}
}
...
`,
```

## Fix

Use variables to define the configuration values within the test setup, allowing greater flexibility and reusability of the test cases.

```go
Config := fmt.Sprintf(`
resource "powerplatform_environment" "example_environment_settings" {
	display_name      = "%s"
	location          = "%s"
	environment_type  = "%s"
	dataverse = {
		language_code     = "%s"
		currency_code     = "%s"
		security_group_id = "%s"
	}
}
...
`, mocks.TestName(), "location_value", "environment_type_value", "language_code_value", "currency_code_value", "security_group_id_value")
```