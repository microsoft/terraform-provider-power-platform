# Inconsistent Error Handling in convertColumnsToState

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go

## Problem

The `convertColumnsToState` function has some places where error values are deliberately ignored using `_` or omitted via no checks. This can hide bugs or introduce hard-to-debug issues.

## Impact

- **Severity:** High
- Can cause subtle or silent state bugs.
- Undermines error visibility and diagnosis.

## Location

`convertColumnsToState`, especially in:

```go
columnField, _ := types.ObjectValue(attributeTypes, attributes)
```

## Code Issue

```go
columnField, _ := types.ObjectValue(attributeTypes, attributes)
result := types.DynamicValue(columnField)
return &result, nil
```

## Fix

Check for the error and propagate or report it via diagnostics:

```go
columnField, err := types.ObjectValue(attributeTypes, attributes)
if err != nil {
    return nil, fmt.Errorf("failed to create object value: %w", err)
}
result := types.DynamicValue(columnField)
return &result, nil
```
