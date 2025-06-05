# Provider Data Type Assertion Safety Issues

This document contains type assertion safety issues related to provider data handling, configuration, and resource provider data type assertions in the codebase.

## ISSUE 1

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

---

## ISSUE 2

# Title

ProviderData type assertion error handling silently continues

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go

## Problem

In the Configure function for the resource, the type assertion for ProviderData to *api.ProviderClient is performed and, if it fails, an error diagnostic is added and return is called. However, for the next type assertion for client.Api, an error is reported that refers to*http.Client, which does not actually match the real client type, causing misleading messages and a possible silent failure if the provider data shape changes. Also, the earlier type assertion for *api.ProviderClient does not panic or halt, but carries on error handling just by message--in future code refactoring, this can make debugging configuration issues hard to diagnose.

## Impact

If ProviderData is not the expected type, resource endpoints will silently misbehave or fail, causing provider malfunctions or hidden state drift. Additionally, the error message may confuse users ("Expected *http.Client..." when it's a nil Api that fails, not a wrong type).

## Location

In the Configure method:

## Code Issue

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
clientApi := client.Api

if clientApi == nil {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )

    return
}
r.ManagedEnvironmentClient = newManagedEnvironmentClient(clientApi)
```

## Fix

Ensure error messages are clear and type-correct. When Api is nil, report that explicitly, not as a wrong type. For example:

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
if client.Api == nil {
    resp.Diagnostics.AddError(
        "Nil Api client",
        "ProviderData contained a *api.ProviderClient but with nil Api. Please check provider initialization and credentials.",
    )
    return
}
r.ManagedEnvironmentClient = newManagedEnvironmentClient(client.Api)
```

This prevents confusion and makes diagnostics cleaner for provider maintainers and users.

---

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
