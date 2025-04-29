# Issue 3: Unnecessary Code in `Update` Method

##

`/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go`

## Problem

The `Update` method contains code that appends diagnostics from `Plan.Get` and `State.Get`. However, the method states that no updates can be made, as all attribute changes would require a delete and recreate. The inclusion of this code section is redundant and does not serve any functional purpose.

## Impact

This redundancy can lead to confusion for readers and developers, who may assume the `Update` function is capable of making updates. It increases the complexity of the code without providing any value.

Severity: **Low**

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

## Fix

Remove the redundant code and simplify the method implementation.

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()

	// No updates can be made; attribute changes require delete and recreate.
	tflog.Info(ctx, fmt.Sprintf("Update operation is not supported for resource %s", r.FullTypeName()))
}
```

This makes clear that the update operation is not supported and reduces unnecessary code overhead.