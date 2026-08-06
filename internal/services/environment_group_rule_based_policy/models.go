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
	RuleBasedPolicyClient client
}

type environmentGroupRuleBasedPolicyResourceModel struct {
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String   `tfsdk:"id"`
	EnvironmentGroupId types.String   `tfsdk:"environment_group_id"`
	RuleSets           types.Object   `tfsdk:"rule_sets"`
}

type advancedConnectorPoliciesOnlyModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type contentSecurityPolicyModel struct {
	Enabled                  types.Bool   `tfsdk:"enabled"`
	EnabledForCanvas         types.Bool   `tfsdk:"enabled_for_canvas"`
	EnabledForCodeApps       types.Bool   `tfsdk:"enabled_for_code_apps"`
	ReportUri                types.String `tfsdk:"report_uri"`
	ReportingEndpoint        types.String `tfsdk:"reporting_endpoint"`
	Configuration            types.Object `tfsdk:"configuration"`
	ConfigurationForCanvas   types.Object `tfsdk:"configuration_for_canvas"`
	ConfigurationForCodeApps types.Object `tfsdk:"configuration_for_code_apps"`
}
