# Title

Improper Handling of `Internal Server Error` in `UpdateEnvironmentSettings`

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/api_environment_settings.go

## Problem

In the `UpdateEnvironmentSettings` function, the handling of the `Internal Server Error` status code is incomplete. The conditional check does not properly re-throw the error (`err`) if the response body indicates a failure. As a result, error details can be lost, and debugging can become difficult.

## Impact

This issue can cause misleading error handling, as the actual error might not propagate back up the stack trace. Only the response body message is reported, missing valuable error chain data. Severity: critical

## Location

UpdateEnvironmentSettings, specifically within the conditional block handling `http.StatusInternalServerError`.

## Code Issue

```go
if resp != nil && resp.HttpResponse.StatusCode == http.StatusInternalServerError {
    return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
}
if err != nil {
    return nil, err
}
```

## Fix

Properly chain the error by passing `err` into the `WrapIntoProviderError` function. This captures the details of both `resp.BodyAsBytes` and the `err`. Example:

```go
if resp != nil && resp.HttpResponse.StatusCode == http.StatusInternalServerError {
    return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
}
if err != nil {
    return nil, err
}
```