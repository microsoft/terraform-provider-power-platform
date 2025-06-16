# Title

Potential Redundancy and Confusion: Both `Required` and `Optional` Tags Set Explicitly

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

In various schema attribute definitions, both `Required` and `Optional` fields are explicitly set, for example:

```go
"scope": schema.StringAttribute{
    MarkdownDescription: "Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)",
    Required:            true,
    Optional:            false,
},
```

According to the Terraform Plugin Framework conventions, a field is either required, optional, or computed. Explicitly setting both can reduce readability and cause maintenance confusion as the mutually exclusive nature is implicit in the framework.

## Impact

While the logic functions as intended (since `Required` and `Optional` are mutually exclusive), it reduces code clarity and could lead to confusion or mistakes during maintenance and future edits. This is a low-severity, but relevant for code readability and maintainability.

## Location

In the map of attributes for `resp.Schema` in the `Schema` method:

## Code Issue

```go
"scope": schema.StringAttribute{
    MarkdownDescription: "Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)",
    Required:            true,
    Optional:            false,
},
// ... same pattern for other String and Set attributes ...
```

## Fix

Only set the field that is applicable. For required fields, set only `Required: true` and omit `Optional`. For optional fields, use `Optional: true` and omit `Required`.

```go
"scope": schema.StringAttribute{
    MarkdownDescription: "Authentication scope for the request. See more: [Authentication Scopes](https://learn.microsoft.com/en-us/entra/identity-platform/scopes-oidc)",
    Required:            true,
},
// for optional:
"body": schema.StringAttribute{
    MarkdownDescription: "Body of the request",
    Optional:            true,
},
```

This improves readability and helps prevent ambiguity in intent.
