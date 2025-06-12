# Title

Potential Error Masking by Not Appending Diagnostics on State Setting

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

In the `Read` function, after setting `state.Output`, you set the state again and append the diagnostics. However, you immediately check for errors and return, which is redundant because diagnostics are already appended, and in Terraform Plugin SDKs, the diagnostics should generally be returned if present, but it's best practice to append diagnostics in-place and return only if a critical condition is hit (not just always immediately after).

## Impact

This statement is not a bug, but is minorly misleading in idiomatic Go/Terraform Plugin SDK code, potentially confusing future maintainers or causing over-defensive exit flow. It is of low severity.

## Location

At the end of the `Read` method:

## Code Issue

```go
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
```

## Fix

Just append the diagnostics. Returning from `Read` after appending diagnostics that might not actually be errors is an overly defensive code pattern. Instead, you may append and only return early in more critical flows. For clarity, this block could be trimmed as:

```go
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
```

If you want to retain the `HasError` early return, consider adding a comment and restructuring for clarity. If this is code style for all your providers, it is not strictly wrong, but be aware that the plugin will handle the diagnostics array, and returning is usually needed only for control flow exit, not always after every set.

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/datasource_rest_query.go_return_after_state_set_low.md`
