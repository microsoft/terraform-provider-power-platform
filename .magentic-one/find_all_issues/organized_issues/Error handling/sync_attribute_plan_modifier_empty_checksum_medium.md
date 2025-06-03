# Title

Unnecessary Separate Handling for Empty Checksum Value

##

/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go

## Problem

After calculating the checksum, the code checks `if value == ""` and sets `resp.PlanValue = types.StringUnknown()`. Normally, a checksum calculation failure (which could result in an empty string) should already be handled by the error check above; reaching this branch is ambiguous and may mask edge cases.

## Impact

This check may hide potential errors in checksum generation or underlying issues, making it harder to diagnose problems. It also doesn't document why an empty hash should be interpreted as “unknown.” Severity: **medium**.

## Location

Within `PlanModifyString`:

```go
if value == "" {
	resp.PlanValue = types.StringUnknown()
} else {
	resp.PlanValue = types.StringValue(value)
}
```

## Code Issue

```go
if value == "" {
	resp.PlanValue = types.StringUnknown()
} else {
	resp.PlanValue = types.StringValue(value)
}
```

## Fix

Document explicitly why an empty hash might occur (if it's a valid state), or treat this as an error/diagnostic. Otherwise, remove this block so that unexpected empty values do not silently propagate:

```go
if value == "" {
	resp.Diagnostics.AddError(fmt.Sprintf("Checksum is empty for %s", d.syncAttribute), "Calculated SHA256 checksum resulted in an empty value, which is unexpected.")
	resp.PlanValue = types.StringUnknown()
} else {
	resp.PlanValue = types.StringValue(value)
}
```

Or, if empty checksums are not possible, simply assign the value without this check.
