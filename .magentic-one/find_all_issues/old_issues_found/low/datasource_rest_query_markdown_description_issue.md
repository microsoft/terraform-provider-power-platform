# Title

Potential Misuse of `MarkdownDescription`

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

In the `Schema` function, `MarkdownDescription` is used excessively and redundantly, in situations where more concise documentation formatting (like plain text) might suffice. This can add unnecessary complexity and make it difficult for contributors to maintain consistency with description texts across different schemas. 

## Impact

The issue is mainly a stylistic concern but can lead to widespread inconsistency across the project, making the codebase harder to read and maintain. The severity is **low**, as this does not directly affect runtime functionality.

## Location

The issue surfaces in the `Schema` function in the `MarkdownDescription` attributes under schema definitions.

## Code Issue

```go
"scope": schema.StringAttribute{
    MarkdownDescription: "Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)",
    Required:            true,
    Optional:            false,
},
"method": schema.StringAttribute{
    MarkdownDescription: "HTTP method",
    Required:            true,
    Optional:            false,
},
"url": schema.StringAttribute{
    MarkdownDescription: "Absolute url of the api call",
    Required:            true,
    Optional:            false,
},
```

## Fix

Using concise documentation frameworks can help reduce redundancy and improve maintainability of code.

```go
"scope": schema.StringAttribute{
    Description: "Authentication scope for the request. See more: Authentication Scopes",
    Required:    true,
    Optional:    false,
},
"method": schema.StringAttribute{
    Description: "HTTP method",
    Required:    true,
    Optional:    false,
},
"url": schema.StringAttribute{
    Description: "Absolute url of the api call",
    Required:    true,
    Optional:    false,
},
```