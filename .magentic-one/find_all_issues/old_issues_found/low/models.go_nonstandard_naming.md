### Title

Nonstandard Naming

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

Certain variable and method names like `environmentSettingsDto` and `auditAndLogsSourceModel` are too verbose or unclear, reducing readability:

```go
environmentSettingsDto := &environmentSettingsDto{}
var auditAndLogsSourceModel AuditSettingsSourceModel
```

### Impact

Severity: Low.
Nonstandard naming increases cognitive load for developers unfamiliar with naming conventions, impacting maintainability.

### Location

File: models.go
Spread across all major function definitions.

### Code Issue

```go
environmentSettingsDto := &environmentSettingsDto{}
var auditAndLogsSourceModel AuditSettingsSourceModel
```

### Fix

Adopt clearer and concise naming conventions:

```go
settings := &SettingsDTO{}
var auditAndLogsModel AuditSourceModel
```

Improves readability and consistency across the codebase.
