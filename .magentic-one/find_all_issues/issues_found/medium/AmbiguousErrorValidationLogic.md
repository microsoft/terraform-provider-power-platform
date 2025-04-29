# Title

Ambiguous Error Validation Logic

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

The error validation used within tests lacks clarity. It employs regular expressions that may not fully validate error messages or handle edge cases correctly. For instance, checks like `regexp.MustCompile(".*moving data across regions is not supported in the unitedstates location.*")` are potentially brittle without explicit error handling.

## Impact

This issue may result in tests that pass incorrectly or fail unexpectedly due to loosely defined regular expressions. The severity is medium because it affects test accuracy, leading to potential misinterpretations of system behavior.

## Location

Examples can be found in the following places:
1. In the test `TestAccEnvironmentsResource_Validate_CreateGenerativeAiFeatures_US_Region_Expect_Fail`, the ExpectError uses overly broad regex patterns.

## Code Issue

```go
ExpectError: regexp.MustCompile(".*moving data across regions is not supported in the unitedstates location.*"),
```

## Fix

Define and assert specific error codes or structured error responses for better validation. This ensures that errors are properly checked without relying on loosely defined patterns.

```go
ExpectError: regexp.MustCompile("^Error\\: moving data across regions is not supported in the unitedstates location$"),
```

Alternatively, refactor the logic to include structured error handling:

```go
ExpectError: func(err error) bool {
    return IsSpecificError(err, "moving data across regions is not supported in the unitedstates location")
}
```

This approach minimizes ambiguity and ensures that errors are correctly identified during tests.