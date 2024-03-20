// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

type ApplicationArrayDto struct {
	Value []ApplicationDto `json:"value"`
}

type ApplicationDto struct {
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

type ApplicationDeleteDto struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ApplicationCreateDto struct {
	Location string `json:"location"`
}

type ApplicationLifecycleCreatedDto struct {
	Name       string                                   `json:"name"`
	Properties ApplicationLifecycleCreatedPropertiesDto `json:"properties"`
}

type ApplicationLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
}

type ApplicationBapi struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ApplicationPropertiesBapi struct {
	TenantID    string `json:"tenantId"`
	DisplayName string `json:"displayName"`
}

type LinkedApplicationMetadataBapi struct {
	Version string `json:"version"`
}

type StatesApplicationBapi struct {
	Management StatesManagementApplicationBapi `json:"management"`
}

type StatesManagementApplicationBapi struct {
	Id string `json:"id"`
}

type ApplicationDtoArray struct {
	Value []ApplicationDto `json:"value"`
}

type ApplicationCreateBapi struct {
	Location string `json:"location"`

	Properties ApplicationPropertiesBapi `json:"properties"`
}

type ApplicationCreatePropertiesBapi struct {
	DisplayName string `json:"displayName"`
}

type ApplicationCreateLinkApplicationMetadataBapi struct {
}

type ApplicationLifecycleDto struct {
	OperationId        string                       `json:"operationId"`
	CreatedDateTime    string                       `json:"createdDateTime"`
	LastActionDateTime string                       `json:"lastActionDateTime"`
	Status             string                       `json:"status"`
	StatusMessage      string                       `json:"statusMessage"`
	Error              ApplicationLifecycleErrorDto `json:"error"`
}

type ApplicationLifecycleErrorDto struct {
	ErrorName  string `json:"errorName"`
	ErrorCode  int    `json:"errorCode"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	StatusCode int    `json:"statusCode"`
	Source     string `json:"source"`
}