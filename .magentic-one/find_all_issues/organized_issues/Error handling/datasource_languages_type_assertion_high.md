# Title

Error-prone type assertion for ProviderData without validation

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages.go

## Problem

In the `Configure` method, the code directly type asserts `req.ProviderData.(*api.ProviderClient)` without validating that `ProviderData` is in fact of this type. If the assertion fails, this will cause a panic and interrupt control flow, instead of gracefully handling configuration errors.

## Impact

A failed type assertion leads to a provider panic instead of a diagnostic error. This impacts user experience and debugging and is considered a **high severity** issue for Terraform providers.

## Location

```go
clientApi := req.ProviderData.(*api.ProviderClient).Api
if clientApi == nil {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )

    return
}
d.LanguagesClient = newLanguagesClient(clientApi)
```

## Code Issue

```go
clientApi := req.ProviderData.(*api.ProviderClient).Api
```

## Fix

First, check that `ProviderData` is of the correct type via an assertion with the `ok` idiom, and report a diagnostic error if not.

```go
providerClient, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
clientApi := providerClient.Api
if clientApi == nil {
    resp.Diagnostics.AddError(
        "Unexpected nil Api in ProviderClient",
        "The 'Api' field on ProviderClient was nil.",
    )
    return
}
d.LanguagesClient = newLanguagesClient(clientApi)
```
