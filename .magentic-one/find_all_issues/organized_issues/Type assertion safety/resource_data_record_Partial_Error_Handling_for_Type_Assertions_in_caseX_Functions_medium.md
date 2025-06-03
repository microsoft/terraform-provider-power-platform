# Partial Error Handling for Type Assertions in caseX Functions

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go

## Problem

The helper functions `caseMapStringOfAny`, `caseArrayOfAny`, etc., only populate values if the type assertion succeeds but do not log or report if assertions fail. The calling code has no idea if the value could not be set. This may result in silent data loss or state drift.

## Impact

- **Severity:** Medium
- Can lead to silent errors or ignored/empty attributes.
- Reduces code robustness and makes debugging more difficult.

## Location

Functions: `caseMapStringOfAny`, `caseArrayOfAny`, etc.

## Code Issue

```go
value, ok := columnValue.(string)
if ok {
	// ...
}
```

## Fix

Consider logging or returning an error (if critical), or at minimum, logging via tflog for unexpected value types:

```go
value, ok := columnValue.(string)
if !ok {
    tflog.Warn(context.TODO(), "caseMapStringOfAny: failed to cast value to string", map[string]interface{}{ "key": key })
    return // or capture error for diagnostic
}
// ... (continue existing logic)
```
