package powerplatform_models

type SolutionSettings struct {
	EnvironmentVariables []SolutionSettingsEnvironmentVariable  `json:"environmentvariables"`
	ConnectionReferences []SolutionSettingsConnectionReferences `json:"connectionreferences"`
}

type SolutionSettingsEnvironmentVariable struct {
	SchemaName string `json:"schemaname"`
	Value      string `json:"value"`
}

type SolutionSettingsConnectionReferences struct {
	LogicalName  string `json:"logicalname"`
	ConnectionId string `json:"connectionid"`
	ConnectorId  string `json:"connectorid"`
}

type SolutionDto struct {
	Id              string `json:"solutionid"`
	EnvironmentName string `json:"environment_name"`
	Name            string `json:"uniquename"`
	DisplayName     string `json:"friendlyname"`
	IsManaged       bool   `json:"ismanaged"`
	CreatedTime     string `json:"createdon"`
	Version         string `json:"version"`
	ModifiedTime    string `json:"modifiedon"`
	InstallTime     string `json:"installedon"`
}

type SolutionDtoArray struct {
	Value []SolutionDto `json:"value"`
}

type StageSolutionImportDto struct {
	CustomizationFile string `json:"CustomizationFile"`
}

type StageSolutionImportResponseDto struct {
	StageSolutionResults StageSolutionImportResultResponseDto `json:"StageSolutionResults"`
}

type StageSolutionImportResultResponseDto struct {
	StageSolutionUploadId     string                          `json:"StageSolutionUploadId"`
	StageSolutionStatus       string                          `json:"StageSolutionStatus"`
	SolutionValidationResults []string                        `json:"SolutionValidationResults"`
	MissingDependencies       []string                        `json:"MissingDependencies"`
	SolutionDetails           StageSolutionSolutionDetailsDto `json:"SolutionDetails"`
}

type StageSolutionSolutionDetailsDto struct {
	SolutionUniqueName   string `json:"SolutionUniqueName"`
	SolutionFriendlyName string `json:"SolutionFriendlyName"`
	IsManaged            bool   `json:"IsManaged"`
	SolutionVersion      string `json:"SolutionVersion"`
}

type ImportSolutionDto struct {
	PublishWorkflows                 bool                                `json:"PublishWorkflows"`
	OverwriteUnmanagedCustomizations bool                                `json:"OverwriteUnmanagedCustomizations"`
	ComponentParameters              []interface{}                       `json:"ComponentParameters"`
	SolutionParameters               ImportSolutionSolutionParametersDto `json:"SolutionParameters"`
}

type ImportSolutionResponseDto struct {
	ImportJobKey     string `json:"ImportJobKey"`
	AsyncOperationId string `json:"AsyncOperationId"`
}

type ImportSolutionSolutionParametersDto struct {
	StageSolutionUploadId string `json:"StageSolutionUploadId"`
}

type ImportSolutionConnectionReferencesDto struct {
	Type                           string `json:"@odata.type"`
	ConnectionReferenceDisplayName string `json:"connectionreferencedisplayname"`
	ConnectionReferenceLogicalName string `json:"connectionreferencelogicalname"`
	Description                    string `json:"description"`
	ConnectorId                    string `json:"connectorid"`
	ConnectionId                   string `json:"connectionid"`
}

type ImportSolutionEnvironmentVariablesDto struct {
	Type       string `json:"@odata.type"`
	SchemaName string `json:"schemaname"`
	Value      string `json:"value"`
}

type AsyncSolutionPullResponseDto struct {
	AsyncOperationId string `json:"AsyncOperationId"`
	CreatedOn        string `json:"createdon"`
	CompletedOn      string `json:"completedon"`
}

type ValidateSolutionImportResponseDto struct {
	SolutionOperationResult ValidateSolutionImportResponseSolutionOperationResultDto `json:"SolutionOperationResult"`
}

type ValidateSolutionImportResponseSolutionOperationResultDto struct {
	Status          string        `json:"Status"`
	WarningMessages []interface{} `json:"WarningMessages"`
	ErrorMessages   []interface{} `json:"ErrorMessages"`
}
