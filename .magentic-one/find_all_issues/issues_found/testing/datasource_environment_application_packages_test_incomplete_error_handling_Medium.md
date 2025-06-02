# Incomplete Test Coverage for Error Handling Paths

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

There are no tests for HTTP error cases (e.g., 403, 404, 500 errors from the mocked API). Current tests only check positive cases or a special error for "No Dataverse exists in environment". Failures due to network/response errors are not considered.

## Impact

Severity: Medium

Lack of negative test cases for error handling can lead to undetected regressions or failures in edge cases, since production code may not be sufficiently exercised.

## Location

All test functions currently only mock successful HTTP responses or application-domain errors, not generic HTTP/network errors.

## Fix

Add additional test cases covering:

- 404 Not Found from key APIs
- 500 Internal Server Error
- Unexpected/unparsable body
- Other network failures, as appropriate

Example:

```go
httpmock.RegisterResponder("GET", "<some endpoint>",
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusInternalServerError, "Internal Server Error"), nil
    })
// ...and write a corresponding test step with ExpectError or asserting the correct error handling path.
```
