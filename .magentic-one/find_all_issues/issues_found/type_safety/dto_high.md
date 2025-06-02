# Issue: Use of float32 for monetary or capacity values

## 
/workspaces/terraform-provider-power-platform/internal/services/capacity/dto.go

## Problem

The `tenantCapacityDto` and `consumptionDto` structs are using the type `float32` for fields representing capacity and consumption values (`TotalCapacity`, `MaxCapacity`, `Actual`, and `Rated`). Using `float32` for quantities that may require precision (such as capacities or monetary values) can result in rounding errors and loss of precision, especially when the values are manipulated or transferred between systems.

## Impact

This issue impacts the precision of capacity-related calculations, potentially affecting business logic relying on precise values. Severity: **High**, especially if the results of these calculations influence billing or resource allocation.

## Location

- tenantCapacityDto: `TotalCapacity float32`, `MaxCapacity float32`
- consumptionDto: `Actual float32`, `Rated float32`

## Code Issue

```go
type tenantCapacityDto struct {
    // ...
    TotalCapacity float32        `json:"totalCapacity"`
    MaxCapacity   float32        `json:"maxCapacity"`
    // ...
}

type consumptionDto struct {
    Actual float32 `json:"actual"`
    Rated  float32 `json:"rated"`
    // ...
}
```

## Fix

Use `float64` instead of `float32` for higher precision and to align with Go's default float type for numeric computations.

```go
type tenantCapacityDto struct {
    // ...
    TotalCapacity float64        `json:"totalCapacity"`
    MaxCapacity   float64        `json:"maxCapacity"`
    // ...
}

type consumptionDto struct {
    Actual float64 `json:"actual"`
    Rated  float64 `json:"rated"`
    // ...
}
```
