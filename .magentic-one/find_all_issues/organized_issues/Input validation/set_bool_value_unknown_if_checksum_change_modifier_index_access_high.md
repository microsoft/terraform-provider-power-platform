# Use of Unchecked Index Access for Attribute Pairs

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go

## Problem

The code accesses the first and second elements of `firstAttributePair` and `secondAttributePair` via `[0]` and `[1]` indexing in `PlanModifyBool`, assuming each pair always has at least two elements. There is no validation or safeguard if fewer elements are provided.

## Impact

If the input slice does not have at least two elements, this will cause a runtime panic due to index out-of-range errors. This is a high-severity issue since it can crash the program unexpectedly.

## Location

```go
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
...
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```

## Code Issue

```go
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```

## Fix

Add validation in the constructor (or before usage) to ensure that all required elements are present in each slice before using them. Example:

```go
func SetBoolValueToUnknownIfChecksumsChangeModifier(firstAttributePair, secondAttributePair []string) planmodifier.Bool {
    if len(firstAttributePair) < 2 || len(secondAttributePair) < 2 {
        panic("Each attribute pair must have at least two elements: attribute name and checksum attribute name")
    }
    return &setBoolValueToUnknownIfChecksumsChangeModifier{
        firstAttributePair:  firstAttributePair,
        secondAttributePair: secondAttributePair,
    }
}
```

Alternatively, add error handling logic instead of panic depending on application requirements.
