# Potential Code Duplication of Context Management Pattern

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go

## Problem

Every CRUD method, as well as Metadata, Schema, Configure, and ImportState, start with an identical two-line pattern for context management:

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
```

This leads to a lot of boilerplate duplication that makes the code harder to maintain and more verbose. Refactoring this repeated logic into a decorator (higher-order function) or a method wrapper could reduce duplication and centralize context-management logic.

## Impact

Severity: **medium**

The impact is mainly on code maintainability, clarity, and future extensibility. Changes to context management would need to be repeated in every method, increasing maintenance overhead and risk of inconsistency.

## Location

- All public methods of `AdminManagementApplicationResource` (Metadata, Schema, Configure, ImportState, Create, Read, Update, Delete)

## Code Issue

The repeated pattern is present in every relevant function header as shown:

```go
func (r *AdminManagementApplicationResource) Create(...) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()
    ...
}
```

## Fix

Consider creating a wrapper/helper function that takes the handler as a parameter, performs the context management, and then calls the handler. Example skeleton:

```go
func withRequestContext(r *AdminManagementApplicationResource, req any, handler func(ctx context.Context)) {
    ctx, exitContext := helpers.EnterRequestContext(context.Background(), r.TypeInfo, req)
    defer exitContext()
    handler(ctx)
}
```

Then use:

```go
func (r *AdminManagementApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    withRequestContext(r, req, func(ctx context.Context) {
        // original function body
    })
}
```

Alternatively, macro generation or code linting could ensure uniformity as well.
