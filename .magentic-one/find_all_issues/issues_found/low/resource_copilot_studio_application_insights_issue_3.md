# Title
Excessive Logging in Metadata Method

##

/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

The `Metadata` method performs a debug-level logging of the resource type name without checking the logger's verbosity level. This can lead to unnecessary, verbose outputs for production builds where minimal logging is preferred.

## Impact

- Unnecessary log entries clutter system logs, especially in production, making debugging harder.
- While unlikely to cause errors, this can reduce system performance slightly due to increased I/O.
- Severity: Low

## Location

Lines: near 55

## Code Issue

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

## Fix

Introduce a verbosity-level or conditional check around the logging statement.

```go
if tflog.SettingEnabled(ctx, tflog.DebugLevel) {
    tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
```
