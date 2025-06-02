# Title:
Missing Error Handling for Application Insights Config

## Path
/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem
The function `createAppInsightsConfigDtoFromSourceModel` initializes and returns an `AppInsightsConfigDto` object but lacks robust input validation and error-checking mechanisms—no checks are performed for invalid or empty `types.String` or other fields before accessing methods like `.ValueString()` and `.ValueBool()`.

## Impact
If any required field in `appInsightsConfigSource` is invalid or empty, calling `.ValueString()` or `.ValueBool()` could result in undefined behavior or runtime panics. This poses a risk to the stability of the application.

## Location
`createAppInsightsConfigDtoFromSourceModel` function

## Code Issue:
```go
func createAppInsightsConfigDtoFromSourceModel(appInsightsConfigSource ResourceModel) (*CopilotStudioAppInsightsDto, error) {
	appInsightsConfigDto := &CopilotStudioAppInsightsDto{
		EnvironmentId:               appInsightsConfigSource.EnvironmentId.ValueString(),
		BotId:                       appInsightsConfigSource.BotId.ValueString(),
		AppInsightsConnectionString: appInsightsConfigSource.AppInsightsConnectionString.ValueString(),
		IncludeSensitiveInformation: appInsightsConfigSource.IncludeSensitiveInformation.ValueBool(),
		IncludeActivities:           appInsightsConfigSource.IncludeActivities.ValueBool(),
		IncludeActions:              appInsightsConfigSource.IncludeActions.ValueBool(),
		NetworkIsolation:            "PublicNetwork",
	}

	return appInsightsConfigDto, nil
}
```

## Fix:
Implement proper checks to validate that required fields are non-empty and valid. If a field is empty or invalid, return an appropriate error.

```go
func createAppInsightsConfigDtoFromSourceModel(appInsightsConfigSource ResourceModel) (*CopilotStudioAppInsightsDto, error) {
	// Validate required fields
	if appInsightsConfigSource.EnvironmentId.IsNull() || appInsightsConfigSource.EnvironmentId.ValueString() == "" {
		return nil, fmt.Errorf("environment ID is required")
	}
	if appInsightsConfigSource.BotId.IsNull() || appInsightsConfigSource.BotId.ValueString() == "" {
		return nil, fmt.Errorf("bot ID is required")
	}
	if appInsightsConfigSource.AppInsightsConnectionString.IsNull() || appInsightsConfigSource.AppInsightsConnectionString.ValueString() == "" {
		return nil, fmt.Errorf("Application Insights connection string is required")
	}

	// Initialize the DTO
	appInsightsConfigDto := &CopilotStudioAppInsightsDto{
		EnvironmentId:               appInsightsConfigSource.EnvironmentId.ValueString(),
		BotId:                       appInsightsConfigSource.BotId.ValueString(),
		AppInsightsConnectionString: appInsightsConfigSource.AppInsightsConnectionString.ValueString(),
		IncludeSensitiveInformation: appInsightsConfigSource.IncludeSensitiveInformation.ValueBool(),
		IncludeActivities:           appInsightsConfigSource.IncludeActivities.ValueBool(),
		IncludeActions:              appInsightsConfigSource.IncludeActions.ValueBool(),
		NetworkIsolation:            "PublicNetwork",
	}

	return appInsightsConfigDto, nil
}
```