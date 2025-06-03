# Repeated Context Management Pattern Should Be Abstracted

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

Each method (`Metadata`, `Schema`, `Configure`, `Read`) includes boilerplate context management with:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
defer exitContext()
```

This pattern is duplicated across the file, leading to unnecessary repetition and possible divergence in future modifications.

## Impact

This duplication hampers maintainability and increases the chance of inconsistencies if changes are required in context management, affecting readability and increasing cognitive load. Severity: **medium**.

## Location

Top of every exported method on `DataLossPreventionPolicyDataSource`.

## Code Issue

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
defer exitContext()
```

## Fix

Encapsulate the context management in a helper, for example:

```go
func withRequestContext(ctx context.Context, ti helpers.TypeInfo, req any, fn func(context.Context)) {
    ctx2, exit := helpers.EnterRequestContext(ctx, ti, req)
    defer exit()
    fn(ctx2)
}
```

And call methods using this wrapper pattern or factor into method receivers.

**Explanation:**  
Reduces duplication, making future updates easier and methods more concise.
