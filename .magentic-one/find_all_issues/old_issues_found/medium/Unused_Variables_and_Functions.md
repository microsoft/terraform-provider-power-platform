# Title

Unused Variables and Functions

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

The function `securityRolesArray()` is defined but appears unused elsewhere within this file or imported to other parts of the application. Similarly, structs such as `environmentIdDto` and `environmentIdPropertiesDto` are defined but do not seem to interact with the rest of the application logic.

## Impact

Unused code introduces unnecessary complexity and confusion, which can lead to maintenance challenges and bloated codebases. It is also a source of potential bugs if future development unintentionally misuses these unverified pieces of the codebase. Severity is **medium**, as it affects maintainability but does not directly break functionality.

## Location

1. Function `securityRolesArray` defined inside `userDto`.
2. Struct `environmentIdDto`.
3. Struct `environmentIdPropertiesDto`.

## Code Issue

```go
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

type environmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties environmentIdPropertiesDto `json:"properties"`
}

type environmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata linkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}
```

## Fix

The unused function and structs should either be:
1. Removed if genuinely not necessary.
2. Integrated with relevant parts of the codebase if they are meant to serve a purpose in the application.

Example removal:

```go
// Remove unused function
// func (u *userDto) securityRolesArray() []string { 
//     if len(u.SecurityRoles) == 0 { 
//         return []string{} 
//     } 
//     var roles []string 
//     for _, role := range u.SecurityRoles { 
//         roles = append(roles, role.RoleId) 
//     } 
//     return roles 
// } 

// Remove unused structs
// type environmentIdDto struct { ... }
// type environmentIdPropertiesDto struct { ... }
```

Example integration:

```go
// Utilize securityRolesArray() in convertDataverseFromUserDto or other places
func convertDataverseFromUserDto(userDto *userDto, disableDelete bool) UserResourceModel {
	model := UserResourceModel{
		Id:                types.StringValue(userDto.Id),
		AadId:             types.StringValue(userDto.AadObjectId),
		SecurityRoles:     userDto.securityRolesArray(), // Add integration
		UserPrincipalName: types.StringValue(userDto.DomainName),
		FirstName:         types.StringValue(userDto.FirstName),
		LastName:          types.StringValue(userDto.LastName),
		BusinessUnitId:    types.StringValue(userDto.BusinessUnitId),
	}
	model.DisableDelete = types.BoolValue(disableDelete)
	return model
}

// Add appropriate usage for `environmentIdDto` based on application context
```