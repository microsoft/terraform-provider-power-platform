# Title

Misleading Markdown Descriptions in Schema Attribute Definitions

##

Path: `/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go`

## Problem

Several attributes in the `Schema` definition have misleading or incorrect `MarkdownDescription` values. For instance:
- `modified_time`, `install_time`, and `version` all share the description `"Created time"`, which is incorrect.

## Impact

This can confuse users implementing the schema in their infrastructure and cause misconfigurations or require unnecessary debugging. Severity is medium due to potential user friction.

## Location

`Schema` method, in attributes `modified_time`, `install_time`, and `version`.

## Code Issue

```go
"modified_time": schema.StringAttribute{
    MarkdownDescription: "Created time", // This is incorrect
    Computed:            true,
},
"install_time": schema.StringAttribute{
    MarkdownDescription: "Created time", // This is incorrect
    Computed:            true,
},
"version": schema.StringAttribute{
    MarkdownDescription: "Created time", // This is incorrect
    Computed:            true,
},
```

## Fix

Update the `MarkdownDescription` to reflect the correct context.

```go
"modified_time": schema.StringAttribute{
    MarkdownDescription: "Last modification time of the solution",
    Computed:            true,
},
"install_time": schema.StringAttribute{
    MarkdownDescription: "Installation time of the solution",
    Computed:            true,
},
"version": schema.StringAttribute{
    MarkdownDescription: "Version of the solution",
    Computed:            true,
},
```