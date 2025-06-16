# Issue: Inconsistent Naming Conventions for Client Field

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/models.go

## Problem

The field names for the client in `TenantSettingsDataSource` and `TenantSettingsResource` structs are inconsistent:
- In `TenantSettingsDataSource`: `TenantSettingsClient client`
- In `TenantSettingsResource`: `TenantSettingClient client`

Additionally, the type `client` is not declared anywhere in this file, which may indicate a missing import or mistyped type.

## Impact

This increases cognitive load and can lead to confusion or bugs during development and maintenance, especially since these types are structurally related. Severity: **medium**.

## Location

- Line 13 (`TenantSettingsDataSource`)
- Bottom of file (`TenantSettingsResource`)

## Code Issue

```go
type TenantSettingsDataSource struct {
	helpers.TypeInfo
	TenantSettingsClient client
}
...
type TenantSettingsResource struct {
	helpers.TypeInfo
	TenantSettingClient client
}
```

## Fix

- Decide upon a consistent field name, such as `TenantSettingsClient`.
- Ensure the custom type `client` is defined and imported.
- Apply the chosen field name to all relevant struct definitions for consistency.

```go
type TenantSettingsDataSource struct {
	helpers.TypeInfo
	TenantSettingsClient client
}

type TenantSettingsResource struct {
	helpers.TypeInfo
	TenantSettingsClient client
}
```
