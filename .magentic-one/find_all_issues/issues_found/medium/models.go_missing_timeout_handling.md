### Title

Missing Timeout Handling

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

Some function blocks, such as `convertFromEnvironmentBehaviorSettings`, do not validate or handle requested timeouts effectively. For instance:

```go
behaviorSettings := environmentSettings.Product.Attributes()["behavior_settings"]
if behaviorSettings != nil && !behaviorSettings.IsNull() && !behaviorSettings.IsUnknown() {
    // Implementation logic
}
```

Timeout edge cases, such as `{"behavior_settings": timeout}` scenarios, are absent from checks.

### Impact

Severity: Medium.
Improper timeout handling could lead to delayed services and reduced performance.

### Location

File: models.go
Function: `convertFromEnvironmentBehaviorSettings`

### Code Issue

```go
behaviorSettings := environmentSettings.Product.Attributes()["behavior_settings"]
if behaviorSettings != nil && !behaviorSettings.IsNull() && !behaviorSettings.IsUnknown() {
    // Implementation logic
}
```

### Fix

Introduce timeout validation conditions within function layers.

```go
if timeoutExceeded(ctx) {
    return nil, errors.New("timeout occurred during behavior settings conversion")
}

behaviorSettings := environmentSettings.Product.Attributes()["behavior_settings"]
if behaviorSettings != nil && !behaviorSettings.IsNull() && !behaviorSettings.IsUnknown() {
    // Implementation logic
}
```

This ensures timeout events trigger appropriate failovers.
