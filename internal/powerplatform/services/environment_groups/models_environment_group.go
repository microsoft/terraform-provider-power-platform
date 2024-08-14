// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

type EnvironmentGroupPrincipalDto struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type EnvironmentGroupDto struct {
	DisplayName string                       `json:"displayName"`
	Description string                       `json:"description"`
	Id          string                       `json:"id,omitempty"`
	CreatedTime string                       `json:"createdTime,omitempty"`
	CreatedBy   EnvironmentGroupPrincipalDto `json:"createdBy,omitempty"`
}
