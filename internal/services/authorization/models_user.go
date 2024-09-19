// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authorization

import "github.com/hashicorp/terraform-plugin-framework/types"

type userDto struct {
	Id             string            `json:"systemuserid"`
	DomainName     string            `json:"domainname"`
	FirstName      string            `json:"firstname"`
	LastName       string            `json:"lastname"`
	AadObjectId    string            `json:"azureactivedirectoryobjectid"`
	BusinessUnitId string            `json:"_businessunitid_value"`
	SecurityRoles  []securityRoleDto `json:"systemuserroles_association,omitempty"`
}

type securityRoleDto struct {
	RoleId         string `json:"roleid"`
	Name           string `json:"name"`
	IsManaged      bool   `json:"ismanaged"`
	BusinessUnitId string `json:"_businessunitid_value"`
}

type securityRoleDtoArray struct {
	Value []securityRoleDto `json:"value"`
}

func (u *userDto) securityRolesArray() []string {
	if len(u.SecurityRoles) == 0 {
		return []string{}
	}
	var roles []string
	for _, role := range u.SecurityRoles {
		roles = append(roles, role.RoleId)
	}
	return roles
}

type userDtoArray struct {
	Value []userDto `json:"value"`
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

func convertFromUserDto(userDto *userDto, disableDelete bool) UserResourceModel {
	model := UserResourceModel{
		Id:                types.StringValue(userDto.Id),
		AadId:             types.StringValue(userDto.AadObjectId),
		SecurityRoles:     userDto.securityRolesArray(),
		UserPrincipalName: types.StringValue(userDto.DomainName),
		FirstName:         types.StringValue(userDto.FirstName),
		LastName:          types.StringValue(userDto.LastName),
		BusinessUnitId:    types.StringValue(userDto.BusinessUnitId),
	}
	model.DisableDelete = types.BoolValue(disableDelete)
	return model
}
