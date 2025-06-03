# Issue 1: Unnecessary Use of log.Default() for Logging

##

/workspaces/terraform-provider-power-platform/internal/modifiers/restore_original_value_modifier.go

## Problem

The code uses `log.Default().Printf` for logging inside plan modifiers. This is not suitable for Terraform provider code, as logs are not properly captured or displayed in Terraform output. Terraform's plugin framework typically provides log contexts or expects no standard logging. Logging to `log.Default()` may pollute test output or miss important integration with Terraform's debugging facilities.

## Impact

Using the standard library's log package can result in logs not appearing where provider users expect them (e.g., in Terraform logs), thus reducing observability and potentially leaking information or adding noise to test runs. **Severity: medium**

## Location

Lines containing:
```go
log.Default().Printf("Storing original value for attribute %s", req.PathExpression.String())
```
and similar log calls.

## Code Issue

```go
log.Default().Printf("Storing original value for attribute %s", req.PathExpression.String())
```

## Fix

Remove standard log statements or replace them with proper diagnostic messages via the Terraform logging context, or (preferably) remove them if they're not essential for user diagnostics. If you need to report diagnostics to the Terraform run, you should use the `resp.Diagnostics` or similar feature of the framework.

```go
// If logging is needed for debugging, use the framework's Diagnostics system.
// Otherwise, simply remove the log line:
if req.State.Raw.IsNull() {
    if (!req.ConfigValue.IsNull()) {
        // Optionally add a diagnostics message here if you wish to surface info to the user.
        resp.Private.SetKey(ctx, req.Path.String(), []byte{1})
    }
}
```
