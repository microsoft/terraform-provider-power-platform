# Title

Inconsistent Naming: Go Types Should Not Be Suffixed with 'Type'

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid_type.go

## Problem

The struct is named `UUIDType`, which is redundant and non-idiomatic in Go. According to Go naming conventions, the “Type” suffix should be avoided unless it specifically disambiguates. Since this code is in a custom types package, a better name would be simply `UUID` or similar.

## Impact

Severity: **Low**

- May reduce readability and burden cross-reference and refactoring processes.
- May complicate code, especially when using code navigation tools or documentation.

## Location

```go
type UUIDType struct {
	basetypes.StringType
}
```

## Code Issue

```go
type UUIDType struct {
	basetypes.StringType
}
```

## Fix

- Rename the struct to `UUID`, and update all associated usages within the codebase.

Example:

```go
type UUID struct {
	basetypes.StringType
}
```
