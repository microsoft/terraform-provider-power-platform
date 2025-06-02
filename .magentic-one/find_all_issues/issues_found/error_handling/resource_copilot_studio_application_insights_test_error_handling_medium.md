# Title

Error handling: Possible Resource Leak if httpmock.DeactivateAndReset is Not Run

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights_test.go

## Problem

`httpmock.DeactivateAndReset()` is deferred immediately after `httpmock.Activate()` in the test functions. If for some reason `Activate()` fails (e.g. panics, though unlikely in this context), this may lead to inconsistent test state in future tests—especially since the provider testing framework sometimes executes tests in parallel or in unpredictable order. Additionally, resource cleanup is not mentioned for other mocks or global state changes.

## Impact

Medium severity. While this is a rare occurrence, improper cleanup in tests might lead to flakiness or cross-test contamination, especially as the suite grows or parallelization is enabled.

## Location

Unit test functions activating/deactivating httpmock, such as:

```go
func TestUnitCopilotStudioApplicationInsights_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	...
}
```

and

```go
func TestUnitCopilotStudioApplicationInsights_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	...
}
```

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix

Consider using subtests (`t.Run`) to scope setup/teardown, or use `t.Cleanup` (Go 1.14+) for better integration with the test runner. This ensures cleanup runs only if activation is successful and properly tracks nested test state.

```go
if err := httpmock.Activate(); err != nil {
    t.Fatalf("failed to activate httpmock: %v", err)
}
t.Cleanup(func() {
    httpmock.DeactivateAndReset()
})
```

---

I will continue scanning for additional issues in this file.
