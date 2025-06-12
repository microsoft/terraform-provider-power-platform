# Title
Error Handling for `resp` Possibly Nil

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/api_environment_settings.go

## Problem

In the `UpdateEnvironmentSettings` method, there is a check on `resp != nil && resp.HttpResponse.StatusCode == http.StatusInternalServerError`. However, `resp` may be nil, but subsequent functions such as `client.Api.HandleForbiddenResponse(resp)` and `client.Api.HandleNotFoundResponse(resp)` are called without verifying this, which may lead to nil pointer dereferences.

## Impact

High. This can cause runtime panics if `resp` is nil.

## Location

```go
if resp != nil && resp.HttpResponse.StatusCode == http.StatusInternalServerError {
    return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
if err != nil {
    return nil, err
}
```

## Code Issue

```go
if resp != nil && resp.HttpResponse.StatusCode == http.StatusInternalServerError {
    return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
if err != nil {
    return nil, err
}
```

## Fix

Return early if `resp` is nil to avoid calling methods on a nil pointer.

```go
if resp == nil {
    return nil, fmt.Errorf("response is nil")
}
if resp.HttpResponse.StatusCode == http.StatusInternalServerError {
    return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
}
if err := client.Api.HandleForbiddenResponse(resp); err != nil {
    return nil, err
}
if err := client.Api.HandleNotFoundResponse(resp); err != nil {
    return nil, err
}
if err != nil {
    return nil, err
}
```
