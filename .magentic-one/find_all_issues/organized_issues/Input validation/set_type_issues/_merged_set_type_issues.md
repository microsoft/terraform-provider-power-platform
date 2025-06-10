# Set Type Issues - Input Validation

This document contains all set type-related input validation issues found in the terraform-provider-power-platform codebase.


## ISSUE 1

# Use of Unchecked Index Access for Attribute Pairs

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_bool_value_unknown_if_checksum_change_modifier.go

## Problem

The code accesses the first and second elements of `firstAttributePair` and `secondAttributePair` via `[0]` and `[1]` indexing in `PlanModifyBool`, assuming each pair always has at least two elements. There is no validation or safeguard if fewer elements are provided.

## Impact

If the input slice does not have at least two elements, this will cause a runtime panic due to index out-of-range errors. This is a high-severity issue since it can crash the program unexpectedly.

## Location

```go
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
...
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```

## Code Issue

```go
firstAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.firstAttributePair[0], d.firstAttributePair[1])
secondAttributeHasChanged := d.hasChecksumChanged(ctx, req, resp, d.secondAttributePair[0], d.secondAttributePair[1])
```

## Fix

Add validation in the constructor (or before usage) to ensure that all required elements are present in each slice before using them. Example:

```go
func SetBoolValueToUnknownIfChecksumsChangeModifier(firstAttributePair, secondAttributePair []string) planmodifier.Bool {
    if len(firstAttributePair) < 2 || len(secondAttributePair) < 2 {
        panic("Each attribute pair must have at least two elements: attribute name and checksum attribute name")
    }
    return &setBoolValueToUnknownIfChecksumsChangeModifier{
        firstAttributePair:  firstAttributePair,
        secondAttributePair: secondAttributePair,
    }
}
```

Alternatively, add error handling logic instead of panic depending on application requirements.


## ISSUE 2

# Missing Type Safety: No Error or Type Validation Handling

##

/workspaces/terraform-provider-power-platform/internal/helpers/set_to_string.go

## Problem

The function quietly ignores any set element that doesn't successfully type-assert to `types.String`, which may hide bugs or input validation problems upstream. This can result in silent data loss if the set contains non-string elements, which can be hard to detect and debug.

## Impact

Silent type-assertion failures can lead to incomplete data transformation, subtle bugs, and reduce reliability of the helper, with a **high** severity in codebases where strict type handling is required.

## Location

```go
for _, v := range set.Elements() {
    if str, ok := v.(types.String); ok {
        result = append(result, str.ValueString())
    }
}
```

## Code Issue

```go
for _, v := range set.Elements() {
    if str, ok := v.(types.String); ok {
        result = append(result, str.ValueString())
    }
}
```

## Fix

Fail fast by either returning an error when an unexpected element is encountered or, if this function must have the same signature, at least document this behavior clearly in the function comment. Preferably, change the signature to return an error:

```go
// SetOfStringsToSlice converts a types.Set of types.String values to a slice of Go strings.
// Returns an error if a non-string element is encountered.
func SetOfStringsToSlice(set types.Set) ([]string, error) {
    var result []string
    for _, v := range set.Elements() {
        str, ok := v.(types.String)
        if !ok {
            return nil, fmt.Errorf("set contains non-string element: %v", v)
        }
        result = append(result, str.ValueString())
    }
    return result, nil
}
```

And update callers to handle the error.


## ISSUE 3

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


## ISSUE 4

# Inefficient Pre-Allocation When Handling Empty Slices

##

/workspaces/terraform-provider-power-platform/internal/helpers/string_to_set.go

## Problem

The function initializes the `values` slice with length equal to `len(slice)`, which is efficient in most cases. However, when `slice` is empty, it still allocates zero-length slices and proceeds to call `types.SetValue`, which may not be necessary if the Set type can be initialized empty more directly.

Additionally, the function does not explicitly handle the case where `slice` is nil. While Go generally handles `nil` slices gracefully, it may improve clarity and explicitness to handle the nil input directly.

## Impact

Severity: **Low**

This has a minimal direct impact on functionality, but handling the empty or nil slices more explicitly could make code more robust and increase maintainability.

## Location

```go
values := make([]attr.Value, len(slice))
for i, v := range slice {
    values[i] = types.StringValue(v)
}
set, diags := types.SetValue(types.StringType, values)
```

## Code Issue

```go
values := make([]attr.Value, len(slice))
for i, v := range slice {
    values[i] = types.StringValue(v)
}
set, diags := types.SetValue(types.StringType, values)
```

## Fix

Explicitly handle the empty or nil slice input case for improved clarity:

```go
if len(slice) == 0 {
    return types.SetValue(types.StringType, []attr.Value{})
}
values := make([]attr.Value, len(slice))
for i, v := range slice {
    values[i] = types.StringValue(v)
}
set, diags := types.SetValue(types.StringType, values)
if diags.HasError() {
    // error handling as above
}
return set, nil
```


# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
