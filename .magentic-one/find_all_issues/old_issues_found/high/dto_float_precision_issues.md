# Title

Use of `float32` for representing currency or units may lead to precision issues.

---

## Path

`/workspaces/terraform-provider-power-platform/internal/services/capacity/dto.go`

---

## Problem

The types `tenantCapacityDto` and `consumptionDto` use the `float32` type to represent fields like `TotalCapacity`, `MaxCapacity`, `Actual`, and `Rated`. Utilizing `float32` can result in precision issues, as floating-point arithmetic is inherently imprecise and can lead to rounding errors when performing calculations, especially in cases involving currency or capacity units requiring high accuracy.

---

## Impact

Using `float32` for calculations that require precision can result in unexpected behavior or inaccuracies in results. For example, operations on capacities or consumption data could lead to minor inconsistencies, which might escalate in aggregate calculations. This issue has **high severity** because it directly impacts data integrity.

---

## Location

**tenantCapacityDto**:
```go
type tenantCapacityDto struct {
	CapacityType  string         `json:"capacityType"`
	CapacityUnits string         `json:"capacityUnits"`
	TotalCapacity float32        `json:"totalCapacity"`
	MaxCapacity   float32        `json:"maxCapacity"`
	Consumption   consumptionDto `json:"consumption"`
	Status        string         `json:"status"`
}
```

**consumptionDto**:
```go
type consumptionDto struct {
	Actual          float32 `json:"actual"`
	Rated           float32 `json:"rated"`
	ActualUpdatedOn string  `json:"actualUpdatedOn"`
	RatedUpdatedOn  string  `json:"ratedUpdatedOn"`
}
```

---

## Fix

Replace `float32` with `decimal.Decimal` or `float64` where appropriate. Libraries like `github.com/shopspring/decimal` can help achieve precision as needed, especially for operations involving monetary values or capacitance.

**Updated Code**:
```go
// Import required package for high precision decimal type.
import (
	"github.com/shopspring/decimal"
)

type tenantCapacityDto struct {
	CapacityType  string         `json:"capacityType"`
	CapacityUnits string         `json:"capacityUnits"`
	TotalCapacity decimal.Decimal        `json:"totalCapacity"`
	MaxCapacity   decimal.Decimal        `json:"maxCapacity"`
	Consumption   consumptionDto `json:"consumption"`
	Status        string         `json:"status"`
}

type consumptionDto struct {
	Actual          decimal.Decimal `json:"actual"`
	Rated           decimal.Decimal `json:"rated"`
	ActualUpdatedOn string  `json:"actualUpdatedOn"`
	RatedUpdatedOn  string  `json:"ratedUpdatedOn"`
}
```

**Explanation**:
- By switching to `decimal.Decimal`, we ensure precise computation and control over rounding and fractional units.
- Alternatively, replace `float32` with `float64` for slightly better precision without involving a new dependency.
