# Unexported Struct Types in DTO Layer

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

Some DTO types are declared with a leading lowercase letter (unexported) such as `environmentArrayDto`, `environmentCreateDto`, `modifySkuDto`, etc. Even for internal packages, best practice is to make Data Transfer Objects exported, because they may need to be used by other packages or in tests.

## Impact

While currently not causing functional bugs, this restricts future extensibility and makes the code less idiomatic. May lead to confusion about struct visibility, usage, and Go code conventions. Severity: **low**.

## Location

Throughout the bottom 1/3 of the file, e.g.,
- `type environmentArrayDto struct { ... }`
- `type environmentCreateDto struct { ... }`
- `type modifySkuDto struct { ... }`

## Code Issue

```go
type environmentArrayDto struct {
    Value []EnvironmentDto `json:"value"`
}
```

## Fix

Export DTOs as a rule unless there is a strong encapsulation reason:

```go
type EnvironmentArrayDto struct {
    Value []EnvironmentDto `json:"value"`
}
```
