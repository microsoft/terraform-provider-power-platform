# Title

Error handling for plan/state Get does not differentiate between error types

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share.go

## Problem

The calls to `req.Plan.Get` and `req.State.Get` in resource lifecycle methods append diagnostics but do not differentiate between errors caused by user input, type incompatibility, or partial failures. There's no separate logging or error return except for a generic check with `HasError()`. Additionally, nil pointer dereferences could occur if a `nil` value slips into the plan or state, since they are immediately dereferenced later.

## Impact

Severity is **medium**. While most standard errors will be caught by the diagnostics check, more granular error handling or protective checks for nil plan/state would make error handling more robust.

## Location

All resource methods (`Create`, `Read`, `Update`, `Delete`) that get plan or state:

```go
	var plan *ShareResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
```

## Code Issue

```go
	var plan *ShareResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
```

## Fix

Check for a nil value on state or plan before using it, in addition to diagnostics:

```go
	var plan *ShareResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || plan == nil {
		return
	}
```

This avoids panics or errors if nil sneaks past diagnostics (e.g., in testing or unexpected framework edge cases). Repeat for `state` wherever it is used.
