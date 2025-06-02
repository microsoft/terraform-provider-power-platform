# Title

DTO conversion has an inconsistent error handling approach

## Path

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

Within the `convertFromTenantSettingsModel` function, some conversions use error handling (e.g., `convertSearchModel`), while others (like `convertTeamsIntegrationModel`) do not. For instance:

```go
err := convertSearchModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
if err != nil {
    return tenantSettingsDto, err
}
convertTeamsIntegrationModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
```

This discrepancy can lead to issues where conversions silently fail without propagating errors.

## Impact

Inconsistent error handling can result in unreliable data transformation processes, making it difficult to diagnose issues. This could lead to runtime failures or unexpected behavior. Severity: High.

## Location

This issue is scattered across the `convertFromTenantSettingsModel` function.

## Code Issue

Example code snippet demonstrating inconsistent handling:

```go
err := convertSearchModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
if err != nil {
    return tenantSettingsDto, err
}
convertTeamsIntegrationModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
```

## Fix

Ensure all conversion functions consistently implement error handling, and update the logic to propagate errors where necessary:

```go
// Add error handling for all conversion functions
err := convertSearchModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
if err != nil {
    return tenantSettingsDto, err
}

err = convertTeamsIntegrationModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
if err != nil {
    return tenantSettingsDto, err
}
```
