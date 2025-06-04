# Title

Incorrect Error Message Refers to MD5 Instead of SHA256

##

/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go

## Problem

The error message in the error handling references "MD5 checksum" while the function actually computes a SHA256 hash using `helpers.CalculateSHA256()`. This introduces confusion and misleads maintainers or users as to which hashing algorithm is being used.

## Impact

This can cause confusion during debugging or auditing, potentially leading to diagnostic errors or incorrect assumptions about security. Severity: **low**.

## Location

Line where the error message is formed in `PlanModifyString`:
```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", d.syncAttribute), err.Error())
```

## Code Issue

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", d.syncAttribute), err.Error())
```

## Fix

Correct the error message to reference SHA256 instead of MD5:

```go
resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", d.syncAttribute), err.Error())
```
