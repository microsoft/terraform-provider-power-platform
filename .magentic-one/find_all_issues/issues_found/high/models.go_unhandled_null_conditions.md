### Title

Unhandled Null Conditions

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

Code blocks checking for `IsNull()` or `IsUnknown()` properties do not handle all branches effectively. For example:

```go
if !auditSettingsObject.IsNull() && !auditSettingsObject.IsUnknown() {
    // Conversion logic
} else {
    // Missing handling for null/unknown branches
}
```

Whenever the object is `Null` or `Unknown`, subsequent logic has no resolution behavior.

### Impact

Severity: High.
This could cause undefined or unpredictable states in dependent logic, with runtime errors surfacing in production failures.

### Location

File: models.go
Functions: `convertFromEnvironmentSettingsModel`, `convertFromEnvironmentBehaviorSettings`

### Code Issue

```go
if !auditSettingsObject.IsNull() && !auditSettingsObject.IsUnknown() {
    // Conversion logic
}
```

### Fix

Add explicit handling logic for null/unknown states to avoid propagating undefined behavior.

```go
if !auditSettingsObject.IsNull() && !auditSettingsObject.IsUnknown() {
    // Conversion logic
} else {
    return nil, errors.New("audit settings object is either null or unknown")
}
```

Provides a safety check with an appropriate fallback mechanism.
