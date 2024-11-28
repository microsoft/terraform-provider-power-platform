// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type EnvironmentSettingsDataSource struct {
	helpers.TypeInfo
	EnvironmentSettingsClient client
}

type EnvironmentSettingsResource struct {
	helpers.TypeInfo
	EnvironmentSettingClient client
}

type EnvironmentSettingsResourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	Id            types.String   `tfsdk:"id"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	AuditAndLogs  types.Object   `tfsdk:"audit_and_logs"`
	Email         types.Object   `tfsdk:"email"`
	Product       types.Object   `tfsdk:"product"`
}

type EnvironmentSettingsDataSourceModel struct {
	Timeouts      timeouts.Value `tfsdk:"timeouts"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
	AuditAndLogs  types.Object   `tfsdk:"audit_and_logs"`
	Email         types.Object   `tfsdk:"email"`
	Product       types.Object   `tfsdk:"product"`
}

type AuditAndLogsSourceModel struct {
	PluginTraceLogSetting types.String `tfsdk:"plugin_trace_log_setting"`
	AuditSettings         types.Object `tfsdk:"audit_settings"`
}

type AuditSettingsSourceModel struct {
	IsAuditEnabled           types.Bool  `tfsdk:"is_audit_enabled"`
	IsUserAccessAuditEnabled types.Bool  `tfsdk:"is_user_access_audit_enabled"`
	IsReadAuditEnabled       types.Bool  `tfsdk:"is_read_audit_enabled"`
	AuditRetentionPeriodV2   types.Int32 `tfsdk:"log_retention_period_in_days"`
}

type EmailSourceModel struct {
	EmailSettings types.Object `tfsdk:"email_settings"`
}

type EmailSettingsSourceModel struct {
	MaxUploadFileSize types.Int64 `tfsdk:"max_upload_file_size_in_bytes"`
}

type ProductSourceModel struct {
	BehaviorSettings types.Object `tfsdk:"behavior_settings"`
	Features         types.Object `tfsdk:"features"`
}

type BehaviorSettingsSourceModel struct {
	ShowDashboardCardsInExpandedState types.Bool `tfsdk:"show_dashboard_cards_in_expanded_state"`
}

type FeaturesSourceModel struct {
	PowerAppsComponentFrameworkForCanvasApps types.Bool `tfsdk:"power_apps_component_framework_for_canvas_apps"`
}

func convertFromEnvironmentSettingsModel(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel) (*environmentSettingsDto, error) {
	environmentSettingsDto := &environmentSettingsDto{}
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
		if !auditAndLogsSourceModel.AuditRetentionPeriodV2.IsNull() && !auditAndLogsSourceModel.AuditRetentionPeriodV2.IsUnknown() {
			environmentSettingsDto.AuditRetentionPeriodV2 = auditAndLogsSourceModel.AuditRetentionPeriodV2.ValueInt32Pointer()
		}

		pluginSettings := environmentSettings.AuditAndLogs.Attributes()["plugin_trace_log_setting"]
		if pluginSettings != nil && !pluginSettings.IsNull() && !pluginSettings.IsUnknown() {
			pluginSettingsValue, ok := pluginSettings.(basetypes.StringValue)
			if !ok {
				return nil, fmt.Errorf("pluginSettings is not of type basetypes.StringValue")
			}
			var v int64
			if pluginSettingsValue.ValueString() == "Off" {
				environmentSettingsDto.PluginTraceLogSetting = &v
			} else if pluginSettingsValue.ValueString() == "Exception" {
				v = 1
				environmentSettingsDto.PluginTraceLogSetting = &v
			} else if pluginSettingsValue.ValueString() == "All" {
				v = 2
				environmentSettingsDto.PluginTraceLogSetting = &v
			}
		}
	}
	convertFromEnvironmentEmailSettings(ctx, environmentSettings, environmentSettingsDto)
	convertFromEnvironmentBehaviorSettings(ctx, environmentSettings, environmentSettingsDto)
	convertFromEnvironmentFeatureSettings(ctx, environmentSettings, environmentSettingsDto)
	return environmentSettingsDto, nil
}

func convertFromEnvironmentEmailSettings(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel, environmentSettingsDto *environmentSettingsDto) {
	emailSettingsObject := environmentSettings.Email.Attributes()["email_settings"]
	if emailSettingsObject != nil && !emailSettingsObject.IsNull() && !emailSettingsObject.IsUnknown() {
		var emailSourceModel EmailSettingsSourceModel
		if err := emailSettingsObject.(basetypes.ObjectValue).As(ctx, &emailSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); err != nil {
			return
		}

		if !emailSourceModel.MaxUploadFileSize.IsNull() && !emailSourceModel.MaxUploadFileSize.IsUnknown() {
			environmentSettingsDto.MaxUploadFileSize = emailSourceModel.MaxUploadFileSize.ValueInt64Pointer()
		}
	}
}

