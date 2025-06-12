# Naming Convention Violation: Struct Names Should Be Exported

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/dto.go

## Problem

Structs intended to be transferred as DTOs (Data Transfer Objects) are named with an unexported (lowercase) identifier (e.g., `dataRecordDto`). According to Go conventions, struct types intended for use outside their defining package (including serialization in APIs and by packages like encoding/json) should start with an uppercase letter.

## Impact

This can negatively impact code readability, maintainability, and, most importantly, usage outside the current package. Downstream consumers cannot reference these types, and this is a medium-severity issue for public or external-facing DTOs.

## Location

- Lines: 3-32, multiple structs (`dataRecordDto`, `environmentIdDto`, `environmentIdPropertiesDto`, `linkedEnvironmentIdMetadataDto`, `entityDefinitionsDto`, `relationApiResponseDto`, `relationApiBodyDto`, `attributesApiResponseDto`, `attributesApiBodyDto`).

## Code Issue

```go
type dataRecordDto struct {
  // ...
}

type environmentIdDto struct {
  // ...
}

type environmentIdPropertiesDto struct {
  // ...
}

// ... etc
```

## Fix

Export struct types by capitalizing the first letter of each struct name for consistency with Go conventions.

```go
type DataRecordDto struct {
  // ...
}

type EnvironmentIdDto struct {
  // ...
}

type EnvironmentIdPropertiesDto struct {
  // ...
}

// ... etc
```
