# Title

Incorrect error message in checksum calculation

##

`/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_value_unknown_if_checksum_change_modifier.go`

## Problem

The error message indicates an MD5 checksum calculation, but the code is clearly using SHA256 for the checksum. This causes confusion and may mislead developers or users debugging the issue.

## Impact

Medium severity. While it doesn't affect the logic of the code, incorrect error messages can lead to debugging inefficiencies and miscommunication.

## Location

Line in the `hasChecksumChanged` method:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())

```

## Code Issue

```go
value, err := helpers.CalculateSHA256(attribute.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
}
```

## Fix

Update the error message to reflect the correct checksum type being calculated:

```go
value, err := helpers.CalculateSHA256(attribute.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", attribute), err.Error())
}
```

This ensures the error message is consistent with the actual functionality of the code.