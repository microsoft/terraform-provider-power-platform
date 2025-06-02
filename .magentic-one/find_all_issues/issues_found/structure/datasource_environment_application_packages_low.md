# Issue 2: Typo in `ApplicaitonId` and Typo in Key

##

Path: /workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages.go

## Problem

There is a typo in the schema attribute's markdown description for `application_id` field: "ApplicaitonId" instead of "ApplicationId".

## Impact

Severity: **Low**

This typo does not affect code execution but may confuse users reading the generated provider documentation and lower the perceived API quality.

## Location

```go
"application_id": schema.StringAttribute{
    MarkdownDescription: "ApplicaitonId",
    Computed:            true,
},
```

## Code Issue

```go
"application_id": schema.StringAttribute{
    MarkdownDescription: "ApplicaitonId",
    Computed:            true,
},
```

## Fix

Fix the typo in the markdown:

```go
"application_id": schema.StringAttribute{
    MarkdownDescription: "ApplicationId",
    Computed:            true,
},
```
