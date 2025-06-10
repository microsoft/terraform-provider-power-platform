# Error Logging Without Propagation - Data Source Configuration

This document consolidates all issues related to error logging without proper propagation found in data source configuration implementations across the Terraform Provider for Power Platform.

## ISSUE 1

# Title

Improper Error Handling: Ambiguous Type Assertion Failure

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies_environments.go

## Problem

When asserting the type of `req.ProviderData` to `*api.ProviderClient`, the code simply adds an error to diagnostics and returns if it fails. However, it lacks any user guidance or structured logging, and provides a generic error message without any potential recovery or visibility for troubleshooting.

## Impact

This only logs the error as a diagnostic, which may not be enough to capture the underlying problem especially if the provider internals or dependencies change, possibly making troubleshooting more difficult. Proper error structuring and logging are best practice for maintainability. **Severity: Medium**

## Location

In the `Configure` method:

## Code Issue

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
 resp.Diagnostics.AddError(
  "Unexpected ProviderData Type",
  fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
 )
 return
}
```

## Fix

Enhance the error handling by providing more specific context, or add structured/error logging for further traceability, such as tflog.Error, and consider other diagnostic actions, e.g., fail fast if a critical path.

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
 tflog.Error(ctx, fmt.Sprintf(
  "Configure failed due to unexpected ProviderData type: expected *api.ProviderClient, got: %T", req.ProviderData,
 ))
 resp.Diagnostics.AddError(
  "Unexpected ProviderData Type",
  fmt.Sprintf(
   "Configuration failed: Expected *api.ProviderClient, but got %T. This is likely a bug—please report this issue to the provider developers.", req.ProviderData,
  ),
 )
 return
}
```

This ensures improved logging and better communication in both diagnostics and logging outputs.

## ISSUE 2

# Direct diagnostic propagation after resp.State.Set

##

/workspaces/terraform-provider-power-platform/internal/services/connection/datasource_connection_shares.go

## Problem

In the `Read` method, diagnostics from `resp.State.Set` are appended to `resp.Diagnostics`. If any errors are present, an early return is executed. While this is common in Terraform plugins, it is a control flow touchpoint that can benefit from more explicit error handling and possible state clean-up or logging to better support debugging and maintainability.

## Impact

**Severity: medium**

If setting the state fails, only a non-specific error will be returned. There is room for improvement by adding contextual information or logging, which would facilitate debugging complex state issues in production.

## Location

```go
diags := resp.State.Set(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
 return
}
```

## Code Issue

```go
diags := resp.State.Set(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
 return
}
```

## Fix

Consider logging or enriching the diagnostic message before returning; at minimum, add a comment describing control flow intent:

```go
diags := resp.State.Set(ctx, &state)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
 // State could not be set—returning to prevent invalid state.
 // Optionally, add contextual logging here if needed for debugging.
 return
}
```

## ISSUE 3

# Missing Diagnostic Handling for State Unmarshalling

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem

In the `Read` function, after retrieving state:

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

This error check is immediate, which is correct. However, there is a missing diagnostic log or report in the event of an error—execution simply returns. This may make debugging hard, as nothing is logged, and the user may see a silent failure.

## Impact

When errors arise during state retrieval, there is no indication in logs as to why `Read` has returned early, which damages debuggability for both maintainers and end-users. Severity: **medium**.

## Location

`Read` function, state unmarshalling block.

## Code Issue

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

## Fix

Log or otherwise inform the user that state retrieval failed:

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    tflog.Error(ctx, "failed to unmarshal current state in Read: "+resp.Diagnostics.Errors()[0].Summary)
    return
}
```

**Explanation:**  
This provides better traceability in logs/debugging, helping developers quickly identify at what point the function has exited.

## ISSUE 4

# Issue 1: Error Handling Missing Return after AddError

##

Path: /workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages.go

## Problem

In the `Read` function, after adding an error diagnostic when `d.ApplicationClient.DataverseExists` returns an error, there is no `return` statement. The function continues to execute even when the error happens, which can lead to further issues or panics as the downstream logic may depend on successful completion of this check.

## Impact

Severity: **High**

Continuing execution after hitting an error means the next code may operate under invalid or unexpected states. This could cause incorrect diagnostics, nil pointer dereferences, further noise in logs, or panics.

## Location

```go
d.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
```

## Code Issue

```go
d.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
```

## Fix

Add a `return` after reporting the diagnostics error, so further processing is aborted:

```go
dvExits, err := d.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
 return
}
```

## ISSUE 5

# Type Safety: Potential Silent Failure in State Reading

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps.go

## Problem

When reading state using `resp.State.Get(ctx, &state)`, diagnostics are appended, and the code returns if errors are present. However, there's no explicit error handling or log for what went wrong if state reading fails. This could make debugging more difficult.

## Impact

Severity: Medium  
Silent failures can make issues hard to diagnose, specially in a production environment.

## Location

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

## Code Issue

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

## Fix

Optionally, log a debug message when an error occurs or add further handling to help debugging.

```go
resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    tflog.Error(ctx, "Failed to read state for EnvironmentPowerAppsDataSource")
    return
}
```

## ISSUE 6

# Title

Potential Missing Return After Error in Client Check

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings.go

## Problem

When calling `DataverseExists`, if there is an error, you log a diagnostic but *do not return*. This may result in further logic running on bad data, and possible misleading or cascading errors.

## Impact

Medium. If an error occurs but you continue, the logic may malfunction, and this could confuse end users and make debugging more difficult.

## Location

Lines:

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
```

