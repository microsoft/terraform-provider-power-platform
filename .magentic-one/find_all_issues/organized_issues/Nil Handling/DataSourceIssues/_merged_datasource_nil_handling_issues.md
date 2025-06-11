# Data Source Nil Handling Issues

This document contains all identified nil handling issues related to data source components in the Terraform Power Platform provider codebase.

## ISSUE 1

<!-- Source: datasource_analytics_data_exports_go_error_handling_high.md -->

# Error Handling on analyticsDataExport `nil` Return

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

When `analyticsDataExport` is `nil` after calling `GetAnalyticsDataExport`, a generic error is reported:

```
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "Unable to find analytics data export with the specified ID",
    )
    return
}
```

There is a lack of context to help diagnose why `nil` was returned (i.e. no ID in request, filtering logic, or an API-level explanation), and it's not clear if returning `nil` means a logical 'not found' or an internal error occurred. The error message gives a potentially misleading impression there is an ID-based lookup (none is present in code shown; possibly a copy-paste mistake in the message).

## Impact

This could confuse users and maintainers, making debugging more difficult and user feedback inadequate.

**Severity:** High.

## Location

```go
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "Unable to find analytics data export with the specified ID",
    )
    return
}
```

## Code Issue

```go
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "Unable to find analytics data export with the specified ID",
    )
    return
}
```

## Fix

Return a more helpful error explaining the likely reasons and removing reference to an ID (if not present). If another layer returns context, surface it here. If `nil` is legal and indicates a non-error empty state, consider returning early or sending a warning instead.

```go
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "No analytics data exports were returned from the downstream API. This may indicate a configuration issue or that no exports exist for this tenant.",
    )
    return
}
```

## ISSUE 2

<!-- Source: datasource_analytics_data_exports_go_resource_management_low.md -->

# Potential Performance Issue: Use of Pointers to Slices

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

The code pattern `len(*analyticsDataExport)` is used, indicating `analyticsDataExport` is a pointer to a slice:

```go
exports := make([]AnalyticsDataModel, 0, len(*analyticsDataExport))
for _, export := range *analyticsDataExport {
    if model := convertDtoToModel(&export); model != nil {
        exports = append(exports, *model)
    }
}
```

In Go, slices are already reference types. Passing and storing them as pointers to slices usually has negligible benefit, and can increase cognitive load and complexity. Unless mutation or a nil distinction is required, functions should accept and return slices directly.

## Impact

Minor performance/readability issue, as it complicates reasoning for memory ownership. Not a leak, but opportunity for simplification and clarity.

**Severity:** Low.

## Location

```go
len(*analyticsDataExport)
for _, export := range *analyticsDataExport {
```

## Code Issue

```go
analyticsDataExport, err := d.analyticsExportClient.GetAnalyticsDataExport(ctx)
// ...
exports := make([]AnalyticsDataModel, 0, len(*analyticsDataExport))
for _, export := range *analyticsDataExport {
```

## Fix

If possible, use a slice instead of pointer to slice:

```go
// In the client, return a slice, not a pointer to a slice
analyticsDataExport, err := d.analyticsExportClient.GetAnalyticsDataExport(ctx)
// ...
exports := make([]AnalyticsDataModel, 0, len(analyticsDataExport))
for _, export := range analyticsDataExport {
    if model := convertDtoToModel(&export); model != nil {
        exports = append(exports, *model)
    }
}
```

## ISSUE 3

<!-- Source: datasource_billing_policies_missing_error_handling_low.md -->

# Missing Error Handling for `LicensingClient` Initialization

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies.go

## Problem

In the `Configure` method, the initialization of `NewLicensingClient` is not checked for errors. If `NewLicensingClient` were to return an error (e.g., due to a nil or invalid API client), this would not be captured and could cause panics or undefined behavior later. Although the current implementation suggests no error is returned, robust code should always account for future changes.

## Impact

