# Issue Report #3

### Title: Deprecated Use of `Logger` Without Proper Context Validation

### Path to the file: `/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps.go`

---

## Problem

The `tflog.Debug` logger is utilized to log debugging messages (`METADATA: %s` in `Metadata` function). However, there is no check for whether the logger is properly initialized or configured for the current context. This could result in a failure to log or raise unnecessary errors during runtime when debugging is enabled.

---

## Impact

Uninitialized or inaccessible logging could lead to loss of critical debug information, especially when diagnosing issues in metadata resolve steps. Severity: **Medium**

---

## Location

**Function:** Metadata function, during logging.

---

## Code Issue

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
```

---

## Fix

Check the validity of logging capabilities within the given context before calling `tflog.Debug`:

```go
if tflog.IsEnabled(ctx) {
    tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
} else {
    fmt.Println("Debugging logger not enabled. Skipping debug information.")
}
```

This ensures graceful handling when the logger is unavailable, avoiding potential runtime errors or application disruptions.
