# Title

Issue with `PowerAppssClient` field in struct definition.

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/models.go

## Problem

The field `PowerAppssClient` in the `EnvironmentPowerAppsDataSource` struct is wrongly named. "PowerApps" is mistakenly spelled with double 's', which could lead to confusion when using or referencing this property later in the code.

## Impact

This typo in the struct field name could confuse contributors working on the codebase, introduce bugs when this name is used for operations, and lead to errors in refactoring or debugging. Severity: **Medium**.

## Location

File: `models.go`

```go
type EnvironmentPowerAppsDataSource struct {
	helpers.TypeInfo
	PowerAppssClient client
}
```

## Code Issue

```go
type EnvironmentPowerAppsDataSource struct {
	helpers.TypeInfo
	PowerAppssClient client
}
```

## Fix

Change `PowerAppssClient` to `PowerAppsClient`.

```go
type EnvironmentPowerAppsDataSource struct {
	helpers.TypeInfo
	PowerAppsClient client
}
```

Explanation:
Correct naming improves readability and ensures consistency throughout the codebase. This fix prevents operational issues resulting from incorrect field referencing.