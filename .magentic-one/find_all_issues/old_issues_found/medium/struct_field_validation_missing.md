# Title

Struct Field Validation Is Missing

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/dto.go

## Problem

The struct `TenantDto` does not enforce any validation on its fields. This means that invalid or unexpected values could be assigned to its properties. For example, if `TenantId` is required to be non-empty or `Location` should follow a specific format, there is no validation mechanism in place to enforce these constraints.

## Impact

The lack of validation could lead to errors during runtime, especially when interacting with external systems or databases that expect specific formats or constraints. This increases the likelihood of bugs or data inconsistencies in the application, making the severity of this issue medium.

## Location

`TenantDto` struct in `dto.go`.

## Code Issue

```go
type TenantDto struct {
	TenantId                         string `json:"tenantId"`
	State                            string `json:"state"`
	Location                         string `json:"location"`
	AadCountryGeo                    string `json:"aadCountryGeo"`
	DataStorageGeo                   string `json:"dataStorageGeo"`
	DefaultEnvironmentGeo            string `json:"defaultEnvironmentGeo"`
	AadDataBoundary                  string `json:"aadDataBoundary"`
	FedRAMPHighCertificationRequired bool   `json:"fedRAMPHighCertificationRequired"`
}
```

## Fix

Introduce validation methods or libraries to ensure correct values for the fields in the `TenantDto` struct. Use packages like `github.com/go-playground/validator` for field validation.

```go
import (
	"github.com/go-playground/validator/v10"
)

type TenantDto struct {
	TenantId                         string `json:"tenantId" validate:"required"`
	State                            string `json:"state" validate:"required,oneof=active inactive"`
	Location                         string `json:"location" validate:"required,min=3"`
	AadCountryGeo                    string `json:"aadCountryGeo" validate:"omitempty,alpha"`
	DataStorageGeo                   string `json:"dataStorageGeo" validate:"required"`
	DefaultEnvironmentGeo            string `json:"defaultEnvironmentGeo"`
	AadDataBoundary                  string `json:"aadDataBoundary"`
	FedRAMPHighCertificationRequired bool   `json:"fedRAMPHighCertificationRequired" validate:"-"`
}

// Use validate struct method to call validations
func ValidateTenantDto(dto TenantDto) error {
	validate := validator.New()
	return validate.Struct(dto)
}
```

This fix ensures that any time `TenantDto` is used, the values are validated before proceeding.
