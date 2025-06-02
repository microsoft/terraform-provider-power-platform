# Title

Lack of Test Isolation and Cleanup for HTTP Mocks Across Functions

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy_test.go

## Problem

Although `httpmock.Activate()` and `httpmock.DeactivateAndReset()` are used within most unit test functions, the `setupTenantHttpMocks` helper and test setup code do not guarantee that mocks from one test do not leak into another (especially since some tests, like acceptance tests, do not isolate state as well). Furthermore, not all tests activate/deactivate the mock, and errors during cleanup (e.g., failure to deactivate) are ignored rather than causing test failure or reporting.

## Impact

If any test fails or mutates global mock state, subsequent tests may become unpredictable—leading to flaky or misleading test outcomes. This can be a significant problem as the suite grows. Severity: **medium**, as intermittent and nondeterministic tests are disruptive but will not always block CI outright.

## Location

Helper function and pattern throughout the file, e.g.:

```go
func setupTenantHttpMocks() {
    // ...does not clean up or enforce isolation between tests
}
```

and scattered call locations for `httpmock.Activate()` and `httpmock.DeactivateAndReset()`, with possible side effects.

## Code Issue

```go
func setupTenantHttpMocks() {
    // ...
}
// Called from several test functions, without ensuring mock state reset/isolation.
```

## Fix

Wrap *every* test that uses HTTP mocking in code that activates/deactivates (and resets) the mock, and handle errors. Ensure `setupTenantHttpMocks` never leaves a responder registered beyond a test's context. Use `t.Cleanup()` to ensure proper teardown. For example:

```go
func TestSomething(t *testing.T) {
    httpmock.Activate()
    t.Cleanup(func() { httpmock.DeactivateAndReset() })

    // Proceed with setup, test, etc.
}
```

Similarly, update `setupTenantHttpMocks` (or eliminate) to avoid any global side effects—pass explicit context or reset after use.
