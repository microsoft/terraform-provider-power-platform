# Lack of Negative Test Cases

##

/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations_test.go

## Problem

The test file only tests successful data retrieval but does not provide any negative test cases (i.e., what happens if the API fails, or data is missing, or malformed). Testing only the happy path does not ensure robustness of error handling or correct error propagation.

## Impact

Reduces the overall test coverage and allows bugs to slip when errors or corner cases occur. Decreases the reliability and confidence in the code's robustness.

**Severity:** Medium

## Location

The entire test file.

## Code Issue

```go
// No error scenarios tested or asserted.
```

## Fix

### Add test case with mock responder returning an error

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
	},
)
// Add a test step asserting resource read fails as expected.
```
