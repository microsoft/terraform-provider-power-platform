# Documentation and Descriptions Issues

This document contains merged issues related to documentation and descriptions in the Power Platform Terraform provider.

## ISSUE 1

**Title:** Inconsistent markdown in Schema descriptions

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go`

**Problem:**
Within the Schema method, the attribute descriptions are set via MarkdownDescription, but some of them include just the simple string and others could benefit from more clear markdown (for example, code ticks for property names, or stronger formatting). This inconsistency can reduce the perceived documentation quality.

**Impact:**
Low: Documentation and visual output in Terraform docs/UI may be less readable or uniform.

**Location:**
Schema attributes block in Schema method.

**Code Issue:**

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

**Fix:**
Ensure all descriptions use consistent formatting, e.g., code ticks for attribute names:

```go
"tenant_id": schema.StringAttribute{
    MarkdownDescription: "`tenant_id`: Tenant ID of the application.",
    Computed:            true,
},
```

And consider bolding or adding more Markdown emphasis for better clarity where appropriate.

## ISSUE 2

**Title:** Inconsistent use of Description vs MarkdownDescription

**File:** `/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

**Problem:**
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

**Impact:**
Inconsistent documentation formatting in Terraform UI and documentation; maintainers and users can be confused as to which fields support markdown and which don't. Severity: low.

**Location:**
Within schema definitions for "teams_integration", "power_apps", etc.

**Code Issue:**

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

**Fix:**
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

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
