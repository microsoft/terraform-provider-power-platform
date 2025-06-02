# Redundant Embedding of helpers.TypeInfo

##

/workspaces/terraform-provider-power-platform/internal/services/connection/models.go

## Problem

Most structs in this file embed `helpers.TypeInfo`, but there is no indication in the file of the purpose, nor whether every struct requires it. If `helpers.TypeInfo` does not provide needed shared functionality or fields for every struct, its inclusion is unnecessary and clutters the model definition.

## Impact

Low. Redundant embedding increases memory usage (minimally) and confuses the design by suggesting there is shared functionality or a common contract even if there is none. It also makes the models harder to read for new maintainers.

## Location

Repeated in:

- `SharesDataSource`
- `ConnectionsDataSource`
- `ShareResource`
- `Resource`

## Code Issue

```go
type SharesDataSource struct {
	helpers.TypeInfo
	ConnectionsClient client
}
```

## Fix

Review whether every struct truly requires `helpers.TypeInfo`. Remove it from structs where it does not convey actual shared behavior or fields:

```go
type SharesDataSource struct {
	ConnectionsClient client // Removed helpers.TypeInfo if not needed
}
```

This improves model clarity and signals intent accurately.
