# Missing Error Handling for os.Getwd in importSolution

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem
In the function `importSolution`, the call to `os.Getwd()` ignores the error return value. The statement:

```go
cwd, _ := os.Getwd()
```

simply discards the error. This method can return an error, and if the current directory isn't available for any reason, the value of `cwd` could be invalid or misleading. Not even logging the error means that diagnostics or troubleshooting is hindered if such a condition arises.

## Impact
- **Severity:** Medium
- Masking possible failures makes troubleshooting harder and hides potential edge cases in filesystem access from the user and maintainers.
- Logging or diagnosing environment issues becomes more difficult if errors are silently ignored.

## Location
`importSolution` function, where `os.Getwd()` is used.

## Code Issue
```go
cwd, _ := os.Getwd()
tflog.Debug(ctx, fmt.Sprintf("Current working directory: %s", cwd))
```

## Fix
Capture and log the error for better diagnostics:

```go
cwd, err := os.Getwd()
if err != nil {
    tflog.Warn(ctx, fmt.Sprintf("Failed to get working directory: %s", err.Error()))
} else {
    tflog.Debug(ctx, fmt.Sprintf("Current working directory: %s", cwd))
}
```
