# Direct diagnostic propagation after resp.State.Set

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

In the `Read` method, diagnostics from `resp.State.Set` are appended to `resp.Diagnostics`. If any errors are present, an early return is executed. While this is common in Terraform plugins, it is a control flow touchpoint that can benefit from more explicit error handling and possible state clean-up or logging to better support debugging and maintainability.

## Impact

**Severity: medium**

If setting the state fails, only a non-specific error will be returned. There is room for improvement by adding contextual information or logging, which would facilitate debugging complex state issues in production.

## Location

```go
diags := resp.State.Set(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	return
}
```

## Code Issue

```go
diags := resp.State.Set(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	return
}
```

## Fix

Consider logging or enriching the diagnostic message before returning; at minimum, add a comment describing control flow intent:

```go
diags := resp.State.Set(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	// State could not be set—returning to prevent invalid state.
	// Optionally, add contextual logging here if needed for debugging.
	return
}
```
