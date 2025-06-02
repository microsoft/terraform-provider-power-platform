# Title

Direct Use of Global Mocks May Cause Test Interference

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings_test.go

## Problem

The test file uses functions such as `mocks.ActivateEnvironmentHttpMocks()` that may set global state or persistent responders in httpmock, and uses default httpmock state, which could affect parallel tests. There is no test isolation, so adding `t.Parallel()` or running with `-parallel` could cause interference.

## Impact

Potential for test interference, non-determinism, and flakiness if tests are run in parallel or re-used in other suites. Severity: Medium, as Go test tooling and best practices encourage t.Parallel() and proper test isolation.

## Location

Example from `TestUnitTestEnvironmentSettingsDataSource_Validate_No_Dataverse`:

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()
// then registers additional responders
```

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()
```

## Fix

Where possible:
- Avoid using global state in mocks.
- Isolate httpmock state per test.
- Avoid shared responders unless scope and reset are well controlled.
- Avoid global initialization in test helpers, or use subtests to isolate responders.

E.g., reset state before/after each test, and ensure all global state is cleaned between tests.

No code snippet required unless helper is implemented, but could resemble:

```go
func setupTest(t *testing.T) func() {
    httpmock.Activate()
    mocks.ActivateEnvironmentHttpMocks()
    return func() { httpmock.DeactivateAndReset() }
}

func TestSomething(t *testing.T) {
    teardown := setupTest(t)
    defer teardown()
    // ...
}
```

---

This file will be saved as:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/datasource_environment_settings_test.go_global_mocks_parallelism_medium.md`
