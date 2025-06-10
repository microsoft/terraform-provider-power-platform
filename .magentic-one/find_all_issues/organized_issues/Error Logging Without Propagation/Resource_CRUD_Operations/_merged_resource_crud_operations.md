# Error Logging Without Propagation - Resource CRUD Operations

This document consolidates all issues related to error logging without proper propagation found in resource CRUD operation implementations across the Terraform Provider for Power Platform.

## ISSUE 1

# Title

Error Handling Omits Logging and Telemetry

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

Throughout error handling in the CRUD methods (`Read`, `Create`, `Update`, `Delete`), errors are added to `resp.Diagnostics` for user-facing reporting, but there is no consistent or explicit use of logging/telemetry hooks (e.g., via `tflog.Error` or similar) to capture these failures for operator or maintainer telemetry. This means provider developers and support engineers may miss error signals, delays in debugging or RCA, and lose important context for production issues.

## Impact

Severity: Medium

Medium supportability risk: the lack of back-end logs or provider telemetry makes diagnosing production or field issues slow and less efficient, particularly if dealing with cloud operational environments (i.e., managed Terraform Cloud or Enterprise).

## Location

```go
if err != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err_client.Error())
 return
}
```

## Code Issue

```go
if err_client != nil {
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err_client.Error())
 return
}
```

(similar in Read, Update, Delete)

## Fix

Add explicit logging before adding diagnostics, e.g.

```go
if err_client != nil {
 tflog.Error(ctx, fmt.Sprintf("Create error: %s", err_client.Error()))
 resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err_client.Error())
 return
}
```

Or, standardize `tflog.Error` calls alongside diagnostic errors throughout the resource.

## ISSUE 2

# Title

Insufficient Error Handling in Read Operation for NotFound Scenario

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

In the `Read` method, when a DLP policy is not found (matches `customerrors.ERROR_OBJECT_NOT_FOUND`), the code calls `resp.State.RemoveResource(ctx)` and immediately returns. However, there is no informational message logged or exposed in diagnostics to indicate to the user why the resource was removed. This could lead to confusion in debugging Terraform states, especially for providers with non-obvious error handling behaviors.

## Impact

Severity: Medium

This impacts user experience and supportability; users may find it unclear why a resource has disappeared from state, leading to confusion or unnecessary troubleshooting of infrastructure drift.

## Location

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
 resp.State.RemoveResource(ctx)
 return
}
```

## Code Issue

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
 resp.State.RemoveResource(ctx)
 return
}
```

## Fix

Add an entry to `resp.Diagnostics` to record that the resource was not found and has been removed from state.

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
 resp.Diagnostics.AddWarning(
  fmt.Sprintf("%s Not Found", r.FullTypeName()),
  fmt.Sprintf("The resource with ID %s was not found and has been removed from the state.", state.Id.ValueString()),
 )
 resp.State.RemoveResource(ctx)
 return
}
```

## ISSUE 3

# Missing Error Handling for os.Getwd in importSolution

##

/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem

In the function `importSolution`, the call to `os.Getwd()` ignores the error return value. The statement:

```go
cwd, _ := os.Getwd()
```

simply discards the error. This method can return an error, and if the current directory isn't available for any reason, the value of `cwd` could be invalid or misleading. Not even logging the error means that diagnostics or troubleshooting is hindered if such a condition arises.

## Impact

- **Severity:** Medium
- Masking possible failures makes troubleshooting harder and hides potential edge cases in filesystem access from the user and maintainers.
- Logging or diagnosing environment issues becomes more difficult if errors are silently ignored.

## Location

`importSolution` function, where `os.Getwd()` is used.

## Code Issue

```go
cwd, _ := os.Getwd()
tflog.Debug(ctx, fmt.Sprintf("Current working directory: %s", cwd))
```

## Fix

Capture and log the error for better diagnostics:

```go
cwd, err := os.Getwd()
if err != nil {
    tflog.Warn(ctx, fmt.Sprintf("Failed to get working directory: %s", err.Error()))
} else {
    tflog.Debug(ctx, fmt.Sprintf("Current working directory: %s", cwd))
}
```

## ISSUE 4

# Redundant Logging: Use of tflog.Warn for Successful Checksum Calculation

##

/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem

Within the `Create` function, after successfully calculating a checksum for the settings or solution file, the code logs this with `tflog.Warn`:

```go
tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of settings file: %s", value))
```

Checksum calculation is a successful and routine operation, so warning level logging is inappropriate. It may create noise in logs, distracting from real warnings. This may reflect miscommunication of severity/intent or copy-paste error.

## Impact

- **Severity:** Low
- Could mislead users/maintainers or clutter logs
- Minor impact, but relevant for operational hygiene and clarity

## Location

In the `Create` function's handling of the settings and solution file checksums.

## Code Issue

```go
if err != nil {
    resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
} else {
    plan.SettingsFileChecksum = types.StringValue(value)
    tflog.Warn(ctx, fmt.Sprintf("CREATE Calculated md5 hash of settings file: %s", value))
}
```

## Fix

Lower the log level to `tflog.Debug` or remove it if not actually useful:

```go
if err != nil {
    resp.Diagnostics.AddWarning("Issue when calculating checksum for settings file", err.Error())
} else {
    plan.SettingsFileChecksum = types.StringValue(value)
    tflog.Debug(ctx, fmt.Sprintf("CREATE calculated SHA256 hash of settings file: %s", value))
}
```

## ISSUE 5

# Error Handling for Unexpected Resource Type in Configure

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

Within the `Configure` function, if the provider data is not of the expected type (`*api.ProviderClient`), an error is appended to diagnostics, but no further control flow actions are taken (such as aborting or returning early). Although the code later returns after appending the error, it is best practice to always follow an error-adding diagnostic block with an explicit `return` to avoid potential future code additions after the error handling that could cause logic bugs.

## Impact

- **Severity: Low**  
  The current design is safe as there is a `return` immediately after, but it is best practice for future maintainability to make error returns explicit after diagnostic error appends that are meant to halt processing.

## Location

```go
if req.ProviderData == nil {
 // ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
 return
}

client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
 resp.Diagnostics.AddError(
  "Unexpected Resource Configure Type",
  fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
 )
 return
}
```

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
```

## Fix

You should continue always placing explicit `return` after `AddError` in error conditions, as is currently done. No modification is strictly needed, but developers should be aware of this flow and avoid code additions after error blocks.

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
 resp.Diagnostics.AddError(
  "Unexpected Resource Configure Type",
  fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
 )
 return // Explicit, to avoid unintended code execution after error
}
```

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
