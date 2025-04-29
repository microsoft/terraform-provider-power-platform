### Title

Invalid Error Propagation

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

When creating errors, the errors.New and fmt.Errorf methods are directly used, but they might lack critical context. For example:

```go
return nil, errors.New("failed to convert audit settings to ObjectValue")
return nil, fmt.Errorf("failed to convert audit settings: %v", err)
```

Error messages do not contain sufficient context that can help track the source of a problem, especially across different service layers.

### Impact

This has a high-severity impact, as insufficient error context can make debugging complicated and lead to wasted developer time. It also reduces readability for log systems and compromises the traceability of issues.

### Location

File: /workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go
Function: convertFromEnvironmentSettingsModel
Block: Error handling for `auditSettings` object conversion.

### Code Issue

```go
return nil, errors.New("failed to convert audit settings to ObjectValue")
return nil, fmt.Errorf("failed to convert audit settings: %v", err)
```

### Fix

Enrich error contexts for debugging.

```go
return nil, errors.New("convertFromEnvironmentSettingsModel: failed to convert audit settings to ObjectValue")
return nil, fmt.Errorf("convertFromEnvironmentSettingsModel: failed to convert audit settings: %v", err)
```
This fix adds the function name as a prefix, enriching the source context of the error.
