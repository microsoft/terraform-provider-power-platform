# Inconsistent markdown in Schema descriptions

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

Within the Schema method, the attribute descriptions are set via MarkdownDescription, but some of them include just the simple string and others could benefit from more clear markdown (for example, code ticks for property names, or stronger formatting). This inconsistency can reduce the perceived documentation quality.

## Impact

Low: Documentation and visual output in Terraform docs/UI may be less readable or uniform.

## Location

Schema attributes block in Schema method.

## Code Issue

```go
"tenant_id": schema.StringAttribute{
    MarkdownDescription: "Tenant ID of the application.",
    Computed:            true,
},
"aad_country_geo": schema.StringAttribute{
    MarkdownDescription: "AAD country geo.",
    Computed:            true,
},
...
```

## Fix

Ensure all descriptions use consistent formatting, e.g., code ticks for attribute names:

```go
"tenant_id": schema.StringAttribute{
    MarkdownDescription: "`tenant_id`: Tenant ID of the application.",
    Computed:            true,
},
```
And consider bolding or adding more Markdown emphasis for better clarity where appropriate.
