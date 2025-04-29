# Title

Improper Error Handling in `Create` Operation

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest.go

## Problem

In the `Create` function, the plan state (`resp.State.Get(ctx, &plan)`) is being fetched before checking for errors in the diagnostics returned from the `req.Plan.Get(ctx, &plan)` operation. This sequence can lead to undefined behavior if an error exists in `req.Plan.Get`.

## Impact

Error diagnostics may not be handled correctly, causing potential resource misconfigurations and runtime issues. Severity: **High**.

## Location

Function: `Create`

## Code Issue

```go
	var plan DataverseWebApiResourceModel
	resp.State.Get(ctx, &plan) // Plan state being fetched first
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // Error diagnostics checked later
	if resp.Diagnostics.HasError() {
		return
	}
```

## Fix

Change the order of error handling to ensure diagnostics are checked before fetching the plan state. For example:

```go
	var plan DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // Check error diagnostics first
	if resp.Diagnostics.HasError() {
		return // Exit if errors exist
	}
	resp.State.Get(ctx, &plan) // Fetch plan state after confirming no errors
```

This ensures proper error handling and avoids undefined behavior.