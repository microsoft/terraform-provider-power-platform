# Title

No Defensive Checks for nil plan or state in CRUD Methods

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

In the CRUD methods (`Read`, `Create`, `Update`, `Delete`), the code directly assigns into fields of potentially nil pointers (`plan`, `state`). If the value from `req.Plan.Get` or `req.State.Get` is nil due to misconfiguration or changes in Terraform core, a nil pointer dereference panic will occur.

## Impact

Severity: Critical

Critical stability and reliability bug. Any runtime change from Terraform, or upstream changes in schema/parser behavior, could cause a panic and crash the provider, leading to loss of in-flight state, partial infrastructure, and requiring manual operator intervention.

## Location

```go
var plan *dataLossPreventionPolicyResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
// Immediately:
plan.Id = types.StringValue(policy.Name)
```

## Code Issue

```go
var plan *dataLossPreventionPolicyResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
// plan may still be nil here!
plan.Id = types.StringValue(policy.Name)
```

## Fix

Add a nil check after unmarshal and before dereferencing/assigning to fields.

```go
var plan *dataLossPreventionPolicyResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
if plan == nil {
	resp.Diagnostics.AddError("Internal Provider Error", "Plan was nil after reading configuration.")
	return
}
```
(Apply the same for state in all CRUD methods.)

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_dlp_policy_nil_plan_critical.md
