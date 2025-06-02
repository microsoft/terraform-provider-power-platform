# Title

Exit context deferred but not always executed due to early returns

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

Throughout the CRUD (`Create`, `Read`, `Update`, `Delete`) and other methods, `exitContext()` is deferred immediately after entering the helper context. However, if the function returns before the defer is registered (e.g., due to early returns for error handling), the deferred cleanup will not be executed, potentially leaking resources or creating incorrect telemetry/trace lifecycles.

## Impact

This can cause subtle bugs, such as memory/resource leaks, incorrect trace/tool context information, or unflushed logs. Severity: **Medium** (may have cumulative resource impact and create hard-to-debug issues in a long running provider/server process).

## Location

Multiple locations, such as:

```go
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext() // <-- If there's any return before this line, exitContext() will not execute
	// ...
}
```

## Code Issue

```go
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
    //...
}
```

## Fix

Ensure `defer exitContext()` is registered as the very first line in the function after variable/assignment declarations. Don't interleave any early returns before it. For patterns where `helpers.EnterRequestContext` must happen after nil-checks (rare) or conditionally, reorganize the method logic.

```go
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext() // now guaranteed to always run if function entered context!

	var plan *UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
    //...
}
```

If it's not possible due to logic (e.g., sometimes context helper is not called), wrap the early checks outside the helper invocation:

```go
func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()
    // ...
}
```
