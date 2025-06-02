# Title

Excessive Debug-Level Log Messages

## Path

`/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go`

## Problem

The `Metadata` function logs debug messages at the root level without distinguishing important diagnostic steps. Example:

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

While debug logs are useful, excessive or unfiltered usage tends to clutter log files without providing actionable insight.

## Impact

Excessive logging can lead to slower performance and bloated log files, impacting readability and operational troubleshooting. Severity: **Low**

## Location

Function: `Metadata`
Line: Debug-level logging statement.
File: `/internal/services/data_record/resource_data_record.go`
Path within the function implementation.

## Code Issue

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

## Fix

Introduce logging conditions or use log levels effectively based on the context:

```go
if DebugLoggingEnabled {
   tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
```