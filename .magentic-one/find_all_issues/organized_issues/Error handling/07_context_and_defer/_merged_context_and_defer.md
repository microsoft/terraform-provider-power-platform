# Context and Defer Issues

This document consolidates all issues related to context handling and defer statement usage in the Terraform Provider for Power Platform.

## ISSUE 1

### Potential Leaked Context Resources with Defer

**File:** `/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go`

**Problem:** Each resource method wraps code with a call to `helpers.EnterRequestContext`, which returns a cleanup function that is deferred. However, there are multiple `return` statements before the logic ends (e.g., after error handling or diagnosis checks). In rare cases, resource cleanup may occur later than optimal (if the context carries large values or goroutines), as the function returns earlier. This can be a problem for resource leaks or lock contention in long-lived operations.

**Impact:**

- Severity: **Low**
- For most cases in Terraform, this may not be a practical issue (since the function scope is returned shortly afterwards), but for context-bound resources, prompt cleanup is a best practice and makes reasoning about program state easier.

**Location:** In every method that calls:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
```

**Code Issue:**

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
if errorCondition {
    return
}
```

**Fix:** Consider calling `exitContext()` directly before any `return` that happens before the end of the function:

```go
if errorCondition {
    exitContext()
    return
}
```

Alternatively, review if `exitContext()` must always run at function return, or if the current usage is non-problematic. If so, document expected lifetimes in comments.

## ISSUE 2

### Defer exitContext() May Extend Resource Lifetimes

**File:** `/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go`

**Problem:** This code uses `defer exitContext()` at the top of most resource methods after entering a request context. However, since these handler methods are not always long-running (may return early on error), the defer will always execute at function return. In most cases this is fine, but in Go, if any context resources are opened/pinned, they could stay open longer than necessary (e.g., until panic or stack unwinding). This is not a memory leak, but can sometimes defer cleanup (especially in high-throughput servers or in benchmarking scenarios), causing unnecessary resource occupation.

**Impact:** Medium. In the general case, this is idiomatic and safe, but in very tight event loops or in extremely rare error/panic cases, resource holding could be longer than needed. Also, it can obscure where cleanup truly happens as code evolves/forks inside the method.

**Code Issue:**

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
```

**Fix:** For more explicit resource handling, you could opt for:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
// Instead of `defer`, explicitly call exitContext() at each early return when critical resources are involved.
```

Or review/ensure that exitContext() is always idempotent and side-effect-free if accidentally called multiple times, and that defer usage is not hiding resource pinning in complex flows.

## ISSUE 3

### Exit Context Deferred but Not Always Executed Due to Early Returns

**File:** `/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go`

**Problem:** Throughout the CRUD (`Create`, `Read`, `Update`, `Delete`) and other methods, `exitContext()` is deferred immediately after entering the helper context. However, if the function returns before the defer is registered (e.g., due to early returns for error handling), the deferred cleanup will not be executed, potentially leaking resources or creating incorrect telemetry/trace lifecycles.

**Impact:** This can cause subtle bugs, such as memory/resource leaks, incorrect trace/tool context information, or unflushed logs. Severity: **Medium** (may have cumulative resource impact and create hard-to-debug issues in a long running provider/server process).

**Location:** Multiple locations, such as:

```go
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
        ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
        defer exitContext() // <-- If there's any return before this line, exitContext() will not execute
        // ...
}
```

**Code Issue:**

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

**Fix:** Ensure `defer exitContext()` is registered as the very first line in the function after variable/assignment declarations. Don't interleave any early returns before it:

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

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
