# Missing Error Handling for Setting State in `Create` and `Read` Methods

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go

## Problem

In both the `Create` and `Read` methods, after assembling the new state, the code sets the state with `resp.State.Set(ctx, &newState)` (or `&state`), and immediately appends any diagnostics to the response. However, if `Set` returns an error, the function does not check for or handle it (e.g., with `HasError()`) and continues, potentially causing confusion or masking errors related to state persistence. This can lead to diagnostics not being surfaced as prominently as necessary or state being partially set.

## Impact

- Possible silent failures when state setting fails.
- Makes diagnosing errors harder for both users and maintainers.
- Can leave the system in an inconsistent state.
- Severity: **Medium**

## Location

In the `Create` method (and similarly in `Read`):

```go
resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
```

## Code Issue

```go
resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
```

## Fix

After appending diagnostics, always check for errors before proceeding, and return if errors exist.

```go
resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
if resp.Diagnostics.HasError() {
	return
}
```

Apply the same pattern to the analogous state setting in the `Read` method as well:

```go
resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
if resp.Diagnostics.HasError() {
	return
}
```

---

This markdown will be saved under:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_enterprise_policy_error_handling_medium.md`
