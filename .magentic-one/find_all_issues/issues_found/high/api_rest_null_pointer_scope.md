# Title

Potential Null Pointer Dereference in `SendOperation`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

In `SendOperation`, the `operation.Scope.ValueStringPointer()` is passed directly to the `ExecuteApiRequest` function without verifying whether it is nil. If `operation.Scope` is missing or not set, the code may trigger an error at runtime, leading to unexpected behavior or crashes.

## Impact

- **High application fragility**: An unverified null pointer can cause unpredictable runtime errors.
- **Security risks**: Exploitable null pointer dereferences could potentially lead to denial of service.
- **Defective functionality**: The `ExecuteApiRequest` function relies on the scope pointer. A missing value may break the logic if not handled properly.

Severity: **High**

## Location

Found in the `SendOperation` function.

## Code Issue

```go
res, err := client.ExecuteApiRequest(ctx, operation.Scope.ValueStringPointer(), url, method, body, headers, expectedStatusCodes)
```

## Fix

Add proper nil checks before passing the `operation.Scope.ValueStringPointer()` to `ExecuteApiRequest`. This ensures the code behaves correctly even in edge cases.

```go
var scopePointer *string
if operation.Scope.ValueStringPointer() != nil {
    scopePointer = operation.Scope.ValueStringPointer()
} else {
    return types.Object{}, errors.New("operation scope is required but missing")
}

res, err := client.ExecuteApiRequest(ctx, scopePointer, url, method, body, headers, expectedStatusCodes)
```

This fix guarantees that the scope pointer is validated and prevents null dereference errors. It also returns an appropriate error if `operation.Scope` is missing.