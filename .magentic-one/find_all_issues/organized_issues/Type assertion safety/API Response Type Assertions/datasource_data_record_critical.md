# Issue with Error Handling on Dynamic Type Assertion

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go

## Problem

In multiple methods, especially in type assertion blocks (e.g., in `convertColumnsToState`, `buildObjectValueFromX`, `buildExpandObject`), the code is using direct type assertions like `columns[key].(bool)` without handling cases where the assertion might fail. If the assertion fails, a panic will be triggered. This could cause the program to crash unexpectedly if the data returned by the API isn't as expected.

## Impact

Crashing due to a failed type assertion is a critical issue for stability and error handling, as it makes the provider brittle in the face of unexpected data. Severity: **critical**.

## Location

Multiple methods, including:
- `convertColumnsToState`
- `buildObjectValueFromX`
- `buildExpandObject`

## Code Issue

```go
switch value.(type) {
case bool:
    caseBool(columns[key].(bool), attributes, attributeTypes, key)
...
// repeated in all functions with type assertions
```

## Fix

You should use the "comma ok" idiom on type assertion to avoid a panic and handle the error gracefully. For example:

```go
switch v := value.(type) {
case bool:
    caseBool(v, attributes, attributeTypes, key)
case int64:
    caseInt64(v, attributes, attributeTypes, key)
case float64:
    caseFloat64(v, attributes, attributeTypes, key)
case string:
    caseString(v, attributes, attributeTypes, key)
case map[string]any:
    typ, val, _ := d.buildObjectValueFromX(v)
    tupleElementType := types.ObjectType{
        AttrTypes: typ,
    }
    objVal, _ := types.ObjectValue(typ, val)
    attributes[key] = objVal
    attributeTypes[key] = tupleElementType
case []any:
    typeObj, valObj := d.buildExpandObject(v)
    attributeTypes[key] = typeObj
    attributes[key] = valObj
default:
    // optionally log or handle unexpected type
}
```
Add similar handling in other relevant methods. This will prevent panics due to unexpected types and let you report or skip unhandled data types without crashing the provider.
