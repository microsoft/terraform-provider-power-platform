# Title

Missing error handling for `req.Config.GetAttribute`

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

In the `Validate` method, the function `req.Config.GetAttribute` is used to get the value of `currentFieldValue` and `otherFieldValue`. However, the code does not check for errors returned by this function when fetching `currentFieldValue`, and it uses `_` to discard the error. Although the error handling exists for `otherFieldValue`, it does not address cases where the first call might fail.

## Impact

If `req.Config.GetAttribute` fails, the value of `currentFieldValue` may be invalid or uninitialized. This can lead to silent incorrect behavior since errors during the fetch process are ignored. This issue has a **high** severity because it can lead to unexpected results, especially in complex configurations.

## Location

Within the `Validate` method:

```go
currentFieldValue := ""
_ = req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)
```

## Code Issue

Snippet with the problematic code:

```go
currentFieldValue := ""
_ = req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)
```

The error returned by `req.Config.GetAttribute` is ignored completely.

## Fix

Proper error handling must be introduced to capture and deal with errors when the `GetAttribute` call fails. Here is the corrected code:

```go
currentFieldValue := ""
if err := req.Config.GetAttribute(ctx, req.Path, &currentFieldValue); err != nil {
    res.Diagnostics.AddError("Error fetching current field value", err.Error())
    return
}
```