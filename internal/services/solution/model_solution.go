// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

type solutionSettings struct {
	EnvironmentVariables []settingsEnvironmentVariable  `json:"environmentvariables"`
	ConnectionReferences []settingsConnectionReferences `json:"connectionreferences"`
}

type settingsEnvironmentVariable struct {
	SchemaName string `json:"schemaname"`
	Value      string `json:"value"`
}

type settingsConnectionReferences struct {
	LogicalName  string `json:"logicalname"`
	ConnectionId string `json:"connectionid"`
	ConnectorId  string `json:"connectorid"`
}

type solutionDto struct {
	Id            string `json:"solutionid"`
	EnvironmentId string `json:"environment_id"`
	Name          string `json:"uniquename"`
	DisplayName   string `json:"friendlyname"`
	IsManaged     bool   `json:"ismanaged"`
	CreatedTime   string `json:"createdon"`
	Version       string `json:"version"`
	ModifiedTime  string `json:"modifiedon"`
	InstallTime   string `json:"installedon"`
}

type solutionArrayDto struct {
	Value []solutionDto `json:"value"`
}

type stageSolutionImportDto struct {
	CustomizationFile string `json:"CustomizationFile"`
}

type stageSolutionImportResponseDto struct {
	StageSolutionResults stageSolutionImportResultResponseDto `json:"StageSolutionResults"`
}

type stageSolutionImportResultResponseDto struct {
	StageSolutionUploadId     string                          `json:"StageSolutionUploadId"`
	StageSolutionStatus       string                          `json:"StageSolutionStatus"`
	SolutionValidationResults []validationResults             `json:"SolutionValidationResults"`
	MissingDependencies       []missingDependenciesDto        `json:"MissingDependencies"`
	SolutionDetails           stageSolutionSolutionDetailsDto `json:"SolutionDetails"`
}

type validationResults struct {
	SolutionValidationResultType string `json:"SolutionValidationResultType"`
	ErrorCode                    int    `json:"ErrorCode"`
	AdditionalInfo               string `json:"AdditionalInfo"`
	Message                      string `json:"Message"`
}

type missingDependenciesDto struct {
	RequiredComponentSchemaName         string `json:"RequiredComponentSchemaName"`
	RequiredComponentDisplayName        string `json:"RequiredComponentDisplayName"`
	RequiredComponentParentSchemaName   string `json:"RequiredComponentParentSchemaName"`
	RequiredComponentParentDisplayName  string `json:"RequiredComponentParentDisplayName"`
	RequiredComponentId                 string `json:"RequiredComponentId"`
	RequiredSolutionName                string `json:"RequiredSolutionName"`
	RequiredComponentType               string `json:"RequiredComponentType"`
	DependentComponentSchemaName        string `json:"DependentComponentSchemaName"`
	DependentComponentDisplayName       string `json:"DependentComponentDisplayName"`
	DependentComponentParentSchemaName  string `json:"DependentComponentParentSchemaName"`
	DependentComponentParentDisplayName string `json:"DependentComponentParentDisplayName"`
	DependentComponentType              string `json:"DependentComponentType"`
	DependentComponentId                string `json:"DependentComponentId"`
}

type stageSolutionSolutionDetailsDto struct {
	SolutionUniqueName   string `json:"SolutionUniqueName"`
	SolutionFriendlyName string `json:"SolutionFriendlyName"`
	IsManaged            bool   `json:"IsManaged"`
	SolutionVersion      string `json:"SolutionVersion"`
}

type importSolutionDto struct {
	PublishWorkflows                 bool                                `json:"PublishWorkflows"`
	OverwriteUnmanagedCustomizations bool                                `json:"OverwriteUnmanagedCustomizations"`
	ComponentParameters              []any                               `json:"ComponentParameters"`
	SolutionParameters               importSolutionSolutionParametersDto `json:"SolutionParameters"`
}

type importSolutionResponseDto struct {
	ImportJobKey     string `json:"ImportJobKey"`
	AsyncOperationId string `json:"AsyncOperationId"`
}

type importSolutionSolutionParametersDto struct {
	StageSolutionUploadId string `json:"StageSolutionUploadId"`
}

type importSolutionConnectionReferencesDto struct {
	Type                           string `json:"@odata.type"`
	ConnectionReferenceDisplayName string `json:"connectionreferencedisplayname"`
	ConnectionReferenceLogicalName string `json:"connectionreferencelogicalname"`
	Description                    string `json:"description"`
	ConnectorId                    string `json:"connectorid"`
	ConnectionId                   string `json:"connectionid"`
}

type importSolutionEnvironmentVariablesDto struct {
	Type       string `json:"@odata.type"`
	SchemaName string `json:"schemaname"`
	Value      string `json:"value"`
}

type asyncSolutionPullResponseDto struct {
	AsyncOperationId string `json:"AsyncOperationId"`
	CreatedOn        string `json:"createdon"`
	CompletedOn      string `json:"completedon"`
}

type validateSolutionImportResponseDto struct {
	SolutionOperationResult validateSolutionImportResponseSolutionOperationResultDto `json:"SolutionOperationResult"`
}

type validateSolutionImportResponseSolutionOperationResultDto struct {
	Status          string `json:"Status"`
	WarningMessages []any  `json:"WarningMessages"`
	ErrorMessages   []any  `json:"ErrorMessages"`
}

type environmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties environmentIdPropertiesDto `json:"properties"`
}

type environmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata linkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type linkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
