# Title

Use of plan/state as pointer type instead of value causes possible nil dereference risk

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/resource_copilot_studio_application_insights.go

## Problem

The code unmarshals plan and state into a pointer to `ResourceModel` (`var plan *ResourceModel`). If unmarshalling fails (for example, if Terraform or the provider framework evolves and breaks the state logic), `plan` or `state` could be `nil`, which will cause a panic on dereference (e.g., `plan.BotId.ValueString()`) rather than a controlled diagnostic.

## Impact

Severity: **Medium**

This creates a race risk, since failure in the Get call only populates Diagnostics but doesn't guarantee a non-nil plan; a dereference without checking pointer for nil carries panic risk. Although the current flow checks `resp.Diagnostics.HasError()` after unmarshalling, this is a subtle and easy-to-break contract and does not provide robust type safety for future maintainers.

## Location

- All CRUD methods: `Create`, `Update`, `Read`, `Delete`

## Code Issue

```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
// plan is pointer, but never checked for nil
```

## Fix

Use `ResourceModel` as a value (not pointer) so unmarshalling always succeeds or returns the zero struct. Optionally, if pointer semantics are needed, explicitly check for nil after unmarshalling.

```go
var plan ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
// safe to use plan as non-nil value now
```

This approach aligns with idiomatic Go, prevents nil-dereference panics, and improves maintainability.

---

File to save:  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_copilot_studio_application_insights_plan_pointer_medium.md`
