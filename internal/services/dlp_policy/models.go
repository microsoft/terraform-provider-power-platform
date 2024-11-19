// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package dlp_policy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type policiesListDataSourceModel struct {
	Timeouts timeouts.Value                            `tfsdk:"timeouts"`
	Policies []dataLossPreventionPolicyDatasourceModel `tfsdk:"policies"`
}

type dataLossPreventionPolicyDatasourceModel struct {
	Id                                types.String `tfsdk:"id"`
	DisplayName                       types.String `tfsdk:"display_name"`
	DefaultConnectorsClassification   types.String `tfsdk:"default_connectors_classification"`
	EnvironmentType                   types.String `tfsdk:"environment_type"`
	CreatedBy                         types.String `tfsdk:"created_by"`
	CreatedTime                       types.String `tfsdk:"created_time"`
	LastModifiedBy                    types.String `tfsdk:"last_modified_by"`
	LastModifiedTime                  types.String `tfsdk:"last_modified_time"`
	Environments                      types.Set    `tfsdk:"environments"`
	NonBusinessConfidentialConnectors types.Set    `tfsdk:"non_business_connectors"`
	BusinessGeneralConnectors         types.Set    `tfsdk:"business_connectors"`
	BlockedConnectors                 types.Set    `tfsdk:"blocked_connectors"`
	CustomConnectorsPatterns          types.Set    `tfsdk:"custom_connectors_patterns"`
}

type dataLossPreventionPolicyResourceModel struct {
	Timeouts                          timeouts.Value `tfsdk:"timeouts"`
	Id                                types.String   `tfsdk:"id"`
	DisplayName                       types.String   `tfsdk:"display_name"`
	DefaultConnectorsClassification   types.String   `tfsdk:"default_connectors_classification"`
	EnvironmentType                   types.String   `tfsdk:"environment_type"`
	CreatedBy                         types.String   `tfsdk:"created_by"`
	CreatedTime                       types.String   `tfsdk:"created_time"`
	LastModifiedBy                    types.String   `tfsdk:"last_modified_by"`
	LastModifiedTime                  types.String   `tfsdk:"last_modified_time"`
	Environments                      types.Set      `tfsdk:"environments"`
	NonBusinessConfidentialConnectors types.Set      `tfsdk:"non_business_connectors"`
	BusinessGeneralConnectors         types.Set      `tfsdk:"business_connectors"`
	BlockedConnectors                 types.Set      `tfsdk:"blocked_connectors"`
	CustomConnectorsPatterns          types.Set      `tfsdk:"custom_connectors_patterns"`
}

type dataLossPreventionPolicyResourceCustomConnectorPattern struct {
	Order          types.Int64  `tfsdk:"order"`
	HostUrlPattern types.String `tfsdk:"host_url_pattern"`
	DataGroup      types.String `tfsdk:"data_group"`
}

var customConnectorPatternSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":            types.Int64Type,
		"host_url_pattern": types.StringType,
		"data_group":       types.StringType,
	},
}

type dataLossPreventionPolicyResourceConnectorModel struct {
	Id                        types.String                                                 `tfsdk:"id"`
	DefaultActionRuleBehavior types.String                                                 `tfsdk:"default_action_rule_behavior"`
	ActionRules               []dataLossPreventionPolicyResourceConnectorActionRuleModel   `tfsdk:"action_rules"`
	EndpointRules             []dataLossPreventionPolicyResourceConnectorEndpointRuleModel `tfsdk:"endpoint_rules"`
}

type dataLossPreventionPolicyResourceConnectorEndpointRuleModel struct {
	Order    types.Int64  `tfsdk:"order"`
	Behavior types.String `tfsdk:"behavior"`
	Endpoint types.String `tfsdk:"endpoint"`
}

type dataLossPreventionPolicyResourceConnectorActionRuleModel struct {
	ActionId types.String `tfsdk:"action_id"`
	Behavior types.String `tfsdk:"behavior"`
}

var connectorSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                           types.StringType,
		"default_action_rule_behavior": types.StringType,
		"action_rules":                 types.ListType{ElemType: actionRuleListObjectType},
		"endpoint_rules":               types.ListType{ElemType: endpointRuleListObjectType},
	},
}

var endpointRuleListObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":    types.Int64Type,
		"behavior": types.StringType,
		"endpoint": types.StringType,
	},
}

var actionRuleListObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"action_id": types.StringType,
		"behavior":  types.StringType,
	},
}

type DataLossPreventionPolicyResource struct {
	helpers.TypeInfo
	DlpPolicyClient client
}

type DataLossPreventionPolicyDataSource struct {
	helpers.TypeInfo
	DlpPolicyClient client
}
