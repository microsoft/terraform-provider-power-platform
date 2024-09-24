// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerapps

type powerAppBapiDto struct {
	Name       string                    `json:"name"`
	Properties powerAppPropertiesBapiDto `json:"properties"`
}

type powerAppPropertiesBapiDto struct {
	DisplayName      string                 `json:"displayName"`
	Owner            powerAppCreatedByDto   `json:"owner"`
	CreatedBy        powerAppCreatedByDto   `json:"createdBy"`
	LastModifiedBy   powerAppCreatedByDto   `json:"lastModifiedBy"`
	LastPublishedBy  powerAppCreatedByDto   `json:"lastPublishedBy"`
	CreatedTime      string                 `json:"createdTime"`
	LastModifiedTime string                 `json:"lastModifiedTime"`
	LastPublishTime  string                 `json:"lastPublishTime"`
	Environment      powerAppEnvironmentDto `json:"environment"`
}

type powerAppEnvironmentDto struct {
	Id       string `json:"id"`
	Location string `json:"location"`
	Name     string `json:"name"`
}

type powerAppCreatedByDto struct {
	DisplayName       string `json:"displayName"`
	Id                string `json:"id"`
	UserPrincipalName string `json:"userPrincipalName"`
}

type powerAppArrayDto struct {
	Value []powerAppBapiDto `json:"value"`
}