func convertFromEnvironmentBehaviorSettings(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel, environmentSettingsDto *environmentSettingsDto) {
	behaviorSettings := environmentSettings.Product.Attributes()["behavior_settings"]
	if behaviorSettings != nil && !behaviorSettings.IsNull() && !behaviorSettings.IsUnknown() {
		var behaviorSettingsSourceModel BehaviorSettingsSourceModel
		behaviorSettings.(basetypes.ObjectValue).As(ctx, &behaviorSettingsSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if !behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState.IsNull() && !behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState.IsUnknown() {
			environmentSettingsDto.BoundDashboardDefaultCardExpanded = behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState.ValueBoolPointer()
		}
	}
}

func convertFromEnvironmentFeatureSettings(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel, environmentSettingsDto *environmentSettingsDto) {
	features := environmentSettings.Product.Attributes()["features"]
	if features != nil && !features.IsNull() && !features.IsUnknown() {
		var featuresSourceModel FeaturesSourceModel
		features.(basetypes.ObjectValue).As(ctx, &featuresSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if !featuresSourceModel.PowerAppsComponentFrameworkForCanvasApps.IsNull() && !featuresSourceModel.PowerAppsComponentFrameworkForCanvasApps.IsUnknown() {
			environmentSettingsDto.PowerAppsComponentFrameworkForCanvasApps = featuresSourceModel.PowerAppsComponentFrameworkForCanvasApps.ValueBoolPointer()
		}
	}
}

func convertFromEnvironmentSettingsDto[T EnvironmentSettingsResourceModel | EnvironmentSettingsDataSourceModel](environmentSettingsDto *environmentSettingsDto, timeout timeouts.Value) T {
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

	logRetentionPeriodTypeValue := types.Int32Value(-1)
	if environmentSettingsDto.AuditRetentionPeriodV2 != nil {
		logRetentionPeriodTypeValue = types.Int32Value(*environmentSettingsDto.AuditRetentionPeriodV2)
	}

	attrValuesAuditSettingsProperties := map[string]attr.Value{
		"is_audit_enabled":             types.BoolValue(*environmentSettingsDto.IsAuditEnabled),
		"is_user_access_audit_enabled": types.BoolValue(*environmentSettingsDto.IsUserAccessAuditEnabled),
		"is_read_audit_enabled":        types.BoolValue(*environmentSettingsDto.IsReadAuditEnabled),
		"log_retention_period_in_days": logRetentionPeriodTypeValue,
	}

	attrAuditSettingsObject := map[string]attr.Type{
		"is_audit_enabled":             types.BoolType,
		"is_user_access_audit_enabled": types.BoolType,
		"is_read_audit_enabled":        types.BoolType,
		"log_retention_period_in_days": types.Int32Type,
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

	attrFeaturesObject := map[string]attr.Type{
		"power_apps_component_framework_for_canvas_apps": types.BoolType,
	}

	attrTypesProductObject := map[string]attr.Type{
		"behavior_settings": types.ObjectType{AttrTypes: attrBahaviorSettingsObject},
		"features":          types.ObjectType{AttrTypes: attrFeaturesObject},
	}

	attrValuesProductProperties := map[string]attr.Value{
		"behavior_settings": types.ObjectValueMust(attrBahaviorSettingsObject, map[string]attr.Value{
			"show_dashboard_cards_in_expanded_state": types.BoolValue(*environmentSettingsDto.BoundDashboardDefaultCardExpanded),
		}),
		"features": types.ObjectValueMust(attrFeaturesObject, map[string]attr.Value{
			"power_apps_component_framework_for_canvas_apps": types.BoolValue(*environmentSettingsDto.PowerAppsComponentFrameworkForCanvasApps),
		}),
	}

	var environmentSettings T
	var ok bool
	switch any(environmentSettings).(type) {
	case EnvironmentSettingsResourceModel:
		environmentSettings, ok = any(EnvironmentSettingsResourceModel{
			Timeouts:     timeout,
			AuditAndLogs: types.ObjectValueMust(attrTypesAuditAndLogsObject, attrValuesAuditAndLogsProperties),
			Email:        types.ObjectValueMust(attrTypesEmailObject, attrValuesEmailProperties),
			Product:      types.ObjectValueMust(attrTypesProductObject, attrValuesProductProperties),
		}).(T)
	case EnvironmentSettingsDataSourceModel:
		environmentSettings, ok = any(EnvironmentSettingsDataSourceModel{
			Timeouts:     timeout,
			AuditAndLogs: types.ObjectValueMust(attrTypesAuditAndLogsObject, attrValuesAuditAndLogsProperties),
			Email:        types.ObjectValueMust(attrTypesEmailObject, attrValuesEmailProperties),
			Product:      types.ObjectValueMust(attrTypesProductObject, attrValuesProductProperties),
		}).(T)
	default:
		panic(fmt.Sprintf("unexpected type %T", environmentSettings))
	}
	if !ok {
		panic(fmt.Sprintf("unexpected type %T", environmentSettings))
	}
	return environmentSettings
}
