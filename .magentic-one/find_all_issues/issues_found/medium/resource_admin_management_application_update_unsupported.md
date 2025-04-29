# Title

Unsupported operation logging in `Update` method

## Path to the file

`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go`

## Problem

The `Update` method of the resource adds a diagnostic error stating "Update not supported". However, it does not provide detailed reasoning or any logs to indicate why this operation is unsupported, nor does it investigate potential alternative behaviors or recovery steps.

## Impact

While this limitation is directly surfaced to users, its lack of descriptive context or investigation into alternatives reduces clarity and usability of the provider. The unexplained limitation might confuse users, resulting in unnecessary support tickets or missed operational insights. Severity is **medium** as the visible lack of explanation impacts clarity and user satisfaction but does not block any functioning feature.

## Location

```go
resp.Diagnostics.AddError("Update not supported", "Update not supported")
```

## Fix

Provide additional context and reasoning in the diagnostic error, with optional links to documentation or suggestions regarding alternative actions.

### Suggested Fix:

```go
resp.Diagnostics.AddError(
    "Update not supported",
    "Updating an admin management application is not currently supported due to limitations in the underlying API." +
    " If you require changes, please delete and recreate the resource. Learn more: [Relevant-Link-Here]",
)
tflog.Warn(ctx, "Update operation was invoked but is unsupported for admin management applications.")
```
