# Error Logging Without Propagation - API Client Issues

This document consolidates all issues related to error logging without proper propagation found in API client implementations across the Terraform Provider for Power Platform.

## ISSUE 1

# Issue 2: Error Not Propagated After Logging

##

/workspaces/terraform-provider-power-platform/internal/services/application/api_application.go

## Problem

When parsing the `operationLocationHeader` fails, the error is only logged but not returned or handled, so the loop proceeds with an invalid URL that could lead to further errors or confusion.

## Impact

This issue has a **medium** severity. Failing to handle or propagate parse errors could result in attempts to make HTTP requests to an invalid URL, leading to confusing errors and wasted resources.

## Location

Within `InstallApplicationInEnvironment`:

## Code Issue

```go
_, err = url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
}
```

## Fix

Return the error immediately after logging it, stopping further execution with an invalid URL:

```go
_, err = url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
    return "", err
}
```

## ISSUE 2

# Issue: Missing Error Handling in `getCurrencies` Function

**File:** `/workspaces/terraform-provider-power-platform/internal/services/currencies/api_currencies.go`

## Problem

The `getCurrencies` function is missing proper error handling for the HTTP client's `Do` method, which could leave users without clear diagnostics when API calls fail.

## Impact

This issue has **medium** severity. Without proper error handling:
- API failures are not properly communicated to users
- Debugging becomes difficult when operations fail
- The function may return incomplete or corrupted data without indication

## Location

Within the `getCurrencies` function:

## Code Issue

```go
func (client *CurrenciesClient) getCurrencies(ctx context.Context) ([]CurrencyDto, error) {
    // ... request setup ...
    resp, err := client.Api.Do(req)
    // Missing: if err != nil { return nil, err }
    
    defer resp.Body.Close()
    // ... rest of function ...
}
```

## Fix

Add proper error handling after the HTTP request:

```go
func (client *CurrenciesClient) getCurrencies(ctx context.Context) ([]CurrencyDto, error) {
    // ... request setup ...
    resp, err := client.Api.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request: %w", err)
    }
    
    defer resp.Body.Close()
    // ... rest of function ...
}
```

## ISSUE 3

# Issue: Header Parsing Error Not Propagated

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go`

## Problem

When adding Dataverse to an environment, header parsing errors are logged but not returned, potentially causing silent failures that are difficult to troubleshoot.

## Impact

This issue has **high** severity. Header parsing failures without proper error propagation can lead to:
- Silent failures that are hard to debug
- Incomplete operations that appear successful
- Resource inconsistency
- User confusion about operation status

## Location

In the `addDataverseToEnvironment` function around line 200:

## Code Issue

```go
addDataverseUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
}
// Continue execution despite parsing error
```

## Root Cause

The error from `url.Parse()` is logged using `tflog.Error()` but not returned or handled, allowing execution to continue with an invalid URL. This could lead to subsequent HTTP operations failing with confusing error messages.

## Fix

Propagate the error immediately after logging:

```go
addDataverseUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
    return "", err
}
```

## Additional Recommendations

1. **Consistent Error Handling**: Ensure all similar header parsing operations throughout the codebase follow the same pattern
2. **Error Context**: Consider adding more context to the error message to help with debugging
3. **Validation**: Add validation for the expected format of the location header before parsing

## ISSUE 4

# Issue: Header Parsing Error Handling in Lifecycle Operations

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go`

## Problem

In the `copyEnvironment` and `backupEnvironment` functions, URL parsing errors from operation location headers are logged but not properly handled, potentially causing silent failures.

## Impact

This issue has **medium** severity. When header parsing fails without proper error propagation:
- Operations may appear to succeed while actually failing
- Debugging becomes difficult due to masked errors
- Resource state may become inconsistent
- Users receive confusing or incomplete error messages

## Location

In multiple lifecycle operation functions:

### `copyEnvironment` function:
```go
operationLocationUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
}
// Execution continues despite error
```

