# Title

Duplicated Null-Check Logic for PlanValue Assignment

##

/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go

## Problem

The assignment logic for `resp.PlanValue` to `types.StringNull()` is duplicated (i.e., both when `settingsFile.IsNull()` and `settingsFile.IsUnknown()`). This could be simplified to reduce redundancy in the code.

## Impact

While not a severe logic bug, it increases maintenance overhead and slightly reduces code clarity. Severity: **low**.

## Location

Within `PlanModifyString`:

```go
if settingsFile.IsNull() {
	resp.PlanValue = types.StringNull()
} else if settingsFile.IsUnknown() {
	resp.PlanValue = types.StringNull()
} else {
	// ...
}
```

## Code Issue

```go
if settingsFile.IsNull() {
	resp.PlanValue = types.StringNull()
} else if settingsFile.IsUnknown() {
	resp.PlanValue = types.StringNull()
} else {
	// ...
}
```

## Fix

Combine the two checks using a logical OR to reduce duplicate code:

```go
if settingsFile.IsNull() || settingsFile.IsUnknown() {
	resp.PlanValue = types.StringNull()
} else {
	// ...
}
```
