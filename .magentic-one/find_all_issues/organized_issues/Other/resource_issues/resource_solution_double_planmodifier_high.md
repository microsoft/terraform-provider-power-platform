# Double Usage of SyncAttributePlanModifier in Schema attributes

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem
The `solution_file_checksum` and `settings_file_checksum` attributes in the resource schema both use `PlanModifiers` with two calls to `modifiers.SyncAttributePlanModifier(...)`, referencing the same target attribute. There's no clear necessity or documentation for having this modifier twice for each attribute.

This makes the code harder to read and could lead to unexpected behavior if the plan modifier is not idempotent or if its implementation ever changes. It also signals potential copy-paste issues and reduces the clarity of the plan modification logic for maintainers.

## Impact
- **Severity:** High
- Unnecessary code duplication can increase maintenance effort and confusion.
- Potential for side effects or performance issues if a modifier is not fully idempotent or if the implementation changes in the future.
- It makes the codebase harder to reason about and may mask intention. 

## Location
Line(s) where `PlanModifiers` for `solution_file_checksum` and `settings_file_checksum` are defined. Example excerpt:

## Code Issue
```go
"solution_file_checksum": schema.StringAttribute{
    MarkdownDescription: "Checksum of the solution file",
    Computed:            true,
    PlanModifiers: []planmodifier.String{
        modifiers.SyncAttributePlanModifier("solution_file"),
        modifiers.SyncAttributePlanModifier("solution_file"),
    },
},
...
"settings_file_checksum": schema.StringAttribute{
    MarkdownDescription: "Checksum of the settings file",
    Computed:            true,
    PlanModifiers: []planmodifier.String{
        modifiers.SyncAttributePlanModifier("settings_file"),
        modifiers.SyncAttributePlanModifier("settings_file"),
    },
},
```

## Fix
Remove the duplicate `SyncAttributePlanModifier` call for each attribute. One instance is sufficient:

```go
"solution_file_checksum": schema.StringAttribute{
    MarkdownDescription: "Checksum of the solution file",
    Computed:            true,
    PlanModifiers: []planmodifier.String{
        modifiers.SyncAttributePlanModifier("solution_file"),
    },
},
...
"settings_file_checksum": schema.StringAttribute{
    MarkdownDescription: "Checksum of the settings file",
    Computed:            true,
    PlanModifiers: []planmodifier.String{
        modifiers.SyncAttributePlanModifier("settings_file"),
    },
},
```
