# Issue: Duplicate Documentation Strings

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connections.go

## Problem

The documentation strings (MarkdownDescription) for `connection_parameters` and `connection_parameters_set` attributes are identical and repeated. This repetition can lead to confusion for users, making it unclear if there's any practical difference between the two fields.

## Impact

Severity: **Low**

Redundant documentation decreases clarity, increases maintenance cost, and may result in user support requests about the difference (if any) between these fields.

## Location

```go
"connection_parameters": schema.StringAttribute{
    MarkdownDescription: "... (for example)[https://learn.microsoft.com...]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
    Computed:            true,
},
"connection_parameters_set": schema.StringAttribute{
    MarkdownDescription: "... (for example)[https://learn.microsoft.com...]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
    Computed:            true,
},
```

## Code Issue

```go
"connection_parameters": schema.StringAttribute{
    MarkdownDescription: "Connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
    Computed:            true,
},
"connection_parameters_set": schema.StringAttribute{
    MarkdownDescription: "Set of connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
    Computed:            true,
},
```

## Fix

Refactor the documentation to clarify the difference, or provide a single, non-redundant documentation string for both:

```go
"connection_parameters": schema.StringAttribute{
    MarkdownDescription: "Connection parameters as a JSON string containing the authentication parameters. Used for interactive connections; see [documentation](https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal) for details.",
    Computed:            true,
},
"connection_parameters_set": schema.StringAttribute{
    MarkdownDescription: "Set of connection parameters as a JSON string used for service principal or automated connections. See [documentation](https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal) for more details.",
    Computed:            true,
},
```

---

This issue will be saved in:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/datasource_connections.go-duplicate-docs-low.md`
