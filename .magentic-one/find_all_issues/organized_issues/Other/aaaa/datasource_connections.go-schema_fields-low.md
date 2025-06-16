# Issue: Snake case naming in struct fields (e.g., `display_name`)

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections.go

## Problem

Several attribute names in the schema definition for the data source use snake_case (e.g., `display_name`, `connection_parameters_set`). Although Go itself does not enforce a schema attribute style, Terraform convention and Go SDK recommendations often prefer lowerCamelCase for schema attributes to improve consistency and cross-provider familiarity.

## Impact

Severity: **Low**

This is a minor consistency issue. Adhering to convention helps maintain user familiarity and eases documentation and automation across providers and tools.

## Location

```go
"connections": schema.ListNestedAttribute{
    ...
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{ ... },
            "name": schema.StringAttribute{ ... },
            "display_name": schema.StringAttribute{ ... },
            "status": schema.SetAttribute{ ... },
            "connection_parameters": schema.StringAttribute{ ... },
            "connection_parameters_set": schema.StringAttribute{ ... },
        },
    },
},
```

## Code Issue

```go
"display_name": schema.StringAttribute{
    MarkdownDescription: "Display name of the connection.",
    Computed:            true,
},
```

## Fix

Change the schema field names to lower camelCase for consistency with Terraform convention:

```go
"displayName": schema.StringAttribute{
    MarkdownDescription: "Display name of the connection.",
    Computed:            true,
},
"connectionParameters": schema.StringAttribute{ ... },
"connectionParametersSet": schema.StringAttribute{ ... },
```
Update all related field lookups and documentation accordingly.

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/datasource_connections.go-schema_fields-low.md`
