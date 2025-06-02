# Issue 1: Error Message Mentions Wrong Hash Algorithm ("MD5" instead of "SHA256")

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go

## Problem

In the `hasChecksumChanged` method, when the error message is generated for a checksum calculation failure, it states "Error calculating MD5 checksum" even though the code actually uses SHA256 via `helpers.CalculateSHA256`.

## Impact

This discrepancy can cause confusion during debugging and mislead engineers regarding the type of hash being calculated. The severity is low as it only affects log/error clarity, but should be corrected for accuracy.

## Location

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

## Code Issue

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

## Fix

Replace "MD5" with "SHA256" in the error string.

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", attribute), err.Error())
```

