# Title

Lack of Negative/Edge Case Tests

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

The file tests only successful (happy-path) scenarios. There are no tests for missing resources, malformed data, or error conditions (400/404/500 responses, bad JSON, etc.), which means error-handling logic isn't being verified.

## Impact

Medium. This has a negative effect on robustness and completeness of the testing suite.

## Location

Globally, both test functions.

## Code Issue

(No specific code — absence of negative cases.)

## Fix

Add additional test cases using the mock responder to simulate error conditions and ensure the provider fails gracefully.

```go
httpmock.RegisterResponder("GET", mockAPIBaseURL+mockCapacityPath,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
	})
// ... Then, assert that the test returns an error.
```

Or test with malformed JSON, or simulate network errors.
