# Title

Direct type assertion on req.ProviderData without type-safety

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set.go

## Problem

The code forcibly asserts the type of `req.ProviderData` to `*api.ProviderClient` without any type assertion check, which could result in a panic if the type is not as expected.

## Impact

High.

- Can cause provider to panic at runtime instead of gracefully handling the error.

## Location

```go
client := req.ProviderData.(*api.ProviderClient).Api
if client == nil {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
r.EnvironmentGroupRuleSetClient = NewEnvironmentGroupRuleSetClient(client, tenant.NewTenantClient(client))
```

## Fix

Add an ok-check prior to accessing fields, and return a helpful error if the assertion fails:

```go
providerClient, ok := req.ProviderData.(*api.ProviderClient)
if !ok || providerClient.Api == nil {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
client := providerClient.Api
r.EnvironmentGroupRuleSetClient = NewEnvironmentGroupRuleSetClient(client, tenant.NewTenantClient(client))
```
