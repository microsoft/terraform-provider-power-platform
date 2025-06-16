# Use of global interface variable assignments for interface compliance

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

The file uses the following global variable for ensuring the DataSource implements the necessary interfaces:

```go
var (
    _ datasource.DataSource              = &DataSource{}
    _ datasource.DataSourceWithConfigure = &DataSource{}
)
```

This is an idiomatic pattern, but it can be moved closer to the type or into an init function for better contextual locality, visibility, and maintainability. Spreading interface compliance assertions across a large file can reduce readability.

## Impact

Low. It’s purely stylistic/organizational, affecting maintainability only.

## Location

Near the top of the file:

## Code Issue

```go
var (
    _ datasource.DataSource              = &DataSource{}
    _ datasource.DataSourceWithConfigure = &DataSource{}
)
```

## Fix

Consider moving these immediately below the DataSource struct definition or consolidate in comments there for better navigation:

```go
// Ensure DataSource implements the required interfaces.
var _ datasource.DataSource = &DataSource{}
var _ datasource.DataSourceWithConfigure = &DataSource{}
```
