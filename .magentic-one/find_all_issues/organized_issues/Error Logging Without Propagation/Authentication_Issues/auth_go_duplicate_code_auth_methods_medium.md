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
