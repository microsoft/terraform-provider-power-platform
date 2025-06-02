# Validators Applied to Computed-Only Schema Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

Attributes such as `default_connectors_classification` and fields nested under `custom_connectors_patterns` have validators defined, but are `Computed: true`, e.g.:

```go
"default_connectors_classification": schema.StringAttribute{
    MarkdownDescription: "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
    Computed:            true,
    Validators: []validator.String{
        stringvalidator.OneOf("General", "Confidential", "Blocked"),
    },
},
```

Validators are ignored for attributes that are not user-supplied (`Computed: true` and not `Optional`/`Required`). Including them is misleading and may signal to maintainers that some user input is being checked, which is not the case.

## Impact

Reduces clarity and can confuse other developers, who may believe the validator has runtime significance. It may also increase provider memory usage slightly and fails static checks in schema lint tools. Severity: **low**.

## Location

- Any attribute that is only `Computed` but includes validators.

## Code Issue

```go
"default_connectors_classification": schema.StringAttribute{
    ...
    Computed:            true,
    Validators: []validator.String{
        stringvalidator.OneOf("General", "Confidential", "Blocked"),
    },
},
// Also under custom_connectors_patterns.data_group, etc.
```

## Fix

Remove the `Validators` from `Computed`-only fields:

```go
"default_connectors_classification": schema.StringAttribute{
    MarkdownDescription: "Default classification for connectors (\"General\", \"Confidential\", \"Blocked\")",
    Computed:            true,
},
```

**Explanation:**  
Validators only affect user-supplied attributes, so don't add them unless the field is user-modifiable.
