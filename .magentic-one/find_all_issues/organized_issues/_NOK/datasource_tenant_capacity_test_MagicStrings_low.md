# Title

Magic Strings for URLs and File Paths

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

The test case uses hardcoded (magic) strings in multiple places—specifically the API endpoint URL and the JSON test file path when registering the mock responder. This makes the code brittle and harder to maintain, as any changes to the endpoint or test data file location would require multiple updates across different test files.

## Impact

Low. While primarily a maintainability and readability issue, magic strings can lead to more significant maintenance burdens, especially as codebases and their testing suites grow. It increases the risk of bugs arising from accidental inconsistencies.

## Location

Lines 13–19 and line 17 specifically.

## Code Issue

```go
	httpmock.RegisterResponder("GET", `https://licensing.powerplatform.microsoft.com/v0.1-alpha/tenants/00000000-0000-0000-0000-000000000001/TenantCapacity`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant_capacity.json").String()), nil
		})
```

## Fix

Store constant strings for URLs and paths as variables or constants at the beginning of the test file. Preferably, define test data file paths and API URLs that are reused.

```go
const (
	mockTenantID      = "00000000-0000-0000-0000-000000000001"
	mockAPIBaseURL    = "https://licensing.powerplatform.microsoft.com/v0.1-alpha"
	mockCapacityPath  = "/tenants/" + mockTenantID + "/TenantCapacity"
	mockCapacityFile  = "tests/datasource/Validate_Read/get_tenant_capacity.json"
)

// ...

httpmock.RegisterResponder("GET", mockAPIBaseURL+mockCapacityPath,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File(mockCapacityFile).String()), nil
	})
```
