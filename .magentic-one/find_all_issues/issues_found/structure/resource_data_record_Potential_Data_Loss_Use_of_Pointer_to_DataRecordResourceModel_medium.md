# Potential Data Loss: Use of Pointer to DataRecordResourceModel in Update/Read

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go

## Problem

The functions `Update` and `Read` use pointers for `*DataRecordResourceModel` when reading state/plan. The rest of the code and Terraform conventions expect passing by value to ensure correct zero-value/unknown behavior, and because the struct is not especially large.

## Impact

- **Severity:** Medium
- Can cause nil dereference panics if the struct is not initialized.
- Can result in partial or incorrect updates as pointer fields may not accurately reflect unknowns or non-set fields.

## Location

`Read`, `Update` function signatures and variable declarations.

## Code Issue

```go
var state *DataRecordResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
// ...
var plan *DataRecordResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
```

## Fix

Use value types instead of pointers, and pass value addresses:

```go
var state DataRecordResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

var plan DataRecordResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
```

This ensures `state` is always properly initialized and compatible with Terraform conventions.
