# Title

Potential for Outdated Mock Responses for PUT/GET/DELETE

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_test.go

## Problem

The unit test registers HTTP mock responses for PUT, GET, and DELETE endpoints using static files (`tests/resource/connections/Validate_Create/put_connection.json`). If the resource schema or response payload changes in code, and the mock response files are not updated, the unit tests may continue to pass deceptively, not covering changed fields or error cases.

## Impact

This introduces a maintainability risk where mocks are no longer representative of real API responses, reducing test coverage and reliability. Subtle mismatches may result in production defects that are not caught in unit testing.

Severity: Medium

## Location

```go
return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
...
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()), nil
```

## Fix

Clearly document/link resource schema and version in the mock files, and regularly sync/update the mock payloads when the API or Terraform resource changes. Advanced: automate delta checks or implement schema-validated mocks.

```go
// In unit tests, add comments and validation steps to ensure mock files match the actual provider/schema/response
// Example:
payload := httpmock.File("tests/resource/connections/Validate_Create/put_connection.json").String()
// Optionally parse/validate payload against current schema before registering
```
This establishes a routine for trustworthy test coverage.
