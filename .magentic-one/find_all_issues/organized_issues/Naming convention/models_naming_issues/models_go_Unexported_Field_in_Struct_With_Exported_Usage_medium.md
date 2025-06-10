# Unexported Field in Struct With Exported Usage

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go

## Problem

The `UserResource` struct has a field named `UserClient` whose type is simply `client`. However, it is not clear from this file whether `client` is an exported or fully qualified type, and there is no error handling or data validation around its assignment or use, which could lead to confusion or poor control flow. 

Also, `UserClient` should follow Go naming conventions (`userClient` for unexported fields) if it is truly intended to be private. If it should be used outside the package, it must be exported, and its type must also be exported and/or fully qualified.

## Impact

Poor naming may reduce code readability and create confusion between exported and unexported uses. If `client` is meant to be accessed outside the package, mismatch between field and type visibility can break API contracts. Severity: **medium**.

## Location

```go
type UserResource struct {
	helpers.TypeInfo
	UserClient client
}
```

## Code Issue

```go
UserClient client
```

## Fix

If `client` is an internal type and `UserClient` should be unexported, rename to follow Go conventions:

```go
userClient client
```

If both should be exported, ensure both are capitalized and imported/exported as needed:

```go
UserClient Client
```

Or, if type is from another package:

```go
UserClient externalpkg.Client
```

Also review code for appropriate error handling or validation where this field is used or set.

