# Title

***Lack of Error Handling in Delete Method***

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go`

## Problem

The `Delete` method calls `resp.State.RemoveResource(ctx)` but does not provide any diagnostics or error handling in case there are runtime issues or the state removal fails unexpectedly. While Terraform handles most state manipulation reliably, a failure in resource state removal could lead to incomplete clean-ups or inconsistent states, and its absence in diagnostics could make debugging challenging.

## Impact

- **Severity:** Medium
- No logs or diagnostics for a failed `RemoveResource` operation can leave users in confusion about the state of the resource after deletion.
- Errors during state removal could propagate silently without being caught, leading to inconsistent behavior or long-term state tracking issues.

## Location

**Function Name:**  
`func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse)`

## Code Issue

```go
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()

    resp.State.RemoveResource(ctx)
}
```

## Fix

Add a diagnostic message or error handling mechanism to log the status of the `RemoveResource` call, ensuring users are informed if there are issues.

### Fixed Code Example:

```go
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()

    err := resp.State.RemoveResource(ctx)

    if err != nil {
        resp.Diagnostics.AddError(
            "Failed to Delete Resource",
            fmt.Sprintf("An error occurred while trying to remove the resource from the state: %s", err),
        )
    } else {
        tflog.Info(ctx, "Successfully deleted resource from state.")
    }
}
```

### Why Fix This Way?
- Ensures proper error handling for unexpected failures during state removal.
- Logs the operation success, helping users and maintainers understand whether deletion was completed successfully.
- Adds fail-safe mechanisms that improve program robustness and reliability.