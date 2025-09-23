// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	Security         types.Object `tfsdk:"security"`
}

type BehaviorSettingsSourceModel struct {
	ShowDashboardCardsInExpandedState types.Bool `tfsdk:"show_dashboard_cards_in_expanded_state"`
}

type FeaturesSourceModel struct {
	PowerAppsComponentFrameworkForCanvasApps             types.Bool `tfsdk:"power_apps_component_framework_for_canvas_apps"`
	PowerAppsMakerBotEnabled                             types.Bool `tfsdk:"enable_powerapps_maker_bot"`
	EnableAccessToSessionTranscriptsForCopilotStudio     types.Bool `tfsdk:"enable_access_to_session_transcripts_for_copilot_studio"`
	EnableTranscriptRecordingForCopilotStudio            types.Bool `tfsdk:"enable_transcript_recording_for_copilot_studio"`
	EnableCopilotStudioShareDataWithVivaInsights         types.Bool `tfsdk:"enable_copilot_studio_share_data_with_viva_insights"`
	EnableCopilotStudioCrossGeoShareDataWithVivaInsights types.Bool `tfsdk:"enable_copilot_studio_cross_geo_share_data_with_viva_insights"`
	EnablePreviewAndExperimentalAIModels                 types.Bool `tfsdk:"enable_preview_and_experimental_ai_models"`
	AiPromptsEnabled                                     types.Bool `tfsdk:"enable_ai_prompts"`
}

type SecuritySourceModel struct {
	EnableIpBasedCookieBinding           types.Bool `tfsdk:"enable_ip_based_cookie_binding"`
	EnableIpBasedFirewallRule            types.Bool `tfsdk:"enable_ip_based_firewall_rule"`
	AllowedIpRangeForFirewall            types.Set  `tfsdk:"allowed_ip_range_for_firewall"`
	AllowedServiceTagsForFirewall        types.Set  `tfsdk:"allowed_service_tags_for_firewall"`
	AllowApplicationUserAccess           types.Bool `tfsdk:"allow_application_user_access"`
	AllowMicrosoftTrustedServiceTags     types.Bool `tfsdk:"allow_microsoft_trusted_service_tags"`
	EnableIpBasedFirewallRuleInAuditMode types.Bool `tfsdk:"enable_ip_based_firewall_rule_in_audit_mode"`
	ReverseProxyIpAddresses              types.Set  `tfsdk:"reverse_proxy_ip_addresses"`
}

