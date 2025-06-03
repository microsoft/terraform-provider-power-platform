# Issue: Unhandled Error After Failure in Create/Update/Conversion

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

In both the `Create`, `Update` and `Delete` methods, after an error is appended to the diagnostics due to a failed operation (like conversion using `createAppInsightsConfigDtoFromSourceModel`), the code does not return immediately. The subsequent code then runs with potentially invalid or nil variables. This can cause further unintended behavior or panics.

## Impact

Severity: **High**

Running code after an error in input conversion can cause runtime panics or unintended behaviors, such as dereferencing nil pointers, performing operations with incomplete data, or overwriting diagnostics with even more confusing errors. In production code, this could lead to crashes or unpredictable API responses.

## Location

- `Create` function, after error from `createAppInsightsConfigDtoFromSourceModel`.
- `Update` function, after error from `createAppInsightsConfigDtoFromSourceModel`.
- `Delete` function, after error from `createAppInsightsConfigDtoFromSourceModel`.

## Code Issue

```go
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*plan)
if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
}
// No return here, code continues after error

// ...
appInsightsConfigDto, err := r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate, plan.BotId.ValueString())
```

## Fix

Add `return` immediately after appending an error to diagnostics whenever the operation cannot continue without valid data.

```go
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*plan)
if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
	return // Return to prevent further processing
}
```

Repeat similar fix in the `Update` and `Delete` methods, right after the error append for input mapping. This ensures the method does not proceed with invalid data.

---

File to save:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_copilot_studio_application_insights_unhandled_error_high.md`
