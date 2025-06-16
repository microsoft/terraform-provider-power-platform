# Function Naming Convention Non-Idiomatic: `UnLinkEnterprisePolicy`

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

The function name `UnLinkEnterprisePolicy` uses PascalCase for `UnLink`, which is not idiomatic Go. The standard is to use `Unlink` (one word).

## Impact

This detracts from Go code idioms, potentially confusing to Go developers, and decreases consistency and maintainability (but is not functionally impactful). Severity: **low**.

## Location

Function declaration and recursive use:

```go
func (client *Client) UnLinkEnterprisePolicy(...)
```

## Code Issue

```go
func (client *Client) UnLinkEnterprisePolicy(ctx context.Context, environmentId, environmentType, systemId string) error {
	// ...
	return client.UnLinkEnterprisePolicy(ctx, environmentId, environmentType, systemId)
}
```

## Fix

Rename everywhere from `UnLinkEnterprisePolicy` to `UnlinkEnterprisePolicy` for compliance with idiomatic naming:

```go
func (client *Client) UnlinkEnterprisePolicy(ctx context.Context, environmentId, environmentType, systemId string) error {
	// ...
	return client.UnlinkEnterprisePolicy(ctx, environmentId, environmentType, systemId)
}
```
