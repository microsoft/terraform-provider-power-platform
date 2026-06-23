// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_based_policy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type environmentGroupRuleBasedPolicyResource struct {
	helpers.TypeInfo
	RuleBasedPolicyClient Client
}

type environmentGroupRuleBasedPolicyResourceModel struct {
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String   `tfsdk:"id"`
	EnvironmentGroupId types.String   `tfsdk:"environment_group_id"`
	RuleSets           types.Object   `tfsdk:"rule_sets"`
}

type environmentGroupRuleBasedPolicyRuleSetsModel struct {
	AdvancedConnectorPoliciesOnly types.Object `tfsdk:"advanced_connector_policies_only"`
	ContentSecurityPolicy         types.Object `tfsdk:"content_security_policy"`
	AdvancedConnectorPolicies     types.Object `tfsdk:"advanced_connector_policies"`
}

type advancedConnectorPoliciesOnlyModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type advancedConnectorPoliciesModel struct {
	AllowedConnectors types.List `tfsdk:"allowed_connectors"`
}

type allowedConnectorModel struct {
	ConnectorId    types.String `tfsdk:"connector_id"`
	ActionsMode    types.String `tfsdk:"actions_mode"`
	AllowedActions types.List   `tfsdk:"allowed_actions"`
}

type contentSecurityPolicyModel struct {
	Enabled                    types.Bool   `tfsdk:"enabled"`
	EnabledForCanvas           types.Bool   `tfsdk:"enabled_for_canvas"`
	EnabledForCodeApps         types.Bool   `tfsdk:"enabled_for_code_apps"`
	ReportUri                  types.String `tfsdk:"report_uri"`
	ReportingEndpoint          types.String `tfsdk:"reporting_endpoint"`
	Configuration              types.Object `tfsdk:"configuration"`
	ConfigurationForCanvas     types.Object `tfsdk:"configuration_for_canvas"`
	ConfigurationForCodeApps   types.Object `tfsdk:"configuration_for_code_apps"`
}

type cspConfigurationModel struct {
	ImgSrc        types.List `tfsdk:"img_src"`
	StyleSrc      types.List `tfsdk:"style_src"`
	FormAction    types.List `tfsdk:"form_action"`
	FrameSrc      types.List `tfsdk:"frame_src"`
	ConnectSrc    types.List `tfsdk:"connect_src"`
	FontSrc       types.List `tfsdk:"font_src"`
	ScriptSrc     types.List `tfsdk:"script_src"`
	FrameAncestor types.List `tfsdk:"frame_ancestor"`
	StrictCsp     types.Bool `tfsdk:"strict_csp"`
}

type cspConfigurationCodeAppsModel struct {
	ImgSrc        types.List `tfsdk:"img_src"`
	StyleSrc      types.List `tfsdk:"style_src"`
	FormAction    types.List `tfsdk:"form_action"`
	FrameSrc      types.List `tfsdk:"frame_src"`
	ConnectSrc    types.List `tfsdk:"connect_src"`
	FontSrc       types.List `tfsdk:"font_src"`
	ScriptSrc     types.List `tfsdk:"script_src"`
	FrameAncestor types.List `tfsdk:"frame_ancestor"`
}
