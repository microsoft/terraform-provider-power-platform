# Unnecessary Re-assignment in Create and Update Methods

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go

## Problem

Inside both `Create` and `Update` functions, the following lines reassign plan fields to themselves, resulting in no effective operation:

```go
plan.Id = types.StringValue(plan.Id.ValueString())
plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
plan.Columns = types.DynamicValue(plan.Columns)
```

These assignments are likely vestigial, left-over from previous implementations or copy-paste. They are redundant since the `plan` fields already hold these values, and this pattern is repeated in both the `Create` and `Update` methods.

## Impact

- **Severity:** Low
- Causes unnecessary confusion and reduces code clarity.
- May cause maintainers to question if a side-effect is expected, leading to misunderstandings.
- Minor performance impact, though negligible.

## Location

In both `Create` (lines ~100-107) and `Update` (lines ~193-200).

## Code Issue

```go
plan.Id = types.StringValue(plan.Id.ValueString())
plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
plan.Columns = types.DynamicValue(plan.Columns)
```

## Fix

Remove these unnecessary reassignments; the fields are already populated via `req.Plan.Get()` and do not need to be reset.

```go
// Remove these lines in both Create and Update methods:
// plan.Id = types.StringValue(plan.Id.ValueString())
// plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
// plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
// plan.Columns = types.DynamicValue(plan.Columns)
```
