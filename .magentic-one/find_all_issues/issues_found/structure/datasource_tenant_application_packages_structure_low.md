# Title

Type Naming Consistency: Obscure/Composite Names

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go

## Problem

The generated types and methods show potential naming inconsistencies such as `TenantApplicationPackagesListDataSourceModel`, `TenantApplicationPackageDataSourceModel`, etc. While this pattern isn't unusual, consider reviewing if these names are the most meaningful and general (e.g., use `TenantApplicationModel`, `TenantApplicationListModel`) for future extensibility and readability.

## Impact

- Low: Naming inconsistency is less severe since it's a local codebase and not a public API; however, it impacts maintainability.

## Location

- All usages of *DataSourceModel types.

## Code Issue

```go
var state TenantApplicationPackagesListDataSourceModel
...
app := TenantApplicationPackageDataSourceModel{
	// ...
}
```

## Fix

**Consider renaming for simplicity, e.g.:**

```go
var state TenantApplicationListModel
...
app := TenantApplicationModel{
	// ...
}
```