## Code Issue

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
}
```

## Fix

Return after adding the diagnostic to avoid further execution.

```go
dvExits, err := d.EnvironmentSettingsClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
 return
}
```

## ISSUE 7

# Possible Return of Partially Created State on Conversion Error

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments.go

## Problem

In the `Read` method, if an error is encountered in `convertSourceModelFromEnvironmentDto` for a specific environment, the function logs an error and returns immediately. However, since the function appends each converted environment to `state.Environments` inside the loop, a partially filled slice (containing only the successful conversions until the error) will be left in `state.Environments`. This could result in an ambiguous or partial state being persisted/used downstream if the `Set` operation is called before error return or if diagnostics do not interrupt subsequent processing.

## Impact

Might result in partially populated state visible to downstream processing (low severity for most scenarios, but could impact environments expecting all-or-nothing). **Severity: Low**.

## Location

```go
for _, env := range envs {
    ...
    env, err := convertSourceModelFromEnvironmentDto(...)
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Error when converting environment %s", env.DisplayName), err.Error())
        return
    }
    state.Environments = append(state.Environments, *env)
}
```

## Code Issue

```go
for _, env := range envs {
    ...
    env, err := convertSourceModelFromEnvironmentDto(...)
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Error when converting environment %s", env.DisplayName), err.Error())
        return
    }
    state.Environments = append(state.Environments, *env)
}
```

## Fix

Reset (nil or empty) `state.Environments` before returning on error, or avoid setting the state at all if a partial result was built before failure:

```go
for _, env := range envs {
    ...
    converted, err := convertSourceModelFromEnvironmentDto(...)
    if err != nil {
        state.Environments = nil // Or: = state.Environments[:0]
        resp.Diagnostics.AddError(fmt.Sprintf("Error when converting environment %s", env.DisplayName), err.Error())
        return
    }
    state.Environments = append(state.Environments, *converted)
}
```

This version ensures that if an error is encountered, no partial list is left in state, preserving an all-or-nothing update strategy and preventing ambiguity.

## ISSUE 8

# Lack of Error Handling After Config Get in Read Method

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go

## Problem

Within the `Read` method of `TenantSettingsDataSource`, the result of `req.Config.Get(ctx, &configuredSettings)` is not checked for diagnostics/errors. This violates robust error handling principles, since if the configuration cannot be decoded (due to, for instance, a schema mismatch or context error), the method continues execution, potentially leading to incorrect logic or even panics downstream.

## Impact

Failing to handle diagnostics can result in unexpected behavior, invalid state propagation, and complicates troubleshooting for users. **Severity: High**

## Location

Line inside `Read` method:

```go
var configuredSettings TenantSettingsDataSourceModel
req.Config.Get(ctx, &configuredSettings)
state, _, err = convertFromTenantSettingsDto[TenantSettingsDataSourceModel](*tenantSettings, state.Timeouts)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error converting tenant settings: %s", d.FullTypeName()), err.Error())
    return
}
```

## Fix

Check and append diagnostics immediately after calling `req.Config.Get`, and abort execution if any errors are present, consistent with prior error checks in this method:

```go
var configuredSettings TenantSettingsDataSourceModel
diags := req.Config.Get(ctx, &configuredSettings)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
state, _, err = convertFromTenantSettingsDto[TenantSettingsDataSourceModel](*tenantSettings, state.Timeouts)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error converting tenant settings: %s", d.FullTypeName()), err.Error())
    return
}
```

This ensures proper error handling, improves code robustness, and aligns with best practices used elsewhere in the provider.

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
