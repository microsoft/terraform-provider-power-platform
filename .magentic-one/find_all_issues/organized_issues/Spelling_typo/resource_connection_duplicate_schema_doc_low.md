# Title
Duplicated schema documentation for "connection_parameters" and "connection_parameters_set"

##
/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection.go

## Problem
Both the `connection_parameters` and `connection_parameters_set` attributes in the schema have essentially identical documentation strings, including the same note with a typo ("requried in-place-update"). This introduces ambiguity about their distinct purposes, hinders readability, and may confuse users and maintainers.

## Impact
Low: Does not affect correctness, but reduces overall documentation clarity for both code maintainers and resource users.

## Location
Resource Schema definition:

## Code Issue
```go
"connection_parameters": schema.StringAttribute{
    MarkdownDescription: "Connection parameters. Json string containing the authentication connection parameters ...",
    ...
},
"connection_parameters_set": schema.StringAttribute{
    MarkdownDescription: "Set of connection parameters. Json string containing the authentication connection parameters ...",
    ...
},
```

## Fix
Review and clarify the documentation for each field so that their specific roles and any differences are apparent, removing typos and duplicated information where not relevant.

```go
"connection_parameters": schema.StringAttribute{
    MarkdownDescription: "Connection parameters. JSON string with authentication details, used when ...",
    ...
},
"connection_parameters_set": schema.StringAttribute{
    MarkdownDescription: "(Advanced) An explicit set of connection parameters (JSON string), for ... [explain use case and distinction].",
    ...
},
```

Also, fix the typo: "requried in-place-update" → "required in-place update".

