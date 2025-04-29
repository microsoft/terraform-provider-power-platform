# Title

Potential null pointer dereference when accessing `location.Properties`

##

`/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go`

## Problem

In the `Read` method, the loop over `currencies.Value` accesses `location.Properties` without prior null checks. If `location.Properties` is `nil`, this will trigger a runtime panic.

## Impact

This could cause the entire provider to crash during execution. Such runtime crashes are categorized as **Critical** severity because they disrupt functionality entirely.

## Location

The issue is in the `Read` method, within the following loop block:

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

## Code Issue

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

Add a null check for `location.Properties` before accessing its members to ensure safe operations and avoid runtime errors.

```go
for _, location := range currencies.Value {
    if location.Properties == nil {
        tflog.Warn(ctx, "Skipping location with null Properties", map[string]interface{}{
            "locationID": location.ID,
        })
        continue
    }
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