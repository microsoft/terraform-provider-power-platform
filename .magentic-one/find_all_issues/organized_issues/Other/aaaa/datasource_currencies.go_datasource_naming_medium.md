# Inconsistent and Non-Descriptive Type Name: `DataSource`

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go

## Problem

The main type for the data source is named `DataSource`. This is overly generic, especially in a package that may contain more than one data source in the future. In the context of the provider codebase, a more descriptive name should be used (e.g., `CurrenciesDataSource`) to improve clarity and avoid conflicts or confusions in larger codebases.

## Impact

Naming ambiguity can result in decreased maintainability and an increased likelihood of referencing or embedding the wrong data source type in the future. The severity increases in larger projects or if code is reused.

**Severity:** Medium

## Location

```go
type DataSource struct {
	helpers.TypeInfo
	CurrenciesClient client
}
```

## Fix

Rename `DataSource` to `CurrenciesDataSource` throughout the file and update all references accordingly.

```go
type CurrenciesDataSource struct {
	helpers.TypeInfo
	CurrenciesClient client
}
``` 
