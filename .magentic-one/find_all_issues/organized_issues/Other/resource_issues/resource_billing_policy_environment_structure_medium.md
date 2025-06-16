# Title

Repeated context management and unclear responsibility for resource context handling

##

internal/services/licensing/resource_billing_policy_environment.go

## Problem

In every resource method (`Create`, `Read`, `Update`, `Delete`, etc.), there is a repeated pattern for entering and exiting request context:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
```

This leads to significant repetition and noise in each resource method. Furthermore, the `EnterRequestContext` return value and its deferred usage is not clearly documented or abstracted for the reader–it’s unclear if this is mandatory in every method, or what the nuanced responsibility of entering/leaving context is for the lifecycle hooks.

Additionally, there is no helper function or struct that encapsulates the pattern of creating, updating, and returning diagnostics from state and plan in CRUD operations. This tends to couple detailed lifecycle logic to each resource's body and impairs readability and maintainability.

## Impact

Severity: medium

Code duplication increases chances for inconsistency and errors, makes it harder to update context management, and reduces maintainability. Lack of abstraction impedes readability for new contributors, leading to potential misuses or bugs if context is not managed in a coordinated way.

## Location

Any resource method entry, e.g.:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
```

## Code Issue

```go
func (r *BillingPolicyEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	// ...
}
```
(and all similar resource method entries)

## Fix

Encapsulate context management in an embedded struct or reusable method, or create decorators or helper templates for standard resource methods. Add brief documentation for when and why to use `EnterRequestContext`. If the pattern is always repeated, consider refactoring so that it is handled automatically by the framework or as a wrapper, making resource bodies more focused on business logic:

```go
// Example pseudo-fix
func withResourceContext(
	ctx context.Context,
	typeInfo TypeInfo,
	req any,
	fn func(ctx context.Context),
) {
	ctx, exit := helpers.EnterRequestContext(ctx, typeInfo, req)
	defer exit()
	fn(ctx)
}

// Usage:
withResourceContext(ctx, r.TypeInfo, req, func(ctx context.Context) {
    // main logic here
})
```

Or encapsulate context entry/exit in a custom resource base struct.
