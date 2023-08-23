package powerplatform_bapi

type DlpPolicyDto struct {
	Name                                 string                                 `json:"Name"`
	DisplayName                          string                                 `json:"DisplayName"`
	DefaultConnectorsClassification      string                                 `json:"DefaultConnectorsClassification"`
	EnvironmentType                      string                                 `json:"EnvironmentType"`
	ETag                                 string                                 `json:"Etag"`
	CreatedBy                            string                                 `json:"CreatedBy"`
	CreatedTime                          string                                 `json:"CreatedTime"`
	LastModifiedBy                       string                                 `json:"LastModifiedBy"`
	LastModifiedTime                     string                                 `json:"LastModifiedTime"`
	Environments                         []DlpEnvironmentDto                    `json:"Environments"`
	ConnectorGroups                      []DlpConnectorGroupsDto                `json:"ConnectorGroups"`
	CustomConnectorUrlPatternsDefinition []DlpConnectorUrlPatternsDefinitionDto `json:"CustomConnectorUrlPatternsDefinition"`
}

type DlpConnectorUrlPatternsDefinitionDto struct {
	Order                       int64  `json:"Order"`
	ConnectorRuleClassification string `json:"ConnectorRuleClassification"`
	Pattern                     string `json:"Pattern"`
}

type DlpEnvironmentDto struct {
	Name string `json:"Name"`
	Id   string `json:"Id"`
	Type string `json:"Type"`
}

type DlpConnectorGroupsDto struct {
	Classification string            `json:"Classification"`
	Connectors     []DlpConnectorDto `json:"Connectors"`
}

type DlpConnectorDto struct {
	Id                        string               `json:"Id"`
	Name                      string               `json:"Name"`
	Type                      string               `json:"Type"`
	DefaultActionRuleBehavior string               `json:"DefaultActionRuleBehavior"`
	ActionRules               []DlpActionRuleDto   `json:"ActionRules"`
	EndpointRules             []DlpEndpointRuleDto `json:"EndpointRules"`
}

type DlpEndpointRuleDto struct {
	Order    int64  `json:"Order"`
	Behavior string `json:"Behavior"`
	Endpoint string `json:"Endpoint"`
}

type DlpActionRuleDto struct {
	ActionId string `json:"ActionId"`
	Behavior string `json:"Behavior"`
}
