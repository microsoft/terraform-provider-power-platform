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
