# Diagnostics May Contain Errors But Processing Continues

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go

## Problem

In `hasChecksumChanged`, after appending diagnostics from `GetAttribute`, there is no check to see if an error occurred before proceeding with the rest of the logic. If `GetAttribute` fails, `attribute` or `attributeChecksum` might not be populated with valid data, which could lead to incorrect calculations or misleading diagnostics.

## Impact

This can result in attempts to calculate a SHA256 value or compare checksums using invalid or empty data, making error detection and debugging harder. The severity is medium, as it may lead to non-obvious erroneous state propagating downstream.

## Location

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
...
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
...
value, err := helpers.CalculateSHA256(attribute.ValueString())
```

## Code Issue

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)

var attributeChecksum types.String
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)

value, err := helpers.CalculateSHA256(attribute.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
}
```

## Fix

After each `Append` of diagnostics, check if there was an error and halt further processing if needed:

```go
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
if diags.HasError() {
    return false
}

var attributeChecksum types.String
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
if diags.HasError() {
    return false
}
```
