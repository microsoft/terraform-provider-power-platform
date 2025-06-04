# Title

Improper error message and mismatched type check in Configure

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go

## Problem

Error message in `Configure` returns `"Expected *http.Client, got: %T"`, though we are expecting `*api.ProviderClient`.

## Impact

Low.

- The message can mislead users and hinder debugging.

## Location

```go
resp.Diagnostics.AddError(
    "Unexpected Resource Configure Type",
    fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
)
```

## Fix

Update `Expected *http.Client` to the correct type. For example:

```go
resp.Diagnostics.AddError(
    "Unexpected Resource Configure Type",
    fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
)
```
