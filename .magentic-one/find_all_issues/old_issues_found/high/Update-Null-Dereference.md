# Title

Potential Null Dereference in `Update` Function

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest.go

## Problem

The `plan` reference in the `Update` function is declared as a pointer (`*DataverseWebApiResourceModel`) without a prior null-check before dereferencing. If the pointer is null, it could result in a runtime panic.

## Impact

This can cause the program to crash unexpectedly during update operations, leading to system downtime or failures. Severity: **High**.

## Location

Function: `Update`

## Code Issue

```go
	var plan *DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // Pointer might not be initialized properly
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Update != nil { // Dereferencing without a null-check
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, plan.Update)
		if err != nil {
			resp.Diagnostics.AddError("Error executing update operation", err.Error())
			return
		}
		plan.Output = bodyWrapped
		// Additional operations...
	}
```

## Fix

Add a null-check for the pointer before dereferencing it:

```go
	var plan *DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // Initialize the pointer
	if plan == nil { // Null-check added to prevent dereferencing a null pointer
		resp.Diagnostics.AddError("Invalid Plan", "Plan received is null.")
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Update != nil { // Proceed after safety checks
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, plan.Update)
		if err != nil {
			resp.Diagnostics.AddError("Error executing update operation", err.Error())
			return
		}
		plan.Output = bodyWrapped
		// Additional operations...
	}
```

This ensures safe usage of the `plan` pointer and avoids runtime panics.