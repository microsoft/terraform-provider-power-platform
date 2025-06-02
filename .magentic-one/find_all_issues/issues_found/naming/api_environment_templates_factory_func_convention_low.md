# Code Structure: Factory Function Naming Should Match Type Name

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/api_environment_templates.go

## Problem

The factory function `newEnvironmentTemplatesClient` creates a value of type `client`, but the function name and returned type do not match Go conventions (type should usually be `Client`, and function should be `NewEnvironmentTemplatesClient` if exported). This makes usages inconsistent and less idiomatic.

## Impact

Non-standard naming makes it harder for other Go developers to follow the code and for auto-completion/reflection tools. Severity: **Low**.

## Location

```go
func newEnvironmentTemplatesClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Code Issue

```go
func newEnvironmentTemplatesClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Fix

Rename the function and type for clarity and consistency:

```go
func NewEnvironmentTemplatesClient(apiClient *api.Client) *Client {
	return &Client{
		Api: apiClient,
	}
}
```
