// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package dlp_policy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dlpPolicyModelDto struct {
	Name                                 string                                  `json:"name"`
	DisplayName                          string                                  `json:"displayName"`
	DefaultConnectorsClassification      string                                  `json:"defaultConnectorsClassification"`
	EnvironmentType                      string                                  `json:"environmentType"`
	ETag                                 string                                  `json:"etag"`
	CreatedBy                            string                                  `json:"createdBy"`
	CreatedTime                          string                                  `json:"createdTime"`
	LastModifiedBy                       string                                  `json:"lastModifiedBy"`
	LastModifiedTime                     string                                  `json:"lastModifiedTime"`
	Environments                         []dlpEnvironmentDto                     `json:"environments"`
	ConnectorGroups                      []dlpConnectorGroupsModelDto            `json:"connectorGroups"`
	CustomConnectorUrlPatternsDefinition []dlpConnectorUrlPatternsDefinitionDto  `json:"customConnectorUrlPatternsDefinition,omitempty"`
	ConnectorConfigurationsDefinition    dlpConnectorConfigurationsDefinitionDto `json:"connectorConfigurationsDefinition,omitempty"`
}

type dlpPolicyDto struct {
	PolicyDefinition                     dlpPolicyDefinitionDto                   `json:"policyDefinition"`
	ConnectorConfigurationsDefinition    *dlpConnectorConfigurationsDefinitionDto `json:"connectorConfigurationsDefinition,omitempty"`
	CustomConnectorUrlPatternsDefinition dlpConnectorUrlPatternsDefinitionDto     `json:"customConnectorUrlPatternsDefinition"`
}

type dlpPolicyDefinitionDto struct {
	Name                            string                  `json:"name,omitempty"`
	DisplayName                     string                  `json:"displayName"`
	DefaultConnectorsClassification string                  `json:"defaultConnectorsClassification"`
	EnvironmentType                 string                  `json:"environmentType"`
	Environments                    []dlpEnvironmentDto     `json:"environments"`
	ConnectorGroups                 []dlpConnectorGroupsDto `json:"connectorGroups"`
	ETag                            string                  `json:"etag,omitempty"`
	CreatedBy                       dlpPolicyLastActionDto  `json:"createdBy,omitempty"`
	CreatedTime                     string                  `json:"createdTime,omitempty"`
	LastModifiedBy                  dlpPolicyLastActionDto  `json:"lastModifiedBy,omitempty"`
	LastModifiedTime                string                  `json:"lastModifiedTime,omitempty"`
}

type dlpPolicyDefinitionDtoArray struct {
	Value []dlpPolicyDto `json:"value"`
}

type dlpPolicyLastActionDto struct {
	DisplayName string `json:"displayName"`
}

type dlpConnectorConfigurationsDefinitionDto struct {
	ConnectorActionConfigurations []dlpConnectorActionConfigurationsDto `json:"connectorActionConfigurations,omitempty"`
	EndpointConfigurations        []dlpEndpointConfigurationsDto        `json:"endpointConfigurations,omitempty"`
}

type dlpConnectorActionConfigurationsDto struct {
	ConnectorId                        string             `json:"connectorId"`
	DefaultConnectorActionRuleBehavior string             `json:"defaultConnectorActionRuleBehavior"`
	ActionRules                        []dlpActionRuleDto `json:"actionRules"`
}

type dlpEndpointConfigurationsDto struct {
	ConnectorId   string               `json:"connectorId"`
	EndpointRules []dlpEndpointRuleDto `json:"endpointRules"`
}

type dlpConnectorUrlPatternsDefinitionDto struct {
	Rules []dlpConnectorUrlPatternsRuleDto `json:"rules"`
}

type dlpConnectorUrlPatternsRuleDto struct {
	Order                       int64  `json:"order"`
	ConnectorRuleClassification string `json:"customConnectorRuleClassification"`
	Pattern                     string `json:"pattern"`
}

type dlpEnvironmentDto struct {
	Name string `json:"name"`
	Id   string `json:"id"`   // $"/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{x.Name}",.
	Type string `json:"type"` // "Microsoft.BusinessAppPlatform/scopes/environments".
}

type dlpConnectorGroupsDto struct {
	Classification string            `json:"classification"`
	Connectors     []dlpConnectorDto `json:"connectors"`
}

type dlpConnectorGroupsModelDto struct {
	Classification string                 `json:"classification"`
	Connectors     []dlpConnectorModelDto `json:"connectors"`
}

type dlpConnectorModelDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

	DefaultActionRuleBehavior string
	ActionRules               []dlpActionRuleDto
	EndpointRules             []dlpEndpointRuleDto
}

type dlpConnectorDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type dlpEndpointRuleDto struct {
	Order    int64  `json:"order"`
	Behavior string `json:"behavior"`
	Endpoint string `json:"endpoint"`
}

type dlpActionRuleDto struct {
	ActionId string `json:"actionId"`
	Behavior string `json:"behavior"`
}

type policiesListDataSourceModel struct {
	Timeouts timeouts.Value                            `tfsdk:"timeouts"`
	Id       types.String                              `tfsdk:"id"`
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
	Environments                      []string     `tfsdk:"environments"`
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
	Environments                      []string       `tfsdk:"environments"`
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

type dataLossPreventionPolicyResourceEnvironmentsModel struct {
	Name types.String `tfsdk:"name"`
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
