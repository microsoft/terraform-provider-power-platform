# Method Name Does Not Conform to Go Naming Conventions

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

The method `securityRolesArray()` defined on `userDto` is not exported, and the name does not follow Go's conventional CamelCase for method naming (should be `SecurityRolesArray`). If this method is only needed internally that's fine, but if it is intended for wider use or for code consistency, its name should be exported and in CamelCase.

## Impact

**Severity: Low**

If intended for internal use only, this is a low-severity issue, but for API consistency, readability, and ensuring that future maintainers adhere to convention, method names should use proper casing. This matters more if/when the receiver type is exported or the method needs usage outside the declaring package.

## Location

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
```

## Code Issue

```go
func (u *userDto) securityRolesArray() []string {
	...
}
```

## Fix

Rename the method to `SecurityRolesArray`. If the receiver struct is also made exported (`UserDto`), this method will be usable consistently outside the package as well.

```go
func (u *UserDto) SecurityRolesArray() []string {
	...
}
```

If the function is meant to be unexported, document its purpose and scope for clarity.
