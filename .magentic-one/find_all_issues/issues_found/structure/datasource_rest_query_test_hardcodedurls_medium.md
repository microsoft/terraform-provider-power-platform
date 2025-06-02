# Hard-Coded HTTP Mock URLs and Filenames

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query_test.go

## Problem

The test uses hard-coded URLs and paths for setting up mocks, including pasting long URLs and file paths as strings inline. This makes it harder to change the mocked API endpoint or test fixture JSON file names if needed, and increases the risk of typos or inconsistencies.

## Impact

This is a **medium** severity issue. Code maintainability and readability is reduced, and there may be duplicated data or difficulty updating resources for future test changes.

## Location

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
	})

httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Web_Apis_WhoAmI/get_whoami.json").String()), nil
	})
```

## Fix

Extract URLs and test file paths into constants at the top of the file, near the other constants, to improve maintainability:

```go
const (
	mockScopeURL  = "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01"
	mockWhoAmIURL = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/whoami"
	mockEnvFile   = "tests/datasource/Web_Apis_WhoAmI/get_environment_00000000-0000-0000-0000-000000000001.json"
	mockWhoAmIFile = "tests/datasource/Web_Apis_WhoAmI/get_whoami.json"
)

// ...
httpmock.RegisterResponder("GET", mockScopeURL,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File(mockEnvFile).String()), nil
	})

httpmock.RegisterResponder("GET", mockWhoAmIURL,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File(mockWhoAmIFile).String()), nil
	})
```
