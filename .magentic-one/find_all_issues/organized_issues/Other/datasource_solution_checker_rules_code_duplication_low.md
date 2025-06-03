# Title

Potential Code Duplication and Boilerplate in Schema Definition

##

/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/datasource_solution_checker_rules.go

## Problem

There is a large, manually-constructed schema tree for the `rules` list attribute, with each nested attribute explicitly enumerated. If this pattern is repeated for many resources/data sources, it increases the risk of code duplication. This could be mitigated by extracting repeated attribute structures into shared helper functions or variables, improving code reuse.

## Impact

Severity is **low** in this file alone, but across the project it can lead to maintainability burdens. If multiple schema attributes' structures change frequently, the effort to apply changes grows and bugs can slip in.

## Location

In the Schema method, definition of the `rules` nested list attributes:

```go
"rules": schema.ListNestedAttribute{
    ...
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            ...
        },
    },
},
```

## Code Issue

```go
"rules": schema.ListNestedAttribute{
    MarkdownDescription: "List of solution checker rules",
    Computed:            true,
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "code": schema.StringAttribute{ ... },
            // many more ...
        },
    },
},
```

## Fix

Extract the rules schema into a shared helper for reuse in other places, for example:

```go
func rulesAttributes() map[string]schema.Attribute {
    return map[string]schema.Attribute{
        "code": schema.StringAttribute{ ... },
        "description": schema.StringAttribute{ ... },
        // ... etc ...
    }
}
// And in the schema definition:
"rules": schema.ListNestedAttribute{
    MarkdownDescription: "List of solution checker rules",
    Computed:            true,
    NestedObject: schema.NestedAttributeObject{
        Attributes: rulesAttributes(),
    },
},
```
Do this particularly if the same model appears in multiple resources or data sources.