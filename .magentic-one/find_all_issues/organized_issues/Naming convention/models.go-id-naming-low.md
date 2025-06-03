# Inconsistent Struct Field Names: "Id" vs "ID"

application/models.go

## Problem

In `EnvironmentApplicationPackageInstallResourceModel`, the field is named `Id` instead of `ID`, which is inconsistent with Go naming conventions (acronyms should be all uppercase, e.g., `ID`). Other similar fields like `EnvironmentId`, `ApplicationId`, and `PublisherId` also do not follow the Go convention of all-uppercase "ID".

## Impact

This breaks Go idiomatic naming conventions, makes the code less readable and maintainable, and can cause subtle issues with automated tools or code generation that expect conventional struct field names.  
Severity: Low

## Location

All affected structs:

- `EnvironmentApplicationPackageInstallResourceModel` (field `Id`)
- `TenantApplicationPackageDataSourceModel` (fields `ApplicationId`, `PublisherId`)

## Code Issue

```go
	Id            types.String   `tfsdk:"id"`
	UniqueName    types.String   `tfsdk:"unique_name"`
	EnvironmentId types.String   `tfsdk:"environment_id"`
```

## Fix

Rename all "Id" suffixes to "ID" to follow Go conventions:

```go
	ID            types.String   `tfsdk:"id"`
	UniqueName    types.String   `tfsdk:"unique_name"`
	EnvironmentID types.String   `tfsdk:"environment_id"`
```

And similarly update other fields throughout the file.
