# Lack of Defensive Nil Check for Set Input

##

/workspaces/terraform-provider-power-platform/internal/helpers/set_to_string.go

## Problem

The function does not check for a nil or unknown set input. If `set` is `types.Set{Unknown: true}` or `nil`, the function will still attempt to operate on it, potentially leading to panics or undefined behavior.

## Impact

This can cause runtime panics and can be a **high severity** reliability issue, particularly when working with Terraform resource data, which can have unknown or nil states.

## Location

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

Add a check for set's zero value, known/unknown state, or nil (depending on the Terraform Plugin Framework's type safety).

```go
// Defensive handling: return empty slice if set is unknown or nil.
func SetOfStringsToSlice(set types.Set) []string {
    if !set.IsKnown() || set.IsNull() {
        return []string{}
    }
    var result []string
    for _, v := range set.Elements() {
        if str, ok := v.(types.String); ok {
            result = append(result, str.ValueString())
        }
    }
    return result
}
```

Use `IsKnown` and `IsNull` methods as provided by the SDK.
