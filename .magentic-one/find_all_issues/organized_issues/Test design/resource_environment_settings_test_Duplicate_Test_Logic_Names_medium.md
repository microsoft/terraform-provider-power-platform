# Duplicate Test Logic/Names

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resource_environment_settings_test.go

## Problem

There is duplicated logic in tests with nearly identical names and structures for unit and acceptance testing (`TestAccTestEnvironmentSettingsResource_Validate_Create_Empty_Settings`, `TestUnitTestEnvironmentSettingsResource_Validate_Create_Empty_Settings`). This can lead to maintenance overhead—if fixing or improving one, the other may get out of sync.

## Impact

Medium, maintainability: Nuanced test behaviors could diverge. Future changes are prone to missing updates in either path.

## Location

Example:

```go
func TestUnitTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) { ... }
func TestAccTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) { ... }
```

## Code Issue

```go
func TestUnitTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) { ... }
func TestAccTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) { ... }
```

## Fix

Abstract shared logic into helper functions or table-driven tests.

```go
func validateCreateEmptySettings(t *testing.T, isUnitTest bool, factories ...any) {
	// common logic
}

func TestUnitTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) {
	validateCreateEmptySettings(t, true, mocks.TestUnitTestProtoV6ProviderFactories)
}
func TestAccTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) {
	validateCreateEmptySettings(t, false, mocks.TestAccProtoV6ProviderFactories)
}
```
