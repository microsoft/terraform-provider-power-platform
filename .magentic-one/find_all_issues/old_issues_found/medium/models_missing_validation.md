# Title

Missing validation or constraints on struct fields

## Path

/workspaces/terraform-provider-power-platform/internal/services/languages/models.go

## Problem

The struct definitions for `DataSourceModel` and `DataModel` lack field-level validation or constraints. For example, it doesn't enforce validation rules for strings like `Name`, `ID`, `DisplayName`, and `LocalizedName`, or `int64` values like `LocaleID`. Similarly, non-allowed default values for `IsTenantDefault` may impact logical operations.

## Impact

This absence of validation can lead to inconsistent or incorrect data being stored or processed, which might cause errors or unexpected behaviors. The severity of this gap is **medium** because while the code could function without validation, the risk of introducing subtle bugs increases with invalid inputs or constraints.

## Location

In the struct definitions:

```go
type DataSourceModel struct {
  Timeouts timeouts.Value `tfsdk:"timeouts"`
  Location types.String   `tfsdk:"location"`
  Value    []DataModel    `tfsdk:"languages"`
}

type DataModel struct {
  Name            string `tfsdk:"name"`
  ID              string `tfsdk:"id"`
  DisplayName     string `tfsdk:"display_name"`
  LocalizedName   string `tfsdk:"localized_name"`
  LocaleID        int64  `tfsdk:"locale_id"`
  IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}
```

## Code Issue

No constraints or validation logic is provided.

```go
type DataModel struct {
  Name            string `tfsdk:"name"`
  ID              string `tfsdk:"id"`
  DisplayName     string `tfsdk:"display_name"`
  LocalizedName   string `tfsdk:"localized_name"`
  LocaleID        int64  `tfsdk:"locale_id"`
  IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}
```

## Fix

Adding validation tags or constraints at the field level, or implementing validation logic in methods associated with these structs:

```go
type DataModel struct {
  Name            string `tfsdk:"name" validate:"required,max=100"`           // Name should be required and limited in length
  ID              string `tfsdk:"id" validate:"required,max=36"`             // ID should be required and valid UUID length
  DisplayName     string `tfsdk:"display_name" validate:"required,max=100"`  // DisplayName limits enforced
  LocalizedName   string `tfsdk:"localized_name" validate:"max=100"`         // LocalizedName is optional, but max constraint recommended
  LocaleID        int64  `tfsdk:"locale_id" validate:"gte=0"`                // LocaleID should be a positive integer
  IsTenantDefault bool   `tfsdk:"is_tenant_default"`                         // Can be directly boolean
}
```

Associated validation calls should also be added where these structs are utilized.
