# Error suppression in Update and Delete: Diagnostics not always propagated

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest.go

## Problem

In the `Update` and `Delete` methods, after executing the operation and on success, the method does not propagate any diagnostics from `resp.State.Set`. Any errors that occur while setting state will not cause an error to be returned, which may lead to silent failures or lost diagnostics. This is inconsistent with other methods which append the result of state saving to diagnostics.

## Impact

Low to Medium. This could result in errors being silenced and not surfaced to the user, hampering debugging and correct Terraform operation.

## Location

`Delete` method (and potentially others). Example lines:
```go
func (r *DataverseWebApiResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var state *DataverseWebApiResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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
}
```

## Code Issue

```go
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

After mutation of the state/output, save to state and propagate diagnostics as in other methods:

```go
	if state.Destroy != nil {
		bodyWrapped, err := r.DataRecordClient.SendOperation(ctx, state.Destroy)
		if err != nil {
			resp.Diagnostics.AddError("Error executing destroy operation", err.Error())
			return
		}
		state.Output = bodyWrapped
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
```
