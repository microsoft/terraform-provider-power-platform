# Title

API client is tightly coupled to authentication logic

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The `Client` struct directly embeds `BaseAuth *Auth`, leading to tight coupling between the API client and authentication logic. This makes it harder to swap out authentication methods, introduces testability challenges (no interfaces/mocking), and weakens single responsibility.

## Impact

Reduced maintainability, extensibility, and testability. Severity: **medium**

## Location

Definition of the `Client` struct:

## Code Issue

```go
type Client struct {
	Config   *config.ProviderConfig
	BaseAuth *Auth
}
```

## Fix

Refactor to depend on an interface for authentication, e.g.:

```go
type Authenticator interface {
	GetTokenForScopes(ctx context.Context, scopes []string) (string, error)
}

type Client struct {
	Config   *config.ProviderConfig
	Auth     Authenticator
}
```

Update usage throughout to refer to `client.Auth.GetTokenForScopes` accordingly.
