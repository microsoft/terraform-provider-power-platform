### Title

Incorrect Type Assertions

### Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

### Problem

Some type assertion blocks do not have fallback mechanisms, directly assuming type correctness without validation:

```go
pluginSettingsValue, ok := pluginSettings.(basetypes.StringValue)
if !ok {
    return nil, errors.New("pluginSettings is not of type basetypes.StringValue")
}
```

While validation exists, assumptions made in other usage portions lack consistency, potentially leading to mismatched object behavior.

### Impact

Severity: High.
Incorrect type assumptions can lead to runtime panics, compromising service reliability and availability.

### Location

File: models.go
Function: convertFromEnvironmentSettingsModel.
Block: `pluginSettings` type assertion logic.

### Code Issue

```go
pluginSettingsValue, ok := pluginSettings.(basetypes.StringValue)
if !ok {
    return nil, errors.New("pluginSettings is not of type basetypes.StringValue")
}
```

### Fix

Ensure broader validity and fallback mechanisms for type assertions.

```go
pluginSettingsValue, ok := pluginSettings.(basetypes.StringValue)
if !ok {
    return nil, fmt.Errorf("pluginSettings: expected basetypes.StringValue, got %T", pluginSettings)
}
```
Adding clear fallback structure improves type safety and prevents crashes.