### `backupEnvironment` function:
```go
operationLocationUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
}
// Execution continues despite error
```

## Root Cause

Both functions log parsing errors using `tflog.Error()` but continue execution instead of returning the error. This masks critical failures that should prevent further processing.

## Fix

Return errors immediately after logging in both functions:

### `copyEnvironment` fix:
```go
operationLocationUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
    return "", err
}
```

### `backupEnvironment` fix:
```go
operationLocationUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
    return "", err
}
```

## Additional Recommendations

1. **Consistency**: Apply similar error handling patterns across all lifecycle operations
2. **Error Context**: Add operation-specific context to error messages
3. **Header Validation**: Consider validating header format before parsing
4. **Testing**: Add unit tests to verify error handling behavior

## ISSUE 5

# Issue: Retry-After Header Parsing Error Handling

**File:** `/workspaces/terraform-provider-power-platform/internal/services/environment/api_environment.go`

## Problem

In the `getEnvironmentStateChangePollingWaitTime` function, parsing errors for the `Retry-After` header are logged but not handled, potentially causing polling issues during environment state changes.

## Impact

This issue has **medium** severity. When `Retry-After` header parsing fails:
- Polling intervals may not be respected, potentially overwhelming the API
- Resource operations may timeout or fail unexpectedly
- API rate limiting may not be properly handled
- Debugging polling issues becomes difficult

## Location

In the `getEnvironmentStateChangePollingWaitTime` function:

## Code Issue

```go
retryAfterSeconds, err := strconv.Atoi(retryAfterHeader)
if err != nil {
    tflog.Warn(ctx, "Error parsing retry after header: "+err.Error())
}
```

## Root Cause

The function logs parsing errors but continues execution without setting an appropriate fallback value or returning an error. This could result in:
- Zero wait time being used (potentially overwhelming the API)
- Inconsistent polling behavior
- Silent degradation of polling functionality

## Fix

Handle the parsing error appropriately with a fallback mechanism:

```go
retryAfterSeconds, err := strconv.Atoi(retryAfterHeader)
if err != nil {
    tflog.Warn(ctx, "Error parsing retry after header, using default wait time: "+err.Error())
    return 30 * time.Second // Use a reasonable default
}
```

Or if the error should be fatal:

```go
retryAfterSeconds, err := strconv.Atoi(retryAfterHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing retry after header: "+err.Error())
    return 0, fmt.Errorf("failed to parse Retry-After header: %w", err)
}
```

## Additional Recommendations

1. **Default Values**: Always provide reasonable fallback values for non-critical parsing failures
2. **Header Validation**: Validate header format before parsing
3. **Consistent Logging**: Use consistent log levels for similar error types
4. **Documentation**: Document the expected header format and fallback behavior

## ISSUE 6

# Issue: Error Only Logged, Not Handled in Managed Environment

**File:** `/workspaces/terraform-provider-power-platform/internal/services/managed_environment/api_managed_environment.go`

## Problem

In the `deleteManagedEnvironment` function, URL parsing errors are only logged using `tflog.Error()` but not returned or handled, allowing execution to continue with potentially invalid URLs.

## Impact

This issue has **low** severity but could lead to:
- Silent failures during managed environment deletion
- Confusing error messages from subsequent operations
- Difficulty in debugging deletion issues
- Potential resource leakage if deletion fails silently

## Location

In the `deleteManagedEnvironment` function:

## Code Issue

```go
operationLocationUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
}
// Execution continues without handling the error
```

## Root Cause

The function logs the parsing error but continues processing, which could lead to attempts to use an invalid URL in subsequent operations.

## Fix

Return the error immediately after logging to prevent further execution with invalid data:

```go
operationLocationUrl, err := url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
    return "", err
}
```

## Additional Recommendations

1. **Consistency**: Apply the same error handling pattern across all managed environment operations
2. **Error Context**: Include more context about which operation was being performed
3. **Validation**: Consider validating the header format before parsing
4. **Testing**: Add unit tests to verify error handling behavior

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
