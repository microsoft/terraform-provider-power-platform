# Title

`Hardcoded index access in slice may cause runtime panic`

##

`/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go`

## Problem

The code directly accesses specific indices of slices `d.firstAttributePair` and `d.secondAttributePair` assuming their length is at least 2. There is no validation to ensure that these slices contain enough elements before accessing these indices.

## Impact

If the slices have fewer than the required elements (e.g., due to incorrect or missing input data), this will cause a runtime panic, which can crash the application. Severity: **high**.

## Location

Lines in the `PlanModifyBool` function:

```go
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```

## Code Issue

```go
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```

## Fix

Validate the slice lengths before accessing specific indices. For example:

```go
if len(d.firstAttributePair) < 2 || len(d.secondAttributePair) < 2 {
    resp.Diagnostics.AddError(
        "Missing attributes for checksum validation",
        "Expected at least two attributes in each attribute pair, but got fewer.",
    )
    return
}

firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```
