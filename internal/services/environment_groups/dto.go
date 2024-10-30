// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_groups

type environmentGroupPrincipalDto struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type environmentGroupDto struct {
	DisplayName string                       `json:"displayName"`
	Description string                       `json:"description"`
	Id          string                       `json:"id,omitempty"`
	CreatedTime string                       `json:"createdTime,omitempty"`
	CreatedBy   environmentGroupPrincipalDto `json:"createdBy,omitempty"`
}

type environmentArrayDto struct {
	Value []environmentDto `json:"value"`
}

type environmentDto struct {
	Name       string                   `json:"name"`
	Properties environmentPropertiesDto `json:"properties"`
}

type environmentPropertiesDto struct {
	ParentEnvironmentGroup environmentGroupPropertiesDto `json:"parentEnvironmentGroup"`
}

type environmentGroupPropertiesDto struct {
	Id string `json:"id"`
}
