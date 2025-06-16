# Title

Incorrect MarkdownDescription for 'description' in Schema

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go

## Problem

The `description` field in the `Schema` has a `MarkdownDescription` that incorrectly states: `"Display name of the environment group"`. This duplicates the description of the `display_name` field and is inaccurate. 

## Impact

Users will be confused about what the `description` attribute actually represents, which could result in incorrect usage or configuration.

**Severity:** low

## Location

Function: `Schema`, attribute `"description"`

## Code Issue

```go
"description": schema.StringAttribute{
    MarkdownDescription: "Display name of the environment group",
    Required:            true,
},
```

## Fix

Update MarkdownDescription to correctly describe the `description` field.

```go
"description": schema.StringAttribute{
    MarkdownDescription: "Description of the environment group",
    Required:            true, // or Optional, as per previous feedback
},
```
