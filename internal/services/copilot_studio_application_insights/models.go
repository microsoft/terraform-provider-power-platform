// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type ResourceModel struct {
	EnvironmentId               types.String `tfsdk:"environmentId"`
	BotId                       types.String `tfsdk:"botId"`
	AppInsightsConnectionString types.String `tfsdk:"applicationInsightsConnectionString"`
	IncludeSensitiveInformation types.Bool   `tfsdk:"includeSensitiveInformation"`
	IncludeActivities           types.Bool   `tfsdk:"includeActivities"`
	IncludeActions              types.Bool   `tfsdk:"includeActions"`
	NetworkIsolation            types.String `tfsdk:"networkIsolation"`
	Errors                      types.String `tfsdk:"errors"`
}

type Resource struct {
	helpers.TypeInfo
	CopilotStudioApplicationInsightsClient client
}

func createAppInsightsConfigDtoFromSourceModel(ctx context.Context, appInsightsConfigSource ResourceModel) (*CopilotStudioAppInsightsDto, error) {
	appInsightsConfigDto := &CopilotStudioAppInsightsDto{
		EnvironmentId:               appInsightsConfigSource.EnvironmentId.ValueString(),
		BotId:                       appInsightsConfigSource.BotId.ValueString(),
		AppInsightsConnectionString: appInsightsConfigSource.AppInsightsConnectionString.ValueString(),
		IncludeSensitiveInformation: appInsightsConfigSource.IncludeSensitiveInformation.ValueBool(),
		IncludeActivities:           appInsightsConfigSource.IncludeActivities.ValueBool(),
		IncludeActions:              appInsightsConfigSource.IncludeActions.ValueBool(),
	}

	return appInsightsConfigDto, nil
}

func convertAppInsightsConfigModelFromDto(appInsightsConfigDto *CopilotStudioAppInsightsDto) (*ResourceModel, error) {
	model := &ResourceModel{
		EnvironmentId:               types.StringValue(appInsightsConfigDto.EnvironmentId),
		BotId:                       types.StringValue(appInsightsConfigDto.BotId),
		AppInsightsConnectionString: types.StringValue(appInsightsConfigDto.AppInsightsConnectionString),
		IncludeSensitiveInformation: types.BoolValue(appInsightsConfigDto.IncludeSensitiveInformation),
		IncludeActivities:           types.BoolValue(appInsightsConfigDto.IncludeActivities),
		IncludeActions:              types.BoolValue(appInsightsConfigDto.IncludeActions),
		NetworkIsolation:            types.StringValue(appInsightsConfigDto.NetworkIsolation),
		Errors:                      types.StringValue(appInsightsConfigDto.Errors),
	}
	return model, nil
}
