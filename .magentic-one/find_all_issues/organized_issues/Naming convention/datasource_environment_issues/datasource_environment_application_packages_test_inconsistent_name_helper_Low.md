# Inconsistent Use of Test Name Helper Utility

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

In one acceptance test the display name is set using `mocks.TestName()`, in others it's hardcoded to "displayname". This leads to inconsistent naming and could impact test deduplication, parallelization, and diagnostics.

## Impact

Severity: Low

Inconsistencies may lead to confusion during debugging, make it difficult to search logs, and may affect the behavior of tests if uniqueness is required.

## Location

For example:

```go
display_name      = "` + mocks.TestName() + `"
```
vs
```go
display_name                              = "displayname"
```

## Fix

Adopt the same approach for all tests. Preferably use the utility for unique, informative naming:

```go
display_name = "` + mocks.TestName() + `"
```
