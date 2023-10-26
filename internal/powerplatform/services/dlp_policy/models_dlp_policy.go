package powerplatform

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DlpPolicyModelDto struct {
	Name                                 string                                  `json:"name"`
	DisplayName                          string                                  `json:"displayName"`
	DefaultConnectorsClassification      string                                  `json:"defaultConnectorsClassification"`
	EnvironmentType                      string                                  `json:"environmentType"`
	ETag                                 string                                  `json:"etag"`
	CreatedBy                            string                                  `json:"createdBy"`
	CreatedTime                          string                                  `json:"createdTime"`
	LastModifiedBy                       string                                  `json:"lastModifiedBy"`
	LastModifiedTime                     string                                  `json:"lastModifiedTime"`
	Environments                         []DlpEnvironmentDto                     `json:"environments"`
	ConnectorGroups                      []DlpConnectorGroupsModelDto            `json:"connectorGroups"`
	CustomConnectorUrlPatternsDefinition []DlpConnectorUrlPatternsDefinitionDto  `json:"customConnectorUrlPatternsDefinition,omitempty"`
	ConnectorConfigurationsDefinition    DlpConnectorConfigurationsDefinitionDto `json:"connectorConfigurationsDefinition,omitempty"`
}

type DlpPolicyDto struct {
	PolicyDefinition                     DlpPolicyDefinitionDto                   `json:"policyDefinition"`
	ConnectorConfigurationsDefinition    *DlpConnectorConfigurationsDefinitionDto `json:"connectorConfigurationsDefinition,omitempty"`
	CustomConnectorUrlPatternsDefinition DlpConnectorUrlPatternsDefinitionDto     `json:"customConnectorUrlPatternsDefinition"`
	//ExemptResourcesDefinition            DlpExemptResourcesDefinitionDto         `json:"exemptResourcesDefinition"`
	//ScopeDefinition                      DlpScopeDefinitionDto                   `json:"scopeDefinition"`
}

type DlpPolicyDefinitionDto struct {
	Name                            string                  `json:"name,omitempty"`
	DisplayName                     string                  `json:"displayName"`
	DefaultConnectorsClassification string                  `json:"defaultConnectorsClassification"`
	EnvironmentType                 string                  `json:"environmentType"`
	Environments                    []DlpEnvironmentDto     `json:"environments"`
	ConnectorGroups                 []DlpConnectorGroupsDto `json:"connectorGroups"`
	ETag                            string                  `json:"etag,omitempty"`
	CreatedBy                       DlpPolicyLastActionDto  `json:"createdBy,omitempty"`
	CreatedTime                     string                  `json:"createdTime,omitempty"`
	LastModifiedBy                  DlpPolicyLastActionDto  `json:"lastModifiedBy,omitempty"`
	LastModifiedTime                string                  `json:"lastModifiedTime,omitempty"`
}

type DlpPolicyDefinitionDtoArray struct {
	Value []DlpPolicyDto `json:"value"`
}

type DlpPolicyLastActionDto struct {
	DisplayName string `json:"displayName"`
}

type DlpConnectorConfigurationsDefinitionDto struct {
	ConnectorActionConfigurations []DlpConnectorActionConfigurationsDto `json:"connectorActionConfigurations,omitempty"`
	EndpointConfigurations        []DlpEndpointConfigurationsDto        `json:"endpointConfigurations,omitempty"`
}

type DlpConnectorActionConfigurationsDto struct {
	ConnectorId                        string             `json:"connectorId"`
	DefaultConnectorActionRuleBehavior string             `json:"defaultConnectorActionRuleBehavior"`
	ActionRules                        []DlpActionRuleDto `json:"actionRules"`
}

type DlpEndpointConfigurationsDto struct {
	ConnectorId   string               `json:"connectorId"`
	EndpointRules []DlpEndpointRuleDto `json:"endpointRules"`
}

type DlpConnectorUrlPatternsDefinitionDto struct {
	Rules []DlpConnectorUrlPatternsRuleDto `json:"rules"`
}

type DlpConnectorUrlPatternsRuleDto struct {
	Order                       int64  `json:"order"`
	ConnectorRuleClassification string `json:"customConnectorRuleClassification"`
	Pattern                     string `json:"pattern"`
}

type DlpEnvironmentDto struct {
	Name string `json:"name"`
	Id   string `json:"id"`   //$"/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/{x.Name}",
	Type string `json:"type"` //"Microsoft.BusinessAppPlatform/scopes/environments"
}

type DlpConnectorGroupsDto struct {
	Classification string            `json:"classification"`
	Connectors     []DlpConnectorDto `json:"connectors"`
}

type DlpConnectorGroupsModelDto struct {
	Classification string                 `json:"classification"`
	Connectors     []DlpConnectorModelDto `json:"connectors"`
}

type DlpConnectorModelDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

	DefaultActionRuleBehavior string
	ActionRules               []DlpActionRuleDto
	EndpointRules             []DlpEndpointRuleDto
}

type DlpConnectorDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type DlpEndpointRuleDto struct {
	Order    int64  `json:"order"`
	Behavior string `json:"behavior"`
	Endpoint string `json:"endpoint"`
}

type DlpActionRuleDto struct {
	ActionId string `json:"actionId"`
	Behavior string `json:"behavior"`
}

type PoliciesListDataSourceModel struct {
	Id       types.String                            `tfsdk:"id"`
	Policies []DataLossPreventionPolicyResourceModel `tfsdk:"policies"`
}

type DataLossPreventionPolicyResourceModel struct {
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

type DataLossPreventionPolicyResourceCustomConnectorPattern struct {
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

type DataLossPreventionPolicyResourceEnvironmentsModel struct {
	Name types.String `tfsdk:"name"`
}

var environmentSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name": types.StringType,
	},
}

type DataLossPreventionPolicyResourceConnectorModel struct {
	Id                        types.String                                                 `tfsdk:"id"`
	DefaultActionRuleBehavior types.String                                                 `tfsdk:"default_action_rule_behavior"`
	ActionRules               []DataLossPreventionPolicyResourceConnectorActionRuleModel   `tfsdk:"action_rules"`
	EndpointRules             []DataLossPreventionPolicyResourceConnectorEndpointRuleModel `tfsdk:"endpoint_rules"`
}

type DataLossPreventionPolicyResourceConnectorEndpointRuleModel struct {
	Order    types.Int64  `tfsdk:"order"`
	Behavior types.String `tfsdk:"behavior"`
	Endpoint types.String `tfsdk:"endpoint"`
}

type DataLossPreventionPolicyResourceConnectorActionRuleModel struct {
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
