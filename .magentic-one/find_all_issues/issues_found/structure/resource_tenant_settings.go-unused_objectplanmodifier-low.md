# Dead/commented code for unused `objectplanmodifier.UseStateForUnknown()`

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

Within the schema for the `power_apps` nested attribute, there's a commented-out plan modifier:

```go
PlanModifiers: []planmodifier.Object{
	// objectplanmodifier.UseStateForUnknown(),
},
```

Leaving dead/commented code in place is not a good practice unless it's for a planned future change and actively communicated. At best, it adds noise; at worst, it confuses about implementation intentions, especially in shared codebases.

## Impact

Detracts from code clarity and cleanliness, especially for collaborators reading or maintaining the file. Severity: low.

## Location

Within schema registration of the `power_apps` attribute.

## Code Issue

```go
PlanModifiers: []planmodifier.Object{
	// objectplanmodifier.UseStateForUnknown(),
},
```

## Fix

Remove the commented line altogether — or, if this modifier should be active, uncomment and ensure proper import and function.

```go
PlanModifiers: []planmodifier.Object{},
```
or (if you actually want it)
```go
PlanModifiers: []planmodifier.Object{
	objectplanmodifier.UseStateForUnknown(),
},
```

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_tenant_settings.go-unused_objectplanmodifier-low.md`
