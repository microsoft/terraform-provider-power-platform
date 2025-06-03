# Function Naming Doesn't Follow Go Conventions

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/api_tenant.go

## Problem

The function `NewTenantClient` does not match the usual Go conventions for constructor naming, which would be `NewClient` for the struct called `Client`. If the package is already called `tenant`, it is common to use `NewClient` so usage from outside would read `tenant.NewClient()`, which is more idiomatic.

## Impact

This can make the codebase harder to use and less idiomatic for Go developers, leading to confusion. **Severity:** Low

## Location

Line 12

## Code Issue

```go
func NewTenantClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}
```

## Fix

Rename the function to `NewClient`:

```go
func NewClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}
```
