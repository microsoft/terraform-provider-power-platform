package powerplatform_models

type DlpPolicyModel struct {
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
	ConnectorGroups                      []DlpConnectorGroupsModel               `json:"connectorGroups"`
	CustomConnectorUrlPatternsDefinition []DlpConnectorUrlPatternsDefinitionDto  `json:"customConnectorUrlPatternsDefinition,omitempty"`
	ConnectorConfigurationsDefinition    DlpConnectorConfigurationsDefinitionDto `json:"connectorConfigurationsDefinition"`
}

type DlpPolicyDto struct {
	PolicyDefinition                     DlpPolicyDefinitionDto                  `json:"policyDefinition"`
	ConnectorConfigurationsDefinition    DlpConnectorConfigurationsDefinitionDto `json:"connectorConfigurationsDefinition,omitempty"`
	CustomConnectorUrlPatternsDefinition DlpConnectorUrlPatternsDefinitionDto    `json:"customConnectorUrlPatternsDefinition"`
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

type DlpPolicyLastActionDto struct {
	DisplayName string `json:"displayName"`
}

type DlpConnectorConfigurationsDefinitionDto struct {
	ConnectorActionConfigurations []DlpConnectorActionConfigurationsDto `json:"connectorActionConfigurations"`
	EndpointConfigurations        []DlpEndpointConfigurationsDto        `json:"endpointConfigurations"`
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

type DlpConnectorGroupsModel struct {
	Classification string              `json:"classification"`
	Connectors     []DlpConnectorModel `json:"connectors"`
}

type DlpConnectorModel struct {
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
