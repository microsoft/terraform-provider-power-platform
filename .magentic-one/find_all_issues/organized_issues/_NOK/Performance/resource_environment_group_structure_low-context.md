# Title

Unnecessary Context Enrichment in All Resource Entry Points

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go

## Problem

Every function enriches the context using `helpers.EnterRequestContext(ctx, r.TypeInfo, req)` and defers an `exitContext()`, but this code does not clarify why such context manipulation is needed on every request, and the benefit is unclear from the current perspective. It may add hidden complexity or performance overhead if it does not add meaningful value.

## Impact

- Maintainers and reviewers may not easily grasp the need and purpose of context enrichment everywhere (low/medium).
- Possible performance impact if the log/context operation is heavy.
- Reduced code clarity: less obvious what is actually required for most CRUD operations.

**Severity:** low

## Location

Most resource methods (`Create`, `Update`, `Read`, `Delete`, `Schema`, `Configure`, `ImportState`)

## Code Issue

For example:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
```

## Fix

Document or centralize the reason for this code pattern's use and assess whether it can be moved into a decorator/wrapper. Otherwise, add comments describing why it should be present on all entry points. If `EnterRequestContext` is essential, ensure it's clear to new contributors.
