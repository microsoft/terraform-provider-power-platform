# Title

Insufficient Isolation/Reset of External HTTP Mock in Unit Tests

## 

/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies_test.go

## Problem

The unit test `TestUnitTestBillingPoliciesDataSource_Validate_Read` globally activates the httpmock library using `httpmock.Activate()` and deactivates it with `defer httpmock.DeactivateAndReset()`. While this is a recommended pattern, it is important to note that the activation is global, affecting all outgoing HTTP requests within the process. If additional tests are executed in parallel (which is common with `go test -parallel` or test suites), this can cause cross-test interference, leading to flaky tests and unpredictable behavior.

## Impact

Severity: **medium**

Global mocking of HTTP can cause tests to be non-deterministic if Go's test runner decides to run tests in parallel. This can potentially disrupt not only the test in this file but any test within the same process causing hard-to-track flakiness. The current test is not isolated and implicitly assumes exclusive access to HTTP mocking.

## Location

Second test function, approximately lines 62–88:

## Code Issue

```go
func TestUnitTestBillingPoliciesDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.powerplatform.com/licensing/billingPolicies?api-version=2022-03-01-preview`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("test/datasource/policies/get_billing_policies.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: ...
	})
}
```

## Fix

Consider adding test annotations to this test/entire file to prevent parallel execution, or better, refactor the code to run HTTP mocking per test (if a less global approach is available via the library or a custom HTTP transport). In Go, the easiest fix is to forcibly opt out of parallel execution for these tests:

```go
func TestUnitTestBillingPoliciesDataSource_Validate_Read(t *testing.T) {
    // Prevent this test from running in parallel with others
    // Optionally, put this at the top of both/all httpmock-using tests.
    // t.Parallel() // <-- deliberately do NOT call this

    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
	// ... rest of the test
}
```

Alternatively, document/test gate that the test suite is not intended to be run in parallel while using global mocks. For high-isolation, consider patching the HTTP client objects instead of using such global hooks if the library supports it.
