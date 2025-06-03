# Excessive Function Length in Schema Function

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors.go

## Problem

The `Schema` method contains a large block for building up the schema attributes inline, making the function long and less readable. Such inline schemas are harder to maintain, refactor, or extend, especially as attribute numbers grow.

## Impact

**Low severity:** Maintainability and readability are negatively impacted, increasing the risk of introducing errors during future modifications and making code more difficult for new contributors to understand.

## Location

```go
func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
    defer exitContext()
    resp.Schema = schema.Schema{
        MarkdownDescription: "...",
        Attributes: map[string]schema.Attribute{
            "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
                Read: true,
            }),
            "connectors": schema.ListNestedAttribute{
                MarkdownDescription: "List of Connectors",
                Computed:            true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        // ... many repeated pattern fields ...
                    },
                },
            },
        },
    }
}
```

## Fix

Move complex attribute/nested schema definitions to variables or dedicated helper functions. Example:

```go
var connectorAttributes = map[string]schema.Attribute{
    "id":   schema.StringAttribute{MarkdownDescription: "Id", Computed: true},
    "name": schema.StringAttribute{MarkdownDescription: "Name", Computed: true},
    // ... other attributes ...
}

var schemaAttributes = map[string]schema.Attribute{
    "timeouts": timeouts.Attributes(ctx, timeouts.Opts{ Read: true }),
    "connectors": schema.ListNestedAttribute{
        MarkdownDescription: "List of Connectors",
        Computed:            true,
        NestedObject: schema.NestedAttributeObject{
            Attributes: connectorAttributes,
        },
    },
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
    defer exitContext()
    resp.Schema = schema.Schema{
        MarkdownDescription: "...",
        Attributes: schemaAttributes,
    }
}
```
