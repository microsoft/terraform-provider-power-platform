package powerplatform_bapi

type ConnectorDto struct {
	Name       string                 `json:"name"`
	Id         string                 `json:"id"`
	Type       string                 `json:"type"`
	Properties ConnectorPropertiesDto `json:"properties"`
}

type ConnectorPropertiesDto struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Tier        string `json:"tier"`
	Publisher   string `json:"publisher"`
}

type ConnectorDtoArray struct {
	Value []ConnectorDto `json:"value"`
}
