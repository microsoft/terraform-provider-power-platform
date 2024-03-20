// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import "github.com/hashicorp/terraform-plugin-framework/types"

type UserDto struct {
	Id            string            `json:"systemuserid"`
	DomainName    string            `json:"domainname"`
	FirstName     string            `json:"firstname"`
	LastName      string            `json:"lastname"`
	AadObjectId   string            `json:"azureactivedirectoryobjectid"`
	SecurityRoles []SecurityRoleDto `json:"systemuserroles_association,omitempty"`
}

type SecurityRoleDto struct {
	RoleId         string `json:"roleid"`
	Name           string `json:"name"`
	IsManaged      bool   `json:"ismanaged"`
	BusinessUnitId string `json:"_businessunitid_value"`
}

type SecurityRoleDtoArray struct {
	Value []SecurityRoleDto `json:"value"`
}

func (u *UserDto) SecurityRolesArray() []string {
	if len(u.SecurityRoles) == 0 {
		return []string{}
	} else {
		var roles []string
		for _, role := range u.SecurityRoles {
			roles = append(roles, role.RoleId)
		}
		return roles
	}
}

type UserDtoArray struct {
	Value []UserDto `json:"value"`
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

func ConvertFromUserDto(userDto *UserDto, disableDelete bool) UserResourceModel {
	model := UserResourceModel{
		Id:                types.StringValue(userDto.Id),
		AadId:             types.StringValue(userDto.AadObjectId),
		SecurityRoles:     userDto.SecurityRolesArray(),
		UserPrincipalName: types.StringValue(userDto.DomainName),
		FirstName:         types.StringValue(userDto.FirstName),
		LastName:          types.StringValue(userDto.LastName),
	}
	model.DisableDelete = types.BoolValue(disableDelete)
	return model
}
