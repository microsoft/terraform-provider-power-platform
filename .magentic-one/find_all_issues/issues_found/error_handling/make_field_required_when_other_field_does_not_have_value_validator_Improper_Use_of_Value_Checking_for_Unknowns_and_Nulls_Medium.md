# Improper Use of Value Checking for Unknowns and Nulls

##

/workspaces/terraform-provider-power-platform/internal/validators/make_field_required_when_other_field_does_not_have_value_validator.go

## Problem

There is an overlap and potential confusion between handling `IsUnknown()` and `IsNull()` for configuration values. In Terraform's plugin SDK, `IsUnknown()` and `IsNull()` have distinct semantic differences. The code currently checks both for unknown or null values in different places, and in one instance uses diagnostics with specific string matching to infer the unknown state, which can be brittle and hard to maintain.

## Impact

Medium severity. This can result in validation logic not behaving as intended, particularly if the SDK changes its error messages, or if inconsistent treatments of value states cause validators to fail under valid scenarios.

## Location

```go
if req.ConfigValue.IsUnknown() || isUnknown {
	tflog.Debug(ctx, "RequiredFieldValidator: Value is unknown, skipping validation")
	return
}
```

## Code Issue

```go
if req.ConfigValue.IsUnknown() || isUnknown {
	tflog.Debug(ctx, "RequiredFieldValidator: Value is unknown, skipping validation")
	return
}
```

## Fix

Unify and clarify the checking for `IsUnknown` and `IsNull`. Move explicit error message reliance to correct API method use and only employ diagnostics for error reporting, not flow control where possible.

```go
if req.ConfigValue.IsUnknown() || isUnknown {
	tflog.Debug(ctx, "RequiredFieldValidator: Value is unknown, skipping validation")
	return
}
if req.ConfigValue.IsNull() {
	// handle the null case if needed
}
```
