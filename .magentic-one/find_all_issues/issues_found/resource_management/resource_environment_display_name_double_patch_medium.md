# Title

Resource Management: Redundant Update Logic for Display Name Causes Extra PATCH Call

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go

## Problem

When updating a resource, if the plan’s display name differs from the state’s display name, the code issues a PATCH request twice in succession with the same display name and other environment properties. The comment acknowledges this as a temporary fix for a backend propagation issue, but it is not optimal as it leads to redundant API calls during a single update operation.

## Impact

- **Severity**: Medium
- May result in double PATCH operations or rate limits on API usage.
- Increases latency for the update operation due to unnecessary network call.
- Might mask the need for an actual resolution of the backend bug, possibly allowing this workaround to persist indefinitely.

## Location

```go
envDto, err := r.EnvironmentClient.UpdateEnvironment(ctx, plan.Id.ValueString(), environmentDto)
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
	return
}

// This is a temporary fix for the issue in BAPI where the display name is not propagated correctly on environment update
if plan.DisplayName.ValueString() != state.DisplayName.ValueString() {
	envDto, err = r.EnvironmentClient.UpdateEnvironment(ctx, plan.Id.ValueString(), environmentDto)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}
}
```

## Code Issue

```go
// PATCH called once above already
// PATCH is called a second time if display name changes
```

## Fix

Add a TODO with reference to a tracking issue or feature flag
OR  
Implement a condition or flag to ensure the PATCH is not repeated if the issue is resolved—or to only perform the double PATCH if a specific configuration or feature flag is enabled, thus minimizing impact.

Example workaround with minimal impact:

```go
// TODO: Remove double PATCH for DisplayName once backend propagation issue #1234 is fixed
if plan.DisplayName.ValueString() != state.DisplayName.ValueString() && workaroundEnabled { // 'workaroundEnabled' could be sourced from config/env
	envDto, err = r.EnvironmentClient.UpdateEnvironment(ctx, plan.Id.ValueString(), environmentDto)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
		return
	}
}
```

---
