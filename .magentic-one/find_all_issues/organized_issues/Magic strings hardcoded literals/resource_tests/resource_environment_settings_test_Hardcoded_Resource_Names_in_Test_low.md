# Hardcoded Resource Names in Test

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resource_environment_settings_test.go

## Problem

Test cases use hardcoded resource names, such as `powerplatform_environment_settings.settings` and environment IDs, which reduces maintainability and flexibility. If changes in resources or configurations happen, many tests must change in parallel.

## Impact

Low severity (maintainability): Changes in resource naming or configuration will require updates to all tests.

## Location

```go
resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", ...)
```

## Code Issue

```go
resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "false"),
```

## Fix

Abstract resource naming and IDs to reusable consts or variables if used in multiple places.

```go
const envSettingsResourceName = "powerplatform_environment_settings.settings"

// Then use:
resource.TestCheckResourceAttr(envSettingsResourceName, "audit_and_logs.audit_settings.is_audit_enabled", "false"),
```
