# Issue: List appending in loop without pre-allocation

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections.go

## Problem

In the `Read` method, when building the `state.Connections` slice, the current implementation uses `append()` within a loop over `connections` without pre-allocating the size of the slice. This can cause multiple memory reallocations, which is suboptimal, especially when handling large lists.

## Impact

Severity: **Low**

While not a functional bug, this practice is less performant for large lists due to frequent memory reallocations and copying when the slice grows.

## Location

```go
for _, connection := range connections {
	connectionModel := ConvertFromConnectionDto(connection)
	state.Connections = append(state.Connections, connectionModel)
}
```

## Code Issue

```go
for _, connection := range connections {
	connectionModel := ConvertFromConnectionDto(connection)
	state.Connections = append(state.Connections, connectionModel)
}
```

## Fix

Pre-allocate the slice with the expected length to avoid unnecessary reallocations:

```go
state.Connections = make([]ConnectionsDataSourceModel, 0, len(connections))
for _, connection := range connections {
	connectionModel := ConvertFromConnectionDto(connection)
	state.Connections = append(state.Connections, connectionModel)
}
```
