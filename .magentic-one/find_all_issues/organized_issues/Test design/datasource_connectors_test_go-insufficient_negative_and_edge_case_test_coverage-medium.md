# Title

Insufficient Negative and Edge Case Test Coverage

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors_test.go

## Problem

The current tests only verify the "happy path" (successful attribute population). There is no verification of error handling, malformed attributes, empty data, or permission issues.

## Impact

Risks undetected bugs in error scenarios; overall code reliability is reduced, real-world edge cases may be missed. Severity: medium.

## Location

Throughout the test steps in both functions

## Code Issue

```go
resource.TestCheckResourceAttr("data.powerplatform_connectors.all", "connectors.#", "4")
// ...
// No test cases for empty responses, missing fields, permission denied, etc.
```

## Fix

Add unit and acceptance tests with intentionally empty, missing, or erroneous data to ensure robust handling, e.g.:

```go
// Mock a bad response:
httpmock.RegisterResponder("GET", someURL, func(req *http.Request) (*http.Response, error) {
    return httpmock.NewStringResponse(http.StatusForbidden, "access denied"), nil
})
// Or create steps without a required attribute and assert correct failure/error diagnostics
```
