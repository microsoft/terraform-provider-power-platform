# Title
Use of pointer fields for primitive types (Count, Value) in OIDC token response struct

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
In the anonymous struct for OIDC response unmarshaling, fields for `Count *int` and `Value *string` are defined as pointers. While this is sometimes useful for optional fields, it increases complexity as the caller must always check for `nil`. If JSON response always returns these fields or `0`/`""` is a valid value, a non-pointer can be more robust and reduces code noise.

## Impact
Medium. Unnecessary pointer indirection can lead to panics (accidentally dereferencing nil), makes code more verbose, and is less idiomatic Go unless explicit nullability is required by the API contract.

## Location
In the `getAssertion` method:

## Code Issue
```go
var tokenRes struct {
    Count *int    `json:"count"`
    Value *string `json:"value"`
}
```

## Fix
If the contract for the endpoint is that `count` and `value` are always present, prefer non-pointer types:

```go
type oidcTokenResponse struct {
    Count int    `json:"count"`
    Value string `json:"value"`
}
```

If true nullability is required, keep pointers but document carefully and always check for nil before dereferencing. Prefer non-pointers for primitives unless required.