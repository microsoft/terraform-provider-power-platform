# Title

Incorrect Test Data File Path in `TestUnitTenantApplicationPackagesDataSource_Validate_Read`

## File Path

`/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages_test.go`

## Problem

The test case `TestUnitTenantApplicationPackagesDataSource_Validate_Read` references a JSON file via `httpmock.File`. This file path appears to be incorrectly specified and does not follow consistent or validated path practices, which could lead to runtime errors if the file is moved or renamed.

## Impact

Incorrect file paths can lead to test failures due to missing or inaccessible test data. This makes tests unreliable and harder to execute on other environments, especially in CI/CD pipelines.

**Severity: High**

## Location

File path reference in `TestUnitTenantApplicationPackagesDataSource_Validate_Read`:

```go
httpmock.File("tests/datasource/tenant_application_packages/Validate_Read/get_tenant_applications.json")
```

## Code Issue

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/tenant_application_packages/Validate_Read/get_tenant_applications.json").String()), nil
```

## Fix

Update the file path to one that is correctly validated and ensure it is relative to the expected test data directory structure. Here's the corrected code:

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../../../../tests/datasource/tenant_application_packages/Validate_Read/get_tenant_applications.json").String()), nil
```

Additionally, consider using a configuration or constants file to manage test data paths centrally, reducing the risk of errors due to hardcoded paths.