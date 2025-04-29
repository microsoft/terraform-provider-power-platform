# Title
Incomplete Error Handling in Delete Method

##

/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

The error from `createAppInsightsConfigDtoFromSourceModel` within the `Delete` method of the resource is handled by adding it to `resp.Diagnostics` but does not prevent further execution. This can lead to subsequent operations using an invalid or incomplete configuration.

## Impact

- Continued execution after an error can lead to undefined behavior and incorrect application state.
- This can cause failed operations down the line or corrupted data.
- Severity: Medium

## Location

Lines: near 297

## Code Issue

```go
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*state)
if err != nil {
    resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
}

appInsightsConfigToCreate.AppInsightsConnectionString = ""
appInsightsConfigToCreate.IncludeSensitiveInformation = false
appInsightsConfigToCreate.IncludeActivities = false
appInsightsConfigToCreate.IncludeActions = false

_, err = r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate, state.BotId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
    return
}
```

## Fix

Include an early return upon encountering an error to stop further execution.

```go
appInsightsConfigToCreate, err := createAppInsightsConfigDtoFromSourceModel(*state)
if err != nil {
    resp.Diagnostics.AddError("Error when converting source model to create Copilot Studio Application Insights configuration dto", err.Error())
    return  // Stop execution on error
}

// Proceed only if no error occurred
appInsightsConfigToCreate.AppInsightsConnectionString = ""
appInsightsConfigToCreate.IncludeSensitiveInformation = false
appInsightsConfigToCreate.IncludeActivities = false
appInsightsConfigToCreate.IncludeActions = false

_, err = r.CopilotStudioApplicationInsightsClient.updateCopilotStudioAppInsightsConfiguration(ctx, *appInsightsConfigToCreate, state.BotId.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
    return
}
```