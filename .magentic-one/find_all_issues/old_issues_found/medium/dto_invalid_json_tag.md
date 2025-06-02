# Title

JSON Tags in `environmentSettingsDto` Don't Match Struct Field Names

# Path

`/workspaces/terraform-provider-power-platform/internal/services/environment_settings/dto.go`

# Problem

The `json` key for the field `PowerAppsComponentFrameworkForCanvasApps` is specified as `"iscustomcontrolsincanvasappsenabled"`, which does not match its corresponding struct field name. This inconsistency can lead to confusion for developers and bugs in cases where automated tooling or strict schema validation is used.

# Impact

1. Misalignment between the Go struct field and JSON schema causes confusion, reducing code readability and maintainability.
2. Potential runtime issues, such as data not being marshaled/unmarshaled correctly.

**Severity**: **Medium**

# Location

```go
PowerAppsComponentFrameworkForCanvasApps *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```

# Fix

Correct the JSON tag to match the field name, or update the field name to align with the intended JSON key.

Option 1: Update the JSON tag.

```go
PowerAppsComponentFrameworkForCanvasApps *bool   `json:"powerAppsComponentFrameworkForCanvasApps,omitempty"`
```

Option 2: Rename the struct field.

```go
IsCustomControlsInCanvasAppsEnabled *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
```

If the existing JSON tag matches an established schema, favor Option 2 to ensure backwards compatibility with external systems.

---