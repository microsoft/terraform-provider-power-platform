# Redundant/Repetitive HTTP Mock Responders Setup

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

The same set of HTTPMock responders is re-registered in multiple test functions, introducing repetitive and redundant code. This makes maintaining the tests unnecessarily difficult and increases the risk of inconsistencies.

## Impact

Severity: Medium

Duplicated code is harder to maintain and synchronize. Bugs or updates in mock setup require multiple coordinated changes, raising risk for inconsistent test results.

## Location

All three main unit test functions, e.g.:

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()

httpmock.RegisterResponder("GET", `https://api.powerplatform.com/...`,
    func(req *http.Request) (*http.Response, error) {
        // ...
    })
```

## Fix

Extract HTTPMock responder setup into a helper function that is reused in each test function.

```go
func setupEnvironmentApplicationPackagesMocks(testVariant string) {
    httpmock.RegisterResponder(..., // use testVariant for test file selection
        func(req *http.Request) (*http.Response, error) {
            // ...
        })
    // Add the other responders
}

// Call this in the tests:
setupEnvironmentApplicationPackagesMocks("Validate_Read")
```
