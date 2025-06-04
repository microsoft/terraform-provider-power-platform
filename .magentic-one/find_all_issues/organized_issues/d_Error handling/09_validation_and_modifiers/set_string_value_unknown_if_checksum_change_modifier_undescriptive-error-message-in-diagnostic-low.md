# Undescriptive Error Message in Diagnostic

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_value_unknown_if_checksum_change_modifier.go

## Problem

The error message for errors in the `hasChecksumChanged` method references MD5 ("Error calculating MD5 checksum for %s"), but the function used is `CalculateSHA256` (SHA256). This is misleading and can confuse users and maintainers.

## Impact

Misleading error messages can slow down debugging and maintenance, leading to confusion about what hashing function is being used. Severity: low.

## Location

Line in `hasChecksumChanged` method where the diagnostic error is added.

## Code Issue

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

## Fix

Correct the error message to reference SHA256 to match the actual implementation.

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", attributeName), err.Error())
```
