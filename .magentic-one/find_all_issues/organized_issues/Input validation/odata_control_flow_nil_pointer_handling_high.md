# Issue 1: Potential Control Flow and Nil Pointer Handling in Query Appending

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go

## Problem

The `appendQuery` function appends non-nil query parts to the OData query, but there is no handling for what would happen if the input pointer (`query`) itself is nil. Additionally, as this function directly operates on the string pointer, improper/misuse or accidentally providing a nil value could lead to panics. 

## Impact

If a nil pointer is passed as the base query to `appendQuery`, a runtime panic will be triggered (`invalid memory address or nil pointer dereference`). Severity is **high** as this is a potential runtime crash.

## Location

Line(s) 81-92

## Code Issue

```go
func appendQuery(query, part *string) {
	if part != nil {
		if len(*query) > 0 {
			*query += "&"
		}
		*query = strings.Join([]string{*query, *part}, "")
	}
}
```

## Fix

Add nil check for `query` pointer and consider returning an error or avoiding mutation if `query` is nil.

```go
func appendQuery(query, part *string) {
	if query == nil {
		// avoid panic and/or log error
		return
	}
	if part != nil {
		if len(*query) > 0 {
			*query += "&"
		}
		*query = strings.Join([]string{*query, *part}, "")
	}
}
```
