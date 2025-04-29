# Title

`Error message incorrectly refers to MD5 checksum instead of SHA256 checksum`

##

`/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go`

## Problem

In the error message located in the `hasChecksumChanged` function, the implementation uses `helpers.CalculateSHA256` for SHA256 checksum calculation, but the error message incorrectly states "Error calculating MD5 checksum".

## Impact

Misleading error message can confuse developers and troubleshooting efforts, leading to delays and misinterpretation of issues. Severity: **high**.

## Location

Line containing:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

## Code Issue

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", attribute), err.Error())
```

## Fix

Update the error message to correctly reference SHA256 checksum:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", attribute), err.Error())
```
