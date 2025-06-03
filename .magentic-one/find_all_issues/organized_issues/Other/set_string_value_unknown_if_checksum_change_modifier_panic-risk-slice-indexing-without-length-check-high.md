# Panic Risk: Slice Indexing Without Length Check

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_value_unknown_if_checksum_change_modifier.go

## Problem

In the `PlanModifyString` method, the code directly indexes into `firstAttributePair[0]`, `firstAttributePair[1]`, `secondAttributePair[0]`, and `secondAttributePair[1]` without checking if these slices have the required length. This can lead to a runtime panic.

## Impact

If the modifier is constructed with slices shorter than 2 elements, the code will panic and possibly crash the program. Severity: high.

## Location

In `PlanModifyString`:

```go
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
...
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```

## Fix

Check the length of slices before accessing elements to avoid panics.

```go
if len(d.firstAttributePair) < 2 || len(d.secondAttributePair) < 2 {
    resp.Diagnostics.AddError("Invalid attribute pair length", "Each attribute pair must have at least two elements.")
    return
}
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```
