# Title

Incorrect Required Attribute on 'description' Field

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go

## Problem

In the `Schema` function, the `description` attribute is marked as `Required: true`, which may be misleading if `description` should actually be optional or computed. Typically, a description is an optional or computed field unless the resource design explicitly demands it.

## Impact

If `description` is not strictly required when creating/updating an environment group, keeping it as `Required: true` would enforce unnecessary constraints and could cause user confusion or failing applies if the field is omitted.

**Severity:** medium

## Location

Function: `Schema`, attribute definition for `"description"`

## Code Issue

```go
"description": schema.StringAttribute{
    MarkdownDescription: "Display name of the environment group",
    Required:            true,
},
```

## Fix

If `description` should be optional (which is standard for descriptions), change `Required: true` to `Optional: true`. Adjust as appropriate given the actual resource requirements.

```go
"description": schema.StringAttribute{
    MarkdownDescription: "Display name of the environment group",
    Optional:            true,
},
```
