# Multiple assignments of plan variable from req.Plan.Get

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

In the `Create` function, `req.Plan.Get(ctx, &plan)` is called and diagnostics appended twice, one shortly after the other. This is redundant and unnecessary. The planned state should only be extracted once in a single code path unless there's a compelling reason to refresh `plan` between two stages (which does not seem to be the case here).

## Impact

This reduces code clarity and can lead to confusion about which instance of `plan` is in use. Severity: low.

## Location

Line ~170 ("resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)..." is called twice.)

## Code Issue

```go
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}

// ... (unrelated code skipped for clarity)

// Get the plan
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
```

## Fix

Remove the redundant second call to `req.Plan.Get(ctx, &plan)`, keeping only the first instance.

```go
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}

// ... continue with rest of function, do not repeat req.Plan.Get(ctx, &plan)
```
