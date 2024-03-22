// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnvironmenttSettingsSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	EnvironmentId                     types.String `tfsdk:"environment_id"`
	MaxUploadFileSize                 types.Int64  `tfsdk:"max_upload_file_size_in_bytes"`
	ShowDashboardCardsInExpandedState types.Bool   `tfsdk:"show_dashboard_cards_in_expanded_state"`
	PluginTraceLogSetting             types.String `tfsdk:"plugin_trace_log_setting"`
	IsAuditEnabled                    types.Bool   `tfsdk:"is_audit_enabled"`
	IsUserAccessAuditEnabled          types.Bool   `tfsdk:"is_user_access_audit_enabled"`
	IsReadAuditEnabled                types.Bool   `tfsdk:"is_read_audit_enabled"`
}

func ConvertFromEnvironmentSettingsModel(environmentSettings EnvironmenttSettingsSourceModel) EnvironmentSettingsDto {
	environmentSettingsDto := EnvironmentSettingsDto{}
	if !environmentSettings.MaxUploadFileSize.IsNull() && !environmentSettings.MaxUploadFileSize.IsUnknown() {
		environmentSettingsDto.MaxUploadFileSize = environmentSettings.MaxUploadFileSize.ValueInt64Pointer()
	}
	if !environmentSettings.IsAuditEnabled.IsNull() && !environmentSettings.IsAuditEnabled.IsUnknown() {
		environmentSettingsDto.IsAuditEnabled = environmentSettings.IsAuditEnabled.ValueBoolPointer()
	}
	if !environmentSettings.IsReadAuditEnabled.IsNull() && !environmentSettings.IsReadAuditEnabled.IsUnknown() {
		environmentSettingsDto.IsReadAuditEnabled = environmentSettings.IsReadAuditEnabled.ValueBoolPointer()
	}
	if !environmentSettings.IsUserAccessAuditEnabled.IsNull() && !environmentSettings.IsUserAccessAuditEnabled.IsUnknown() {
		environmentSettingsDto.IsUserAccessAuditEnabled = environmentSettings.IsUserAccessAuditEnabled.ValueBoolPointer()
	}
	if !environmentSettings.ShowDashboardCardsInExpandedState.IsNull() && !environmentSettings.ShowDashboardCardsInExpandedState.IsUnknown() {
		environmentSettingsDto.BoundDashboardDefaultCardExpanded = environmentSettings.ShowDashboardCardsInExpandedState.ValueBoolPointer()
	}
	if !environmentSettings.PluginTraceLogSetting.IsNull() && !environmentSettings.PluginTraceLogSetting.IsUnknown() {
		var v int64 = 0
		if *environmentSettings.PluginTraceLogSetting.ValueStringPointer() == "Off" {
			environmentSettingsDto.PluginTraceLogSetting = &v
		} else if *environmentSettings.PluginTraceLogSetting.ValueStringPointer() == "Exception" {
			v = 1
			environmentSettingsDto.PluginTraceLogSetting = &v
		} else if *environmentSettings.PluginTraceLogSetting.ValueStringPointer() == "All" {
			v = 2
			environmentSettingsDto.PluginTraceLogSetting = &v
		}
	}
	return environmentSettingsDto
}

func ConvertFromEnvironmentSettingsDto(environmentSettingsDto *EnvironmentSettingsDto) EnvironmenttSettingsSourceModel {
	environmentSettings := EnvironmenttSettingsSourceModel{
		Id:                                types.StringValue(uuid.New().String()),
		MaxUploadFileSize:                 types.Int64Value(*environmentSettingsDto.MaxUploadFileSize),
		ShowDashboardCardsInExpandedState: types.BoolValue(*environmentSettingsDto.BoundDashboardDefaultCardExpanded),
		IsAuditEnabled:                    types.BoolValue(*environmentSettingsDto.IsAuditEnabled),
		IsUserAccessAuditEnabled:          types.BoolValue(*environmentSettingsDto.IsUserAccessAuditEnabled),
		IsReadAuditEnabled:                types.BoolValue(*environmentSettingsDto.IsReadAuditEnabled),
	}

	if environmentSettingsDto.PluginTraceLogSetting != nil {
		switch *environmentSettingsDto.PluginTraceLogSetting {
		case 0:
			environmentSettings.PluginTraceLogSetting = types.StringValue("Off")
		case 1:
			environmentSettings.PluginTraceLogSetting = types.StringValue("Exception")
		case 2:
			environmentSettings.PluginTraceLogSetting = types.StringValue("All")
		}
	}

	return environmentSettings
}

type EnvironmentSettingsValueDto struct {
	Value []EnvironmentSettingsDto `json:"value"`
}

type EnvironmentSettingsDto struct {
	MaxUploadFileSize                 *int64  `json:"maxuploadfilesize,omitempty"`
	PluginTraceLogSetting             *int64  `json:"plugintracelogsetting,omitempty"`
	IsAuditEnabled                    *bool   `json:"isauditenabled,omitempty"`
	IsUserAccessAuditEnabled          *bool   `json:"isuseraccessauditenabled,omitempty"`
	IsReadAuditEnabled                *bool   `json:"isreadauditenabled,omitempty"`
	BoundDashboardDefaultCardExpanded *bool   `json:"bounddashboarddefaultcardexpanded,omitempty"`
	OrganizationId                    *string `json:"organizationid,omitempty"`
}

type EnvironmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties EnvironmentIdPropertiesDto `json:"properties"`
}

type EnvironmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type LinkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
