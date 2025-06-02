# Test File Uses Hardcoded String Paths

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/datasource_environment_templates_test.go

## Problem

The test setup uses hardcoded string paths (`tests/datasource/Validate_Read/get_environment_templates.json`). This reduces maintainability, making file renames/refactors harder and error-prone.

## Impact

Reduces maintainability and may cause tests to break if files are moved. Severity: Medium.

## Location

```go
httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_templates.json").String())
```

## Code Issue

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_templates.json").String()), nil
```

## Fix

Use a constant for the file path, declared at the top of the file:

```go
const environmentTemplatesJSONPath = "tests/datasource/Validate_Read/get_environment_templates.json"
...
return httpmock.NewStringResponse(http.StatusOK, httpmock.File(environmentTemplatesJSONPath).String()), nil
```
