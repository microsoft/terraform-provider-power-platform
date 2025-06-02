# Title
Inconsistent Variable Naming for Checksum Values

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution_test.go

## Problem
Local variables for file checksums (e.g., `solution_checksum`, `settings_checksum`, `solutionFileChecksum`, `settingsFileChecksum`) use inconsistent naming styles — some use snake_case, others CamelCase or lowerCamelCase. Consistency in variable naming helps code readability and maintainability.

## Impact
Mixed conventions within the same file can confuse developers, reduce code clarity, and make code harder to scan or refactor. Severity: **medium** (style, maintainability).

## Location
Throughout the file, e.g., lines 86, 93, 159, 171, etc.

## Code Issue
```go
solution_checksum := createFile("test_solution.zip", "test_solution")
settings_checksum := createFile("test_solution_settings.json", "")
...
solutionFileChecksum, _ := helpers.CalculateSHA256(SOLUTION_1_NAME)
settingsFileChecksum, _ := helpers.CalculateSHA256(solutionSettingsFileName)
```

## Fix
Choose one naming convention and use it throughout the file. Go typically prefers camelCase or MixedCaps for variables. Example:

```go
solutionFileChecksum := createFile("test_solution.zip", "test_solution")
settingsFileChecksum := createFile("test_solution_settings.json", "")
// Or, if you prefer shorter:
solFileChecksum := ...
setFileChecksum := ...
```
Apply the convention consistently to both `createFile` results and direct calculations.
