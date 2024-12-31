// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// type CopilotStudioAppInsightsDataSource struct {
// 	helpers.TypeInfo
// 	CopilotStudioAppInsightsDatasourceClient client
// }

// type CopilotStudioAppInsightsDataSourceModel struct {
// 	AppInsightsConnectionString types.String `tfsdk:"applicationInsightsConnectionString"`
// 	IncludeSensitiveInformation types.Bool   `tfsdk:"includeSensitiveInformation"`
// 	IncludeActivities           types.Bool   `tfsdk:"includeActivities"`
// 	IncludeActions              types.Bool   `tfsdk:"includeActions"`
// 	NetworkIsolation            types.String `tfsdk:"networkIsolation"`
// 	Errors                      types.String `tfsdk:"errors"`
// }

type ResourceModel struct {
	AppInsightsConnectionString types.String `tfsdk:"applicationInsightsConnectionString"`
	IncludeSensitiveInformation types.Bool   `tfsdk:"includeSensitiveInformation"`
	IncludeActivities           types.Bool   `tfsdk:"includeActivities"`
	IncludeActions              types.Bool   `tfsdk:"includeActions"`
	NetworkIsolation            types.String `tfsdk:"networkIsolation"`
	Errors                      types.String `tfsdk:"errors"`
}

type Resource struct {
	helpers.TypeInfo
	CopilotStudioAppInsightsDatasourceClient client
}
