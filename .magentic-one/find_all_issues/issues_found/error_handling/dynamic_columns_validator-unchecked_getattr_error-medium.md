# Issue 2: Unchecked Error from config.GetAttribute

## 

/workspaces/terraform-provider-power-platform/internal/services/data_record/dynamic_columns_validator.go

## Problem

The returned error from `config.GetAttribute` is ignored (assigned to `_`). If `GetAttribute` fails, the subsequent code may operate on invalid or undefined data.

## Impact

Can potentially lead to data inconsistency or mask actual configuration errors the user should be aware of. Severity: **medium**.

## Location

Within the `Validate` function:

## Code Issue

```go
_ = config.GetAttribute(ctx, matchedPaths[0], &dynamicColumns)
```

## Fix

Check the error and add to diagnostics if present:

```go
if err := config.GetAttribute(ctx, matchedPaths[0], &dynamicColumns); err != nil {
    diags.AddError("Failed to get dynamic columns attribute", err.Error())
    return diags
}
```