func convertFromEnvironmentSettingsModel(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel) (*environmentSettingsDto, error) {
	environmentSettingsDto := &environmentSettingsDto{}
	auditSettingsObject := environmentSettings.AuditAndLogs.Attributes()["audit_settings"]
	if auditSettingsObject != nil && helpers.IsKnown(auditSettingsObject) {
		objectValue, ok := auditSettingsObject.(basetypes.ObjectValue)
		if !ok {
			return nil, fmt.Errorf("failed to convert audit settings to ObjectValue, got %T: %+v", auditSettingsObject, auditSettingsObject)
		}

		var auditAndLogsSourceModel AuditSettingsSourceModel
		if diags := objectValue.As(ctx, &auditAndLogsSourceModel, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		}); diags != nil {
			return nil, fmt.Errorf("failed to convert audit settings: %v", diags)
		}

		if helpers.IsKnown(auditAndLogsSourceModel.IsAuditEnabled) {
			environmentSettingsDto.IsAuditEnabled = auditAndLogsSourceModel.IsAuditEnabled.ValueBoolPointer()
		}
		if helpers.IsKnown(auditAndLogsSourceModel.IsUserAccessAuditEnabled) {
			environmentSettingsDto.IsUserAccessAuditEnabled = auditAndLogsSourceModel.IsUserAccessAuditEnabled.ValueBoolPointer()
		}
		if helpers.IsKnown(auditAndLogsSourceModel.IsReadAuditEnabled) {
			environmentSettingsDto.IsReadAuditEnabled = auditAndLogsSourceModel.IsReadAuditEnabled.ValueBoolPointer()
		}
		if helpers.IsKnown(auditAndLogsSourceModel.AuditRetentionPeriodV2) {
			environmentSettingsDto.AuditRetentionPeriodV2 = auditAndLogsSourceModel.AuditRetentionPeriodV2.ValueInt32Pointer()
		}

		pluginSettings := environmentSettings.AuditAndLogs.Attributes()["plugin_trace_log_setting"]
		if pluginSettings != nil && helpers.IsKnown(pluginSettings) {
			pluginSettingsValue, ok := pluginSettings.(basetypes.StringValue)
			if !ok {
				return nil, errors.New("pluginSettings is not of type basetypes.StringValue")
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
	if err := convertFromEnvironmentEmailSettings(ctx, environmentSettings, environmentSettingsDto); err != nil {
		return nil, err
	}
	if err := convertFromEnvironmentBehaviorSettings(ctx, environmentSettings, environmentSettingsDto); err != nil {
		return nil, err
	}
	if err := convertFromEnvironmentFeatureSettings(ctx, environmentSettings, environmentSettingsDto); err != nil {
		return nil, err
	}
	if err := convertFromEnvironmentSecuritySettings(ctx, environmentSettings, environmentSettingsDto); err != nil {
		return nil, err
	}
	return environmentSettingsDto, nil
}

func convertFromEnvironmentEmailSettings(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel, environmentSettingsDto *environmentSettingsDto) error {
	emailSettingsObject := environmentSettings.Email.Attributes()["email_settings"]
	if emailSettingsObject != nil && helpers.IsKnown(emailSettingsObject) {
		objectValue, ok := emailSettingsObject.(basetypes.ObjectValue)
		if !ok {
			return errors.New("failed to convert email settings to ObjectValue")
		}

		var emailSourceModel EmailSettingsSourceModel
		if diags := objectValue.As(ctx, &emailSourceModel, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		}); diags != nil {
			return fmt.Errorf("failed to convert email settings: %v", diags)
		}

		if helpers.IsKnown(emailSourceModel.MaxUploadFileSize) {
			environmentSettingsDto.MaxUploadFileSize = emailSourceModel.MaxUploadFileSize.ValueInt64Pointer()
		}
	}
	return nil
}

func convertFromEnvironmentBehaviorSettings(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel, environmentSettingsDto *environmentSettingsDto) error {
	behaviorSettings := environmentSettings.Product.Attributes()["behavior_settings"]
	if behaviorSettings != nil && helpers.IsKnown(behaviorSettings) {
		var behaviorSettingsSourceModel BehaviorSettingsSourceModel
		if diags := behaviorSettings.(basetypes.ObjectValue).As(ctx, &behaviorSettingsSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
			return fmt.Errorf("failed to convert behavior settings: %v", diags)
		}

		if helpers.IsKnown(behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState) {
			environmentSettingsDto.BoundDashboardDefaultCardExpanded = behaviorSettingsSourceModel.ShowDashboardCardsInExpandedState.ValueBoolPointer()
		}
	}
	return nil
}

func convertFromEnvironmentFeatureSettings(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel, environmentSettingsDto *environmentSettingsDto) error {
	features := environmentSettings.Product.Attributes()["features"]
	if features != nil && helpers.IsKnown(features) {
		var featuresSourceModel FeaturesSourceModel
		if diags := features.(basetypes.ObjectValue).As(ctx, &featuresSourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
			return fmt.Errorf("failed to convert feature settings: %v", diags)
		}

		if helpers.IsKnown(featuresSourceModel.PowerAppsComponentFrameworkForCanvasApps) {
			environmentSettingsDto.PowerAppsComponentFrameworkForCanvasApps = featuresSourceModel.PowerAppsComponentFrameworkForCanvasApps.ValueBoolPointer()
		}
		if helpers.IsKnown(featuresSourceModel.PowerAppsMakerBotEnabled) {
			environmentSettingsDto.PowerAppsMakerBotEnabled = featuresSourceModel.PowerAppsMakerBotEnabled.ValueBoolPointer()
		}
		if helpers.IsKnown(featuresSourceModel.EnableAccessToSessionTranscriptsForCopilotStudio) {
			val := featuresSourceModel.EnableAccessToSessionTranscriptsForCopilotStudio.ValueBoolPointer()
			if val != nil {
				negated := !(*val)
				environmentSettingsDto.BlockAccessToSessionTranscriptsForCopilotStudio = &negated
			} else {
				environmentSettingsDto.BlockAccessToSessionTranscriptsForCopilotStudio = nil
			}
		}
		if helpers.IsKnown(featuresSourceModel.EnableTranscriptRecordingForCopilotStudio) {
			val := featuresSourceModel.EnableTranscriptRecordingForCopilotStudio.ValueBoolPointer()
			if val != nil {
				negated := !(*val)
				environmentSettingsDto.BlockTranscriptRecordingForCopilotStudio = &negated
			} else {
				environmentSettingsDto.BlockTranscriptRecordingForCopilotStudio = nil
			}
		}
		if helpers.IsKnown(featuresSourceModel.EnableCopilotStudioShareDataWithVivaInsights) {
			environmentSettingsDto.EnableCopilotStudioShareDataWithVivaInsights = featuresSourceModel.EnableCopilotStudioShareDataWithVivaInsights.ValueBoolPointer()
		}
		if helpers.IsKnown(featuresSourceModel.EnableCopilotStudioCrossGeoShareDataWithVivaInsights) {
			environmentSettingsDto.EnableCopilotStudioCrossGeoShareDataWithVivaInsights = featuresSourceModel.EnableCopilotStudioCrossGeoShareDataWithVivaInsights.ValueBoolPointer()
		}
		if helpers.IsKnown(featuresSourceModel.EnablePreviewAndExperimentalAIModels) {
			environmentSettingsDto.PaiPreviewScenarioEnabled = featuresSourceModel.EnablePreviewAndExperimentalAIModels.ValueBoolPointer()
		}
		if helpers.IsKnown(featuresSourceModel.AiPromptsEnabled) {
			environmentSettingsDto.AiPromptsEnabled = featuresSourceModel.AiPromptsEnabled.ValueBoolPointer()
		}
	}
	return nil
}

func convertFromEnvironmentSecuritySettings(ctx context.Context, environmentSettings EnvironmentSettingsResourceModel, environmentSettingsDto *environmentSettingsDto) error {
	security := environmentSettings.Product.Attributes()["security"]
	if security != nil && helpers.IsKnown(security) {
		var securitySourceModel SecuritySourceModel
		if diags := security.(basetypes.ObjectValue).As(ctx, &securitySourceModel, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); diags != nil {
			return fmt.Errorf("failed to convert security settings: %v", diags)
		}

		if helpers.IsKnown(securitySourceModel.EnableIpBasedCookieBinding) {
			environmentSettingsDto.EnableIpBasedCookieBinding = securitySourceModel.EnableIpBasedCookieBinding.ValueBoolPointer()
		}
		if helpers.IsKnown(securitySourceModel.EnableIpBasedFirewallRule) {
			environmentSettingsDto.EnableIpBasedFirewallRule = securitySourceModel.EnableIpBasedFirewallRule.ValueBoolPointer()
		}
		if helpers.IsKnown(securitySourceModel.AllowedIpRangeForFirewall) {
			value := strings.Join(helpers.SetToStringSlice(securitySourceModel.AllowedIpRangeForFirewall), ",")
			environmentSettingsDto.AllowedIpRangeForFirewall = &value
		}
		if helpers.IsKnown(securitySourceModel.AllowedServiceTagsForFirewall) {
			value := strings.Join(helpers.SetToStringSlice(securitySourceModel.AllowedServiceTagsForFirewall), ",")
			environmentSettingsDto.AllowedServiceTagsForFirewall = &value
		}
		if helpers.IsKnown(securitySourceModel.AllowApplicationUserAccess) {
			environmentSettingsDto.AllowApplicationUserAccess = securitySourceModel.AllowApplicationUserAccess.ValueBoolPointer()
		}
		if helpers.IsKnown(securitySourceModel.AllowMicrosoftTrustedServiceTags) {
			environmentSettingsDto.AllowMicrosoftTrustedServiceTags = securitySourceModel.AllowMicrosoftTrustedServiceTags.ValueBoolPointer()
		}
		if helpers.IsKnown(securitySourceModel.EnableIpBasedFirewallRuleInAuditMode) {
			environmentSettingsDto.EnableIpBasedFirewallRuleInAuditMode = securitySourceModel.EnableIpBasedFirewallRuleInAuditMode.ValueBoolPointer()
		}
		if helpers.IsKnown(securitySourceModel.ReverseProxyIpAddresses) {
			value := strings.Join(helpers.SetToStringSlice(securitySourceModel.ReverseProxyIpAddresses), ",")
			environmentSettingsDto.ReverseProxyIpAddresses = &value
		}
	}
	return nil
}

func convertFromEnvironmentSettingsDto[T EnvironmentSettingsResourceModel | EnvironmentSettingsDataSourceModel](environmentSettingsDto *environmentSettingsDto, timeout timeouts.Value) (T, error) {
	pluginTraceSettings := "Unknown"
	if environmentSettingsDto.PluginTraceLogSetting != nil {
		switch *environmentSettingsDto.PluginTraceLogSetting {
		case 0:
			pluginTraceSettings = "Off"
		case 1:
			pluginTraceSettings = "Exception"
		case 2:
			pluginTraceSettings = "All"
		default:
			pluginTraceSettings = "Unknown"
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
		"power_apps_component_framework_for_canvas_apps":                types.BoolType,
		"enable_powerapps_maker_bot":                                    types.BoolType,
		"enable_access_to_session_transcripts_for_copilot_studio":       types.BoolType,
		"enable_transcript_recording_for_copilot_studio":                types.BoolType,
		"enable_copilot_studio_share_data_with_viva_insights":           types.BoolType,
		"enable_copilot_studio_cross_geo_share_data_with_viva_insights": types.BoolType,
		"enable_preview_and_experimental_ai_models":                     types.BoolType,
		"enable_ai_prompts":                                             types.BoolType,
	}

	attrTypesSecurityObject := map[string]attr.Type{
		"enable_ip_based_cookie_binding":              types.BoolType,
		"enable_ip_based_firewall_rule":               types.BoolType,
		"allowed_ip_range_for_firewall":               types.SetType{ElemType: types.StringType},
		"allowed_service_tags_for_firewall":           types.SetType{ElemType: types.StringType},
		"allow_application_user_access":               types.BoolType,
		"allow_microsoft_trusted_service_tags":        types.BoolType,
		"enable_ip_based_firewall_rule_in_audit_mode": types.BoolType,
		"reverse_proxy_ip_addresses":                  types.SetType{ElemType: types.StringType},
	}

	attrTypesProductObject := map[string]attr.Type{
		"behavior_settings": types.ObjectType{AttrTypes: attrBahaviorSettingsObject},
		"features":          types.ObjectType{AttrTypes: attrFeaturesObject},
		"security":          types.ObjectType{AttrTypes: attrTypesSecurityObject},
	}

	reverseProxyAdresses := []attr.Value{}
	if environmentSettingsDto.ReverseProxyIpAddresses != nil {
		for _, proxy := range strings.Split(*environmentSettingsDto.ReverseProxyIpAddresses, ",") {
			reverseProxyAdresses = append(reverseProxyAdresses, types.StringValue(proxy))
		}
	}

	allowedIpRangeForFirewall := []attr.Value{}
	if environmentSettingsDto.AllowedIpRangeForFirewall != nil {
		for _, ip := range strings.Split(*environmentSettingsDto.AllowedIpRangeForFirewall, ",") {
			allowedIpRangeForFirewall = append(allowedIpRangeForFirewall, types.StringValue(ip))
		}
	}

	allowedServiceTags := []attr.Value{}
	if environmentSettingsDto.AllowedServiceTagsForFirewall != nil {
		for _, tag := range strings.Split(*environmentSettingsDto.AllowedServiceTagsForFirewall, ",") {
			allowedServiceTags = append(allowedServiceTags, types.StringValue(tag))
		}
	}

	attrValuesProductProperties := map[string]attr.Value{
		"behavior_settings": types.ObjectValueMust(attrBahaviorSettingsObject, map[string]attr.Value{
			"show_dashboard_cards_in_expanded_state": types.BoolValue(*environmentSettingsDto.BoundDashboardDefaultCardExpanded),
		}),
		"features": types.ObjectValueMust(attrFeaturesObject, map[string]attr.Value{
			"power_apps_component_framework_for_canvas_apps":                types.BoolValue(*environmentSettingsDto.PowerAppsComponentFrameworkForCanvasApps),
			"enable_powerapps_maker_bot":                                    types.BoolPointerValue(environmentSettingsDto.PowerAppsMakerBotEnabled),
			"enable_access_to_session_transcripts_for_copilot_studio":       types.BoolValue(!*environmentSettingsDto.BlockAccessToSessionTranscriptsForCopilotStudio),
			"enable_transcript_recording_for_copilot_studio":                types.BoolValue(!*environmentSettingsDto.BlockTranscriptRecordingForCopilotStudio),
			"enable_copilot_studio_share_data_with_viva_insights":           types.BoolValue(*environmentSettingsDto.EnableCopilotStudioShareDataWithVivaInsights),
			"enable_copilot_studio_cross_geo_share_data_with_viva_insights": types.BoolValue(*environmentSettingsDto.EnableCopilotStudioCrossGeoShareDataWithVivaInsights),
			"enable_preview_and_experimental_ai_models":                     types.BoolValue(*environmentSettingsDto.PaiPreviewScenarioEnabled),
			"enable_ai_prompts":                                             types.BoolValue(*environmentSettingsDto.AiPromptsEnabled),
		}),
		"security": types.ObjectValueMust(attrTypesSecurityObject, map[string]attr.Value{
			"enable_ip_based_cookie_binding":              types.BoolValue(*environmentSettingsDto.EnableIpBasedCookieBinding),
			"enable_ip_based_firewall_rule":               types.BoolValue(*environmentSettingsDto.EnableIpBasedFirewallRule),
			"allowed_ip_range_for_firewall":               types.SetValueMust(types.StringType, allowedIpRangeForFirewall),
			"allowed_service_tags_for_firewall":           types.SetValueMust(types.StringType, allowedServiceTags),
			"allow_application_user_access":               types.BoolValue(*environmentSettingsDto.AllowApplicationUserAccess),
			"allow_microsoft_trusted_service_tags":        types.BoolValue(*environmentSettingsDto.AllowMicrosoftTrustedServiceTags),
			"enable_ip_based_firewall_rule_in_audit_mode": types.BoolValue(*environmentSettingsDto.EnableIpBasedFirewallRuleInAuditMode),
			"reverse_proxy_ip_addresses":                  types.SetValueMust(types.StringType, reverseProxyAdresses),
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
		return environmentSettings, fmt.Errorf("unexpected type %T", environmentSettings)
	}
	if !ok {
		return environmentSettings, fmt.Errorf("unexpected type %T", environmentSettings)
	}
	return environmentSettings, nil
}
