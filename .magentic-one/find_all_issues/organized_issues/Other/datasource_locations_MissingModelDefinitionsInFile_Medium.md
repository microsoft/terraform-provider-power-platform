# Data and Schema Models Not Defined in This File

##

/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations.go

## Problem

The file makes use of `DataSourceModel` and `DataModel` types:

```go
var state DataSourceModel
// ...
state.Value = append(state.Value, DataModel{ ... })
```

However, these types are not defined in this file, nor are they imported explicitly. There is no documentation in this file about their structure, their location, or their intended usage. This can make it difficult for future maintainers or reviewers to quickly verify the shape and correctness of the state expected for the data source.

## Impact

This affects maintainability, code navigation, and onboarding. If their definitions drift from usage, subtle bugs could result. This is of **Medium** severity for code structure/maintainability.

## Location

Methods in this file (`Read`) operate on these models, but there are no type definitions for them nearby or clear references.

## Code Issue

```go
var state DataSourceModel
...
state.Value = append(state.Value, DataModel{
    // fields...
})
```

## Fix

- If models are defined elsewhere, import and reference them clearly, and (optionally) add file-level documentation with comments or import-qualified references.
- If they belong to this package, consider keeping their definitions in this or a closely adjacent file (e.g., `models.go`).
- Add documentation comments with model definitions or usage.

```go
// DataSourceModel and DataModel should be imported or defined here, e.g.:
// type DataSourceModel struct {
//     Value []DataModel
// }
// ...
```

This will improve structure, maintainability, and cross-file understanding.
