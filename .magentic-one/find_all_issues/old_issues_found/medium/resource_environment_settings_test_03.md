# Title

Repeated Configuration with Similar Test Steps

## Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resource_environment_settings_test.go

## Problem

There are several instances in the file where similar configurations are repeated across test cases (`Config:` blocks for `powerplatform_environment_settings` and `powerplatform_environment` resources). This introduces redundancy within the codebase, causing inefficiency and making maintenance harder with future changes.

## Impact

Medium — Repeated configuration can lead to increased maintenance efforts and larger testing files. Additionally, bugs in common configurations could propagate to multiple test cases, making debugging more difficult.

## Location

Examples:
1. Configuration in `TestAccTestEnvironmentSettingsResource_Validate_Create_Empty_Settings`
2. Configuration in `TestAccTestEnvironmentSettingsResource_Validate_Read`
3. Configuration in `TestUnitTestEnvironmentSettingsResource_Validate_Read`

## Code Issue

```go
Config: `
  resource "powerplatform_environment_settings" "settings" {
    environment_id                         = powerplatform_environment.example_environment_settings.id
  }`,
```

This style is repeated across multiple test cases.

## Fix

Refactor the test configuration to use shared helper functions or constants for common repetitive configurations. For example:

```go
func getEnvironmentConfigExample(environmentID string) string {
	return fmt.Sprintf(`
    resource "powerplatform_environment_settings" "settings" {
      environment_id = "%s"
    }
  `, environmentID)
}

func TestAccTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) {
    ...
    Config: getEnvironmentConfigExample("00000000-0000-0000-0000-000000000001"),
    ...
}
```