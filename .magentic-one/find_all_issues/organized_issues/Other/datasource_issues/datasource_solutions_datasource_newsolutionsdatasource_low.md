# Naming Consistency: Function Name Does Not Reflect DataSource Type

##

/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go

## Problem

The constructor function is named `NewSolutionsDataSource` but returns a `DataSource` struct, which is a generic name. There may be confusion if multiple datasource types exist (e.g., for Solution vs other entities). It's more conventional for factory/constructor functions to match the struct or to be explicit about their purpose for clarity and maintainability.

## Impact

**Low**. Primarily an issue of maintainability and clarity for future contributors. It can make the codebase less discoverable and less consistent, particularly as the project scales or more data sources are introduced.

## Location

Line 16-23:

## Code Issue

```go
func NewSolutionsDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "solutions",
		},
	}
}
```

## Fix

Either:

- Rename the struct to `SolutionsDataSource`, or
- (Recommended) Keep the struct as `DataSource` but add a comment to clarify its generality, or
- If this pattern is already normed in the project (single generic DataSource struct), leave as-is but document.

Example comment addition for clarity:

```go
// NewSolutionsDataSource returns a new DataSource configured for "solutions".
// The DataSource struct is generic and reused for various entity types.
func NewSolutionsDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "solutions",
		},
	}
}
```

