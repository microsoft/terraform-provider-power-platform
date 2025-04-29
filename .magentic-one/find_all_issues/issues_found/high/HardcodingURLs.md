# Title
Hardcoding URLs

## 
/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem
The code contains hardcoded URLs such as `https://api.bap.microsoft.com` and `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami`. Hardcoding URLs can make the application inflexible to changes in environments or configurations and introduces risks, such as exposing sensitive information.

## Impact
Hardcoding URLs reduces the ability to adapt to changes in services or configurations without modifying the source code. If the URLs need to change, several code locations need updating, increasing the likelihood of errors. It's also harder to manage environments like development, testing, staging, and production using the same code.

Severity: **High**

## Location
- Line 17: `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`
- Line 22: `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami`

## Code Issue
The code currently hardcodes URLs in the HTTP mock responses:

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
    })
```

```go
httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_whoami.json").String()), nil
    })
```

## Fix
The hardcoded URLs should be replaced with constants or configuration variables. Use environment variables or a centralized configuration file to store the URLs, making it easier to update and manage across environments. Here is an example implementation:

```go
const (
    baseAPIURL         = "https://api.bap.microsoft.com"
    dynamicsBaseURL    = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2"
)

httpmock.RegisterResponder("GET", baseAPIURL+"/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
    })

httpmock.RegisterResponder("GET", dynamicsBaseURL+"/whoami",
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_whoami.json").String()), nil
    })
```

Alternately, you can use environment variables:
```go
baseAPIURL := os.Getenv("BASE_API_URL")
dynamicsBaseURL := os.Getenv("DYNAMICS_BASE_URL")

httpmock.RegisterResponder("GET", baseAPIURL+"/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
    })

httpmock.RegisterResponder("GET", dynamicsBaseURL+"/whoami",
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_whoami.json").String()), nil
    })
```

This approach makes it easier to manage configurations across different environments and improves maintainability.
