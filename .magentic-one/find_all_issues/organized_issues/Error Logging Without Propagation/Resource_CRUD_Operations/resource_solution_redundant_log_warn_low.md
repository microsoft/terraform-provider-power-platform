# Redundant Logging: Use of tflog.Warn for Successful Checksum Calculation

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem
Within the `Create` function, after successfully calculating a checksum for the settings or solution file, the code logs this with `tflog.Warn`:

```go
tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of settings file: %s", value))
```

Checksum calculation is a successful and routine operation, so warning level logging is inappropriate. It may create noise in logs, distracting from real warnings. This may reflect miscommunication of severity/intent or copy-paste error.

## Impact
- **Severity:** Low
- Could mislead users/maintainers or clutter logs
- Minor impact, but relevant for operational hygiene and clarity

## Location
In the `Create` function's handling of the settings and solution file checksums.

## Code Issue
```go
if err != nil {
    resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
} else {
    plan.SettingsFileChecksum = types.StringValue(value)
    tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of settings file: %s", value))
}
```

## Fix
Lower the log level to `tflog.Debug` or remove it if not actually useful:

```go
if err != nil {
    resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
} else {
    plan.SettingsFileChecksum = types.StringValue(value)
    tflog.Debug(ctx, fmt.Sprintf("CREATE calculated SHA256 hash of settings file: %s", value))
}
```
