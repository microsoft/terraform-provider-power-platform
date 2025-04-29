# Title:
Hardcoded Property Value in DTO Configuration

## Path
/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem
The `NetworkIsolation` property in the DTO is hardcoded to `"PublicNetwork"` inside the `createAppInsightsConfigDtoFromSourceModel` function without explanation or allowance for customization.

## Impact
Hardcoding values limits the configurability of the function and makes it difficult to adapt to changes in business requirements without modifying the code. Additionally, this can lead to operational problems if the assumed hardcoded value becomes invalid or needs adjustment.

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
		NetworkIsolation:            "PublicNetwork", // Hardcoded value
	}

	return appInsightsConfigDto, nil
}
```

## Fix:
Introduce a mechanism to allow this value to be dynamically set, either by using a parameter or coupling it with other configurations.

```go
func createAppInsightsConfigDtoFromSourceModel(appInsightsConfigSource ResourceModel, networkIsolation string) (*CopilotStudioAppInsightsDto, error) {
	// Validate network isolation input
	if networkIsolation == "" {
		return nil, fmt.Errorf("network isolation cannot be empty")
	}

	// Initialize the DTO
	appInsightsConfigDto := &CopilotStudioAppInsightsDto{
		EnvironmentId:               appInsightsConfigSource.EnvironmentId.ValueString(),
		BotId:                       appInsightsConfigSource.BotId.ValueString(),
		AppInsightsConnectionString: appInsightsConfigSource.AppInsightsConnectionString.ValueString(),
		IncludeSensitiveInformation: appInsightsConfigSource.IncludeSensitiveInformation.ValueBool(),
		IncludeActivities:           appInsightsConfigSource.IncludeActivities.ValueBool(),
		IncludeActions:              appInsightsConfigSource.IncludeActions.ValueBool(),
		NetworkIsolation:            networkIsolation, // Dynamically set
	}

	return appInsightsConfigDto, nil
}
```