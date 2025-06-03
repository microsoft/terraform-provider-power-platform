# Inadequate Error Handling on State Get in Read

##

/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations.go

## Problem

In the `Read` method, the code retrieves state with:
```go
var state DataSourceModel
resp.State.Get(ctx, &state)
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
	return
}
```
This calls `resp.State.Get` twice: once without checking its error and again with appending diagnostics. If the first call populates `state` or fails, the second call could yield different results or diagnostics could be duplicated/misleading. It also violates the pattern of always handling and appending errors for state operations.

## Impact

This could lead to errors being missed, diagnostics being inconsistent or duplicated, and makes the control flow harder to understand and debug. **Severity: Medium**

## Location

Method: `Read`, lines around:

```go
var state DataSourceModel
resp.State.Get(ctx, &state)
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
	return
}
```

## Code Issue

```go
var state DataSourceModel
resp.State.Get(ctx, &state)
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
	return
}
```

## Fix

Only call `Get` once, append any diagnostics, and proceed based on error state:

```go
var state DataSourceModel
diags := resp.State.Get(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	return
}
```
