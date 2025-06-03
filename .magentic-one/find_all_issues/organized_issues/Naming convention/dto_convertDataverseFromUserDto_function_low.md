# Overly Broad Function Naming and Usage for Conversion

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

The function `convertDataverseFromUserDto` uses a non-idiomatic (Java/C#-style) naming ("convertDataverseFromUserDto") and accepts a pointer to `userDto` as input, which may cause confusion in understanding the conversion direction and the resulting type. Moreover, neither `UserResourceModel` nor what constitutes “Dataverse” is explained or clear, making the intent ambiguous and potentially hard to maintain.

## Impact

**Severity: Low**

Readability, maintainability, and discoverability are affected. Understanding what is being converted to what, and under what conditions, becomes less clear. This could hinder onboarding and correct usage in the future.

## Location

```go
func convertDataverseFromUserDto(userDto *userDto, disableDelete bool) UserResourceModel {
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
```

## Code Issue

```go
func convertDataverseFromUserDto(userDto *userDto, disableDelete bool) UserResourceModel {
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
```

## Fix

Rename the function to clearly describe the conversion direction—ideally: `UserDtoToResourceModel` or `NewUserResourceModelFromDto`.

```go
func UserDtoToResourceModel(dto *UserDto, disableDelete bool) UserResourceModel {
	model := UserResourceModel{
		Id:                types.StringValue(dto.Id),
		AadId:             types.StringValue(dto.AadObjectId),
		SecurityRoles:     dto.SecurityRolesArray(),
		UserPrincipalName: types.StringValue(dto.DomainName),
		FirstName:         types.StringValue(dto.FirstName),
		LastName:          types.StringValue(dto.LastName),
		BusinessUnitId:    types.StringValue(dto.BusinessUnitId),
	}
	model.DisableDelete = types.BoolValue(disableDelete)
	return model
}
```

This makes the intent and ownership of the function clear, and improves discoverability and maintainability.
