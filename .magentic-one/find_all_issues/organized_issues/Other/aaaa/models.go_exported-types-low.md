# Structure: Lack of Documentation for Exported Types

##

/workspaces/terraform-provider-power-platform/internal/services/rest/models.go

## Problem

Exported types and structs (`DataverseWebApiDatasource`, `DataverseWebApiDatasourceModel`) lack doc comments. In Go, all exported elements should be documented for public usage and tooling support.

## Impact

**Severity: Low**

This does not directly affect functionality but hampers both human understanding and automated documentation tools.

## Location

```go
type DataverseWebApiDatasource struct {
    helpers.TypeInfo
    DataRecordClient client
}

type DataverseWebApiDatasourceModel struct {
    ...
}
```

## Code Issue

```go
type DataverseWebApiDatasource struct {
    helpers.TypeInfo
    DataRecordClient client
}

type DataverseWebApiDatasourceModel struct {
    ...
}
```

## Fix

Add appropriate doc comments above each exported type.

```go
// DataverseWebApiDataSource represents a data source for interacting with the Dataverse Web API.
type DataverseWebApiDataSource struct {
    helpers.TypeInfo
    DataRecordClient Client
}

// DataverseWebApiDataSourceModel defines the schema for Dataverse Web API data source in Terraform.
type DataverseWebApiDataSourceModel struct {
    ...
}
```
