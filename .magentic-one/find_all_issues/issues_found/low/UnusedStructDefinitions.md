# Title

Unused Struct Definitions

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/models.go`

## Problem

Some structs in this file are defined but not utilized anywhere in the code. For example, the `TenantSettingsResource` struct appears unused. Code should avoid retaining redundant definitions unless intended for future use or reference.

## Impact

Unused struct definitions can lead to confusion for developers maintaining the code, unnecessarily increase the file's complexity and size, and add to compilation overhead. Severity: **low**.

## Location

Struct definition for `TenantSettingsResource`.

## Code Issue

```go
type TenantSettingsResource struct {
	helpers.TypeInfo
	TenantSettingClient client
}
```

## Fix

Remove or comment out the unused `TenantSettingsResource` struct. If its use is planned for the future, add a comment indicating its intended purpose:

```go
// TenantSettingsResource struct definition - currently unused, kept for future use.
type TenantSettingsResource struct {
	helpers.TypeInfo
	TenantSettingClient client
}

// Alternatively, remove the struct:
```