# Title

Repeated Initialization Pattern without Consolidation

##

`/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go`

## Problem

The pattern `helpers.EnterRequestContext(ctx, d.TypeInfo, req)` followed by `defer exitContext()` is used in multiple methods (`Metadata`, `Schema`, `Configure`, `Read`) with identical logic. This repetition reduces maintainability and can introduce redundant errors when updates to the pattern are necessary.

## Impact

Repetition increases code duplication and makes maintenance harder. Any update to the logic will require changes in multiple places, increasing the likelihood of introducing regressions. **Severity is low**, as the current functionality is not negatively affected, but future scalability and maintainability are compromised.

## Location

Across the methods `Metadata`, `Schema`, `Configure`, and `Read` in the file.

## Code Issue

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
defer exitContext()
```

This code repeats in different methods.

## Fix

Refactor the repeated logic into a utility function that consolidates context initialization and error validation. For instance:

```go
func initializeRequestContext(ctx context.Context, typeInfo helpers.TypeInfo, req interface{}) (context.Context, func(), error) {
    newCtx, exitContext := helpers.EnterRequestContext(ctx, typeInfo, req)
    if newCtx == nil {
        return nil, nil, fmt.Errorf("Failed to initialize request context")
    }
    return newCtx, exitContext, nil
}
```

Then update the original methods to use this utility function:

```go
ctx, exitContext, err := initializeRequestContext(ctx, d.TypeInfo, req)
if err != nil {
    resp.Diagnostics.AddError("Context Initialization Error", err.Error())
    return
}
defer exitContext()
```

This approach centralizes the logic, reduces redundancy, and ensures consistent error validation across methods.
