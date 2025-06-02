# Title

Lack of helper-based field extraction or resource model validation causes repetitive code

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

In all CRUD methods, code that extracts values from the plan or state is repeated, and type/semantic validation and error reporting is manual. This creates repetitive, verbose, and less maintainable code, and makes it easy for developers to miss validations or introduce subtle bugs.

## Impact

Severity: **Low**

This is a maintainability and readability problem. Without a field validation or extraction helper or resource model validator, the codebase is harder to audit for correctness and harder to update/typesafe in the future.

## Location

- `Create`, `Read`, `Update`, `Delete` methods.

## Code Issue

```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*plan)
if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
	return
}
```

## Fix

Centralize validation and extraction logic for resource models. For example, add a `Validate()` method to `ResourceModel`, or write a resource model validation helper:

```go
func (m *ResourceModel) Validate() error {
	if m.BotId.IsUnknown() || m.BotId.ValueString() == "" {
		return fmt.Errorf("bot_id is required")
	}
	// Check other fields as necessary
	return nil
}

// Usage:
if err := plan.Validate(); err != nil {
	resp.Diagnostics.AddError("Invalid resource model", err.Error())
	return
}
```

This reduces code repetition, increases clarity, and ensures that all validation paths are handled consistently, simplifying later refactorings.

---

File to save:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_copilot_studio_application_insights_model_validation_low.md`
