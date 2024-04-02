// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type EnvironmenttSettingsSourceModel struct {
	Id            types.String `tfsdk:"id"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	AuditAndLogs  types.Object `tfsdk:"audit_and_logs"`
	Email         types.Object `tfsdk:"email"`
	Product       types.Object `tfsdk:"product"`
}

type AuditAndLogsSourceModel struct {
	PluginTraceLogSetting types.String `tfsdk:"plugin_trace_log_setting"`
	AuditSettings         types.Object `tfsdk:"audit_settings"`
}

type AuditSettingsSourceModel struct {
	IsAuditEnabled           types.Bool `tfsdk:"is_audit_enabled"`
	IsUserAccessAuditEnabled types.Bool `tfsdk:"is_user_access_audit_enabled"`
	IsReadAuditEnabled       types.Bool `tfsdk:"is_read_audit_enabled"`
}

type EmailSourceModel struct {
	EmailSettings types.Object `tfsdk:"email_settings"`
}

type EmailSettingsSourceModel struct {
	MaxUploadFileSize types.Int64 `tfsdk:"max_upload_file_size_in_bytes"`
}

type ProductSourceModel struct {
	BehaviorSettings types.Object `tfsdk:"behavior_settings"`
}

type BehaviorSettingsSourceModel struct {
	ShowDashboardCardsInExpandedState types.Bool `tfsdk:"show_dashboard_cards_in_expanded_state"`
}

func SetDefaultValuesForEnvironmentSettings(environmentSettings *EnvironmentSettingsDto) {
	if environmentSettings.MaxUploadFileSize == nil {
		defaultValue := int64(5242880)
		environmentSettings.MaxUploadFileSize = &defaultValue
	}
	if environmentSettings.PluginTraceLogSetting == nil {
		environmentSettings.PluginTraceLogSetting = new(int64)
		*environmentSettings.PluginTraceLogSetting = 0
	}
	if environmentSettings.IsAuditEnabled == nil {
		environmentSettings.IsAuditEnabled = new(bool)
		*environmentSettings.IsAuditEnabled = false
	}
	if environmentSettings.IsUserAccessAuditEnabled == nil {
		environmentSettings.IsUserAccessAuditEnabled = new(bool)
		*environmentSettings.IsUserAccessAuditEnabled = false
	}
	if environmentSettings.IsReadAuditEnabled == nil {
		environmentSettings.IsReadAuditEnabled = new(bool)
		*environmentSettings.IsReadAuditEnabled = false
	}
	if environmentSettings.BoundDashboardDefaultCardExpanded == nil {
		environmentSettings.BoundDashboardDefaultCardExpanded = new(bool)
		*environmentSettings.BoundDashboardDefaultCardExpanded = false
	}
}

func ConvertFromEnvironmentSettingsModel(ctx context.Context, environmentSettings EnvironmenttSettingsSourceModel) EnvironmentSettingsDto {
	environmentSettingsDto := EnvironmentSettingsDto{}
	auditSettingsObject := environmentSettings.AuditAndLogs.Attributes()["audit_settings"]
	if auditSettingsObject != nil && !auditSettingsObject.IsNull() && !auditSettingsObject.IsUnknown() {
		var auditAndLogsSourceModel AuditSettingsSourceModel
		auditSettingsObject.(basetypes.ObjectValue).As(ctx, &auditAndLogsSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if !auditAndLogsSourceModel.IsAuditEnabled.IsNull() && !auditAndLogsSourceModel.IsAuditEnabled.IsUnknown() {
			environmentSettingsDto.IsAuditEnabled = auditAndLogsSourceModel.IsAuditEnabled.ValueBoolPointer()
		}
		if !auditAndLogsSourceModel.IsUserAccessAuditEnabled.IsNull() && !auditAndLogsSourceModel.IsUserAccessAuditEnabled.IsUnknown() {
			environmentSettingsDto.IsUserAccessAuditEnabled = auditAndLogsSourceModel.IsUserAccessAuditEnabled.ValueBoolPointer()
		}
		if !auditAndLogsSourceModel.IsReadAuditEnabled.IsNull() && !auditAndLogsSourceModel.IsReadAuditEnabled.IsUnknown() {
			environmentSettingsDto.IsReadAuditEnabled = auditAndLogsSourceModel.IsReadAuditEnabled.ValueBoolPointer()
		}

		pluginSettings := environmentSettings.AuditAndLogs.Attributes()["plugin_trace_log_setting"]
		if pluginSettings != nil && !pluginSettings.IsNull() && !pluginSettings.IsUnknown() {
			pluginSettings := pluginSettings.(basetypes.StringValue)
			var v int64 = 0
			if pluginSettings.ValueString() == "Off" {
				environmentSettingsDto.PluginTraceLogSetting = &v
			} else if pluginSettings.ValueString() == "Exception" {
				v = 1
				environmentSettingsDto.PluginTraceLogSetting = &v
			} else if pluginSettings.ValueString() == "All" {
				v = 2
				environmentSettingsDto.PluginTraceLogSetting = &v
			}
		}
	}
	emailSettingsObject := environmentSettings.Email.Attributes()["email_settings"]
	if emailSettingsObject != nil && !emailSettingsObject.IsNull() && !emailSettingsObject.IsUnknown() {
		var emailSourceModel EmailSettingsSourceModel
		emailSettingsObject.(basetypes.ObjectValue).As(ctx, &emailSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if !emailSourceModel.MaxUploadFileSize.IsNull() && !emailSourceModel.MaxUploadFileSize.IsUnknown() {
			environmentSettingsDto.MaxUploadFileSize = emailSourceModel.MaxUploadFileSize.ValueInt64Pointer()
		}

	}

	behaviorSettings := environmentSettings.Product.Attributes()["behavior_settings"]
	if behaviorSettings != nil && !behaviorSettings.IsNull() && !behaviorSettings.IsUnknown() {
		var behaviorSettingsSourceModel BehaviorSettingsSourceModel
		behaviorSettings.(basetypes.ObjectValue).As(ctx, &behaviorSettingsSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if !behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState.IsNull() && !behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState.IsUnknown() {
			environmentSettingsDto.BoundDashboardDefaultCardExpanded = behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState.ValueBoolPointer()
		}
	}

	return environmentSettingsDto
}

func ConvertFromEnvironmentSettingsDto(environmentSettingsDto *EnvironmentSettingsDto) EnvironmenttSettingsSourceModel {
	environmentSettings := EnvironmenttSettingsSourceModel{
		Id: types.StringValue(uuid.New().String()),
	}

	pluginTraceSettings := "Unknown"
	if environmentSettingsDto.PluginTraceLogSetting != nil {
		switch *environmentSettingsDto.PluginTraceLogSetting {
		case 0:
			pluginTraceSettings = "Off"
		case 1:
			pluginTraceSettings = "Exception"
		case 2:
			pluginTraceSettings = "All"
		}
	}

	attrValuesAuditSettingsProperties := map[string]attr.Value{
		"is_audit_enabled":             types.BoolValue(*environmentSettingsDto.IsAuditEnabled),
		"is_user_access_audit_enabled": types.BoolValue(*environmentSettingsDto.IsUserAccessAuditEnabled),
		"is_read_audit_enabled":        types.BoolValue(*environmentSettingsDto.IsReadAuditEnabled),
	}

	attrAuditSettingsObject := map[string]attr.Type{
		"is_audit_enabled":             types.BoolType,
		"is_user_access_audit_enabled": types.BoolType,
		"is_read_audit_enabled":        types.BoolType,
	}

	attrTypesAuditAndLogsObject := map[string]attr.Type{
		"plugin_trace_log_setting": types.StringType,
		"audit_settings":           types.ObjectType{AttrTypes: attrAuditSettingsObject},
	}

	attrValuesAuditAndLogsProperties := map[string]attr.Value{
		"plugin_trace_log_setting": types.StringValue(pluginTraceSettings),
		"audit_settings":           types.ObjectValueMust(attrAuditSettingsObject, attrValuesAuditSettingsProperties),
	}

	attrEmailSettingsObject := map[string]attr.Type{
		"max_upload_file_size_in_bytes": types.Int64Type,
	}

	attrValuesEmailProperties := map[string]attr.Value{
		"email_settings": types.ObjectValueMust(attrEmailSettingsObject, map[string]attr.Value{
			"max_upload_file_size_in_bytes": types.Int64Value(*environmentSettingsDto.MaxUploadFileSize),
		}),
	}

	attrTypesEmailObject := map[string]attr.Type{
		"email_settings": types.ObjectType{AttrTypes: attrEmailSettingsObject},
	}

	attrBahaviorSettingsObject := map[string]attr.Type{
		"show_dashboard_cards_in_expanded_state": types.BoolType,
	}

	attrTypesProductObject := map[string]attr.Type{
		"behavior_settings": types.ObjectType{AttrTypes: attrBahaviorSettingsObject},
	}

	attrValuesProductProperties := map[string]attr.Value{
		"behavior_settings": types.ObjectValueMust(attrBahaviorSettingsObject, map[string]attr.Value{
			"show_dashboard_cards_in_expanded_state": types.BoolValue(*environmentSettingsDto.BoundDashboardDefaultCardExpanded),
		}),
	}

	environmentSettings.AuditAndLogs = types.ObjectValueMust(attrTypesAuditAndLogsObject, attrValuesAuditAndLogsProperties)
	environmentSettings.Email = types.ObjectValueMust(attrTypesEmailObject, attrValuesEmailProperties)
	environmentSettings.Product = types.ObjectValueMust(attrTypesProductObject, attrValuesProductProperties)

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
