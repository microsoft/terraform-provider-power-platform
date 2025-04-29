# Title

Hardcoded API Endpoint Without Adequate Testing

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings_test.go`

## Problem

In the test function `TestUnitTestTenantSettingsDataSource_Validate_Read`, the hardcoded URL `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01` is directly embedded in the `httpmock.RegisterResponder`. This creates a strong coupling with the API structure and version, which can lead to breaks in functionality if the endpoint or version changes.

## Impact

Hardcoding the API endpoint increases fragility in test cases. If the endpoint changes, all tests relying on this function will fail. It also makes it harder to maintain the code and reuse parts of the test across different API versions or endpoints. Severity: **High**.

## Location

Within the `TestUnitTestTenantSettingsDataSource_Validate_Read` function.

## Code Issue

```go
httpmock.RegisterResponder("POST", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings?api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/post_list_tenant_settings.json").String()), nil
    })
```

## Fix

Abstract the hardcoded URL into a configurable constant or test parameter so it can be easily updated or customized without altering multiple test functions.

```go
const TenantSettingsAPIEndpoint = "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/listTenantSettings"

httpmock.RegisterResponder("POST", fmt.Sprintf("%s?api-version=2023-06-01", TenantSettingsAPIEndpoint),
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/post_list_tenant_settings.json").String()), nil
    })
```

This approach reduces direct coupling and makes the code easier to maintain and scale across different API versions or endpoints.