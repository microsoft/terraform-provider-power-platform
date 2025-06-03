# Inconsistent Resource Naming for New Client Constructor

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go

## Problem

The function `newConnectorsClient` returns a value of type `client`, but the type it returns is simply named `client`, which is unexported, very generic, and could clash or confuse with idiomatic usage or other `client` types in the codebase.

## Impact

Poor naming for exported entities (even unexported in this file) can decrease codebase maintainability and readability, leading to confusion during development or future refactoring. Severity: **low**.

## Location

At the top of the file:

```go
func newConnectorsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}
```

## Code Issue

```go
func newConnectorsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}
```

## Fix

Rename both the constructor and struct to `connectorsClient` or a similar descriptive name:

```go
type connectorsClient struct {
	Api *api.Client
}

func newConnectorsClient(apiClient *api.Client) connectorsClient {
	return connectorsClient{
		Api: apiClient,
	}
}
```

This makes the code more descriptive and less prone to naming conflicts.
