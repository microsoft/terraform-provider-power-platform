# Title

Inconsistent casing of initialisms in function names (e.g., `NewApiClientBase` should be `NewAPIClientBase`)

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The function `NewApiClientBase` uses `Api` instead of the Go convention `API`. This creates inconsistency and goes against Go best practices regarding common initialisms.

## Impact

Reduces codebase consistency and makes code less idiomatic. Severity: **low**

## Location

Constructor for API client base

## Code Issue

```go
func NewApiClientBase(providerConfig *config.ProviderConfig, baseAuth *Auth) *Client {
	return &Client{
		Config:   providerConfig,
		BaseAuth: baseAuth,
	}
}
```

## Fix

Rename to use conventional casing and update usages accordingly:

```go
func NewAPIClientBase(providerConfig *config.ProviderConfig, baseAuth *Auth) *Client {
	return &Client{
		Config:   providerConfig,
		BaseAuth: baseAuth,
	}
}
```
