# Naming: Struct Name Does Not Convey Intention

##

/workspaces/terraform-provider-power-platform/internal/services/rest/models.go

## Problem

The struct `DataverseWebApiDatasource` may lead to confusion as it uses "Datasource" (non-standard spelling, standard is "DataSource"). Furthermore, unless this convention is enforced project-wide, using a non-standard spelling can decrease code maintainability and readability for other developers.

## Impact

**Severity: Low**

While not functionally problematic, inconsistent or incorrect naming makes the code less accessible, harder to search and can lead to misunderstandings about the nature of the type.

## Location

```go
type DataverseWebApiDatasource struct {
    helpers.TypeInfo
    DataRecordClient client
}
```

## Code Issue

```go
type DataverseWebApiDatasource struct {
    helpers.TypeInfo
    DataRecordClient client
}
```

## Fix

Use the standard naming convention by renaming the struct to `DataverseWebApiDataSource` throughout the codebase and update references accordingly.

```go
type DataverseWebApiDataSource struct {
    helpers.TypeInfo
    DataRecordClient client
}
```

