# Title

Fields Without Validation

## 

/workspaces/terraform-provider-power-platform/internal/services/tenant/models.go

## Problem

The fields in `DataSourceModel` lack field-level validation or constraints, which can lead to invalid or empty data being processed. Additionally, it doesn't specify any validation methods for required fields, lengths, or formats (e.g., valid Tenant ID format for `TenantId`).

## Impact

- It may lead to runtime errors and unreliable data processing.
- Missed validations can increase debugging time and effort.

**Severity: High**

## Location

```go
type DataSourceModel struct {
	TenantId                         types.String `tfsdk:"tenant_id"`
	State                            types.String `tfsdk:"state"`
	Location                         types.String `tfsdk:"location"`
	AadCountryGeo                    types.String `tfsdk:"aad_country_geo"`
	DataStorageGeo                   types.String `tfsdk:"data_storage_geo"`
	DefaultEnvironmentGeo            types.String `tfsdk:"default_environment_geo"`
	AadDataBoundary                  types.String `tfsdk:"aad_data_boundary"`
	FedRAMPHighCertificationRequired types.Bool   `tfsdk:"fed_ramp_high_certification_required"`
}
```

## Fix

Introduce proper validation for fields such as `TenantId`, `State`, and `Location`. This can be achieved by wrapping the fields in a method or incorporating a validation mechanism.

Example:
```go
type DataSourceModel struct {
	TenantId                         types.String `tfsdk:"tenant_id" validate:"required,min=1,max=36"`
	State                            types.String `tfsdk:"state" validate:"oneof=enabled disabled"`
	Location                         types.String `tfsdk:"location" validate:"required,min=2"`
	AadCountryGeo                    types.String `tfsdk:"aad_country_geo" validate:"country"`
	DataStorageGeo                   types.String `tfsdk:"data_storage_geo" validate:"geo"`
	DefaultEnvironmentGeo            types.String `tfsdk:"default_environment_geo"`
	AadDataBoundary                  types.String `tfsdk:"aad_data_boundary"`
	FedRAMPHighCertificationRequired types.Bool   `tfsdk:"fed_ramp_high_certification_required"`
}
```

You may also use Go validators like `go-playground/validator` to handle such validations systematically.