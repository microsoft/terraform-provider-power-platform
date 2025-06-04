# Inconsistent error handling in Read and Update methods

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest.go

## Problem

In the `Update` and `Read` methods, error handling is not fully consistent. In the `Update` method, a pointer (`plan *DataverseWebApiResourceModel`) is used, and it is set via `req.Plan.Get(ctx, &plan)`, and similar for the Read method with `state`. However, in the `Read` method, after executing the "read" operation, any new output received is not assigned to the state. In contrast, in `Create` and `Update`, state/output is always replaced.

Additionally, after calling SendOperation (e.g., in `Read`), the new output value `bodyWrapped` is never actually set to the new state, and thus the state saved might be stale. This may lead to out-of-sync resource output and missed drift detection.

## Impact

Medium. This may cause state inconsistencies and prevent correct drift detection and reporting in Terraform, leading to potentially confusing user experience or state corruption if resource data changes remotely.

## Location

`Read` method, lines around:
```go
	if state.Read != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Read)
		if err != nil {
			resp.Diagnostics.AddError("Error executing read operation", err.Error())
			return
		}

		if state.Output.String() != bodyWrapped.String() {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("true"))
		} else {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("false"))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
```

## Code Issue

```go
	if state.Read != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Read)
		if err != nil {
			resp.Diagnostics.AddError("Error executing read operation", err.Error())
			return
		}

		if state.Output.String() != bodyWrapped.String() {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("true"))
		} else {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("false"))
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
```

## Fix

Assign the `bodyWrapped` value to `newState.Output` for proper state update after a successful read. This will ensure that the state is always in sync with the actual read operation results.

```go
	if state.Read != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Read)
		if err != nil {
			resp.Diagnostics.AddError("Error executing read operation", err.Error())
			return
		}

		if state.Output.String() != bodyWrapped.String() {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("true"))
		} else {
			resp.Private.SetKey(ctx, "force_value_unknown", []byte("false"))
		}

		newState.Output = bodyWrapped // <-- Update the output here
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
```
