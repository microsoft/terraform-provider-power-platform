# Issue: Lack of Unit Tests for Edge Cases and Error Handling

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/api_tenant_isolation_policy.go

## Problem

There is no evidence of unit tests validating the edge cases and error handling of functions such as `getTenantIsolationPolicy`, `createOrUpdateTenantIsolationPolicy`, and `doWaitForLifecycleOperationStatus`. Critical paths such as retries, error wrapping, status code handling, and asynchronous polling are not demonstrated to be tested, making it hard to guarantee the code’s reliability and regression safety.

## Impact

Medium to High. The absence of such testing risks undetected regressions, insufficient error handling coverage, and accidental breaking of asynchronous logic or edge case handling—especially around status codes and HTTP header parsing.

## Location

Applies to the whole file, especially the public API client methods and polling logic.

## Code Issue

No code block to show; refers to lack of tests surrounding:

- Error wrapping and propagation
- Async flow (`202 Accepted`)
- Header parsing (`Retry-After`)
- Failure modes (missing headers, context cancellation, etc.)

## Fix

Create dedicated unit tests validating:

- Error propagation when the API responds with unexpected errors or missing headers.
- Polling and sleep logic upon asynchronous operations.
- Correct parsing of `Retry-After` header.
- Handling of 404 Not Found and similar status codes gracefully.
- Context cancellation.

Example (in a test file, e.g. `api_tenant_isolation_policy_test.go`):

```go
func TestGetTenantIsolationPolicy_NotFound(t *testing.T) {
	// Arrange mock API client that returns 404
	// Call getTenantIsolationPolicy and validate output is (nil, nil)
}

func TestCreateOrUpdateTenantIsolationPolicy_AsyncSuccess(t *testing.T) {
	// Arrange mock client with 202 Accepted then 200 OK
	// Test proper polling and final fetch
}
```
