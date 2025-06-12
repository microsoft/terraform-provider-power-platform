# Code Structure: Redundant Assignments in Create Method

##

/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install.go

## Problem

In the `Create` method, after fetching the plan values, the code redundantly reassigns `EnvironmentId` and `UniqueName` to themselves, already as `types.StringValue(...)`. This is unnecessary unless the values are being normalized, which does not appear to be the case here.

## Impact

Severity: **Low**

Redundant assignments increase noise and can mislead readers to believe that the values are being processed when they're simply copied, slightly impacting code readability and maintainability.

## Location

Within `Create`:
```go
state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), strings.ReplaceAll(strings.ToLower(state.UniqueName.ValueString()), " ", "_")))
state.EnvironmentId = types.StringValue(state.EnvironmentId.ValueString())
state.UniqueName = types.StringValue(state.UniqueName.ValueString())
```

## Code Issue

```go
state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), strings.ReplaceAll(strings.ToLower(state.UniqueName.ValueString()), " ", "_")))
state.EnvironmentId = types.StringValue(state.EnvironmentId.ValueString())
state.UniqueName = types.StringValue(state.UniqueName.ValueString())
```

## Fix

Remove the redundant assignments unless normalization or transformation is required. Only set these values if actual conversion, validation, or business logic is needed.

```go
state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.EnvironmentId.ValueString(), strings.ReplaceAll(strings.ToLower(state.UniqueName.ValueString()), " ", "_")))
// Remove unnecessary copies of EnvironmentId and UniqueName unless transformation is needed
```

---

This output will be saved as:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_environment_application_package_install.go_redundant_assignments-low.md`
