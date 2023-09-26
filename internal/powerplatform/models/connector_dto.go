package powerplatform_models

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
	Unblockable bool
}

type ConnectorDtoArray struct {
	Value []ConnectorDto `json:"value"`
}

type UnblockableConnectorDto struct {
	Id       string                          `json:"id"`
	Metadata UnblockableConnectorMetadataDto `json:"metadata"`
}

type UnblockableConnectorMetadataDto struct {
	Unblockable bool `json:"unblockable"`
}

type VirtualConnectorDto struct {
	Id       string                      `json:"id"`
	Metadata VirtualConnectorMetadataDto `json:"metadata"`
}

type VirtualConnectorMetadataDto struct {
	VirtualConnector bool   `json:"virtualConnector"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	DisplayName      string `json:"displayName"`
}
