# Issue Title

Unnecessary Type Assertion in Loop

##

`/workspaces/terraform-provider-power-platform/internal/helpers/set_to_string.go`

## Problem

The code contains a type assertion within the loop (`if str, ok := v.(types.String)`) to check whether `v` is of type `types.String`. However, based on the context, the input `set` should already only contain elements of type `types.String`. If the input contains mixed or unexpected types, this may lead to unpredictable behavior.

## Impact

This introduces unnecessary complexity and potential inefficiencies in the code, as every iteration involves a type assertion. The impact is **Medium** because this inefficiency can grow as the size of the set increases, and improper type assertions can lead to runtime errors if elements of unexpected types are encountered.

## Location

Function: `SetToStringSlice`

## Code Issue

```go
func SetToStringSlice(set types.Set) []string {
	var result []string
	for _, v := range set.Elements() {
		if str, ok := v.(types.String); ok {
			result = append(result, str.ValueString())
		}
	}
	return result
}
```

## Fix

To enhance code efficiency and prevent potential runtime issues, validate that the input `set` contains only `types.String` elements prior to the loop, or pass a set with a defined generic type. The following modification directly converts the elements without type assertions:

```go
func SetToStringSlice(set types.Set) ([]string, error) {
	var result []string
	for _, v := range set.Elements() {
		str, ok := v.(types.String)
		if !ok {
			return nil, fmt.Errorf("unexpected type in the set: expected types.String, got %T", v)
		}
		result = append(result, str.ValueString())
	}
	return result, nil
}
```

Alternatively, refactor the logic to avoid reliance on dynamic type checking during processing. Ensure the types of elements in the set are validated prior to calling this function.