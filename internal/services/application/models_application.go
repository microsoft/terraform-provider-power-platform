// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package application

type tenantApplicationDto struct {
	ApplicationDescription string                         `json:"applicationDescription"`
	ApplicationId          string                         `json:"applicationId"`
	ApplicationName        string                         `json:"applicationName"`
	ApplicationVisibility  string                         `json:"applicationVisibility"`
	CatalogVisibility      string                         `json:"catalogVisibility"`
	LastError              *tenantApplicationErrorDetails `json:"errorDetails,omitempty"`
	LearnMoreUrl           string                         `json:"learnMoreUrl"`
	LocalizedDescription   string                         `json:"localizedDescription"`
	LocalizedName          string                         `json:"localizedName"`
	PublisherId            string                         `json:"publisherId"`
	PublisherName          string                         `json:"publisherName"`
	UniqueName             string                         `json:"uniqueName"`
}

type tenantApplicationErrorDetails struct {
	ErrorCode  string `json:"errorCode"`
	ErrorName  string `json:"errorName"`
	Message    string `json:"message"`
	Source     string `json:"source"`
	StatusCode int64  `json:"statusCode"`
	Type       string `json:"type"`
}

type tenantApplicationArrayDto struct {
	Value []tenantApplicationDto `json:"value"`
}

type environmentApplicationArrayDto struct {
	Value []environmentApplicationDto `json:"value"`
}

type environmentApplicationDto struct {
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

type environmentApplicationDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type environmentApplicationCreateDto struct {
	Location string `json:"location"`
}

type environmentApplicationLifecycleCreatedDto struct {
	Name       string                                              `json:"name"`
	Properties environmentApplicationLifecycleCreatedPropertiesDto `json:"properties"`
}

type environmentApplicationLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}

type environmentApplicationBapi struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type environmentApplicationPropertiesBapi struct {
	TenantID    string `json:"tenantId"`
	DisplayName string `json:"displayName"`
}

type linkedEnvironmentApplicationMetadataBapi struct {
	Version string `json:"version"`
}

type statesEnvironmentApplicationBapi struct {
	Management statesManagementApplicationBapi `json:"management"`
}

type statesManagementApplicationBapi struct {
	Id string `json:"id"`
}

type environmentApplicationDtoArray struct {
	Value []environmentApplicationDto `json:"value"`
}

type environmentApplicationCreateBapi struct {
	Location string `json:"location"`

	Properties environmentApplicationPropertiesBapi `json:"properties"`
}

type environmentApplicationCreatePropertiesBapi struct {
	DisplayName string `json:"displayName"`
}

type environmentApplicationCreateLinkApplicationMetadataBapi struct {
}

type environmentApplicationLifecycleDto struct {
	OperationId        string                                  `json:"operationId"`
	CreatedDateTime    string                                  `json:"createdDateTime"`
	LastActionDateTime string                                  `json:"lastActionDateTime"`
	Status             string                                  `json:"status"`
	StatusMessage      string                                  `json:"statusMessage"`
	Error              environmentApplicationLifecycleErrorDto `json:"error"`
}

type environmentApplicationLifecycleErrorDto struct {
	ErrorName  string `json:"errorName"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode int    `json:"statusCode"`
	Source     string `json:"source"`
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
