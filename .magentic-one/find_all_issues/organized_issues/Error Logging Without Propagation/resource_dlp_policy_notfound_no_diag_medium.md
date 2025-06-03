# Title

Insufficient Error Handling in Read Operation for NotFound Scenario

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

In the `Read` method, when a DLP policy is not found (matches `customerrors.ERROR_OBJECT_NOT_FOUND`), the code calls `resp.State.RemoveResource(ctx)` and immediately returns. However, there is no informational message logged or exposed in diagnostics to indicate to the user why the resource was removed. This could lead to confusion in debugging Terraform states, especially for providers with non-obvious error handling behaviors.

## Impact

Severity: Medium

This impacts user experience and supportability; users may find it unclear why a resource has disappeared from state, leading to confusion or unnecessary troubleshooting of infrastructure drift.

## Location

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
	resp.State.RemoveResource(ctx)
	return
}
```

## Code Issue

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
	resp.State.RemoveResource(ctx)
	return
}
```

## Fix

Add an entry to `resp.Diagnostics` to record that the resource was not found and has been removed from state.

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
	resp.Diagnostics.AddWarning(
		fmt.Sprintf("%s Not Found", r.FullTypeName()),
		fmt.Sprintf("The resource with ID %s was not found and has been removed from the state.", state.Id.ValueString()),
	)
	resp.State.RemoveResource(ctx)
	return
}
```
---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_dlp_policy_notfound_no_diag_medium.md
