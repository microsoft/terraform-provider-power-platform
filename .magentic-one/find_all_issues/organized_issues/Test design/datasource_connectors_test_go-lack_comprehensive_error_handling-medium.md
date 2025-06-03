# Title

Lack of Comprehensive Error Handling in HTTP Mock Responders

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors_test.go

## Problem

The test file’s HTTP mock responders do not simulate error conditions such as unexpected status codes, malformed JSON, or HTTP client failures. All test responders only return HTTP 200 responses with valid JSON, missing negative test cases that are key to robust testing.

## Impact

This issue reduces test coverage and reliability. Code that handles error responses and exceptions is not exercised, risking unhandled edge cases in production. Severity: medium.

## Location

Lines 27–74 (HTTP mock responders in `TestUnitConnectorsDataSource_Validate_Read`)

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Read/get_virtual.json").String()), nil
    })
```

## Fix

Add test steps to simulate error responses and ensure error-handling code paths are covered. For example, register a responder that returns an error code and assert that the provider deals with it as expected.

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual`,
    func(req *http.Request) (*http.Response, error) {
        // Simulate an HTTP 500 Internal Server Error
        return httpmock.NewStringResponse(http.StatusInternalServerError, "Internal error"), nil
    })

// Add another test step to check behavior on error
resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.#", "0")
// Or expect a specific error, depending on provider's contract
```
