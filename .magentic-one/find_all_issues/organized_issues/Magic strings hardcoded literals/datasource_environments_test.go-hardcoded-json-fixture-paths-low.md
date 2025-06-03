# Title

Hardcoded JSON Fixture Paths Used in Unit Tests

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments_test.go

## Problem

The test file uses hardcoded string paths for JSON fixtures (e.g., `tests/datasource/Validate_Read/get_environments.json`). If these files are renamed, moved, or missing, the test will break. Since fixture location conventions or load errors are not handled, test maintainability is reduced. Code navigation and migration/refactoring are also harder if paths are not constructed from variables or constants.

## Impact

Severity: Low

- Test refactoring or moving fixtures becomes more error-prone.
- Test failures may not indicate missing/moved data in a clear way.
- New contributors may struggle to locate or update referenced fixture data.

## Location

Within `TestUnitEnvironmentsDataSource_Validate_Read` in:
- `httpmock.RegisterResponder(...httpmock.File("tests/datasource/Validate_Read/get_environments.json").String()), nil`
- and similar lines for related files

## Code Issue

```go
httpmock.RegisterResponder("GET", ..., func(req *http.Request) (*http.Response, error) {
    return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environments.json").String()), nil
})
```

## Fix
- Extract paths to appropriately-named constants
- Consider centralizing fixture file discovery and error handling logic.
- Optional: Add logic to verify existence/failure and improve error output if a referenced fixture is missing.

Example:

```go
const getEnvironmentsFixture = "tests/datasource/Validate_Read/get_environments.json"
const getEnvironment1Fixture = "tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json"

httpmock.RegisterResponder("GET", ..., func(req *http.Request) (*http.Response, error) {
    return httpmock.NewStringResponse(http.StatusOK, httpmock.File(getEnvironmentsFixture).String()), nil
})
```
