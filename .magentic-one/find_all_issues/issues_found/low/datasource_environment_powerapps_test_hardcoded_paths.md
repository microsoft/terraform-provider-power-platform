# Issue Report #2

## Title

Hardcoded API Response File Paths in Unit Tests

##

File Path:

`/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps_test.go`

## Problem

In the function `TestUnitEnvironmentPowerAppsDataSource_Validate_Read`, the paths of mock API response files (e.g., `tests/Validate_Read/get_environments.json`) are hardcoded directly into the test code. This practice can lead to maintenance difficulties when paths change or need to be updated.

## Impact

- Hardcoding file paths makes tests less maintainable and less reusable.
- Makes it harder to cleanly refactor test data organization.
- Introduces additional work if paths require updates.

Severity: **Low**

## Location

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_environments.json").String()), nil
    })

httpmock.RegisterResponder("GET", `=~^https://api\\.powerapps\\.com/providers/Microsoft\\.PowerApps/scopes/admin/environments/([\\d-]+)/apps`,
    func(req *http.Request) (*http.Response, error) {
        id := httpmock.MustGetSubmatch(req, 1)
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_apps_"+id+".json").String()), nil
    })
```

## Fix

Abstract the file paths into reusable constants or configuration values. This will make maintenance and refactoring easier if file locations ever change.

Code example for the fix:

```go
// Define constants for file paths
const (
    EnvironmentsResponsePath = "tests/Validate_Read/get_environments.json"
    AppsResponsePathTemplate = "tests/Validate_Read/get_apps_%s.json"
)

httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File(EnvironmentsResponsePath).String()), nil
    })

httpmock.RegisterResponder("GET", `=~^https://api\\.powerapps\\.com/providers/Microsoft\\.PowerApps/scopes/admin/environments/([\\d-]+)/apps`,
    func(req *http.Request) (*http.Response, error) {
        id := httpmock.MustGetSubmatch(req, 1)
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf(AppsResponsePathTemplate, id)).String()), nil
    })
```
