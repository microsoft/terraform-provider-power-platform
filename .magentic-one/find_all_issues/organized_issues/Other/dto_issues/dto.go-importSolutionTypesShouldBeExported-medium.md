# Title

Import Solution DTO Types Should Be Exported

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

The types `importSolutionDto`, `importSolutionResponseDto`, `importSolutionSolutionParametersDto`, `importSolutionConnectionReferencesDto`, and `importSolutionEnvironmentVariablesDto` are all unexported. For DTOs exchanged between packages, they should be exported (type names should begin with uppercase).

## Impact

Not exporting these types inhibits their use outside the defining package and prevents consistent and extensible API or data model definitions. Severity: medium.

## Location

Lines 68–98

## Code Issue

```go
type importSolutionDto struct { ... }
type importSolutionResponseDto struct { ... }
type importSolutionSolutionParametersDto struct { ... }
type importSolutionConnectionReferencesDto struct { ... }
type importSolutionEnvironmentVariablesDto struct { ... }
```

## Fix

Capitalize the first letter of each type name:

```go
type ImportSolutionDto struct { ... }
type ImportSolutionResponseDto struct { ... }
type ImportSolutionSolutionParametersDto struct { ... }
type ImportSolutionConnectionReferencesDto struct { ... }
type ImportSolutionEnvironmentVariablesDto struct { ... }
```
