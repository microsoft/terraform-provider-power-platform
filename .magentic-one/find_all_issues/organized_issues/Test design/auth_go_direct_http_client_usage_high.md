# Title
Direct use of http.DefaultClient without timeouts or customization can cause issues

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The code directly invokes `http.DefaultClient.Do(req)` in the `getAssertion` method of `OidcCredential`. This default client has no timeouts and is shared globally, which can cause resource starvation or unbounded wait in case of issues (e.g., network stalls). In high-availability or cloud tools, always use a custom http.Client with strict timeout policies, and avoid global state for improved testability and reliability.

## Impact
High severity for API clients. This can cause hangs in production, and may inadvertently be affected by changes to the global client elsewhere in the process. It also reduces testability.

## Location
In the `getAssertion` method:

## Code Issue
```go
resp, err := http.DefaultClient.Do(req)
```

## Fix
Declare and use a custom http.Client with a sensible timeout, e.g.:

```go
client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)
```

In production code, consider making the timeout and client customizable (injected via DI or config), or using a shared, non-default client for all requests, set up at package level.