# Inconsistent Naming for Type and Field Names

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/models.go

## Problem

The code uses a mixture of naming conventions across type and field names. Some structs use the "Model" or "DataSource" suffix (e.g., `DataRecordResourceModel`, `DataRecordListDataSourceModel`, `ExpandModel`), while others use slightly different forms (e.g., `DataRecordResource`, `DataRecordDataSource`). This inconsistency can confuse maintainers and users, and makes the code less readable and predictable.

## Impact

Medium. Inconsistent naming increases the likelihood of confusion when using or extending the codebase, especially for new contributors.

## Location

Multiple locations throughout the file:

```go
type DataRecordDataSource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type ExpandModel struct {
    ...
}

type DataRecordListDataSourceModel struct {
    ...
}

type DataRecordResource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataRecordResourceModel struct {
    ...
}
```

## Fix

Adopt a consistent naming convention for all data model types and field names, such as always using `Model` for structs that represent Terraform model schemas, and always suffixing resource types with `Resource` or `DataSource` clearly.

```go
// Consistent naming using Model and Resource suffixes:
type DataRecordResourceModel struct { ... }
type DataRecordListDataSourceModel struct { ... }
type ExpandModel struct { ... }
```

Further, ensure that corresponding fields across types consistently use the same naming, such as always referring to table names, record IDs, or similar concepts the same way in all types.

