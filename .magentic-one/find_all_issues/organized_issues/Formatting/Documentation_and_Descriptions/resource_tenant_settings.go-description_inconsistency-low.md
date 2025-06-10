# Inconsistent use of Description vs MarkdownDescription

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

Within the resource schema, some attributes use `Description` while others use `MarkdownDescription`. For example:

```go
"teams_integration": schema.SingleNestedAttribute{
    Description: "Teams Integration",
    Optional:    true,
    ...
},
"power_apps": schema.SingleNestedAttribute{
    Description:   "Power Apps",
    Optional:      true,
    ...
},
"power_platform": schema.SingleNestedAttribute{
    MarkdownDescription: "Power Platform",
    Optional:            true,
    ...
},
```

Framework conventions favor consistent documentation strings, with `MarkdownDescription` preferred because it supports richer formatting and is rendered better in Terraform documentation and UIs.

## Impact

Inconsistent documentation formatting in Terraform UI and documentation; maintainers and users can be confused as to which fields support markdown and which don’t. Severity: low.

## Location

Within schema definitions for "teams_integration", "power_apps", etc.

## Code Issue

```go
"teams_integration": schema.SingleNestedAttribute{
    Description: "Teams Integration",
    Optional:    true,
    ...
},
"power_apps": schema.SingleNestedAttribute{
    Description:   "Power Apps",
    Optional:      true,
    ...
},
```

## Fix

Use `MarkdownDescription` for all attribute/documentation strings for consistency.

```go
"teams_integration": schema.SingleNestedAttribute{
    MarkdownDescription: "Teams Integration",
    Optional:            true,
    ...
},
"power_apps": schema.SingleNestedAttribute{
    MarkdownDescription: "Power Apps",
    Optional:            true,
    ...
},
```

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_tenant_settings.go-description_inconsistency-low.md`
