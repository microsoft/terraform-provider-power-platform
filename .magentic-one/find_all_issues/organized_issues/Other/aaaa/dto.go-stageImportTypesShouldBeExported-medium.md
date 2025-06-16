# Title

Stage Solution Import DTO Types Should Be Exported

##

/workspaces/terraform-provider-power-platform/internal/services/solution/dto.go

## Problem

The types `stageSolutionImportDto`, `stageSolutionImportResponseDto`, `stageSolutionImportResultResponseDto`, and `stageSolutionSolutionDetailsDto` are all unexported (names start with lowercase). If these types are required from outside the `solution` package (as is likely for DTOs representing structured data across APIs), they must be exported.

## Impact

Prevents reuse and data exchange involving the types from outside the package, causing maintainability and reuse issues. Severity: medium.

## Location

Lines 29–66

## Code Issue

```go
type stageSolutionImportDto struct { ... }
type stageSolutionImportResponseDto struct { ... }
type stageSolutionImportResultResponseDto struct { ... }
type stageSolutionSolutionDetailsDto struct { ... }
```

## Fix

Capitalize the first letter of each type:

```go
type StageSolutionImportDto struct { ... }
type StageSolutionImportResponseDto struct { ... }
type StageSolutionImportResultResponseDto struct { ... }
type StageSolutionSolutionDetailsDto struct { ... }
```
