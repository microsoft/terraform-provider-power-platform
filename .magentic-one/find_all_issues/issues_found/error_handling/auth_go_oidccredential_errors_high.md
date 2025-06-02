# Title
Error message details are inconsistent and could be improved for OIDC credential instantiation

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The error checks in `NewOidcCredential` method return hardcoded, low-information errors (e.g., `"tenant is required for OIDC credential"`, `"request Token is required for OIDC credential"`). These are not wrapped or annotated, and do not use constants, leading to less consistency and not enabling callers to distinguish error types. Furthermore, there are places where more contextual information can be given, or a specific error type (possibly a typed error) could be beneficial to facilitate control flow for the caller.

## Impact
Medium severity. Poor error context and inconsistent error reporting can make troubleshooting difficult for users and maintainers, and hinder programmatic error handling downstream.

## Location
Within the implementation of `NewOidcCredential`:

## Code Issue
```go
if c.requestToken == "" {
    return nil, errors.New("request Token is required for OIDC credential")
}
if c.requestUrl == "" {
    return nil, errors.New("request URL is required for OIDC credential")
}
if options.TenantID == "" {
    return nil, errors.New("tenant is required for OIDC credential")
}
if options.ClientID == "" {
    return nil, errors.New("client is required for OIDC credential")
}
```

## Fix
Provide more actionable and detailed error messages, use error wrapping where applicable, and consider custom error types if errors need to be programmatically distinguished. Example:

```go
if c.requestToken == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: requestToken")
}
if c.requestURL == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: requestURL")
}
if options.TenantID == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: TenantID")
}
if options.ClientID == "" {
	return nil, fmt.Errorf("NewOIDCCredential: missing required parameter: ClientID")
}
```
Or, if strict error handling is needed for the API, define specific error types.
