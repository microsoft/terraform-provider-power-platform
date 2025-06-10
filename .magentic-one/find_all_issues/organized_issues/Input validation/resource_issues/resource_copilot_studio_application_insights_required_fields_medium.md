# Title

No validation of required resource attributes before use

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

Throughout the CRUD methods, the plan and state models are accessed via pointer dereference (e.g., `plan.BotId.ValueString()` and `state.BotId.ValueString()`), but there is no explicit validation that these required fields are non-empty or valid before use, especially just after fetching from the plan or state.

If an invalid Terraform configuration, a provider bug, or a state migration issue leads to these being empty, subsequent downstream client calls could fail in an uncontrolled way.

## Impact

Severity: **Medium**

Lack of validation can lead to cryptic API errors or panics that do not give a user-friendly error message in Diagnostics. A more robust approach that validates all required attributes before invoking an API improves maintainability, user experience, and robustness.

## Location

- Methods: `Create`, `Read`, `Update`, `Delete`
- Usage: Immediately after plan/state is loaded and before API client methods

## Code Issue

```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

// No checks here!
// Directly using plan.BotId.ValueString()
```

## Fix

Explicitly check all required fields in the plan (or state) after loading and before use, and return a helpful error if they are not set.

```go
if plan.BotId.IsUnknown() || plan.BotId.ValueString() == "" {
	resp.Diagnostics.AddError("Missing Bot ID", "The bot_id field is required but was not provided.")
	return
}
if plan.EnvironmentId.IsUnknown() || plan.EnvironmentId.ValueString() == "" {
	resp.Diagnostics.AddError("Missing Environment ID", "The environment_id field is required but was not provided.")
	return
}
// Repeat as needed for other critical fields
```

Repeat for the `state` variable in Read and Delete, and for all critical variables wherever appropriate.

---

File to save:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_copilot_studio_application_insights_required_fields_medium.md`
