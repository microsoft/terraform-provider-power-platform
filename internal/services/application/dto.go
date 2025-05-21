// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
package application

type tenantApplicationDto struct {
	ApplicationDescription string                            `json:"applicationDescription"`
	ApplicationId          string                            `json:"applicationId"`
	ApplicationName        string                            `json:"applicationName"`
	ApplicationVisibility  string                            `json:"applicationVisibility"`
	CatalogVisibility      string                            `json:"catalogVisibility"`
	LastError              *tenantApplicationErrorDetailsDto `json:"errorDetails,omitempty"`
	LearnMoreUrl           string                            `json:"learnMoreUrl"`
	LocalizedDescription   string                            `json:"localizedDescription"`
	LocalizedName          string                            `json:"localizedName"`
	PublisherId            string                            `json:"publisherId"`
	PublisherName          string                            `json:"publisherName"`
	UniqueName             string                            `json:"uniqueName"`
}

type tenantApplicationErrorDetailsDto struct {
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

type environmentApplicationLifecycleCreatedDto struct {
	Name       string                                              `json:"name"`
	Properties environmentApplicationLifecycleCreatedPropertiesDto `json:"properties"`
}

type environmentApplicationLifecycleCreatedPropertiesDto struct {
	ProvisioningState string `json:"provisioningState"`
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

type applicationUsersResponseDto struct {
	Value []applicationUserDto `json:"value"`
}

type applicationUserDto struct {
	FullName      string `json:"fullname"`
	ApplicationId string `json:"applicationid"`
	Id            string `json:"applicationuserid"`
	SystemUserId  string `json:"systemuserid"`
}
