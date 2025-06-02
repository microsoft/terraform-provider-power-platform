# Title

Potential Panic Due to Unchecked resp.State.Get Error

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_tenant_application_packages.go

## Problem

In the Read method, `resp.State.Get(ctx, &state)` is called, but the error returned (diagnostics) is not checked. In unusual conditions (malformed state, framework bug, etc.), this could lead to panics or undefined data usage.

## Impact

- High: Possible panic or silent failure with bad input or Terraform state bugs.
- Error handling best practice.

## Location

- Read method, at start of method.

## Code Issue

```go
var state TenantApplicationPackagesListDataSourceModel
resp.State.Get(ctx, &state)
```

## Fix

**Capture diagnostics and handle errors:** 

```go
var state TenantApplicationPackagesListDataSourceModel
diags := resp.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	return
}
```
