# Title

Unnecessary nil check on EnvironmentGroupClient in `Delete` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go`

## Problem

The `Delete` function checks and attempts repeated calls on the `EnvironmentGroupClient` even if the initial call to `DeleteEnvironmentGroup` fails. The error handling logic could benefit from consolidation to avoid redundant code execution.

## Impact

Unnecessary function calls or excessive resource checks reduce the maintainability of the code and can lead to performance inefficiencies when managing resources. Severity: High

## Location

```go
// Delete function.
func (r *EnvironmentGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    ...
    err := r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
        return
    }

    if customerrors.Code(err) == customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND || customerrors.Code(err) == customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP {
        envs, err := r.EnvironmentGroupClient.GetEnvironmentsInEnvironmentGroup(ctx, state.Id.ValueString())
        ...
        err = r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
        if err != nil {
            resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
            return
        }
    }
}
```

## Fix

Simply exit early after displaying the error in the first occurrence, rather than relying on redundant calls.

```go
// Delete function.
func (r *EnvironmentGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    ...
    err := r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
    if err != nil {
        if customerrors.Code(err) == customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND || customerrors.Code(err) == customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP {
           ...
        } else {
           resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
           return
        }
	}
}
```