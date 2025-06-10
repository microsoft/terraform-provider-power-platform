# Inconsistent Struct Naming Conventions

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

Several types in the file use different leading capitalization and struct naming conventions (e.g., `environmentArrayDto`, `environmentCreateDto`, `modifySkuDto`, etc.) versus the majority which use capitalized names (e.g., `EnvironmentDto`). In Go, exported types (usable outside the package) should always use capitalized, CamelCase names.

## Impact

While all types in this file may be internal, this inconsistency confuses both users and maintainers, as some types are exported while others are package-private. Go style recommends using capitalized names for exported types for consistency and clarity. Severity: **low**.

## Location

Throughout the entire file, e.g.
- `environmentArrayDto`
- `environmentCreateDto`
- `modifySkuDto`

## Code Issue

```go
type environmentArrayDto struct {
    Value []EnvironmentDto `json:"value"`
}

// ...

type modifySkuDto struct {
    EnvironmentSku string `json:"environmentSku,omitempty"`
}
```

## Fix

Rename these types to use capitalized CamelCase style, e.g.:

```go
type EnvironmentArrayDto struct {
    Value []EnvironmentDto `json:"value"`
}

type EnvironmentCreateDto struct {
    Location   string                         `json:"location"`
    Properties EnvironmentCreatePropertiesDto `json:"properties"`
}

type ModifySkuDto struct {
    EnvironmentSku string `json:"environmentSku,omitempty"`
}
```
