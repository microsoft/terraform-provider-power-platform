# Title

`Improper diagnostic error handling with silent failures`

##

`/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go`

## Problem

In the `hasChecksumChanged` function, where attribute diagnostics are appended using `resp.Diagnostics.Append(diags...)`, there is no error check or handling for the case when these diagnostics contain issues.

## Impact

If any of the diagnostics added contain errors, they may go unnoticed, resulting in silent failures or incorrect behavior downstream. Severity: **medium**.

## Location

Lines in the `hasChecksumChanged` function:

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)

diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
```

## Code Issue

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)

diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
```

## Fix

Check if the diagnostics have errors and handle them accordingly. For example:

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return false // Or another appropriate action
}

diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return false // Or another appropriate action
}
```
