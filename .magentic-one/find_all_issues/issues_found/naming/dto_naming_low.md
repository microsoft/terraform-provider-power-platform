# Title

Struct Name Does Not Follow Go Naming Conventions

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/dto.go

## Problem

The struct `adminManagementApplicationDto` uses camelCase, which does not conform to Go's convention of using PascalCase (UpperCamelCase) for type names. This can affect readability, maintainability, and may lead to confusion when exporting types.

## Impact

Low. This is primarily a readability and maintainability concern, but consistent application of naming conventions is important for overall code quality and cooperation in a Go codebase.

## Location

Line 4-6

## Code Issue

```go
type adminManagementApplicationDto struct {
	ClientId string `json:"applicationId"`
}
```

## Fix

Rename the struct to use PascalCase, following Go conventions. That is, change `adminManagementApplicationDto` to `AdminManagementApplicationDTO` (acronyms in Go types are conventionally all-caps, but can be just `Dto` as well if you wish).

```go
type AdminManagementApplicationDTO struct {
	ClientId string `json:"applicationId"`
}
```
