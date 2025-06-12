# Redundant Combination of `ExpectError` and Empty Check Function

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

In `TestUnitEnvironmentApplicationPackagesDataSource_Validate_No_Dataverse()`, the test step sets both an `ExpectError` and an empty `Check` via `resource.ComposeAggregateTestCheckFunc()`. Specifying an empty check function is redundant when only error validation is needed.

## Impact

Severity: Low

This adds a minor amount of unnecessary cognitive overhead to the test; the purpose of the test step would be clearer if non-applicable code was omitted.

## Location

```go
ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
Check:       resource.ComposeAggregateTestCheckFunc(),
```

## Fix

You may remove the `Check` entirely when the test's only assertion is that an error occurs:

```go
ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
// Remove `Check`
```
Or simply avoid the empty function where no checks are required.
