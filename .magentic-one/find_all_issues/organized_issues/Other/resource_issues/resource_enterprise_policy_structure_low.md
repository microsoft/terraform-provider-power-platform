# Issue: Noop Update Method with Incomplete Comment

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go

## Problem

The `Update` method contains only a comment stating that there's nothing to update and that any attribute change requires recreate. However, the method still retrieves both plan and state, appends their diagnostics, and shows a confusing structure with dead code. This is misleading, adds noise, and could create confusion for future maintainers regarding whether an update is ever allowed.

## Impact

- Reduces code clarity and maintainability.
- May lead to incorrect assumptions by future contributors regarding the resource's capabilities.
- **Severity:** Low

## Location

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	var plan *sourceModel
	var state *sourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// There is nothing to update in this resource, as any attribute change would require a delete and create.
}
```

## Code Issue

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// ... snip ...
	// There is nothing to update in this resource, as any attribute change would require a delete and create.
}
```

## Fix

If the resource is not updatable (i.e., all updates require recreation), simplify the `Update` to explicitly state this by just returning, and add a comment explaining why:
```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This resource does not support in-place updates. All changes require resource replacement.
	return
}
```
Alternatively, if plan/state are required for logging or diagnostics, document that purpose.

---

This markdown will be saved under:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_enterprise_policy_structure_low.md`
