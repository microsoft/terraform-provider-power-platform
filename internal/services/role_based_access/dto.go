// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

// DTOs for the Power Platform Authorization RBAC API

type roleAssignmentRequestDto struct {
	PrincipalObjectId string `json:"principalObjectId"`
	PrincipalType     string `json:"principalType"`
	RoleDefinitionId  string `json:"roleDefinitionId"`
	Scope             string `json:"scope"`
}

type roleAssignmentDto struct {
	RoleAssignmentId           string  `json:"roleAssignmentId"`
	Scope                      string  `json:"scope"`
	PrincipalType              string  `json:"principalType"`
	PrincipalObjectId          string  `json:"principalObjectId"`
	RoleDefinitionId           string  `json:"roleDefinitionId"`
	CreatedByPrincipalType     string  `json:"createdByPrincipalType"`
	CreatedByPrincipalObjectId string  `json:"createdByPrincipalObjectId"`
	CreatedOn                  string  `json:"createdOn"`
	ExpiresOn                  *string `json:"expiresOn"`
}

type roleAssignmentsListDto struct {
	Value []roleAssignmentDto `json:"value"`
}

type roleDefinitionDto struct {
	RoleDefinitionId   string   `json:"roleDefinitionId"`
	RoleDefinitionName string   `json:"roleDefinitionName"`
	Permissions        []string `json:"permissions"`
}

type roleDefinitionsListDto struct {
	Value []roleDefinitionDto `json:"value"`
}
