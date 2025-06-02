# Title

Inconsistent Use of Composite Fields in Structs

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/models.go`

## Problem

The `ExpandModel` struct uses a composite field (`[]ExpandModel`), which can recursively define its structure. However, other structs like `DataRecordListDataSourceModel` and `DataRecordResourceModel` have non-recursive composite fields and do not maintain consistency. This lack of standardization might lead to confusion and incorrect usage.

## Impact

Inconsistent use of composite fields may introduce challenges in scaling the codebase. It could result in bugs during data serialization/deserialization, affecting the reliability of API responses and data integrity. Severity: **Medium**

## Location

### Found in the following field declaration:

```go
type ExpandModel struct {
	NavigationProperty types.String  `tfsdk:"navigation_property"`
	Select             []string      `tfsdk:"select"`
	Filter             types.String  `tfsdk:"filter"`
	OrderBy            types.String  `tfsdk:"order_by"`
	Top                types.Int64   `tfsdk:"top"`
	Expand             []ExpandModel `tfsdk:"expand"`
}
```

## Code Issue

```go
Expand             []ExpandModel `tfsdk:"expand"`
```

## Fix

Standardize the definition and usage of composite fields across all models. If recursion is required, document the rationale and provide safe ways to handle deep structures (e.g., loop limits or breaking conditions). Alternatively, flatten `ExpandModel` to avoid recursive structures, and maintain uniformity across other model declarations.

### Adjusted Code Example:

```go
type ExpandModel struct {
	NavigationProperty types.String  `tfsdk:"navigation_property"`
	Select             []string      `tfsdk:"select"`
	Filter             types.String  `tfsdk:"filter"`
	OrderBy            types.String  `tfsdk:"order_by"`
	Top                types.Int64   `tfsdk:"top"`
	Expand             []string      `tfsdk:"expand"` // Flatten the composite field
}

// Adding nested-reference handling for deeper ExpandModel iterations
func HandleExpandField(field []string) error {
	for _, expand := range field {
		// Validate and process the values or hierarchy
		// Implement limits to avoid infinite recursive depth
		fmt.Println(expand)
	}
	return nil
}
```

Implementing a flat structure or stricter control mechanisms ensures better maintainability and scalability.

