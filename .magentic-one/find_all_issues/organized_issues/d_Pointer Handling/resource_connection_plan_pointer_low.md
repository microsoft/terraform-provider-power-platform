# Title
Unnecessary use of double-pointer for ResourceModel variable `plan`

##
/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go

## Problem
The code uses double-pointer syntax for `plan` (i.e., `var plan *ResourceModel`) and then passes its address around (e.g., `&plan`), but Go's Terraform Plugin Framework typically expects a value or a single pointer. This can lead to confusion and inconsistent state handling.

## Impact
Low: This is mostly a readability/maintainability issue. It may introduce subtle bugs if pointer semantics are misunderstood, but does not currently break the code.

## Location
Create, Update methods:

## Code Issue
```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// ... then plan usage
```

## Fix
Use a non-pointer or a single pointer appropriately for `plan` with Terraform plugin calls.

```go
var plan ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // pass address of value, not address of pointer
// ... use plan fields as plan.Foo
```
If you need `plan` to be a pointer, then don't take `&plan` (it's already a pointer):

```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, plan)...) // just pass plan
```

Choose the style based on framework expectations and keep it consistent throughout the file.

