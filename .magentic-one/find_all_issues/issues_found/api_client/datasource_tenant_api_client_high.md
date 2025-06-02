# No nil check for TenantClient in Read method

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

In the Read method, d.TenantClient is used without first checking whether it is nil. If Configure was never called, or if there was an error in configuration/setup, this could lead to a panic during operation.

## Impact

Severity: High. If d.TenantClient is nil, this will cause a runtime panic, causing the provider to crash and breaking the entire Terraform operation.

## Location

In Read method:

## Code Issue

```go
tenant, err := d.TenantClient.GetTenant(ctx)
```

## Fix

Add a nil check prior to use, and fail gracefully:

```go
if d.TenantClient == nil {
    resp.Diagnostics.AddError("Tenant client is not configured", "The TenantClient is nil. This is likely a bug in the provider initialization.")
    return
}

tenant, err := d.TenantClient.GetTenant(ctx)
```
