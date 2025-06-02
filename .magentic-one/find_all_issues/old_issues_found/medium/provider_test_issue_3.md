# Title

Hardcoded API URLs in Mocked HTTP Responses

## Path

`/workspaces/terraform-provider-power-platform/internal/provider/provider_test.go`

## Problem

Throughout the test file, hardcoded API URLs are directly used when defining mocked HTTP responders for `httpmock`. For example:

```go
httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments.*api-version=2023-06-01$`),
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
	})
```

## Impact

- **Maintenance Overhead**: Changes in actual API URLs can cause tests to fail and require manual updates of these hardcoded strings across multiple places.
- **Code Duplication**: Repeated hardcoding of the same URLs leads to duplicated code.
- **Medium Severity**: While tests work under current conditions, scalability and manageability are negatively impacted.

## Location

Lines 105–112 and 131–138 in the file `provider_test.go`.

## Code Issue

```go
httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments.*api-version=2023-06-01$`),
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
	})
```

## Fix

Define constant variables for the API URLs and reuse them throughout the file. This reduces duplication and enhances maintainability.

```go
const (
	bapAPIBaseURL              = "https://api.bap.microsoft.com"
	environmentsEndpointRegexp = `^` + bapAPIBaseURL + `/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments.*api-version=2023-06-01$`
)

httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(environmentsEndpointRegexp),
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
	})
```