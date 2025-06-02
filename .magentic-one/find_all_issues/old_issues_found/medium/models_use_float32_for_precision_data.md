# Title

Use of `types.Float32` for handling financial or capacity data in `TenantCapacityDataSourceModel`.

##

`/workspaces/terraform-provider-power-platform/internal/services/capacity/models.go`

## Problem

The `types.Float32` type is used for `total_capacity`, `max_capacity`, `actual`, and `rated`. Float32 has limited precision and may result in rounding errors for values that require finer granularity.

## Impact

Potential inaccuracies in capacity calculations can propagate errors in platform configurations. Severity: **medium**.

## Location

`TenantCapacityDataSourceModel` and `ConsumptionDataSourceModel`.

## Code Issue

### Problematic Code:

```go
TotalCapacity types.Float32 `tfsdk:"total_capacity"`
MaxCapacity   types.Float32 `tfsdk:"max_capacity"`
Actual        types.Float32 `tfsdk:"actual"`
Rated         types.Float32 `tfsdk:"rated"`
```

## Fix

Replace `types.Float32` with a higher precision type, such as `types.Float64`.

```go
TotalCapacity types.Float64 `tfsdk:"total_capacity"`
MaxCapacity   types.Float64 `tfsdk:"max_capacity"`
Actual        types.Float64 `tfsdk:"actual"`
Rated         types.Float64 `tfsdk:"rated"`
```