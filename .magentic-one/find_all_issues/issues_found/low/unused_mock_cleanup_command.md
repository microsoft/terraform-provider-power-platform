# Title

Unused HTTP Mock Cleanup Command

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles_test.go

## Problem

The `defer httpmock.DeactivateAndReset()` statement is placed but in some cases isn't required because mocks are automatically cleared after tests.

## Impact

While this doesn't directly break the tests, it introduces unnecessary overhead and slightly hinders readability. Severity: Low.

## Location

TestUnitSecurityDataSource_Validate_No_Dataverse

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix

Verify whether the cleanup is required and remove it if redundant.

```go
httpmock.Activate()
// Remove `defer httpmock.DeactivateAndReset()` if unnecessary.
```