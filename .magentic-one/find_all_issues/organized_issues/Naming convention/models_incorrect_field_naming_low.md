# Incorrect Field Naming in Struct: `PowerAppssClient`

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/models.go

## Problem

The field `PowerAppssClient` in the `EnvironmentPowerAppsDataSource` struct appears to have a naming error. The double "s" ("PowerAppss") is likely unintentional and does not conform to Go naming conventions, which prefer clear, concise, and consistently named fields. This could cause confusion and maintenance issues.

## Impact

Incorrect or inconsistent field names reduce code readability and maintainability, increasing the risk of bugs or misunderstandings (severity: low).

## Location

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

Change the field name from `PowerAppssClient` to `PowerAppsClient` for consistency and clarity:

```go
type EnvironmentPowerAppsDataSource struct {
	helpers.TypeInfo
	PowerAppsClient client
}
```
