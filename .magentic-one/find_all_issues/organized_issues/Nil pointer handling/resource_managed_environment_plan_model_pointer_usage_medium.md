# Title

Use of pointer-to-struct for resource plan/state model may risk nil dereference

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

The code repeatedly uses a pointer (`var plan *ManagedEnvironmentResourceModel`) for the resource plan/state model throughout Create, Update, and Read. In most cases, this is populated by `.Get()` methods, but any future code changes, test scaffolding, or framework changes could result in nil being assigned or improperly mocked, yielding nil pointer panics on field access. Best practice is to prefer value structs unless mutation of the struct pointer itself is needed (rare in Terraform provider logic) or nil is a valid state. Type safety is improved by avoiding unnecessary indirection.

## Impact

Medium. Could cause panics or hidden bugs, especially in future test or code refactoring where plan is accidentally left nil after .Get or during mock initialization.

## Location

Example in Create, Update, Read:

## Code Issue

```go
var plan *ManagedEnvironmentResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// use plan immediately: plan.EnvironmentId.ValueString()
```

## Fix

Declare the struct as a value, not a pointer:

```go
var plan ManagedEnvironmentResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// use plan.EnvironmentId.ValueString(), etc.
```

This avoids nil dereferences while using zero-values for safety throughout logic and error handling.
