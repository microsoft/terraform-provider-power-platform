# Large Schema Function, Highly Nested Structure

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go

## Problem

The `Schema` method of `TenantSettingsDataSource` is extremely large and consists of deeply nested attribute declarations to describe all possible tenant settings. This results in a method that is very difficult to read, scan, debug, and extend. It increases the risk of unintentional errors (omitting commas, incorrect nesting, hard-to-locate settings) and makes consistency and future modifications daunting for maintainers.

## Impact

Complicates maintenance and readability, increases cognitive overhead for new contributors, and heightens the risk of subtle bugs due to difficult-to-review code. **Severity: Medium**

## Location

Main body of the `Schema` method, especially

```go
func (d *TenantSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
   ...
        Attributes: map[string]schema.Attribute{
            "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
                Create: false,
                Update: false,
                Delete: false,
                Read:   false,
            }),
            ...
            "power_platform": schema.SingleNestedAttribute{
                ...
                Attributes: map[string]schema.Attribute{
                    ... // deep nesting for dozens of settings
                },
            },
        },
   ...
}
```

## Fix

Factor the deepest/nested groups of attributes into their own helper functions that return `map[string]schema.Attribute` values. For example:

```go
func powerPlatformAttributes() map[string]schema.Attribute {
    return map[string]schema.Attribute{
        // place Power Platform nested attributes here...
    }
}

// Then in Schema:
"power_platform": schema.SingleNestedAttribute{
    MarkdownDescription: "Power Platform",
    Computed:            true,
    Attributes:          powerPlatformAttributes(),
},
```

This reduces function size and nesting, and helps enforce separation of concerns. Repeat for other major top-level groups ("governance", "power_apps", etc.) where applicable.
