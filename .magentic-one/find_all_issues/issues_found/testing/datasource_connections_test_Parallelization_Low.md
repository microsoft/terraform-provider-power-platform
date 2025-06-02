# Issue 5

No Parallelization of Unit Tests (t.Parallel)

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections_test.go

## Problem

The unit test function `TestUnitConnectionsDataSource_Validate_Read` does not indicate `t.Parallel()`, which would allow Go to run this test concurrently with others. Adding this directive speeds up test runs and prevents accidental state bleed-through between tests (when safe).

## Impact

- **Severity:** Low
- Potentially slower test execution during development/CI.
- Might miss out on parallel test discovery benefits that can surface data races or side effects when tests are run together.

## Location

```go
func TestUnitConnectionsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
    ...
}
```

## Fix

Call `t.Parallel()` at the beginning of unit test functions, provided the setup and teardown logic supports parallelism and does not have shared state.

```go
func TestUnitConnectionsDataSource_Validate_Read(t *testing.T) {
    t.Parallel()
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
}
```
This allows Go's test runner to execute this and compatible tests more efficiently.
