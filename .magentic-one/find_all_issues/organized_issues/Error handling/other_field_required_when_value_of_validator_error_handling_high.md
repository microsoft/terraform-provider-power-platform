# Issue: Use of error return values without error handling

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

There are multiple instances in the code where functions that return an error (or a diagnostic error) are used, but their error return values are ignored or not handled properly. For example, the calls to `GetAttribute` and `PathMatches` do not check for errors, which can lead to unnoticed issues, panics, or confusion during debugging.

## Impact

If errors are not handled properly, this could lead to faulty logic, program panics, or incorrect validation behavior. Severity: **high**.

## Location

Relevant lines include (but are not limited to):

```go
	currentFieldValue := ""
	_ = req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)
	paths, _ := req.Config.PathMatches(ctx, av.OtherFieldExpression)
```

## Code Issue

```go
	currentFieldValue := ""
	_ = req.Config.GetAttribute(ctx, req.Path, &currentFieldValue)

	// ...
	paths, _ := req.Config.PathMatches(ctx, av.OtherFieldExpression)
```

## Fix

Check and handle the errors properly, such as:

```go
	currentFieldValue := ""
	if err := req.Config.GetAttribute(ctx, req.Path, &currentFieldValue); err != nil {
		res.Diagnostics.AddError(
			"Failed to get attribute for current field",
			err.Error(),
		)
		return
	}

	paths, err := req.Config.PathMatches(ctx, av.OtherFieldExpression)
	if err != nil {
		res.Diagnostics.AddError(
			"Failed to match path expression for other field",
			err.Error(),
		)
		return
	}
```

---

This issue impacts error handling and should be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/other_field_required_when_value_of_validator_error_handling_high.md`
