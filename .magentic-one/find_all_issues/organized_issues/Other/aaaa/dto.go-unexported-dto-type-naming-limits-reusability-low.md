# Unexported DTO Type Naming Limits Reusability

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/dto.go

## Problem

All DTO types in this file are defined with unexported (lowercase) names (e.g., `powerAppBapiDto`, `powerAppPropertiesBapiDto`, etc.). DTOs (Data Transfer Objects) are typically used for structuring data across packages, such as for API interaction. If they need to be used outside this package, they should be exported (with uppercase names). If they are intentionally internal, this is fine, but Go convention suggests DTO types might be reused.

## Impact

Severity: Low

If DTOs are unexported, they can't be used in other packages, reducing flexibility and reusability. It also makes them less consistent with typical Go naming conventions for structural types.

## Location

All struct type definitions in this file.

## Code Issue

```go
type powerAppBapiDto struct {
    // ...
}
type powerAppPropertiesBapiDto struct {
    // ...
}
type powerAppEnvironmentDto struct {
    // ...
}
type powerAppCreatedByDto struct {
    // ...
}
type powerAppArrayDto struct {
    // ...
}
```

## Fix

If these DTOs should be used outside of this package, export them by renaming with uppercase first letters.

```go
type PowerAppBapiDto struct {
    // ...
}
type PowerAppPropertiesBapiDto struct {
    // ...
}
type PowerAppEnvironmentDto struct {
    // ...
}
type PowerAppCreatedByDto struct {
    // ...
}
type PowerAppArrayDto struct {
    // ...
}
```
