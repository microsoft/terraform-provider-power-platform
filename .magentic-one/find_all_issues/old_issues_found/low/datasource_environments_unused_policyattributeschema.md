# Title
Unused or Redundant `policyAttributeSchema` Variable in `Schema` Method

##
/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments.go

## Problem

The variable `policyAttributeSchema` is defined and used once within the `Schema` method to set the attributes for `enterprise_policies`. However, this abstraction is unnecessary as it is used directly and does not add clarity or reduce repetition.

## Impact

While this does not affect the functionality, it reduces the code clarity, adds unnecessary variable initialization, and slightly impacts performance and code readability when this pattern is repeated across a codebase.

Severity: Low

## Location

```go  
policyAttributeSchema := map[string]schema.Attribute{  
    "type": schema.StringAttribute{  
        MarkdownDescription: "Type of the policy according to [schema definition](https://learn.microsoft.com/en-us/azure/templates/microsoft.powerplatform/enterprisepolicies?pivots=deployment-language-terraform#enterprisepolicies-2)",  
        Computed:            true,  
    },  
    "id": schema.StringAttribute{  
        MarkdownDescription: "Id (guid)",  
        Computed:            true,  
    },  
    "location": schema.StringAttribute{  
        MarkdownDescription: "Location of the policy",  
        Computed:            true,  
    },  
    "system_id": schema.StringAttribute{  
        MarkdownDescription: "System id (guid)",  
        Computed:            true,  
    },  
    "status": schema.StringAttribute{  
        MarkdownDescription: "Link status of the policy",  
        Computed:            true,  
    },  
}  
```

## Fix

Eliminate the unnecessary variable declaration by directly assigning a map literal inside the attribute `enterprise_policies`.

```go  
"enterprise_policies": schema.SetNestedAttribute{  
    MarkdownDescription: "Enterprise policies for the environment. See [Enterprise policies](https://learn.microsoft.com/en-us/power-platform/admin/enterprise-policies) for more details.",  
    Computed:            true,  
    NestedObject: schema.NestedAttributeObject{  
        Attributes: map[string]schema.Attribute{  
            "type": schema.StringAttribute{  
                MarkdownDescription: "Type of the policy according to [schema definition](https://learn.microsoft.com/en-us/azure/templates/microsoft.powerplatform/enterprisepolicies?pivots=deployment-language-terraform#enterprisepolicies-2)",  
                Computed:            true,  
            },  
            "id": schema.StringAttribute{  
                MarkdownDescription: "Id (guid)",  
                Computed:            true,  
            },  
            "location": schema.StringAttribute{  
                MarkdownDescription: "Location of the policy",  
                Computed:            true,  
            },  
            "system_id": schema.StringAttribute{  
                MarkdownDescription: "System id (guid)",  
                Computed:            true,  
            },  
            "status": schema.StringAttribute{  
                MarkdownDescription: "Link status of the policy",  
                Computed:            true,  
            },  
        },  
    },  
},  
```

Explanation:

- Reduces unnecessary variable initialization (`policyAttributeSchema`).  
- Ensures attributes are set directly within the logical context of the parent object.  
- Simplifies readability by removing indirection.