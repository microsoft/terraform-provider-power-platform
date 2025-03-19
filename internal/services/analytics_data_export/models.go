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

// convertDtoToModel converts an AnalyticsDataDto to AnalyticsDataModel.
func convertDtoToModel(dto *AnalyticsDataDto) *AnalyticsDataModel {
	if dto == nil {
		return nil
	}

	// Map environments
	environments := make([]EnvironmentModel, 0, len(dto.Environments))
	for _, env := range dto.Environments {
		environments = append(environments, EnvironmentModel{
			EnvironmentId:  types.StringValue(env.EnvironmentId),
			OrganizationId: types.StringValue(env.OrganizationId),
		})
	}

	// Map status
	status := make([]StatusModel, 0, len(dto.Status))
	for _, s := range dto.Status {
		message := types.StringNull()
		if s.Message != nil {
			message = types.StringValue(*s.Message)
		}
		status = append(status, StatusModel{
			Name:      types.StringValue(s.Name),
			State:     types.StringValue(s.State),
			LastRunOn: types.StringValue(s.LastRunOn),
			Message:   message,
		})
	}

	// Map scenarios
	scenarios := make([]types.String, 0, len(dto.Scenarios))
	for _, s := range dto.Scenarios {
		scenarios = append(scenarios, types.StringValue(s))
	}

	return &AnalyticsDataModel{
		ID:           types.StringValue(dto.ID),
		Source:       types.StringValue(dto.Source),
		Environments: environments,
		Status:       status,
		Sink: SinkModel{
			ID:                types.StringValue(dto.Sink.ID),
			Type:              types.StringValue(dto.Sink.Type),
			SubscriptionId:    types.StringValue(dto.Sink.SubscriptionId),
			ResourceGroupName: types.StringValue(dto.Sink.ResourceGroupName),
			ResourceName:      types.StringValue(dto.Sink.ResourceName),
			Key:               types.StringValue(dto.Sink.Key),
		},
		PackageName:      types.StringValue(dto.PackageName),
		Scenarios:        scenarios,
		ResourceProvider: types.StringValue(dto.ResourceProvider),
		AiType:           types.StringValue(dto.AiType),
	}
}
