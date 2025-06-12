# Title

Inconsistent Naming Conventions for JSON Tags vs. Struct Field Names

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

The DTO struct field names use camelCase (e.g., `TenantId`, `PowerPlatform`), but the JSON tags and maps in conversion functions use snake_case (e.g., `tenantId`, `power_platform`). However, in several conversion functions and attribute maps, inconsistencies exist where some keys use camelCase and others snake_case. For example, keys like "power_platform" (snake_case in maps), but the field is `PowerPlatform`, and some referenced strings (such as map keys in function calls) do not always match their respective struct tags.

## Impact

This inconsistency could lead to serialization/deserialization bugs, confusion when mapping HTTP API responses/requests, and increased cognitive load for future maintainers. If the backend or integration expects strict casing, this may lead to deployment or runtime errors. Severity: medium.

## Location

- Throughout DTO struct definitions and conversion functions.
- Example: `powerPlatformAttributes["power_apps"]` vs. JSON tag `powerApps`, etc.

## Code Issue

```go
type powerPlatformSettingsDto struct {
    ...
    PowerApps *powerAppsSettingsDto `json:"powerApps,omitempty"`
    ...
}
...
func convertPowerAppsModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
    powerAppsObject := powerPlatformAttributes["power_apps"]
    ...
}
```

## Fix

Unify the naming conventions for map keys, struct fields, and JSON tags. Either:
- Always use snake_case everywhere (preferred for Go, for API consistency).
- Or, only diverge for struct tags but never for map keys.

Example (move to snake_case consistently):

```go
type powerPlatformSettingsDto struct {
    PowerApps *powerAppsSettingsDto `json:"power_apps,omitempty"` // snake_case
    // ...
}

// Now, always use "power_apps" in map lookups and JSON tags
powerAppsObject := powerPlatformAttributes["power_apps"]
```

This change should be applied throughout all DTOs and relevant conversion utilities for consistency.