This issue is a "low" severity as the current implementation of `NewLicensingClient` may not fail, but if refactored, fail-safe patterns should be in place. Not handling potential errors could lead to panics or non-obvious failures if the client is ever changed to error.

## Location

`Configure` method, line:

```go
d.LicensingClient = NewLicensingClient(client.Api)
```

## Code Issue

```go
d.LicensingClient = NewLicensingClient(client.Api)
```

## Fix

Check for error on client construction — if `NewLicensingClient` can return an error, it should be handled.

```go
licensingClient, err := NewLicensingClient(client.Api)
if err != nil {
    resp.Diagnostics.AddError(
        "Failed to create Licensing client",
        fmt.Sprintf("Could not initialize licensing client: %s", err.Error()),
    )
    return
}
d.LicensingClient = licensingClient
```

## ISSUE 4

<!-- Source: datasource_connectors_connectorsclient-nil-validation-medium.md -->

# No Validation or Checks on ConnectorsClient Construction

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors.go

## Problem

The `Configure` method assigns `d.ConnectorsClient = newConnectorsClient(client.Api)` but does not validate whether the result is non-nil or if the provided API client is valid. If `newConnectorsClient` can return a nil value or an improperly initialized client (depending on API evolutions or failures), later API calls could panic or fail unclearly.

## Impact

**Medium severity:** Possible panics or unclear errors at later points in execution, especially if any dependency involved in constructing the client fails or changes behavior in the future.

## Location

```go
d.ConnectorsClient = newConnectorsClient(client.Api)
```

## Fix

Validate that `newConnectorsClient` returns a non-nil instance (or proper values), and add a check to append a diagnostic error and return if the client was not constructed properly:

```go
d.ConnectorsClient = newConnectorsClient(client.Api)
if d.ConnectorsClient == nil {
    resp.Diagnostics.AddError(
        "Failed to create connectors client",
        "Connectors client returned nil. Check provider configuration and upstream client logic.",
    )
    return
}
```

## ISSUE 5

<!-- Source: datasource_connectors_nil-slice-initialization-low.md -->

# Lack of State Initialization Before Read in Read Method

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/datasource_connectors.go

## Problem

In the `Read` method, the `state` object is decoded from the request's state, but there is no initialization to ensure `state.Connectors` is a non-nil (empty) slice prior to appending connector models. If `state.Connectors` is nil, repeated Read invocations or certain state transitions could result in a nil slice being marshalled to Terraform, which may not be handled consistently by the framework or downstream consumers.

## Impact

**Low to Medium severity:** Possible risk of data inconsistencies or unexpected nil slices propagating to the Terraform state, which could lead to intermittent deserialization issues or subtle schema mismatches.

## Location

```go
var state ListDataSourceModel
resp.State.Get(ctx, &state)

for _, connector := range connectors {
    connectorModel := convertFromConnectorDto(connector)
    state.Connectors = append(state.Connectors, connectorModel)
}
```

## Fix

Ensure `state.Connectors` is initialized to an empty slice if nil, before appending new elements:

```go
var state ListDataSourceModel
resp.State.Get(ctx, &state)
if state.Connectors == nil {
    state.Connectors = []ConnectorModel{}
}

for _, connector := range connectors {
    connectorModel := convertFromConnectorDto(connector)
    state.Connectors = append(state.Connectors, connectorModel)
}
```

## ISSUE 6

<!-- Source: datasource_environments_nil_slice_low.md -->

# Nil Pointer Risk on Append to Slice

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments.go

## Problem

In the `Read` function, `state.Environments` is being appended to without checking if it is nil. If `state.Environments` has not been initialized (is nil), appending to it will still work in Go, but it may cause issues with serialization or downstream expectations if the slice is expected to be non-nil (e.g., always return an empty list instead of `null`). This is especially important for TF state models.

## Impact

