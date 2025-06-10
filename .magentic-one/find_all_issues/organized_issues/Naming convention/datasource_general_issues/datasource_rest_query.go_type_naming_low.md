# Title

Type Naming: DataverseWebApiDatasource Has Inconsistent Casing

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

The main type in this file is named `DataverseWebApiDatasource`, which mixes "Api" and "Datasource" casing. In Go, the recommended convention based on acronym usage is either "API" or "Datasource"/"DataSource" but not partial. Additionally, the naming is inconsistent with the commonly used "DataSource" in the Terraform Plugin SDK community.

## Impact

Inconsistent or unconventional naming can cause confusion and decrease maintainability, as contributors might be unsure whether the correct form is "API", "Api", "Datasource", or "DataSource". This is a low-severity issue related to code clarity.

## Location

Throughout the file as the type name and references.

## Code Issue

```go
type DataverseWebApiDatasource struct {
  //...
}
```

## Fix

Rename the struct and relevant usages to `DataverseWebAPIDatasource` (or, if following the Go/TF convention strongly, `DataverseWebAPIDataSource`):

```go
type DataverseWebAPIDatasource struct {
  //...
}
```

And update all references accordingly. This will bring clarity and consistency to the codebase.

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/datasource_rest_query.go_type_naming_low.md`
