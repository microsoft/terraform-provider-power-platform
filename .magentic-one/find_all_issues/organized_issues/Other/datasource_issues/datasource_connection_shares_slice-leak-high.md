# Appending shares to a possibly non-empty state slice

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

On every `Read`, the current implementation appends to `state.Shares` without clearing or re-instantiating the slice. If a `Read` runs multiple times for the same Terraform resource (in runs, re-plans, or refreshes), the shares slice may accumulate/duplicate data.

## Impact

**Severity: high**

Leads to state bloat, possible provider inconsistencies, and bugs in state reconciliation.

## Location

```go
for _, connection := range connectionsList.Value {
	connectionModel := ConvertFromConnectionSharesDto(connection)
	state.Shares = append(state.Shares, connectionModel)
}
```

## Code Issue

```go
for _, connection := range connectionsList.Value {
	connectionModel := ConvertFromConnectionSharesDto(connection)
	state.Shares = append(state.Shares, connectionModel)
}
```

## Fix

Clear the `state.Shares` at the start of the function or just before the loop:

```go
state.Shares = make([]SharesDataSourceModel, 0, len(connectionsList.Value))
for _, connection := range connectionsList.Value {
	connectionModel := ConvertFromConnectionSharesDto(connection)
	state.Shares = append(state.Shares, connectionModel)
}
```
