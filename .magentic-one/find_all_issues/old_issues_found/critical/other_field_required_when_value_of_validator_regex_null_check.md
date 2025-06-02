# Title

Improper regex null checks

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

Regular expressions (`CurrentFieldValueRegex` and `OtherFieldValueRegex`) may be nil based on current construction logic, and although the code follows null-check conditions in some cases, there are locations where nil access could lead to runtime panics. For example, `av.OtherFieldValueRegex.MatchString(*otherFieldValue)` doesn't explicitly check if `OtherFieldValueRegex` is nil.

## Impact

Potential runtime panics which can crash the application in production-failure cases. This is a **critical** severity issue because the application can abruptly terminate if `OtherFieldValueRegex` is accessed while nil.

## Location

Within the `Validate` method:

```go
if av.OtherFieldValueRegex != nil && !av.OtherFieldValueRegex.MatchString(*otherFieldValue) || av.OtherFieldValueRegex == nil && (otherFieldValue == nil || *otherFieldValue == "") && !isUnknown {
```

## Code Issue

Snippet with the problematic code:

```go
if av.OtherFieldValueRegex != nil && !av.OtherFieldValueRegex.MatchString(*otherFieldValue) || av.OtherFieldValueRegex == nil && (otherFieldValue == nil || *otherFieldValue == "") && !isUnknown {
```

## Fix

Ensure comprehensive nil checks before using regular expressions. Rewrite the condition as follows:

```go
if av.OtherFieldValueRegex != nil {
    if !av.OtherFieldValueRegex.MatchString(*otherFieldValue) {
        res.Diagnostics.AddError(av.ErrorMessage, av.ErrorMessage)
    }
} else if (otherFieldValue == nil || *otherFieldValue == "") && !isUnknown {
    res.Diagnostics.AddError(av.ErrorMessage, av.ErrorMessage)
}
```

This ensures no nil dereference, simplifying readability and improving safety.