# Title

Issue with inappropriate naming convention for constant variables

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/constants.go

## Problem

The constant variables `ROLE_ENVIRONMENT_ADMIN` and `ROLE_ENVIRONMENT_MAKER` do not follow standard Go naming conventions. According to Go style guidelines, constant names should be written in CamelCase if they are intended to be exported or lowerCamelCase if they are not exported. Using all uppercase letters with underscores is against Go's idiomatic practices.

## Impact

- Code readability and maintainability may be affected because the naming convention is inconsistent with Go standards.
- Developers working on the project may face difficulty understanding and maintaining the code due to unfamiliar naming styles.
- Medium severity because this does not directly break any functionality but impacts code consistency and readability.

## Location

Line numbers:
- Line 6: `ROLE_ENVIRONMENT_ADMIN`
- Line 7: `ROLE_ENVIRONMENT_MAKER`

## Code Issue

```go
const (
	ROLE_ENVIRONMENT_ADMIN = "Environment Admin"
	ROLE_ENVIRONMENT_MAKER = "Environment Maker"
)
```

## Fix

The constants should be renamed following Go standard naming conventions:

```go
const (
	RoleEnvironmentAdmin = "Environment Admin"
	RoleEnvironmentMaker = "Environment Maker"
)
```

Renaming these constants improves consistency and aligns with best practices in Go programming.
