# Issue with Data Consistency When Building Rows

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

When building the `rows` attribute in the `Read` method, the code constructs a slice of elements from records and builds a `TupleValue` using a parallel slice of their types. However, if `elements` has a length of zero (no records returned), the code calls `types.TupleValue(elementsTypes, elements)`, which may not behave as expected. Additionally, error return values from `types.TupleValue` and similar constructors are ignored, which could result in subtle or silent failures, especially if the input data has inconsistencies.

## Impact

- If the row count is zero, an invalid or nil value may be stored (unclear how framework handles zero length).
- Any internal error in attribute/type construction will be silently ignored (since errors are discarded).
- This weakens data consistency guarantees for consumers. 
- Severity: **medium**.

## Location

In `Read`:
```go
elementTypes := []attr.Type{}
for range elements {
	elementTypes = append(elementTypes, types.DynamicType)
}
rows, _ := types.TupleValue(elementTypes, elements)
state.Rows = types.DynamicValue(rows)
```

## Code Issue

```go
rows, _ := types.TupleValue(elementTypes, elements)
state.Rows = types.DynamicValue(rows)
```

## Fix

- Always check the error returned by `types.TupleValue` and handle it (at a minimum, surface the error to diagnostics and return).
- If there are no elements, set the value to empty appropriately using the Terraform SDK type helpers.

```go
rows, err := types.TupleValue(elementTypes, elements)
if err != nil {
	resp.Diagnostics.AddError("Failed to build tuple value for rows", err.Error())
	return
}
state.Rows = types.DynamicValue(rows)
```
This ensures that any error when constructing the attribute is reported, and consumers won't receive silently broken state.
