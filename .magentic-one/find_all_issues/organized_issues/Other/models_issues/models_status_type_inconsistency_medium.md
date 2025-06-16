# Wrong Slice Element Type for Status

##

/workspaces/terraform-provider-power-platform/internal/services/connection/models.go

## Problem

The `ConnectionsDataSourceModel` struct defines `Status` as `[]string`, but its peer, `ResourceModel`, defines `Status` as `types.Set`. This type inconsistency may result in issues when mapping between the models or interfacing with Terraform SDK, which expects the correct type defined in the schema. The discrepancy is likely unintentional and may lead to confusion or subtle runtime bugs.

## Impact

Medium. This inconsistency hurts maintainability and could introduce runtime issues, especially in conversion, serialization, or when relying on uniform interfaces/types across resource/data source models.

## Location

- Structs:
  - `ConnectionsDataSourceModel`
    - `Status []string`
  - `ResourceModel`
    - `Status types.Set`
- File: `/workspaces/terraform-provider-power-platform/internal/services/connection/models.go`

## Code Issue

```go
type ConnectionsDataSourceModel struct {
	// ...
	Status                  []string     `tfsdk:"status"`
	// ...
}

type ResourceModel struct {
	// ...
	Status                  types.Set      `tfsdk:"status"`
	// ...
}
```

## Fix

Choose a single, consistent type for `Status` (likely `types.Set` for interoperability with the Terraform framework) and update all models to match:

```go
type ConnectionsDataSourceModel struct {
	// ...
	Status                  types.Set     `tfsdk:"status"`
	// ...
}
```

This ensures consistency, simplifies conversions, and aligns with the SDK's expectations.
