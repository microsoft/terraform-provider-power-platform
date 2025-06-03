# Title

No Unit Tests for API Methods

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/api_admin_management_application.go

## Problem

There are no unit tests in this file for the API client methods, meaning no assurance of correctness, regression protection, or contract validation between provider and API.

## Impact

Medium to High. Affects code reliability and maintainability.

## Location

This file; absence of any `_test.go` or test functions.

## Code Issue

_No code present, as tests are missing._

## Fix

Write unit tests using Go's testing framework, possibly using a mock for `*api.Client`.

```go
// Example proto-test
func TestGetAdminApplication_Success(t *testing.T) {
    // setup mock api.Client
    // simulate response
    // call GetAdminApplication
    // assert result
}
```
