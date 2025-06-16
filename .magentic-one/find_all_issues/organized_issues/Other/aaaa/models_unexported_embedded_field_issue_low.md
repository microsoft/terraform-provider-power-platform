# Unexported Embedded Field Issue

##

/workspaces/terraform-provider-power-platform/internal/services/connection/models.go

## Problem

The structs `SharesDataSource`, `ConnectionsDataSource`, `ShareResource`, and `Resource` embed an unexported field `ConnectionsClient client` without exporting the `client` type or providing a way for users in other packages to utilize these clients, unless the fields are intentionally package-private. The use of the lowercase identifier `client` for the type is also inconsistent and could cause confusion, especially if the actual type definition is not visible here.

## Impact

Low to Medium. This reduces extensibility if the intention is for these structs or their embedded clients to be used outside the package. Additionally, the lack of clarity as to what `client` refers to can be confusing, especially for maintainers or consumers of this package.

## Location

- Structs:
  - `SharesDataSource`
  - `ConnectionsDataSource`
  - `ShareResource`
  - `Resource`

```go
type SharesDataSource struct {
	helpers.TypeInfo
	ConnectionsClient client
}
```

## Fix

1. If `client` is an internal type, clarify the naming by making it exported if needed (e.g., `Client`).
2. Consider exporting the field if it is meant to be accessible outside the package.
3. If it is intentionally private, add documentation to clarify the intent.

```go
type SharesDataSource struct {
	helpers.TypeInfo
	ConnectionsClient Client // Use exported type if it should be accessible
}
```

Or, add a comment:

```go
type SharesDataSource struct {
	helpers.TypeInfo
	ConnectionsClient client // package-private client for internal use
}
```

This clarification aids package consumers and reduces confusion.
