# Title

Missing nil checks for nested structs in response handling

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

When mapping values from `capacity.Consumption` in the Read function, the code assumes that `capacity.Consumption` is never nil. If the API response contains a capacity object without a `Consumption` object, this will cause a runtime panic due to dereferencing a nil pointer.

## Impact

High severity. This can cause the entire provider operation to fail with a panic if upstream data does not guarantee non-nil `Consumption` fields. This is both a stability and safety concern.

## Location

Function: `Read`, during state mapping

## Code Issue

```go
Consumption: ConsumptionDataSourceModel{
    Actual:          types.Float32Value(capacity.Consumption.Actual),
    Rated:           types.Float32Value(capacity.Consumption.Rated),
    ActualUpdatedOn: types.StringValue(capacity.Consumption.ActualUpdatedOn),
    RatedUpdatedOn:  types.StringValue(capacity.Consumption.RatedUpdatedOn),
},
```

## Fix

Add a nil check before accessing fields of `capacity.Consumption`:

```go
var consumption ConsumptionDataSourceModel
if capacity.Consumption != nil {
    consumption = ConsumptionDataSourceModel{
        Actual:          types.Float32Value(capacity.Consumption.Actual),
        Rated:           types.Float32Value(capacity.Consumption.Rated),
        ActualUpdatedOn: types.StringValue(capacity.Consumption.ActualUpdatedOn),
        RatedUpdatedOn:  types.StringValue(capacity.Consumption.RatedUpdatedOn),
    }
}
// ...
Consumption: consumption,
```

This ensures safe handling in the absence of nested data.
