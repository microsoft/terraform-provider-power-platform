// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package analytics_data_export

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AnalyticsDataExportModel represents the model for Analytics Data Export.
type AnalyticsDataExportModel struct {
	Exports []AnalyticsDataModel `tfsdk:"exports"`
}

// AnalyticsDataModel represents the model for individual Analytics Data entry.
type AnalyticsDataModel struct {
	ID               types.String       `tfsdk:"id"`
	Source           types.String       `tfsdk:"source"`
	Environments     []EnvironmentModel `tfsdk:"environments"`
	Status           []StatusModel      `tfsdk:"status"`
	Sink             SinkModel          `tfsdk:"sink"`
	PackageName      types.String       `tfsdk:"package_name"`
	Scenarios        []types.String     `tfsdk:"scenarios"`
	ResourceProvider types.String       `tfsdk:"resource_provider"`
	AiType           types.String       `tfsdk:"ai_type"`
}

// StatusModel represents the model for Status entry.
type StatusModel struct {
	Name      types.String `tfsdk:"name"`
	State     types.String `tfsdk:"state"`
	LastRunOn types.String `tfsdk:"last_run_on"`
	Message   types.String `tfsdk:"message"`
}

// EnvironmentModel represents the model for Environment configuration.
type EnvironmentModel struct {
	EnvironmentId  types.String `tfsdk:"environment_id"`
	OrganizationId types.String `tfsdk:"organization_id"`
}

// SinkModel represents the model for Sink configuration.
type SinkModel struct {
	ID                types.String `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	SubscriptionId    types.String `tfsdk:"subscription_id"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	ResourceName      types.String `tfsdk:"resource_name"`
	Key               types.String `tfsdk:"key"`
}

// AnalyticsDataCreateModel represents the model for creating Analytics Data Export.
type AnalyticsDataCreateModel struct {
	Source           types.String       `tfsdk:"source"`
	Environments     []EnvironmentModel `tfsdk:"environments"`
	Sink             SinkModel          `tfsdk:"sink"`
	PackageName      types.String       `tfsdk:"package_name"`
	Scenarios        []types.String     `tfsdk:"scenarios"`
	ResourceProvider types.String       `tfsdk:"resource_provider"`
}
