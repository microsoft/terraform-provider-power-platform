# Redundant State Removal During Resource Read

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

In the `Read` method, the state is removed both when the feature lookup returns a not-found error and when the returned feature is `nil`. This may be correct defensive programming if both cases occur in practice (for example, if the client returns `nil, nil`), but having to check both may indicate unclear API contracts or create subtle redundancy. This pattern can hide bugs in the client function or inflate code.

## Impact

Medium severity, as this can create confusion or hide upstream/client problems, and might result in maintainers missing cases where improper `nil, nil` is returned from the client, making error diagnosis harder.

## Location

```go
feature, err := r.EnvironmentWaveClient.GetFeature(ctx, state.EnvironmentId.ValueString(), state.FeatureName.ValueString())
if err != nil {
	if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
	return
}

if feature == nil {
	resp.State.RemoveResource(ctx)
	return
}
```

## Code Issue

```go
feature, err := r.EnvironmentWaveClient.GetFeature(ctx, state.EnvironmentId.ValueString(), state.FeatureName.ValueString())
if err != nil {
	if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
	return
}

if feature == nil {
	resp.State.RemoveResource(ctx)
	return
}
```

## Fix

Ensure the upstream API contract for `GetFeature` is clear:
- Returns `nil, error` for not found, and NEVER `nil, nil`.
- Or
- Returns `nil, nil` for not found, but never error for not found.
- Adjust the state removal to match the documented contract, and add an in-code comment explaining the reason for both removals if both are required for safety.

For example:

```go
// Defensive: client may return (nil, nil) for not found as well as error
if feature == nil {
	resp.State.RemoveResource(ctx)
	return
}
```

Or document or fix the client, so only one check is needed.