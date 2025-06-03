# Lack of Negative/Edge Case Tests

##
/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages_test.go

## Problem

The test file solely covers expected/positive test scenarios for both acceptance and unit tests. It does not include negative or edge case tests, such as invalid API responses, empty responses, failed HTTP calls, configuration errors, or boundary conditions for attributes.

## Impact

Severity: **Medium**  
Not having negative or edge case tests reduces test coverage and leaves the provider susceptible to undetected failures in real-world scenarios such as bad input, schema changes, or unexpected API issues.

## Location

- Applies throughout the test functions in the file; there are no test cases for negative or edge case behaviors.

## Code Issue

```go
// No tests for when the HTTP API returns errors, invalid JSON, or attributes are missing/invalid.
// All test cases assume success responses.
```

## Fix

Add test steps to simulate API errors (e.g., return non-200 response codes or malformed JSON), check for missing/invalid configurations, and validate that the datasource fails as expected (with proper error messages or no data).

```go
httpmock.RegisterResponder("GET", "...",
    func(req *http.Request) (*http.Response, error) {
        // Simulate API error
        return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
    })

resource.TestStep{
    Config: `data "powerplatform_tenant_application_packages" "invalid" {}`,
    ExpectError: regexp.MustCompile("internal server error"),
},
```
