# Title
Schema Definition Duplication

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem
The `connectorSchema` is defined multiple times within nested attributes. While this approach works, it leads to code duplication and reduces maintainability. If changes are required in the schema structure, they would need to be modified in multiple locations.

## Impact
The duplicated schema definitions make the code harder to maintain and increase the likelihood of inconsistencies or errors during updates. This is a medium-severity issue as it impacts code quality and maintainability but does not result in runtime errors.

## Location
Function: Schema

## Code Issue
```go
connectorSchema := schema.NestedAttributeObject{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            MarkdownDescription: "ID of the connector",
            Optional:            true,
        },
        "default_action_rule_behavior": schema.StringAttribute{
            MarkdownDescription: "Default action rule behavior for the connector (\"Allow\", \"Block\")",
            Optional:            true,
            Validators: []validator.String{
                stringvalidator.OneOf("Allow", "Block", ""),
            },
        },
        // Rest of the schema omitted for brevity
    },
}

// Appears in multiple nested attributes such as `business_connectors` and `non_business_connectors`
```

## Fix
Define the `connectorSchema` as a reusable constant or function that can be invoked whenever needed. This improves maintainability and consistency.

```go
var connectorSchema = schema.NestedAttributeObject{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            MarkdownDescription: "ID of the connector",
            Optional:            true,
        },
        "default_action_rule_behavior": schema.StringAttribute{
            MarkdownDescription: "Default action rule behavior for the connector (\"Allow\", \"Block\")",
            Optional:            true,
            Validators: []validator.String{
                stringvalidator.OneOf("Allow", "Block", ""),
            },
        },
        // Rest of the schema omitted for brevity
    },
}

// Reuse `connectorSchema` in nested attributes
"business_connectors": schema.SetNestedAttribute{
    MarkdownDescription: "Connectors for sensitive data",
    Computed:            true,
    NestedObject:        connectorSchema,
},

"non_business_connectors": schema.SetNestedAttribute{
    MarkdownDescription: "Connectors for non-sensitive data",
    Computed:            true,
    NestedObject:        connectorSchema,
},
```