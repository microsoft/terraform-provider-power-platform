# Structure and Maintainability: Large Method for Schema Definition

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

The `Schema` method contains a deeply nested and lengthy literal schema definition with 8+ levels of indentation and many attributes. This can be difficult to read, maintain, and extend, as it is easy for errors or copy-paste mistakes to occur. It is also hard to quickly grasp the schema structure, especially for more complex objects.

## Impact

- Maintenance is challenging; adding or removing fields is error prone.
- Code review and comprehension are harder for complex, deeply nested object schemas.
- Can introduce subtle inconsistencies in style, docs, or required/computed flags over time.

**Severity:** Medium

## Location

The method `func (d *AnalyticsExportDataSource) Schema(...)` and its schema literal.

## Code Issue

```go
func (d *AnalyticsExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
    defer exitContext()
    resp.Schema = schema.Schema{
        // ... huge nested literal ...
    }
}
```

## Fix

Move repeated or deeply nested attributes into helper functions or variables. For example:

```go
var sinkAttribute = schema.SingleNestedAttribute{
    MarkdownDescription: "The sink configuration for analytics data",
    Required: true,
    Attributes: map[string]schema.Attribute{
        // ...
    },
}

// In Schema method:
"sink": sinkAttribute,
```

This modularizes the schema definition, improves readability, enables reuse, and reduces risk of copy-paste or indentation errors.
