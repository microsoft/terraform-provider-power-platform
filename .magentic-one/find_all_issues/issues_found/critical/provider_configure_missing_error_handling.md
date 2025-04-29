# Title

Error Handling in `provider.Configure` Missing when `getUrls` Returns Undefined

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

In the `Configure` function, the results of `getPublicCloudUrls`, `getGccUrls`, etc., are returned by the `getCloudUrls` methods but are not validated against `nil`. As a result, the provider health and configuration may fail further downstream.

For example:
- No fallback for `nil` cloud URL.
- This results in breaking variable-handling logic further down.

## Impact

Severity: Critical

- Prevents the provider from functioning correctly when the cloud type is invalid or unknown.
- Can cause runtime failures, leading to poor user experience and untrustworthy state.

## Location

Function `Configure`

## Code Issue

```go
switch cloudType {
case string(config.CloudTypePublic):
    providerConfigUrls, cloudConfiguration = getCloudPublicUrls()
case string(config.CloudTypeGcc):
    providerConfigUrls, cloudConfiguration = getGccUrls()
case string(config.CloudTypeGccHigh):
    providerConfigUrls, cloudConfiguration = getGccHighUrls()
case string(config.CloudTypeDod):
    providerConfigUrls, cloudConfiguration = getDodUrls()
case string(config.CloudTypeChina):
    providerConfigUrls, cloudConfiguration = getChinaUrls()
case string(config.CloudTypeEx):
    providerConfigUrls, cloudConfiguration = getExUrls()
case string(config.CloudTypeRx):
    providerConfigUrls, cloudConfiguration = getRxUrls()
default:
    resp.Diagnostics.AddAttributeError(
        path.Root("cloud"),
        "Unknown cloud",
        fmt.Sprintf("The provider cannot create the API client as there is an unknown configuration value for `cloud`. Either set the value in the provider configuration or use the '%s' environment variable.", constants.ENV_VAR_POWER_PLATFORM_CLOUD),
    )
}
```

## Fix

Add nil checks to the `getUrls` functions and default to safe configurations where these are not resolved:

```go
providerConfigUrls, cloudConfiguration := getCloudPublicUrls()
if providerConfigUrls == nil || cloudConfiguration == nil {
    resp.Diagnostics.AddAttributeError(
        path.Root("urls"),
        "Nil Cloud Configuration",
        "The provider cannot configure URLs as cloud configuration returned nil values."
    )
    return
}
```