# Title

Partial Field Plan Modifier coverage for Computed Fields

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

Some computed fields (for example: `"last_modified_by"` and `"last_modified_time"`) do not use `PlanModifiers` such as `UseStateForUnknown()`, while others do. This inconsistency can lead to cases where the Terraform plan is unable to preserve values that the backend manages, causing state/plan mismatches after apply or updates.

## Impact

Severity: Medium

Can lead to incorrect Terraform plan output, state drift, and unintuitive user experience when fields change only on the backend.

## Location

```go
"created_by": schema.StringAttribute{
	MarkdownDescription: "User who created the policy",
	Computed:            true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	},
},
"last_modified_by": schema.StringAttribute{
	MarkdownDescription: "User who last modified the policy",
	Computed:            true,
},
"last_modified_time": schema.StringAttribute{
	MarkdownDescription: "Time when the policy was last modified",
	Computed:            true,
},
```

## Code Issue

```go
// Only some computed fields use PlanModifiers
"created_by": schema.StringAttribute{
	MarkdownDescription: "User who created the policy",
	Computed:            true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	},
},
"last_modified_by": schema.StringAttribute{
	MarkdownDescription: "User who last modified the policy",
	Computed:            true,
	// missing PlanModifiers
},
```

## Fix

Add `PlanModifiers: []planmodifier.String{ stringplanmodifier.UseStateForUnknown() },` to all computed fields that have backend controlled values to preserve state and avoid unnecessary drift.

```go
"last_modified_by": schema.StringAttribute{
	MarkdownDescription: "User who last modified the policy",
	Computed:            true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	},
},
"last_modified_time": schema.StringAttribute{
	MarkdownDescription: "Time when the policy was last modified",
	Computed:            true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	},
},
```
---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_dlp_policy_missing_plan_modifiers_medium.md
