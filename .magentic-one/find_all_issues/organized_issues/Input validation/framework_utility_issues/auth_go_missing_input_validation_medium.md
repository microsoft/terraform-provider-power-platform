# Title
Missing input validation for scopes argument throughout authentication methods

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
Many authentication methods accept `scopes []string` as input, but there is no verification that `scopes` is non-nil and contains valid, non-empty scope URIs. Passing an empty or malformed scope list may lead to hard-to-debug failures downstream or requests for access tokens without correct audience.

## Impact
Medium. If invalid scope input is accepted, may result in confusing Azure SDK/network/server errors or incorrect authentication behavior, reducing robustness for users and making debugging harder.

## Location
Throughout all auth methods including:
- AuthenticateClientCertificate
- AuthenticateUsingCli
- AuthenticateClientSecret
- AuthenticateOIDC
- AuthenticateUserManagedIdentity
- AuthenticateSystemManagedIdentity
- AuthenticateAzDOWorkloadIdentityFederation
- and indirectly through `GetTokenForScopes`

## Code Issue
No validation for argument:
```go
func (client *Auth) AuthenticateClientCertificate(ctx context.Context, scopes []string) (string, time.Time, error) {
    // ...
}
```

## Fix
Validate input at public API boundaries:

```go
if len(scopes) == 0 {
    return "", time.Time{}, errors.New("at least one scope is required for token request")
}
```
And document the behavior in GoDoc comments for each relevant method.