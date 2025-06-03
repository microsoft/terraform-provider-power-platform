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
