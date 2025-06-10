# Title

Async, Validate, and Environment ID DTOs Should Be Exported

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

The types `asyncSolutionPullResponseDto`, `validateSolutionImportResponseDto`, `validateSolutionImportResponseSolutionOperationResultDto`, `environmentIdDto`, `environmentIdPropertiesDto`, and `linkedEnvironmentIdMetadataDto` are all unexported, likely limiting their usability for consumers of the package with type embedding, function parameters, or API responses. DTOs should generally be exported for broad usability.

## Impact

Not exporting these types limits code extensibility, clarity, and reuse. Severity: medium.

## Location

Lines 100–127

## Code Issue

```go
type asyncSolutionPullResponseDto struct { ... }
type validateSolutionImportResponseDto struct { ... }
type validateSolutionImportResponseSolutionOperationResultDto struct { ... }
type environmentIdDto struct { ... }
type environmentIdPropertiesDto struct { ... }
type linkedEnvironmentIdMetadataDto struct { ... }
```

## Fix

Capitalize each type to make them exported:

```go
type AsyncSolutionPullResponseDto struct { ... }
type ValidateSolutionImportResponseDto struct { ... }
type ValidateSolutionImportResponseSolutionOperationResultDto struct { ... }
type EnvironmentIdDto struct { ... }
type EnvironmentIdPropertiesDto struct { ... }
type LinkedEnvironmentIdMetadataDto struct { ... }
```
