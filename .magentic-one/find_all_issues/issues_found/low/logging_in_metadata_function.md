# Title
Logging in Metadata Function

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem
The `Metadata` function uses `tflog.Debug` for logging metadata information. While this is acceptable, there is no consideration for log verbosity levels, which might lead to excessive debug logs in some environments.

## Impact
The excessive usage of debug logs can clutter log files and reduce readability, making it harder to debug real issues in production or higher verbosity environments. The severity of this issue is low.

## Location
Function: Metadata

## Code Issue
```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

## Fix
Use log verbosity levels to differentiate between information, warnings, and debug logs. Enhance the logging mechanism accordingly.
```go
if helpers.IsLogLevelDebug(ctx) {
    tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
```
