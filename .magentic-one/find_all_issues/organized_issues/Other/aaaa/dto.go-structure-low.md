# Unnecessary Struct Nesting Increases Complexity

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/dto.go

## Problem

Many DTO structs are deeply nested for relatively simple cases (e.g., `environmentIdPropertiesDto` contains only one field). Over-nesting can make the code harder to read and maintain without providing value.

## Impact

Low severity: This impedes code readability and maintainability but is not a critical error.

## Location

- Lines 10-22 (DTO nesting for environment and properties)

## Code Issue

```go
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

Consider flattening struct hierarchy where possible or grouping fields logically to minimize pointless nesting.

```go
type EnvironmentIdDto struct {
	Id                        string                        `json:"id"`
	Name                      string                        `json:"name"`
	LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}
```
