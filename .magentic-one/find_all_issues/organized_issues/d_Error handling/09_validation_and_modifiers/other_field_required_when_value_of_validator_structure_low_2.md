# Duplicate error message in AddError usage

##

/workspaces/terraform-provider-power-platform/internal/validators/other_field_required_when_value_of_validator.go

## Problem

In the line where diagnostics are added for an invalid "other field" value, the `ErrorMessage` is used as both the summary and the detail, causing redundancy. Ideally, the detail section should provide extra context or actionable information, not just reiterate the summary.

## Impact

This reduces clarity and effectiveness of error diagnostics, especially for users and maintainers debugging complex configurations. Severity: **low**.

## Location

```go
if (av.OtherFieldValueRegex != nil && otherFieldValue != nil && !av.OtherFieldValueRegex.MatchString(*otherFieldValue)) ||
	(av.OtherFieldValueRegex == nil && (otherFieldValue == nil || *otherFieldValue == "") && !isUnknown) {
	res.Diagnostics.AddError(av.ErrorMessage, av.ErrorMessage)
}
```

## Code Issue

```go
	res.Diagnostics.AddError(av.ErrorMessage, av.ErrorMessage)
```

## Fix

Change to something like:

```go
	res.Diagnostics.AddError(
		av.ErrorMessage,
		"Field \""+paths[0].String()+"\" does not meet required value conditions.",
	)
```
Or provide an actionable/contextual detail:

```go
	res.Diagnostics.AddError(
		av.ErrorMessage,
		"The other required field did not match the expected value or was empty.",
	)
```

---

This issue impacts maintainability/readability and should be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/other_field_required_when_value_of_validator_structure_low_2.md`
