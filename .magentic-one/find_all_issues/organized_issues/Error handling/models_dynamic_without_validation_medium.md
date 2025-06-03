# Use of `types.Dynamic` for Heterogeneous Data Without Explicit Validation

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/models.go

## Problem

The `Rows` and `Columns` fields are typed as `types.Dynamic`, which allows for great flexibility but also opens the door to runtime errors and data inconsistency, especially since there is no documentation or code comment to describe the expected structure for dynamic content.

## Impact

Medium. Unvalidated and undocumented dynamic types increase the risk of runtime panics, unexpected results, or bugs in downstream code that expects a particular shape. This also makes the code harder to maintain over time.

## Location

```go
type DataRecordListDataSourceModel struct {
    // ...
    Rows types.Dynamic  `tfsdk:"rows"`
}

type DataRecordResourceModel struct {
    // ...
    Columns types.Dynamic  `tfsdk:"columns"`
}
```

## Fix

Document the expected shape and content of any `types.Dynamic` fields clearly in code comments or documentation. Introduce validation logic in the code path that manipulates these fields to ensure they contain the required keys and value types, or consider replacing with explicit types if the structure is known and stable.

```go
// Example code comment:
// Columns is expected to be a JSON object with keys...
//
// In resource logic, add schema or runtime validation.

if err := validateColumns(model.Columns); err != nil {
    return fmt.Errorf("invalid columns format: %w", err)
}
```
