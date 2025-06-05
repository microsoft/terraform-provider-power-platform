# Title
Missing retry or backoff logic on HTTP request for OIDC assertion

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The `getAssertion` method makes an HTTP request to acquire a token assertion but does not attempt any retry or exponential backoff for robustness. Any temporary network hiccup, throttling, or transient error results in immediate failure. Since this is part of authentication logic, reliability is critical.

## Impact
Medium severity, as this could cause authentication failures due to transient network conditions, cloud API throttling, or hiccups in the OIDC endpoint, reducing reliability of the provider in real-world scenarios.

## Location
In `getAssertion` method:

## Code Issue
```go
resp, err := client.Do(req)
if err != nil {
    return "", fmt.Errorf("getAssertion: cannot request token: %w", err)
}
```

## Fix
Implement retry logic with exponential backoff for network-related errors and server-side temporary errors (such as 429, 500, 502, 503, 504). You can use a helper or a for-loop, e.g.:

```go
for i := 0; i < 3; i++ {
    resp, err := client.Do(req)
    if transientErr(err, resp) {
        time.Sleep(backoff(i))
        continue
    }
    // ... handle normal case ...
    break
}
```

Or integrate a robust HTTP client library with built-in retry mechanisms.