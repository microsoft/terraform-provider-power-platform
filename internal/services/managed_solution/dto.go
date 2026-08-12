// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managedsolution

import "encoding/xml"

type solutionPackageXML struct {
	XMLName          xml.Name                   `xml:"ImportExportXml"`
	SolutionManifest solutionPackageManifestXML `xml:"SolutionManifest"`
}

type solutionPackageManifestXML struct {
	UniqueName          string                         `xml:"UniqueName"`
	Version             string                         `xml:"Version"`
	Managed             int                            `xml:"Managed"`
	LocalizedNames      []solutionLocalizedNameXML     `xml:"LocalizedNames>LocalizedName"`
	MissingDependencies solutionMissingDependenciesXML `xml:"MissingDependencies"`
}

type solutionLocalizedNameXML struct {
	Description string `xml:"description,attr"`
}

type solutionMissingDependenciesXML struct {
	MissingDependencies []solutionMissingDependencyXML `xml:"MissingDependency"`
}

type solutionMissingDependencyXML struct {
	Required []solutionRequiredDependencyXML `xml:"Required"`
}

type solutionRequiredDependencyXML struct {
	Type       string `xml:"type,attr"`
	SchemaName string `xml:"schemaName,attr"`
	Solution   string `xml:"solution,attr"`
}

type solutionCustomizationsXML struct {
	XMLName              xml.Name                         `xml:"ImportExportXml"`
	ConnectionReferences []solutionConnectionReferenceXML `xml:"connectionreferences>connectionreference"`
}

type solutionConnectionReferenceXML struct {
	LogicalName string `xml:"connectionreferencelogicalname,attr"`
	ConnectorID string `xml:"connectorid"`
}

type environmentVariableDefinitionXML struct {
	XMLName      xml.Name `xml:"environmentvariabledefinition"`
	SchemaName   string   `xml:"schemaname,attr"`
	DefaultValue *string  `xml:"defaultvalue"`
}

type environmentVariableValuesJSON struct {
	EnvironmentVariableValues environmentVariableValueContainerJSON `json:"environmentvariablevalues"`
}

type environmentVariableValueContainerJSON struct {
	EnvironmentVariableValue *environmentVariableValueJSON `json:"environmentvariablevalue"`
}

type environmentVariableValueJSON struct {
	Value *string `json:"value"`
}

type stageSolutionImportDto struct {
	CustomizationFile string `json:"CustomizationFile"`
}

type stageSolutionImportResponseDto struct {
	StageSolutionResults stageSolutionImportResultResponseDto `json:"StageSolutionResults"`
}

type stageSolutionImportResultResponseDto struct {
	StageSolutionUploadId     string                   `json:"StageSolutionUploadId"`
	StageSolutionStatus       string                   `json:"StageSolutionStatus"`
	SolutionValidationResults []validationResultsDto   `json:"SolutionValidationResults"`
	MissingDependencies       []missingDependenciesDto `json:"MissingDependencies"`
	SolutionDetails           stageSolutionDetailsDto  `json:"SolutionDetails"`
}

type validationResultsDto struct {
	Message string `json:"Message"`
}

type missingDependenciesDto struct {
	RequiredComponentSchemaName string `json:"RequiredComponentSchemaName"`
}

type stageSolutionDetailsDto struct {
	SolutionUniqueName string `json:"SolutionUniqueName"`
}

type importSolutionDto struct {
	PublishWorkflows                 bool                        `json:"PublishWorkflows"`
	OverwriteUnmanagedCustomizations bool                        `json:"OverwriteUnmanagedCustomizations"`
	SkipProductUpdateDependencies    bool                        `json:"SkipProductUpdateDependencies"`
	ComponentParameters              []any                       `json:"ComponentParameters,omitempty"`
	SolutionParameters               importSolutionParametersDto `json:"SolutionParameters"`
}

type importSolutionParametersDto struct {
	StageSolutionUploadId string `json:"StageSolutionUploadId"`
}

type importSolutionResponseDto struct {
	ImportJobKey     string `json:"ImportJobKey"`
	AsyncOperationId string `json:"AsyncOperationId"`
}

type importSolutionConnectionReferenceDto struct {
	Type                           string `json:"@odata.type"`
	ConnectionReferenceDisplayName string `json:"connectionreferencedisplayname"`
	ConnectionReferenceLogicalName string `json:"connectionreferencelogicalname"`
	Description                    string `json:"description"`
	ConnectorId                    string `json:"connectorid"`
	ConnectionId                   string `json:"connectionid"`
}

type asyncSolutionPullResponseDto struct {
	CompletedOn string `json:"completedon"`
}

type validateSolutionImportResponseDto struct {
	SolutionOperationResult validateSolutionImportResultDto `json:"SolutionOperationResult"`
}

type validateSolutionImportResultDto struct {
	Status        string `json:"Status"`
	ErrorMessages []any  `json:"ErrorMessages"`
}

type connectionArrayDto struct {
	Value []connectionDto `json:"value"`
}

type connectionDto struct {
	Name       string                  `json:"name"`
	Properties connectionPropertiesDto `json:"properties"`
}

type connectionPropertiesDto struct {
	ApiId       string `json:"apiId"`
	DisplayName string `json:"displayName"`
}

type environmentVariableValuesResponseDto struct {
	Value []environmentVariableValueDto `json:"value"`
}

type environmentVariableValueDto struct {
	Id         string  `json:"environmentvariablevalueid"`
	SchemaName string  `json:"schemaname"`
	Value      *string `json:"value"`
}

type packageEnvironmentVariable struct {
	SchemaName            string
	HasDefaultValue       bool
	ContainsPackagedValue bool
}

type packageConnectionReference struct {
	LogicalName string
	ConnectorID string
}

type solutionPackage struct {
	UniqueName           string
	DisplayName          string
	Version              string
	IsManaged            bool
	ConnectionReferences map[string]packageConnectionReference
	EnvironmentVariables map[string]packageEnvironmentVariable
	Dependencies         map[string]string
}
