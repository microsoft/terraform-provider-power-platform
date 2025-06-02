# Title

Locally Defined Dynamic Validator Usage Without Safe State Validation

## Path

`/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go`

## Problem

In the method `ConfigValidators`, the validator `DynamicColumns` is invoked without ensuring that the associated attribute `columns` has been safely validated against unsupported states or ranged values:

```go
func (d *DataRecordResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        DynamicColumns(
            path.Root("columns").Expression(),
        ),
    }
}
```

A stricter validation mechanism is required to handle non-conforming `columns`. Without this, malformed or improperly defined columns may cause ambiguous downstream behavior.

## Impact

This omission risks misconfigurations propagating into runtime or affecting data-oriented operations. Furthermore, debugging these edge-case misconfigurations could become complex. Severity: **High**

## Location

Method: `ConfigValidators`
File: `/internal/services/data_record/resource_data_record.go`
Line Location: DynamicColumns(Path Usage)

## Code Issue

```go
func (d *DataRecordResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        DynamicColumns(
            path.Root("columns").Expression(),
        ),
    }
}
```

## Fix

Enhance input validations for dynamic column definitions prior to invoking `DynamicColumns`. Example:

```go
func (d *DataRecordResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
    columnValidator = ValidateForEmptyOrUnexpectedType(path.Root("columns"))

    if columnValidator.HasIssues() {
       throw UserDiagnostics.Append(errorStatement);
    }

    return MainLogicValidator(UseExtension DynamicColumns);
}
```