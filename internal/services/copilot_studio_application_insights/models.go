// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type ResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	EnvironmentId               types.String `tfsdk:"environment_id"`
	BotId                       types.String `tfsdk:"bot_id"`
	AppInsightsConnectionString types.String `tfsdk:"application_insights_connection_string"`
	IncludeSensitiveInformation types.Bool   `tfsdk:"include_sensitive_information"`
	IncludeActivities           types.Bool   `tfsdk:"include_activities"`
	IncludeActions              types.Bool   `tfsdk:"include_actions"`
}

type Resource struct {
	helpers.TypeInfo
	CopilotStudioApplicationInsightsClient client
}

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

func convertAppInsightsConfigModelFromDto(appInsightsConfigDto *CopilotStudioAppInsightsDto) (*ResourceModel, error) {
	model := &ResourceModel{
		Id:                          types.StringValue(fmt.Sprintf("%s_%s", appInsightsConfigDto.EnvironmentId, appInsightsConfigDto.BotId)),
		EnvironmentId:               types.StringValue(appInsightsConfigDto.EnvironmentId),
		BotId:                       types.StringValue(appInsightsConfigDto.BotId),
		AppInsightsConnectionString: types.StringValue(appInsightsConfigDto.AppInsightsConnectionString),
		IncludeSensitiveInformation: types.BoolValue(appInsightsConfigDto.IncludeSensitiveInformation),
		IncludeActivities:           types.BoolValue(appInsightsConfigDto.IncludeActivities),
		IncludeActions:              types.BoolValue(appInsightsConfigDto.IncludeActions),
	}
	return model, nil
}
