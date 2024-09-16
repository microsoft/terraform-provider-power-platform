// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

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
	SolutionValidationResults []SolutionValidationResults     `json:"SolutionValidationResults"`
	MissingDependencies       []MissingDependenciesDto        `json:"MissingDependencies"`
	SolutionDetails           StageSolutionSolutionDetailsDto `json:"SolutionDetails"`
}

type SolutionValidationResults struct {
	SolutionValidationResultType string `json:"SolutionValidationResultType"`
	ErrorCode                    int    `json:"ErrorCode"`
	AdditionalInfo               string `json:"AdditionalInfo"`
	Message                      string `json:"Message"`
}

type MissingDependenciesDto struct {
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
