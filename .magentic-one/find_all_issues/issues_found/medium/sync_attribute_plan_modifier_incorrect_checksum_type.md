# Title

Incorrect Checksum Type in Error Message

##

/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go

## Problem

In the `PlanModifyString` method, the error message specifies "Error calculating MD5 checksum for...". However, the actual checksum calculation function being used is `helpers.CalculateSHA256`, which computes a SHA256 hash instead of an MD5 checksum. This discrepancy can lead to confusion and inaccurate debugging or logging.

## Impact

- **Severity**: Medium
- Incorrect communication in error logs can mislead developers or system analysts trying to debug issues involving checksum calculations.
- Errors related to checksum calculations might be overlooked or wrong measures applied, impacting system reliability.

## Location

Line: `resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", d.syncAttribute), err.Error())`

## Code Issue

```go
value, err := helpers.CalculateSHA256(settingsFile.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Error calculating MD5 checksum for %s", d.syncAttribute), err.Error())
	return
}
```

## Fix

Fix the error message to accurately reflect the checksum type being calculated.

```go
value, err := helpers.CalculateSHA256(settingsFile.ValueString())
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Error calculating SHA256 checksum for %s", d.syncAttribute), err.Error())
	return
}
}
```
