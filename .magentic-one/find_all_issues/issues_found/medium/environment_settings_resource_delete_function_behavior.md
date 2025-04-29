# Title

Lack of Implementation in Delete Function

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resources_environment_settings.go

## Problem

The `Delete` function lacks a clear implementation or rationale for its "Do nothing on purpose" comment. This raises concerns about whether this behavior is intentional and appropriate for all use cases.

## Impact

The missing implementation can lead to confusion, as users and maintainers may not understand why the resource deletion logic is omitted intentionally. If the resource deletion is needed, this could result in inconsistency or leave orphaned resources, potentially affecting system integrity.

Severity: **Medium**

## Location

Line ~297: Function `Delete`.

## Code Issue

```go
func (r *EnvironmentSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()
    // Do nothing on purpose
}
```

## Fix

Add a meaningful implementation or provide clear documentation justifying the omission, ensuring that users and maintainers understand its purpose.

```go
func (r *EnvironmentSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()

    // Currently, no deletion logic is implemented as resources may be immutable.
    // If the resource needs to be deleted physically, implement the appropriate cleanup logic.
    tflog.Info(ctx, "Delete function executed. No deletion logic implemented intentionally.")
}
```
