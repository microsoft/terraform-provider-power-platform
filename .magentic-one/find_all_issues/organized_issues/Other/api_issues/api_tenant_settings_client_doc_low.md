# Lack of Documentation and Exported Types for Client Struct

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/api_tenant_settings.go

## Problem

The `client` struct and its methods are not documented and are unexported (lowercase `client`). While this is not a critical issue for internal packages, adding documentation (GoDoc comments) and considering eventual export (`Client`) prepares the code for reuse and improves clarity. Interfaces that this struct satisfies could also be declared to ease mocking and future extensibility.

## Impact

- **Severity: Low**
- The code is less approachable for new contributors.
- Limits future extensibility, testing/mocking, and clarity.
- Reduces the benefits of GoDoc and static analysis tooling.

## Location

```go
type client struct {
    Api *api.Client
}

func (client *client) GetTenant(ctx context.Context) (*tenantDto, error) {
    ...
}
```

## Code Issue

```go
type client struct {
    Api *api.Client
}
```

## Fix

Add GoDoc comments for the type and its methods. Consider exporting the struct if future use outside the current package is intended and define the corresponding interface.

```go
// Client provides access to tenant and tenant settings endpoints for the Power Platform API.
type Client struct {
    Api *api.Client
}

// GetTenant returns the tenant details for the current API context.
func (c *Client) GetTenant(ctx context.Context) (*tenantDto, error) {
    ...
}
```

If mocking is intended, add an interface like:

```go
// TenantSettingsAPI defines methods for interacting with tenant settings.
type TenantSettingsAPI interface {
    GetTenant(ctx context.Context) (*tenantDto, error)
    GetTenantSettings(ctx context.Context) (*tenantSettingsDto, error)
    UpdateTenantSettings(ctx context.Context, tenantSettings tenantSettingsDto) (*tenantSettingsDto, error)
}
```
This improves maintainability and clarity.
