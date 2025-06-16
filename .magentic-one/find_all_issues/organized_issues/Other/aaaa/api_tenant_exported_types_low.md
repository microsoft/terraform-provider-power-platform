# Potential Overexposure with Exported Struct and Method

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/api_tenant.go

## Problem

Both the `Client` struct and its method `GetTenant` are exported (capitalized names), although it's not clear from the context that they are intended to be part of the public API. In Go, only types and functions meant to be used outside the package should be exported; otherwise, they should be unexported (start with lowercase letters).

## Impact

Unnecessarily exporting types and functions increases the public API surface area, hurting encapsulation and possibly leading to misuse of internal components. **Severity:** Low

## Location

```go
type Client struct {
	Api *api.Client
}

func (client *Client) GetTenant(ctx context.Context) (*TenantDto, error) {
```

## Code Issue

```go
type Client struct {
	Api *api.Client
}

func (client *Client) GetTenant(ctx context.Context) (*TenantDto, error) {
```

## Fix

If these are not intended for use outside `tenant` package, make them unexported:

```go
type client struct {
	API *api.Client
}

func (c *client) getTenant(ctx context.Context) (*TenantDto, error) {
```

Be sure to update all references accordingly.
