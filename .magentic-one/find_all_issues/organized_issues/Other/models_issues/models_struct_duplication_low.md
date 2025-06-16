# Struct Duplication for `DataRecordDataSource` and `DataRecordResource`

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/models.go

## Problem

There are two structs, `DataRecordDataSource` and `DataRecordResource`, that have identical fields:

```go
type DataRecordDataSource struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataRecordResource struct {
	helpers.TypeInfo
	DataRecordClient client
}
```

This results in code duplication, making maintenance harder. If a change is needed (e.g., addition or modification of a field), it must be repeated in both structs.

## Impact

Low to Medium. While this duplication is not currently problematic if all usages are separate and justified, it introduces technical debt and increases the risk of drifting definitions in the future.

## Location

As shown above.

## Fix

If possible, consolidate these into a single struct, or embed a common struct:

```go
type DataRecordBase struct {
	helpers.TypeInfo
	DataRecordClient client
}

type DataRecordDataSource struct {
	DataRecordBase
}

type DataRecordResource struct {
	DataRecordBase
}
```

This approach reduces duplication, makes maintenance easier, and clarifies how shared functionality is structured.
