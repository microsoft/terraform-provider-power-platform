# Resource Test Issues - Merged Issues

## ISSUE 1

# Title
Use of ALL_CAPS for Constant Naming Is Unidiomatic in Go

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution_test.go

## Problem
Go convention for constants naming uses CamelCase or mixedCaps, not ALL_CAPS with underscores as is common in other languages. Constants such as `SOLUTION_1_NAME`, `SOLUTION_1_RELATIVE_PATH`, `SOLUTION_2_NAME`, and `SOLUTION_2_RELATIVE_PATH` do not follow Go naming conventions.

## Impact
Using non-idiomatic Go naming conventions decreases code readability and makes code look inconsistent to Go developers. Severity: **high**, since idiomatic naming impacts maintainability and team adoption.

## Location
Lines 18-24

## Code Issue
```go
const (
    SOLUTION_1_NAME          = "TerraformTestSolution_Complex_1_1_0_0.zip"
    SOLUTION_1_RELATIVE_PATH = "tests/resource/Test_Files/" + SOLUTION_1_NAME

    SOLUTION_2_NAME          = "TerraformSimpleTestSolution_1_0_0_1_managed.zip"
    SOLUTION_2_RELATIVE_PATH = "tests/resource/Test_Files/" + SOLUTION_2_NAME
)
```

## Fix
Change constants to use CamelCase or mixedCaps format. For example:

```go
const (
    Solution1Name          = "TerraformTestSolution_Complex_1_1_0_0.zip"
    Solution1RelativePath = "tests/resource/Test_Files/" + Solution1Name

    Solution2Name          = "TerraformSimpleTestSolution_1_0_0_1_managed.zip"
    Solution2RelativePath = "tests/resource/Test_Files/" + Solution2Name
)
```


---

## ISSUE 2

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


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
