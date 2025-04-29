# Title

Missing diagnostic error specificity for `attribute` and `attributeChecksum`

##

`/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_value_unknown_if_checksum_change_modifier.go`

## Problem

In the `hasChecksumChanged` method, if `attribute` or `attributeChecksum` values are missing, no specific diagnostic error is logged. This omission can lead to debugging challenges, as the underlying reason for checksum mismatch remains unclear.

## Impact

Medium severity. The lack of diagnostic errors for missing attribute values makes troubleshooting difficult, as there is insufficient detail about what caused the checksum calculation to fail.

## Location

Code in the `hasChecksumChanged` method:

```go
var attribute types.String
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)

var attributeChecksum types.String
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
```

## Code Issue

```go
var attribute types.String
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

Add diagnostic errors specifically to check if the attribute values are missing or invalid before proceeding with the checksum calculation:

```go
var attribute types.String
diags := req.Plan.GetAttribute(ctx, path.Root(attributeName), &attribute)
resp.Diagnostics.Append(diags...)
if !attribute.IsUnknown() && attribute.IsNull() {
	resp.Diagnostics.AddError("Missing attribute value for "+attributeName, fmt.Sprintf("The attribute %s is either unknown or null.", attributeName))
	return false
}

var attributeChecksum types.String
diags = req.State.GetAttribute(ctx, path.Root(checksumAttributeName), &attributeChecksum)
resp.Diagnostics.Append(diags...)
if !attributeChecksum.IsUnknown() && attributeChecksum.IsNull() {
	resp.Diagnostics.AddError("Missing checksum value for "+checksumAttributeName, fmt.Sprintf("The checksum attribute %s is either unknown or null.", checksumAttributeName))
	return false
}

value, err := helpers.CalculateSHA256(attribute.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", attributeName), err.Error())
}
```

This improvement provides clear error messages that help pinpoint the specific cause of checksum issues.