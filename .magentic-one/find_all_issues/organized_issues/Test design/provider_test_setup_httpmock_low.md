# Title

Duplication of HTTP Mock Responders in Multiple Test Functions

##

internal/provider/provider_test.go

## Problem

The two tests for telemetry opt-out (`TestUnitPowerPlatformProvider_Validate_Telementry_Optout_Is_False` and `_Is_True`) duplicate mock registrations and setup logic. This increases maintenance costs and risks the mocks going out of sync if test logic or URLs change.

## Impact

Low. Duplication makes tests harder to maintain.

## Location

Repeated `httpmock.RegisterRegexpResponder` and `httpmock.RegisterResponder` code.

## Fix

Extract out repeated mock setup into a helper function used by both tests:

```go
func setupEnvironmentHttpMocks(t *testing.T) {
    mocks.ActivateEnvironmentHttpMocks()
    httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(...), ...)
    httpmock.RegisterResponder("GET", "...", ...)
}
```

Then call `setupEnvironmentHttpMocks(t)` in both test cases.
