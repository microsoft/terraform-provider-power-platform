# Title

Redundant Code Block in Create Method

## Path

`/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go`

## Problem

In the `Create` method of the file, the plan object is being redundantly set to itself:

```go
plan.Id = types.StringValue(plan.Id.ValueString())
plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
plan.Columns = types.DynamicValue(plan.Columns)
```

This redundancy does not contribute any meaningful functionality since modifying the `plan` object does not change its functionality.

## Impact

Leaving redundant code reduces code readability and might lead to confusion among developers. Additionally, it adds unnecessary complexity. Severity: **Critical**

## Location

Line located in the `Create` implementation.
File: `/internal/services/data_record/resource_data_record.go`

## Code Issue

```go
plan.Id = types.StringValue(plan.Id.ValueString())
plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
plan.TableLogicalName = types.StringValue(plan.TableLogicalName.ValueString())
plan.Columns = types.DynamicValue(plan.Columns)
```

## Fix

Remove the redundant code as it has no effect and does not modify the `plan` object meaningfully:

```go
// Removed redundant lines
```