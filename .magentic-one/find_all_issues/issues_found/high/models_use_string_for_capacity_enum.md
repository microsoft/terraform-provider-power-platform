# Title

Using `types.String` for `CapacityType` and `CapacityUnits` fields for structured data.

##

`/workspaces/terraform-provider-power-platform/internal/services/capacity/models.go`

## Problem

Capacity types and units are stored as `types.String`. While this provides flexibility, it does not validate against expected data formats (e.g., enums or specific predefined constants). This can lead to input errors and ambiguity in usage.

## Impact

Incorrect or invalid capacity types or units may introduce issues during runtime interpretation. Severity: **high**.

## Location

`TenantCapacityDataSourceModel` struct.

## Code Issue

### Problematic Code:

```go
CapacityType  types.String `tfsdk:"capacity_type"`
CapacityUnits types.String `tfsdk:"capacity_units"`
```

## Fix

Replace `types.String` with a specific type or validation mechanism to ensure accuracy.

```go
type CapacityType string

const (
	CapacityTypeStandard CapacityType = "standard"
	CapacityTypePremium  CapacityType = "premium"
)

type TenantCapacityDataSourceModel struct {
	CapacityType  CapacityType `tfsdk:"capacity_type"`
	CapacityUnits types.String `tfsdk:"capacity_units"`
}
```