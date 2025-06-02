# Direct Use of `append` Without Resetting `state.Value` in DataSource `Read`

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go

## Problem

In the `Read` method of the data source, the slice `state.Value` is directly appended to in a loop without first resetting or truncating it. In situations where `Read` is called multiple times within the same lifecycle (for example, with partial state), appending could result in duplicate or stale data.

## Impact

This can cause duplicated currency entries and inconsistent data to be returned to the user, violating expectations and potentially breaking infrastructure operations or causing confusion.

**Severity:** Medium

## Location

```go
for _, location := range currencies.Value {
	state.Value = append(state.Value, DataModel{
		ID:              location.ID,
		Name:            location.Name,
		Type:            location.Type,
		Code:            location.Properties.Code,
		Symbol:          location.Properties.Symbol,
		IsTenantDefault: location.Properties.IsTenantDefault,
	})
}
```

## Fix

Always reset `state.Value = nil` before appending data during each read to ensure idempotency.

```go
state.Value = nil // Reset before populating
for _, location := range currencies.Value {
	state.Value = append(state.Value, DataModel{
		ID:              location.ID,
		Name:            location.Name,
		Type:            location.Type,
		Code:            location.Properties.Code,
		Symbol:          location.Properties.Symbol,
		IsTenantDefault: location.Properties.IsTenantDefault,
	})
}
```
