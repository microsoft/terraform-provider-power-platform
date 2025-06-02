# Title

Missing Validations for Struct Fields in `DataModel`

##

`/workspaces/terraform-provider-power-platform/internal/services/currencies/models.go`

## Problem

The `DataModel` struct contains fields like `ID`, `Name`, `Type`, `Code`, and `Symbol` which appear to represent critical attributes of a currency entity. However, there are no validations or constraints (e.g., checking for non-empty strings or specific format requirements) on these fields. This could lead to invalid data being processed.

## Impact

Without proper validations, data inconsistencies and potential runtime errors may occur. For example:
- If `ID` is empty, it might break downstream applications assuming valid IDs.
- If `Code` doesn't follow standard currency code formats (e.g., ISO 4217), it could cause interoperability issues.

Severity: **Critical**

## Location

The issue affects the `DataModel` struct in the following file:

`/workspaces/terraform-provider-power-platform/internal/services/currencies/models.go`

## Code Issue

```go
type DataModel struct {
	ID              string `tfsdk:"id"`
	Name            string `tfsdk:"name"`
	Type            string `tfsdk:"type"`
	Code            string `tfsdk:"code"`
	Symbol          string `tfsdk:"symbol"`
	IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}
```

## Fix

Add structure field validation logic to ensure necessary constraints are applied:

```go
type DataModel struct {
	ID              string `tfsdk:"id" validate:"required"`
	Name            string `tfsdk:"name" validate:"required"`
	Type            string `tfsdk:"type" validate:"required,oneof=physical digital"`
	Code            string `tfsdk:"code" validate:"required,min=3,max=3"`
	Symbol          string `tfsdk:"symbol"`
	IsTenantDefault bool   `tfsdk:"is_tenant_default"`
}
```

Explanation:
- `validate:"required"` ensures the field is mandatory.
- `validate:"min=3,max=3"` ensures that the `Code` field conforms to the standard ISO 4217 format.

Additionally, you may want to use a validation library (e.g., `github.com/go-playground/validator`) or custom validation methods for better control over data integrity.
