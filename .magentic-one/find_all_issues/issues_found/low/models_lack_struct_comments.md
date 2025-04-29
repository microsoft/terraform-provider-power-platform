# Title

Struct definitions lack comprehensive documentation.

##

`/workspaces/terraform-provider-power-platform/internal/services/capacity/models.go`

## Problem

The struct definitions do not include comments explaining their purpose or fields. This reduces code readability and increases onboarding time for developers.

## Impact

Undocumented code can lead to misunderstandings of its functionality, introducing human errors during usage. Severity: **low**.

## Location

All struct definitions.

## Code Issue

### Problematic Code:

```go
type TenantCapacityDataSourceModel struct {
	CapacityType  types.String               `tfsdk:"capacity_type"`
	CapacityUnits types.String               `tfsdk:"capacity_units"`
	Consumption   ConsumptionDataSourceModel `tfsdk:"consumption"`
}
```

## Fix

Add meaningful comments describing the purpose of structs and fields.

```go
// TenantCapacityDataSourceModel holds information about a tenant's capacity usage and limits.
type TenantCapacityDataSourceModel struct {
	// CapacityType represents the type of capacity (standard, premium, etc.).
	CapacityType types.String `tfsdk:"capacity_type"`
	// CapacityUnits indicate the measurement unit for capacity.
	CapacityUnits types.String `tfsdk:"capacity_units"`
	// Consumption provides detailed consumption metrics for the tenant.
	Consumption ConsumptionDataSourceModel `tfsdk:"consumption"`
}
```