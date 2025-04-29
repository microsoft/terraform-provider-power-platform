# Title

Error-prone usage of `fmt.Sprintf` in logging statements

##

`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go`

## Problem

Multiple logging statements in the code utilize `fmt.Sprintf` in combination with `tflog.Debug` for logging at runtime. This practice can lead to unnecessary performance overhead. `tflog.Debug` supports formatting arguments directly, which would eliminate the intermediate processing performed by `fmt.Sprintf`.

## Impact

Using `fmt.Sprintf` in logging statements increases computational overhead and can affect performance, especially in heavily used or concurrent environments. This issue is a **medium severity** concern as it undermines logging efficiency, but does not lead to critical breakdowns in functionality.

## Location

Instances found in logging statements such as:

```go
tflog.Debug(ctx, fmt.Sprintf("Dataverse exist in eviroment %t", hasEnvDataverse))
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
tflog.Debug(ctx, fmt.Sprintf("created a resource with ID %s", plan.Id.ValueString()))
tflog.Debug(ctx, fmt.Sprintf("READ: %s with id %s", r.FullTypeName(), state.Id.ValueString()))
```

## Code Issue

### Example

```go
tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
tflog.Debug(ctx, fmt.Sprintf("Dataverse exist in eviroment %t", hasEnvDataverse))
```

## Fix

Use `tflog.Debug`'s built-in argument formatting feature to simplify logging statements.

### Corrected Code

```go
tflog.Debug(ctx, "METADATA:", resp.TypeName)
tflog.Debug(ctx, "Dataverse exist in environment:", hasEnvDataverse)
```

### Explanation

`tflog.Debug` accepts multiple arguments, enabling direct formatting without relying on `fmt.Sprintf`. This approach reduces computational overhead and improves readability.
