# Field Name and JSON Tag Mismatch (PowerAppsComponentFrameworkForCanvasApps)

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/dto.go

## Problem

The Go struct field `PowerAppsComponentFrameworkForCanvasApps` does not match the JSON tag `iscustomcontrolsincanvasappsenabled`, reducing maintainability and readability. The field name should more closely reflect the JSON tag or domain vocabulary for clarity and consistency.

## Impact

This inconsistency can confuse developers and cause maintenance issues, especially when generating or mapping from API documentation. **Severity:** low.

## Location

```go
PowerAppsComponentFrameworkForCanvasApps *bool `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```

## Code Issue

```go
PowerAppsComponentFrameworkForCanvasApps *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```

## Fix

Rename the field so it matches the domain concept and JSON tag:

```go
IsCustomControlsInCanvasAppsEnabled *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```