Can cause confusion for consumers or downstream code expecting a non-nil list (`[]`) rather than `null`. Severity: **Low** (Go handles nil-slice appends, but empty list is generally preferred for model consistency).

## Location

```go
var state ListDataSourceModel
...
for _, env := range envs {
    ...
    state.Environments = append(state.Environments, *env)
}
```

## Code Issue

```go
var state ListDataSourceModel

for _, env := range envs {
    ...
    state.Environments = append(state.Environments, *env)
}
```

## Fix

Pre-initialize the slice or ensure it is never nil for better consistency in returned data.

```go
var state ListDataSourceModel
state.Environments = make([]EnvironmentModel, 0, len(envs)) // EnvironmentModel is a placeholder; use actual type.

for _, env := range envs {
    ...
    state.Environments = append(state.Environments, *env)
}
```

## ISSUE 7

<!-- Source: datasource_locations_NoNilCheckForLocationsClient_High.md -->

# No Nil Check for LocationsClient Before Method Call

##

/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations.go

## Problem

In the `Read` method, code calls `d.LocationsClient.GetLocations(ctx)` without ever verifying that `d.LocationsClient` is non-nil. If the `Configure` method was not called or failed to set up the client (or if it was set to nil by some external event or mistake), calling a method on a nil struct would result in a runtime panic.

## Impact

This oversight could cause panics during provider execution, leading to plugin crashes and poor user/developer experience. This is a **High** severity control flow and error handling issue.

## Location

In the `Read` method:

```go
locations, err := d.LocationsClient.GetLocations(ctx)
```

## Code Issue

```go
locations, err := d.LocationsClient.GetLocations(ctx)
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
 return
}
```

## Fix

Check `d.LocationsClient` for `nil` before usage, and add a meaningful diagnostic if it is not configured:

```go
if d.LocationsClient == nil {
 resp.Diagnostics.AddError(
  "Locations client not configured",
  "The locations client was not configured. This may indicate a problem with the provider setup or authentication.",
 )
 return
}

locations, err := d.LocationsClient.GetLocations(ctx)
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
 return
}
```

## ISSUE 8

<!-- Source: datasource_tenant_api_client_high.md -->

# No nil check for TenantClient in Read method

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

In the Read method, d.TenantClient is used without first checking whether it is nil. If Configure was never called, or if there was an error in configuration/setup, this could lead to a panic during operation.

## Impact

Severity: High. If d.TenantClient is nil, this will cause a runtime panic, causing the provider to crash and breaking the entire Terraform operation.

## Location

In Read method:

## Code Issue

```go
tenant, err := d.TenantClient.GetTenant(ctx)
```

## Fix

Add a nil check prior to use, and fail gracefully:

```go
if d.TenantClient == nil {
    resp.Diagnostics.AddError("Tenant client is not configured", "The TenantClient is nil. This is likely a bug in the provider initialization.")
    return
}

tenant, err := d.TenantClient.GetTenant(ctx)
```

## ISSUE 9

<!-- Source: datasource_tenant_type_safety_medium.md -->

# Type assertion without error handling in Configure

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

In the `Configure` method, the code asserts that `req.ProviderData` is of type `*api.ProviderClient` via:

```go
client, ok := req.ProviderData.(*api.ProviderClient)
```

If the assertion fails, it adds an error to diagnostics, but does not actually return early or halt execution. Instead, the code proceeds, which may lead to further logic executing with a nil or incorrect client, resulting in nil reference errors or unexpected behaviors elsewhere.

## Impact

Medium: Adds error diagnostics but possible logic flow may continue with invalid state, risking subtle bugs or panics later in the lifecycle.

## Location

In `Configure`:

## Code Issue

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    // Should return here to avoid further usage
}

d.TenantClient = NewTenantClient(client.Api)
```

## Fix

Return immediately after adding diagnostics for the assertion failure:

```go
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return // <-- ensure no further logic executes
}
```

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

Apply this fix to the whole codebase
