# Magic Static API Response File Path in Test

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant_test.go

## Problem

The test `TestUnitTenantDataSource_Validate_Read` uses a hard-coded static path to `"tests/datasource/Validate_Read/get_tenant.json"`. This file must be present and readable in the test environment. If the file is missing, moved, or altered, the unit test will either fail or behave unpredictably. This introduces fragility and reduces portability of the tests. Test data should be bundled or handled via test fixtures or embedded resources when necessary.

## Impact

**Severity: Medium**

- This can lead to broken CI/CD pipelines if the referenced file is missing.
- Test reliability becomes tied to file system state, reducing reproducibility.
- Contributors may miss adding this file, introducing accidental test failures.

## Location

Line where `httpmock.File("tests/datasource/Validate_Read/get_tenant.json")` is referenced in the `TestUnitTenantDataSource_Validate_Read` function.

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01`,
	func(_ *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant.json").String()), nil
	})
```

## Fix

Embed the file into the test binary using Go's embed feature, or check for file existence prior to running the test and provide a clear error if missing. With the embed package (Go 1.16+):

```go
import (
    _ "embed"
    // ...
)

//go:embed tests/datasource/Validate_Read/get_tenant.json
var tenantJSON string

httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01`,
	func(_ *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, tenantJSON), nil
	})
```

This guarantees the test data is always present when tests are run.
