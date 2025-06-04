# Potential Type Safety Issue: Usage of Unspecified `client` Type

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go

## Problem

The field `CurrenciesClient` in `DataSource` is typed as `client`, but there is no import or definition provided in this file or through visible imports. This could lead to ambiguity or errors if the wrong type is used, or if it conflicts with other `client` types in the project.

## Impact

This may cause compile-time errors, reduced code readability, and type confusion if multiple `client` types exist. Future maintainers may import the wrong type or misunderstand the dependency. This introduces a risk for refactoring and integration.

**Severity:** Medium

## Location

```go
type DataSource struct {
	helpers.TypeInfo
	CurrenciesClient client
}
```

## Fix

Use a more descriptive and package-qualified type, such as `api.CurrenciesClient` if that is the intended type, and ensure it is imported correctly.

```go
type DataSource struct {
	helpers.TypeInfo
	CurrenciesClient api.CurrenciesClient // Or the correct fully qualified type
}
```
