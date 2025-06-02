# Title

Missing input validation for DTO fields can lead to unintended behavior.

---

## Path

`/workspaces/terraform-provider-power-platform/internal/services/capacity/dto.go`

---

## Problem

The DTO structures (`capacityDto`, `tenantCapacityDto`, `consumptionDto`) lack mechanisms or embedded validation logic to ensure that fields contain valid or expected values. For instance:
- `CapacityType`, `LicenseModelType`, `TenantId`, and other string fields are not validated for acceptable formats or values.
- Numeric fields like `TotalCapacity`, `MaxCapacity`, `Actual`, and `Rated` do not have constraints for minimum or maximum values.
- The absence of validation could result in processing invalid data (e.g., negative capacities or incorrect tenant IDs).

---

## Impact

Without validation, the system is prone to accepting and processing unexpected or invalid data, leading to possible downstream errors and unreliable computations. This issue has **medium severity**, as it jeopardizes data reliability but does not directly impact system stability unless combined with other vulnerabilities.

---

## Location

**CapacityDto**:
```go
type capacityDto struct {
	TenantId         string              `json:"tenantId"`
	LicenseModelType string              `json:"licenseModelType"`
	TenantCapacities []tenantCapacityDto `json:"tenantCapacities"`
}
```

**TenantCapacityDto**:
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

**ConsumptionDto**:
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

Introduce validations for the DTO fields either in a separate validation layer or within constructors/methods if applicable. For example:

**Updated Code**:
```go
import (
	"errors"
	"strings"
	"github.com/shopspring/decimal"
)

type tenantCapacityDto struct {
	CapacityType  string         `json:"capacityType"`
	CapacityUnits string         `json:"capacityUnits"`
	TotalCapacity decimal.Decimal `json:"totalCapacity"`
	MaxCapacity   decimal.Decimal `json:"maxCapacity"`
	Consumption   consumptionDto  `json:"consumption"`
	Status        string          `json:"status"`
}

func (t *tenantCapacityDto) Validate() error {
	if strings.TrimSpace(t.CapacityType) == "" {
		return errors.New("CapacityType cannot be empty")
	}
	if t.TotalCapacity.LessThan(decimal.Zero) {
		return errors.New("TotalCapacity must be non-negative")
	}
	// Add validation for other fields...
	return nil
}
```

**Explanation**:
- The `Validate` method ensures that fields conform to expected formats and business rules.
- Place similar validation logic for other DTO types (`capacityDto`, `consumptionDto`) to verify integrity.