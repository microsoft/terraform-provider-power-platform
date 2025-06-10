# Title

Struct fields may not conform to Go naming conventions

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

The struct fields used in the `state.TenantCapacities` slice and the corresponding `TenantCapacityDataSourceModel` and `ConsumptionDataSourceModel` structs (presumably defined elsewhere) use names like `CapacityType`, `CapacityUnits`, and `TotalCapacity`, which match the upstream JSON keys but should be capitalized and formatted according to Go naming conventions when used in Go struct definitions. 

Additionally, field names like `TenantId` should ideally be `TenantID` (per Go's initialism guidelines).

## Impact

Medium severity. While the code functions, improper naming can hinder readability and long-term maintainability, and doesn't conform to Go style guidelines, reducing clarity for future maintainers.

## Location

Throughout the read function (struct assignments)

## Code Issue

```go
state.TenantId = types.StringValue(tenantCapacityDto.TenantId)
// ...
TenantCapacityDataSourceModel{
    CapacityType:  types.StringValue(capacity.CapacityType),
    CapacityUnits: types.StringValue(capacity.CapacityUnits),
    // ...
}
```

## Fix

Follow Go's naming conventions for initialisms. For example, rename `TenantId` to `TenantID`, and similarly update struct field definitions and variable assignments. Ensure all references match the revised names.

```go
state.TenantID = types.StringValue(tenantCapacityDto.TenantID)
// ...
TenantCapacityDataSourceModel{
    CapacityType:  types.StringValue(capacity.CapacityType),
    CapacityUnits: types.StringValue(capacity.CapacityUnits),
    // ...
}
```

Note: This example assumes that the corresponding struct definitions are updated to match the revised naming.
