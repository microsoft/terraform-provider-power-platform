# Title

Repetitive If-Null-And-Unknown Pattern Leads to Boilerplate

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

Throughout the DTO conversion logic, code is repeatedly written for null and unknown checks before assigning or converting values:
```go
if !value.IsNull() && !value.IsUnknown() {
    target = value.ValueBoolPointer()
}
```
This leads to excessive boilerplate and makes the code harder to read and maintain.

## Impact

- Reduced code maintainability.
- More places for subtle bugs if new fields are added or checks are missed.
- Contributes to unnecessarily bloated and less readable functions.
- Increases technical debt over time. Severity: medium.

## Location

- All major conversion functions: e.g., `convertFromTenantSettingsModel`, `convertTeamsIntegrationModel`, `convertPowerAppsModel`, etc.

## Code Issue

```go
if !tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers.IsNull() && !tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers.IsUnknown() {
    tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers = tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers.ValueBoolPointer()
}
```

## Fix

Abstract common check-use patterns into reusable helper functions, such as:

```go
func getBoolPointer(v basetypes.BoolValue) *bool {
    if !v.IsNull() && !v.IsUnknown() {
        return v.ValueBoolPointer()
    }
    return nil
}

// Usage:
tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers = getBoolPointer(tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers)
```

This enables more concise and reliable code, and helps enforce best practices for null/unknown-safe conversions.

