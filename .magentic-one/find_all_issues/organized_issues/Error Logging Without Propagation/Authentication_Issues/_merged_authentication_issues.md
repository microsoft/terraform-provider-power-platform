# Error Logging Without Propagation - Authentication Issues

This document consolidates all issues related to error logging without proper propagation found in authentication implementations across the Terraform Provider for Power Platform.

## ISSUE 1

# Title

Duplicated logic across multiple authentication methods violates DRY principle

##

/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem

All authentication methods repeat similar logic: initialize credential, call GetToken, handle errors. The structure of methods like `AuthenticateClientCertificate`, `AuthenticateUsingCli`, `AuthenticateClientSecret`, etc., leads to boilerplate code, increased surface for maintenance, and risk of inconsistent handling or logging.

## Impact

Medium. Duplication increases cognitive load, promotes bugs and inconsistencies if one method changes, and makes adding new authentication methods error-prone.

## Location

Across all the `Authenticate*` methods, e.g.:

```go
func (client *Auth) AuthenticateClientCertificate(ctx context.Context, scopes []string) (string, time.Time, error) {
    // ...
    accessToken, err := azureCertCredentials.GetToken(ctx, client.createTokenRequestOptions(ctx, scopes))
    if err != nil {
        return "", time.Time{}, err
    }
    return accessToken.Token, accessToken.ExpiresOn, nil
}
```

And similarly for other auth methods.

## Fix

Abstract common token acquisition logic into a helper function:

```go
func getToken(ctx context.Context, cred azcore.TokenCredential, opts policy.TokenRequestOptions) (string, time.Time, error) {
    token, err := cred.GetToken(ctx, opts)
    if err != nil {
        return "", time.Time{}, err
    }
    return token.Token, token.ExpiresOn, nil
}
```

Then call this helper from each auth method after credential construction. This reduces code repetition and simplifies maintenance.

## ISSUE 2

# Title

Missing retry or backoff logic on HTTP request for OIDC assertion

##

/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem

The `getAssertion` method makes an HTTP request to acquire a token assertion but does not attempt any retry or exponential backoff for robustness. Any temporary network hiccup, throttling, or transient error results in immediate failure. Since this is part of authentication logic, reliability is critical.

## Impact

Medium severity, as this could cause authentication failures due to transient network conditions, cloud API throttling, or hiccups in the OIDC endpoint, reducing reliability of the provider in real-world scenarios.

## Location

In `getAssertion` method:

## Code Issue

```go
resp, err := client.Do(req)
if err != nil {
    return "", fmt.Errorf("getAssertion: cannot request token: %w", err)
}
```

## Fix

Implement retry logic with exponential backoff for network-related errors and server-side temporary errors (such as 429, 500, 502, 503, 504). You can use a helper or a for-loop, e.g.:

```go
for i := 0; i < 3; i++ {
    resp, err := client.Do(req)
    if transientErr(err, resp) {
        time.Sleep(backoff(i))
        continue
    }
    // ... handle normal case ...
    break
}
```

Or integrate a robust HTTP client library with built-in retry mechanisms.

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
