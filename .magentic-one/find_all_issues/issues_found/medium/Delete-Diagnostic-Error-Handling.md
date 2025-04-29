# Title

Improper Handling of Diagnostic Errors in `Delete` Function

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest.go

## Problem

The `Delete` function does not check for existing diagnostic errors (`resp.Diagnostics.HasError`) after appending diagnostics from `req.State.Get(ctx, &state)`. Failing to validate diagnostics might lead to continuing execution when critical issues exist.

## Impact

This issue can cause undefined behavior, such as attempting to dereference a corrupted or invalid `state` object. Severity: **Medium**.

## Location

Function: `Delete`

## Code Issue

```go
	var state *DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	// Missing validation for diagnostics
	if state.Destroy != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Destroy)
		if err != nil {
			resp.Diagnostics.AddError("Error executing destroy operation", err.Error())
			return
		}
		state.Output = bodyWrapped
	}
```

## Fix

Add a check to validate diagnostics before proceeding with any `state` operations. For example:

```go
	var state *DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	
	// Check diagnostics for errors before continuing
	if resp.Diagnostics.HasError() {
		return
	}
	
	if state.Destroy != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Destroy)
		if err != nil {
			resp.Diagnostics.AddError("Error executing destroy operation", err.Error())
			return
		}
		state.Output = bodyWrapped
	}
```

This ensures that the function stops execution in case of critical diagnostic errors, reducing the risk of undefined behavior.