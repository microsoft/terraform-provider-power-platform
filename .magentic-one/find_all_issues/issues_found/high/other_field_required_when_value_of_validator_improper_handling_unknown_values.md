# Title

Improper handling of unknown values

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

The code determines whether values are unknown by matching error details or summaries with specific strings, such as "Received unknown value". If the error messages change or are localized, this brittle logic may fail.

## Impact

Using string-based checks for unknown values makes the code fragile and dependent on implementation details of the `terraform-plugin-framework`. This issue has **high** severity due to the risk of breaking functionality if the library's error messages change.

## Location

Within the `Validate` method:

```go
isUnknown := false
if d.HasError() {
    isUnknown = strings.Contains(d.Errors()[0].Detail(), "Received unknown value") || strings.Contains(d.Errors()[0].Summary(), "Received unknown value")
}
```

## Code Issue

Snippet with the problematic code:

```go
isUnknown := false
if d.HasError() {
    isUnknown = strings.Contains(d.Errors()[0].Detail(), "Received unknown value") || strings.Contains(d.Errors()[0].Summary(), "Received unknown value")
}
```

This approach is error-prone and difficult to maintain.

## Fix

Introduce a dedicated method or library-supported functionality to check for unknown values. If such a method is not available, clearly document assumptions and encapsulate the logic in reusable functions.

Example fix:

```go
isUnknown := false
if d.HasError() {
    for _, diagError := range d.Errors() {
        if diagError.IsOfType(diag.ErrorTypeUnknownValue) { // Assuming framework provides error type constants
            isUnknown = true
            break
        }
    }
}
```

Alternatively, use a helper function to encapsulate string-based checks, making it easier to update and replace globally if the framework changes.