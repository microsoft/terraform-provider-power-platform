// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package dlp_policy

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

type dlpPolicyDefinitionArrayDto struct {
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
