# Title
Naming: Struct and Function Parameter Names Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/api_environment_settings.go

## Problem

The `client` struct and its receiver variable are both named `client`. By convention, struct names in Go should be capitalized (`Client`) when exported, and receiver variable names should preferably be small (e.g., `c`). The current naming can be confusing and does not follow Go idioms.

## Impact

Low/Medium. This impacts code readability and maintainability, especially for people reading or maintaining the code in the future.

## Location

```go
type client struct {
    Api *api.Client
}

func (client *client) ...
```

## Code Issue

```go
type client struct {
    Api *api.Client
}

func (client *client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
    //...
}
```

## Fix

Rename the struct to `Client` (capitalize if you want to export), and use a shorter receiver name, typically `c`.

```go
type Client struct {
    Api *api.Client
}

func (c *Client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
    //...
}
```
