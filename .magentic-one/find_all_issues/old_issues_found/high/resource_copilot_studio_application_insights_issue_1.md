# Title

Error Handling Missing in Update and Create Methods

##

/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

In both the `Update` and `Create` methods, the error from `createAppInsightsConfigDtoFromSourceModel` function is captured and added to `resp.Diagnostics` but the function continues execution. This leads to additional operations executing on an invalid object, increasing the risk of panic or undefined behavior.

## Impact

- This can lead to undefined behavior, obscure stack traces, and potential panics in production.
- Faulty execution paths might corrupt or store invalid state as well.
- Severity: High

## Location

Lines:
- `Update` method: near line 214
- `Create` method: near line 135

## Code Issue

```go
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*plan)
if err != nil {
    resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
}

// You can't really create a config, so treat a create as an update
appInsightsConfigDto, err := r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate, plan.BotId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating/updating %s", r.FullTypeName()), err.Error())
    return
}
```

## Fix

Add return immediately after adding the error to `resp.Diagnostics` to ensure the function exits cleanly upon handling the error.

```go
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*plan)
if err != nil {
    resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
    return  // Ensure early return on error
}

// Proceed with further processing only if no error occurred
appInsightsConfigDto, err := r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate, plan.BotId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating/updating %s", r.FullTypeName()), err.Error())
    return
}
```
