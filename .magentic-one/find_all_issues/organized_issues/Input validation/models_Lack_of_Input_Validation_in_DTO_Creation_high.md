# Issue 1: Lack of Input Validation in DTO Creation

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem

The function `createAppInsightsConfigDtoFromSourceModel` directly converts values from the `ResourceModel` to the DTO without any validation for required fields (such as `EnvironmentId`, `BotId`, or `AppInsightsConnectionString`). If any of these fields are empty or invalid, the DTO may be created with incomplete or invalid data, potentially causing runtime errors further along in the workflow.

## Impact

Severity: **High**

Unvalidated input may lead to the creation of DTOs with missing or malformed data, which can cause downstream API errors, unexpected behavior, or subtle bugs that are difficult to trace.

## Location

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

## Fix

Add validation to ensure all required fields are present and valid before creating the DTO. Return an error if validation fails.

```go
func createAppInsightsConfigDtoFromSourceModel(appInsightsConfigSource ResourceModel) (*CopilotStudioAppInsightsDto, error) {
	envId := appInsightsConfigSource.EnvironmentId.ValueString()
	botId := appInsightsConfigSource.BotId.ValueString()
	connStr := appInsightsConfigSource.AppInsightsConnectionString.ValueString()

	if envId == "" {
		return nil, fmt.Errorf("EnvironmentId cannot be empty")
	}
	if botId == "" {
		return nil, fmt.Errorf("BotId cannot be empty")
	}
	if connStr == "" {
		return nil, fmt.Errorf("ApplicationInsightsConnectionString cannot be empty")
	}

	appInsightsConfigDto := &CopilotStudioAppInsightsDto{
		EnvironmentId:               envId,
		BotId:                       botId,
		AppInsightsConnectionString: connStr,
		IncludeSensitiveInformation: appInsightsConfigSource.IncludeSensitiveInformation.ValueBool(),
		IncludeActivities:           appInsightsConfigSource.IncludeActivities.ValueBool(),
		IncludeActions:              appInsightsConfigSource.IncludeActions.ValueBool(),
		NetworkIsolation:            "PublicNetwork",
	}

	return appInsightsConfigDto, nil
}
```
