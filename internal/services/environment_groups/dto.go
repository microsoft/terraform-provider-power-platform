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
