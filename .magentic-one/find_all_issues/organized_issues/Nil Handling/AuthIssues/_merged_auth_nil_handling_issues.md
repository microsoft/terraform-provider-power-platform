# Authentication Nil Handling Issues

This document contains all identified nil handling issues related to authentication components in the Terraform Power Platform provider codebase.

## ISSUE 1

<!-- Source: auth_go_defer_close_panic_high.md -->

# Title

defer resp.Body.Close() not nil-safe, risk of panic

##

/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem

The code in `getAssertion` uses:

```go
defer resp.Body.Close()
```

immediately after a nil-check on error, but does not check if `resp` or `resp.Body` are `nil`. If the HTTP client or net/http returns a non-nil error but constructs a partially-initialized response (which can happen in rare network cases or due to custom clients), then deferring `resp.Body.Close()` can cause a panic if `resp` or `resp.Body` are `nil`.

## Impact

High severity. This could cause a runtime panic and crash the program or provider, especially under certain HTTP/network error conditions, instead of returning an error as expected.

## Location

In `getAssertion` method:

## Code Issue

```go
resp, err := http.DefaultClient.Do(req)
if err != nil {
    return "", fmt.Errorf("getAssertion: cannot request token: %w", err)
}
defer resp.Body.Close()
```

## Fix

Guard the defer statement so it is only called if `resp != nil && resp.Body != nil`, for example:

```go
resp, err := client.Do(req)
if err != nil {
    return "", fmt.Errorf("getAssertion: cannot request token: %w", err)
}
if resp != nil && resp.Body != nil {
    defer resp.Body.Close()
}
```

This ensures no panic if the underlying library misbehaves or returns partial data.

## ISSUE 2

<!-- Source: auth_go_pointer_fields_in_structs_medium.md -->

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

---

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

Apply this fix to the whole codebase
