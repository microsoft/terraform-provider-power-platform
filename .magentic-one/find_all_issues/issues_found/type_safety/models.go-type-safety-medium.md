# Potential Type-Safety Issues with Loose Use of types.String and types.Int64

application/models.go

## Problem

All fields use the generic `types.String` and `types.Int64` types from the Terraform framework. This prevents the Go compiler from helping you enforce stricter typing or enum-like safety for fields that have specific, limited allowed values (for example, `ApplicationVisibility`, `CatalogVisibility`, `ErrorCode`, etc.).

## Impact

Loose typing can result in invalid or unexpected values being assigned to these struct fields at runtime. This could introduce bugs, hurt maintainability, and make validation harder, as there is no compile-time enforcement of allowed values.  
Severity: Medium

## Location

Example locations (many fields):

- `TenantApplicationPackageDataSourceModel.ApplicationVisibility`
- `TenantApplicationPackageDataSourceModel.CatalogVisibility`
- `TenantApplicationErrorDetailsDataSourceModel.ErrorCode`
- etc.

## Code Issue

```go
	ApplicationVisibility  types.String                                   `tfsdk:"application_visibility"`
	CatalogVisibility      types.String                                   `tfsdk:"catalog_visibility"`
	LastError              []TenantApplicationErrorDetailsDataSourceModel `tfsdk:"last_error"`
```

## Fix

For fields with specific allowed values, consider defining custom Go types or enums and using them internally, with conversion/validation functions when reading from `types.String`. For example:

```go
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

type TenantApplicationPackageDataSourceModel struct {
	// ...
	ApplicationVisibility Visibility `tfsdk:"application_visibility"`
	// ...
}
```

Provide validation/marshal/unmarshal helpers as needed.
