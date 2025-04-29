# Title

Insufficient error message in diagnostics

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

In the case where `req.Config.PathMatches` finds zero or more than one match, the diagnostics error simply states that *"Other field required when value of validator should have exactly one match"*. However, this error is too generic, and users may find it difficult to discern how this error arose and how to resolve it.

## Impact

A vague error message makes debugging difficult for end users, which can lead to frustration and ineffective problem resolution. This is a **medium** severity issue because it primarily affects usability but does not influence the correctness of the validator logic.

## Location

Within the `Validate` method:

```go
if paths == nil && len(paths) != 1 {
    res.Diagnostics.AddError("Other field required when value of validator should have exactly one match", "")
    return
}
```

## Code Issue

Snippet with the problematic code:

```go
res.Diagnostics.AddError("Other field required when value of validator should have exactly one match", "")
```

The second argument to `AddError`, which is designed to provide additional details, is left empty.

## Fix

Provide a more descriptive error message that includes what went wrong and what was expected. Here’s the corrected suggestion:

```go
if paths == nil && len(paths) != 1 {
    res.Diagnostics.AddError(
        "Invalid configuration: Exactly one match required",
        "The configuration specified a condition that requires exactly one match for the other field path expression, but either none or multiple matches were found. Verify the path expression syntax and ensure only one match exists."
    )
    return
}
```