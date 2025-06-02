# Error Handling Incomplete for Unknown Error Types

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments.go

## Problem

In the `Read` function, the error handling for `d.EnvironmentClient.GetDefaultCurrencyForEnvironment` only adds a warning if the error code matches a known value. However, other unexpected error types might occur, and these are currently ignored, potentially leading to silent failures.

## Impact

Unexpected errors that are not warnings or the known error (`ERROR_ENVIRONMENT_URL_NOT_FOUND`) are ignored, resulting in suppressed diagnostics and more difficult troubleshooting for users. Severity: **Medium**.

## Location

```go
defaultCurrency, err := d.EnvironmentClient.GetDefaultCurrencyForEnvironment(ctx, env.Name)
if err != nil {
    if customerrors.Code(err) != customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND {
        resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", env.Name), err.Error())
    }
} else {
    currencyCode = defaultCurrency.IsoCurrencyCode
}
```

## Code Issue

```go
if err != nil {
    if customerrors.Code(err) != customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND {
        resp.Diagnostics.AddWarning(fmt.Sprintf("Error when reading default currency for environment %s", env.Name), err.Error())
    }
} else {
    currencyCode = defaultCurrency.IsoCurrencyCode
}
```

## Fix

Add an explicit branch to handle truly unexpected errors, perhaps with a proper error diagnostic instead of a warning. You may also consider logging unexpected error types for debugging.

```go
if err != nil {
    switch customerrors.Code(err) {
    case customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND:
        // Non-critical, just skip currency.
    default:
        resp.Diagnostics.AddError(
            fmt.Sprintf("Unexpected error when reading default currency for environment %s", env.Name),
            err.Error(),
        )
        return
    }
} else {
    currencyCode = defaultCurrency.IsoCurrencyCode
}
```

This approach helps catch truly unexpected issues and makes debugging easier for consumers.
