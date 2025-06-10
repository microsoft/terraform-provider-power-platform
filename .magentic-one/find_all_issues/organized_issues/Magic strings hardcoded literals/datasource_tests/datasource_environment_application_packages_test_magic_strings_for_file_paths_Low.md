# Magic Strings for File Paths in Mocks

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

File paths to test fixtures are hardcoded as raw strings repeatedly, risking typos and making refactoring harder. Multiple places use the string "tests/datasource/environment_application_packages/..." directly, which is also duplicated across test functions.

## Impact

Severity: Low

Centralizing and reusing constants makes renaming and refactoring easier (single source of truth), reduces risk of typos, and improves readability.

## Location

Example:

```go
httpmock.File("tests/datasource/environment_application_packages/Validate_Read/get_applications.json")
```

## Fix

Introduce a constant or helper at the top of the test file:

```go
const testDataDir = "tests/datasource/environment_application_packages"

// then use:
httpmock.File(fmt.Sprintf("%s/Validate_Read/get_applications.json", testDataDir))
```
