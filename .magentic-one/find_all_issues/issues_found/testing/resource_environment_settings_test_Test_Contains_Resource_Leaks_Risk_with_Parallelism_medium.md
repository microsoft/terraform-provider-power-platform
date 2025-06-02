# Test Contains Resource Leaks Risk with Parallelism

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resource_environment_settings_test.go

## Problem

Tests use global HTTP mocking (`httpmock.Activate`) and mutate a package-global variable (`getOrgInx`). This variable is not protected against concurrent access if tests are ever run in parallel (via `t.Parallel()`), creating a data race risk.

## Impact

Medium to high: If future sub-tests or parallel test execution are enabled, global variables produce nondeterministic and flaky tests (hard to debug).

## Location

All tests using:

```go
var getOrgInx = 0
// usage:
getOrgInx++
```

## Code Issue

```go
var getOrgInx = 0

// Within responder:
getOrgInx++
```

## Fix

Move `getOrgInx` to be a local variable and pass via closure for each test, and (optionally) protect with a mutex if truly needed globally.

```go
test := func(t *testing.T) {
    var getOrgInx int
    httpmock.RegisterResponder("GET", "...",
        func(req *http.Request) (*http.Response, error) {
            getOrgInx++
            // ...
        },
    )
    // ...
}
```

Or use:

```go
var (
    getOrgInx int
    mu        sync.Mutex
)
httpmock.RegisterResponder("GET", "...",
    func(req *http.Request) (*http.Response, error) {
        mu.Lock()
        getOrgInx++
        mu.Unlock()
        // ...
    },
)
```

However, keep mutable state as local as possible for test reliability.
