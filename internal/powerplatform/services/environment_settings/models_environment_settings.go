package powerplatform

type EnvironmentSettingsValueDto struct {
	Value []EnvironmentSettingsDto `json:"value"`
}

type EnvironmentSettingsDto struct {
	MaxUploadFileSize int64 `json:"maxuploadfilesize"`
}

type EnvironmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties EnvironmentIdPropertiesDto `json:"properties"`
}

type EnvironmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type LinkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
