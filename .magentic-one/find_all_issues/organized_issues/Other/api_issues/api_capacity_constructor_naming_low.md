# Issue: Construction Function Name Not Idiomatic

## 
/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go

## Problem

The factory/construction function for the `client` type is named `newCapacityClient`. In Go, such constructors are usually named `new<Type>` (e.g. `newClient`), unless there is a need to distinguish among multiple client types. If the context is unique enough in its package, `newClient` would be more idiomatic and concise.

## Impact

Severity: **low**

This issue is about code idioms and package cohesion. Using a redundant or overly verbose constructor name may reduce readability and break Go conventions.

## Location

```go
func newCapacityClient(clientApi *api.Client) client {
    return client{
        Api: clientApi,
    }
}
```

## Code Issue

```go
func newCapacityClient(clientApi *api.Client) client {
    return client{
        Api: clientApi,
    }
}
```

## Fix

If this is the principal client for this package, rename to:

```go
func newClient(apiClient *api.Client) client {
    return client{
        Api: apiClient,
    }
}
```
And apply previously recommended struct renamings (`Api` → `api`).
