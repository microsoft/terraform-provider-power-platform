### Title

Inconsistent Error Checking

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

Some error checks in the code are skipped or incompletely handled, particularly when calling certain methods that involve conversions like:

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
    return nil, errors.New("failed to convert audit settings to ObjectValue")
}
```

Certain external calls might fail due to runtime issues, but the impact contexts or follow-up checks are absent.

### Impact

Severity: Critical.
This could lead to runtime errors where corrupted, incomplete or unexpected data fails conversions but continues execution. Such errors are catastrophic in production environments and compromise data integrity.

### Location

File: /workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go
Function: convertFromEnvironmentSettingsModel.
Block: `auditSettingsObject` conversion.

### Code Issue

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
    return nil, errors.New("failed to convert audit settings to ObjectValue")
}
```

### Fix

Always validate external object interactions and handle follow-up conditions properly.

```go
objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
if !ok {
    return nil, fmt.Errorf("convertFromEnvironmentSettingsModel: invalid audit settings object conversion")
}
```
This ensures structured debugging without arbitrary assumptions.
