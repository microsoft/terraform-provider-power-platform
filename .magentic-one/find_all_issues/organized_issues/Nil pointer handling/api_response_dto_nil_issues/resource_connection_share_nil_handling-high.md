# Title

Missing error return after detecting nil `share` in Create, Read, and Update methods

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share.go

## Problem

In the `Create`, `Read`, and `Update` methods, when a nil `share` or `newShare` is detected, an error is logged to `resp.Diagnostics` but the control flow continues to run and may attempt to access properties on a nil pointer. For example, after checking `if share == nil`, it should immediately return rather than continuing on to `convertFromConnectionResourceSharesDto(plan, share)` which will panic if `share` is `nil`.

## Impact

This introduces a risk of panics due to nil pointer dereferences if the `share` or `newShare` objects are missing. The severity is **high** since it will cause the provider to crash during runtime.

## Location

- `Create` method: after `if share == nil { ... }`
- `Read` method: after `if share == nil { ... }`
- `Update` method: after `if newShare == nil { ... }`

## Code Issue

```go
	if share == nil {
		resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
	}
	// code continues and uses "share" without returning

	// ... also in Read and Update, same problem with "share" and "newShare"
```

## Fix

Add a `return` statement immediately after adding the error to diagnostics for these nil checks:

```go
	if share == nil {
		resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
		return
	}
```

Apply the above adjustment in all affected places (Create, Read, Update):

```go
	// Example in Create:
	if share == nil {
		resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
		return
	}

	// Example in Update:
	if newShare == nil {
		resp.Diagnostics.AddError("Error getting connection share", "Connection share not found")
		return
	}
```
