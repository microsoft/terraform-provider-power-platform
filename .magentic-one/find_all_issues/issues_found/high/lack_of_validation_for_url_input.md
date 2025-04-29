# Title

Lack of Validation for URL Input in Execute Method.

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The `Execute` method fails to adequately validate the URL input before parsing it. While some validation is performed (`neturl.Parse(url)`), no checks ensure the URL's scheme or host adheres to the expected format.

## Impact

An improperly formatted URL may pass initial validation but cause runtime errors when interacting with external APIs. Additionally, this can lead to security concerns, such as exploiting weak URL validation mechanisms.

Severity: **High**

## Location

The URL validation block in the `Execute` method:

```go
if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
    return nil, customerrors.NewUrlFormatError(url, e)
}
```

## Code Issue

The current validation mechanism doesn't ensure the URL contains an expected scheme (e.g., HTTPS). Additionally, error messages don't provide actionable details.

```go
if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
    return nil, customerrors.NewUrlFormatError(url, e)
}
```

## Fix

Introduce stricter validation checks for URL input and provide clearer error logs.

```go
if u, e := neturl.Parse(url); e != nil || !u.IsAbs() {
    return nil, customerrors.NewUrlFormatError(url, fmt.Errorf("invalid URL '%s': %w", url, e))
}

if u.Scheme != "https" {
    return nil, fmt.Errorf("URL scheme '%s' is not secure. HTTPS is required.", u.Scheme)
}
```

This fix ensures:
- The URL's scheme is validated as `HTTPS`.
- Error messages include actionable context.