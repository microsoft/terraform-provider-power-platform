# Struct Embedding: Unclear Purpose of Embedded `helpers.TypeInfo`

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/models.go

## Problem

`helpers.TypeInfo` is embedded in the `EnvironmentPowerAppsDataSource` struct, but it is unclear why. If this is not documented or immediately useful, such embedding can reduce code clarity.

## Impact

Unintentional or unclear struct embedding can hinder codebase readability and may cause unexpected issues with struct field/method shadowing (severity: low).

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

Add a comment explaining the purpose of embedding `helpers.TypeInfo`. If not needed, remove the embedding:

```go
// TypeInfo provides metadata about the datasource type
helpers.TypeInfo
```
or
```go
type EnvironmentPowerAppsDataSource struct {
	PowerAppsClient client
}
```
