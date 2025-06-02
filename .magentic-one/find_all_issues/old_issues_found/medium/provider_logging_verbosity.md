# Title

Logging Verbosity Issue

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

Logging in the provider uses `tflog.Info` and `tflog.Warn` without adjustable verbosity levels. There is no mechanism for users to configure the verbosity level of logs to match their preference or environment, particularly important for debugging or production settings.

## Impact

Severity: Medium

- Unnecessarily verbose logs might inundate the user with information, making debugging harder.
- In production, verbose logs could leak sensitive operational information.

## Location

Throughout logging statements in the file.

Example:

1. Logging for test mode:

```go
tflog.Info(ctx, "Test mode enabled. Authentication requests will not be sent to the backend APIs.")
```

2. Logging for CLI mode:

```go
tflog.Info(ctx, "Using CLI for authentication")
```

## Code Issue

```go
tflog.Info(ctx, "Test mode enabled. Authentication requests will not be sent to the backend APIs.")
...
tflog.Info(ctx, "Using CLI for authentication")
```

## Fix

Introduce a logging verbosity level and adjust the log statements to adhere to this configured level:

```go
var verbosityLevel int

func log(ctx context.Context, level int, message string, params ...map[string]any) {
    if level <= verbosityLevel {
        tflog.Info(ctx, message, params...)
    }
}

// Usage:
log(ctx, 1, "Using CLI for authentication")
log(ctx, 2, "Test mode enabled.")
```