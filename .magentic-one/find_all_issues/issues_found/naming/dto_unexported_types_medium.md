# Unexported Struct Types Used as DTOs

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/dto.go

## Problem

The structs `userDto`, `securityRoleDto`, `securityRoleArrayDto`, etc., are named with an initial lowercase letter, making them unexported. Since these are Data Transfer Objects (DTOs), they are likely to be used across multiple packages, particularly for marshaling/unmarshaling JSON. Keeping them unexported restricts their usage, reduces potential reusability, and goes against Go naming conventions for types meant for sharing.

## Impact

**Severity: Medium**

DTOs that are unexported but are expected to be used outside of the current package cannot be accessed, resulting in the need for unnecessary wrapper types or copy-pasted structures in other packages. It also reduces testability from external packages and contradicts idiomatic Go naming guidelines for types expected for broader-use.

## Location

Multiple locations:
- Definition of `userDto`
- Definition of `securityRoleDto`
- Definition of `securityRoleArrayDto`
- ... (others follow this pattern)

## Code Issue

```go
type userDto struct {
  ...
}

type securityRoleDto struct {
  ...
}
```

## Fix

Capitalize the struct type names that are intended to be used outside this package. This will make them exported and accessible in other packages.

```go
type UserDto struct {
  ...
}

type SecurityRoleDto struct {
  ...
}
```

This applies to all relevant DTO structs in this file. If some types are intentionally kept private, please comment on the intention for clarity.
