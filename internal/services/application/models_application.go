// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package application

type TenantApplicationDto struct {
	ApplicationDescription string                         `json:"applicationDescription"`
	ApplicationId          string                         `json:"applicationId"`
	ApplicationName        string                         `json:"applicationName"`
	ApplicationVisibility  string                         `json:"applicationVisibility"`
	CatalogVisibility      string                         `json:"catalogVisibility"`
	LastError              *TenantApplicationErrorDetails `json:"errorDetails,omitempty"`
	LearnMoreUrl           string                         `json:"learnMoreUrl"`
	LocalizedDescription   string                         `json:"localizedDescription"`
	LocalizedName          string                         `json:"localizedName"`
	PublisherId            string                         `json:"publisherId"`
	PublisherName          string                         `json:"publisherName"`
	UniqueName             string                         `json:"uniqueName"`
}

type TenantApplicationErrorDetails struct {
	ErrorCode  string `json:"errorCode"`
	ErrorName  string `json:"errorName"`
	Message    string `json:"message"`
	Source     string `json:"source"`
	StatusCode int64  `json:"statusCode"`
	Type       string `json:"type"`
}

type TenantApplicationArrayDto struct {
	Value []TenantApplicationDto `json:"value"`
}

type EnvironmentApplicationArrayDto struct {
	Value []EnvironmentApplicationDto `json:"value"`
}

type EnvironmentApplicationDto struct {
	ApplicationId         string `json:"applicationId"`
	Name                  string `json:"applicationName"`
	UniqueName            string `json:"uniqueName"`
	Version               string `json:"version"`
	Description           string `json:"localizedDescription"`
	PublisherId           string `json:"publisherId"`
	PublisherName         string `json:"publisherName"`
	LearnMoreUrl          string `json:"learnMoreUrl"`
	State                 string `json:"state"`
	ApplicationVisibility string `json:"applicationVisibility"`
}

type EnvironmentApplicationDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type EnvironmentApplicationCreateDto struct {
	Location string `json:"location"`
}

type EnvironmentApplicationLifecycleCreatedDto struct {
	Name       string                                              `json:"name"`
	Properties EnvironmentApplicationLifecycleCreatedPropertiesDto `json:"properties"`
}

type EnvironmentApplicationLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}

type EnvironmentApplicationBapi struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type EnvironmentApplicationPropertiesBapi struct {
	TenantID    string `json:"tenantId"`
	DisplayName string `json:"displayName"`
}

type LinkedEnvironmentApplicationMetadataBapi struct {
	Version string `json:"version"`
}

type StatesEnvironmentApplicationBapi struct {
	Management StatesManagementApplicationBapi `json:"management"`
}

type StatesManagementApplicationBapi struct {
	Id string `json:"id"`
}

type EnvironmentApplicationDtoArray struct {
	Value []EnvironmentApplicationDto `json:"value"`
}

type EnvironmentApplicationCreateBapi struct {
	Location string `json:"location"`

	Properties EnvironmentApplicationPropertiesBapi `json:"properties"`
}

type EnvironmentApplicationCreatePropertiesBapi struct {
	DisplayName string `json:"displayName"`
}

type EnvironmentApplicationCreateLinkApplicationMetadataBapi struct {
}

type EnvironmentApplicationLifecycleDto struct {
	OperationId        string                                  `json:"operationId"`
	CreatedDateTime    string                                  `json:"createdDateTime"`
	LastActionDateTime string                                  `json:"lastActionDateTime"`
	Status             string                                  `json:"status"`
	StatusMessage      string                                  `json:"statusMessage"`
	Error              EnvironmentApplicationLifecycleErrorDto `json:"error"`
}

type EnvironmentApplicationLifecycleErrorDto struct {
	ErrorName  string `json:"errorName"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode int    `json:"statusCode"`
	Source     string `json:"source"`
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
