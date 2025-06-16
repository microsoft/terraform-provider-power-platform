# Redundant context management with EnterRequestContext and defer

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

In every exported method (`Metadata`, `Schema`, `Read`, `Configure`), there's a call to `helpers.EnterRequestContext(ctx, d.TypeInfo, req)`, with the returned function deferred as `exitContext()`. While this is likely part of a logging or resource management pattern, if the EnterRequestContext or deferred exitContext manage resources (e.g. spans, log scopes, etc.), they could stack up or be used incorrectly if future code is added or exceptions occur. It is always important to ensure their use is really needed in every function, and that any panics or early returns still guarantee exitContext() executes.

## Impact

Low: Currently, usage seems safe (because of defer), but this can lead to maintenance risk if copy-pasted elsewhere without care or if EnterRequestContext's semantics change.

## Location

At the top of each main DataSource method:

## Code Issue

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
defer exitContext()
```

## Fix

Document the reason for using EnterRequestContext/exitContext in a code comment at each location, and ensure their use is truly warranted. If not mandatory, refactor or centralize resource/context management.

```go
// EnterRequestContext tracks/logs operation context for debugging
ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
defer exitContext()
```
