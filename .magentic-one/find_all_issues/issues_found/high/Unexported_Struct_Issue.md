# Title

Unexported Structs Used in Interface

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

Several structs such as `userDto`, `securityRoleDto`, etc., have been declared as unexported (use lowercase first letter). However, these structs are clearly being used for JSON serialization through tags, which suggests that they are part of an external interface. Unexported structs cannot be accessed if needed for deserialization or manipulation outside this package.

## Impact

This issue limits the usability of these structs across other packages and could lead to serialization/deserialization errors if these structs or fields need to be referenced outside the package. The severity is **high**, as the unexported nature of these structs creates scalability and integration issues for consumers of the codebase.

## Location

Struct issues were found starting with the declarations of `userDto`, `securityRoleDto`, and others.

## Code Issue

```go
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
```

## Fix

Export these structs by changing their names to start with an uppercase letter, making them accessible across packages and resolving JSON serialization/deserialization issues.

```go
type UserDto struct {
	Id             string            `json:"systemuserid"`
	DomainName     string            `json:"domainname"`
	FirstName      string            `json:"firstname"`
	LastName       string            `json:"lastname"`
	AadObjectId    string            `json:"azureactivedirectoryobjectid"`
	BusinessUnitId string            `json:"_businessunitid_value"`
	SecurityRoles  []SecurityRoleDto `json:"systemuserroles_association,omitempty"`
}

type SecurityRoleDto struct {
	RoleId         string `json:"roleid"`
	Name           string `json:"name"`
	IsManaged      bool   `json:"ismanaged"`
	BusinessUnitId string `json:"_businessunitid_value"`
}
```