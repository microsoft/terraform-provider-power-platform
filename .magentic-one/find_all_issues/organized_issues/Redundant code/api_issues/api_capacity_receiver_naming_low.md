# Issue: Function Receiver Naming is Too Generic

## 
/workspaces/terraform-provider-power-platform/internal/services/capacity/api_capacity.go

## Problem

Throughout the file, the receiver for the `client` struct methods is named `client`, which is both verbose and redundant. Go convention is to use a one- or two-letter receiver name (commonly derived from the struct type). Long or generic receiver names can reduce code readability, especially when distinguishing between the receiver and local variables or types that have similar names.

## Impact

Severity: **low**

This is primarily a readability/naming issue but can make code harder to read and follow, especially as the struct grows or if there are naming clashes with local variables.

## Location

All method receivers for `client`:

```go
func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    // ...
}
```

## Code Issue

```go
func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    // ...
}
```

## Fix

Use a concise, idiomatic receiver name (such as `c`):

```go
func (c *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
    // replace all occurrences of 'client.' in the method body with 'c.'
}
```
