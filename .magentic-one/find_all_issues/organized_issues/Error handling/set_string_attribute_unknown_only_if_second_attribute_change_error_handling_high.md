# Unhandled Error Diagnostics in PlanModifyString

##

/workspaces/terraform-provider-power-platform/internal/modifiers/set_string_attribute_unknown_only_if_second_attribute_change.go

## Problem

Error diagnostics returned from `GetAttribute` calls in the `PlanModifyString` method are appended to `resp.Diagnostics`, but the code does not check or handle whether errors exist before proceeding to work with the potentially invalid data.

## Impact

- **Error Handling**: The function may execute further logic using invalid or uninitialized variables if errors exist in diagnostics, leading to unreliable behavior.
- **Severity**: High

## Location

```go
	var planSecondAttribute types.String
	diags := req.Plan.GetAttribute(ctx, d.secondAttributePath, &planSecondAttribute)
	resp.Diagnostics.Append(diags...)

	var stateSecondAttribute types.String
	diags = req.State.GetAttribute(ctx, d.secondAttributePath, &stateSecondAttribute)
	resp.Diagnostics.Append(diags...)

	if planSecondAttribute.ValueString() != stateSecondAttribute.ValueString() && !planSecondAttribute.IsUnknown() && !planSecondAttribute.IsNull() {
		resp.PlanValue = types.StringUnknown()
	}
```

## Fix

After appending diagnostics, check if any errors are present and return early if so to prevent use of potentially invalid data.

```go
	var planSecondAttribute types.String
	diags := req.Plan.GetAttribute(ctx, d.secondAttributePath, &planSecondAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateSecondAttribute types.String
	diags = req.State.GetAttribute(ctx, d.secondAttributePath, &stateSecondAttribute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if planSecondAttribute.ValueString() != stateSecondAttribute.ValueString() && !planSecondAttribute.IsUnknown() && !planSecondAttribute.IsNull() {
		resp.PlanValue = types.StringUnknown()
	}
```
