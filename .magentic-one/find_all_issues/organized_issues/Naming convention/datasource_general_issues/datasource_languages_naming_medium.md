# Title

Inconsistent and unclear naming for DataSource and DataSourceModel

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go

## Problem

The struct `DataSource` (referenced in this file) is named ambiguously; its use does not clarify it is specifically for the Languages data source, and neither do the models like `DataSourceModel` and `DataModel`. They are generic and could lead to confusion, notably as the codebase expands or is read by others. Naming should reflect the domain and use case for clarity and maintainability.

## Impact

Poor naming conventions decrease codebase maintainability and readability, increasing the cognitive load for new collaborators and risk for naming clashes or mistakes. Severity: **Medium**.

## Location

Usage throughout the file, especially:

```go
var (
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)
func NewLanguagesDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "languages",
		},
	}
}
...
var state DataSourceModel
...
for _, language := range languages.Value {
	state.Value = append(state.Value, DataModel{
		ID:              language.ID,
		Name:            language.Name,
		DisplayName:     language.Properties.DisplayName,
		LocalizedName:   language.Properties.LocalizedName,
		LocaleID:        language.Properties.LocaleID,
		IsTenantDefault: language.Properties.IsTenantDefault,
	})
}
```

## Code Issue

```go
type DataSource struct { ... }
type DataSourceModel struct { ... }
type DataModel struct { ... }
...
func NewLanguagesDataSource() datasource.DataSource { ... }
```

## Fix

Rename the types and functions to specifically reference `Languages` so their intent is clear and unambiguous throughout the codebase.

```go
type LanguagesDataSource struct { ... }
type LanguagesDataSourceModel struct { ... }
type LanguagesDataModel struct { ... }

func NewLanguagesDataSource() datasource.DataSource {
	return &LanguagesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "languages",
		},
	}
}
```
