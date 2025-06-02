# Title
`OidcCredential` and related fields do not follow official Go initialism convention

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
Go code conventionally uses "OIDC" as a capitalized initialism ("OIDC", not "Oidc"). Several types/fields (e.g., `OidcCredential`, `OidcCredentialOptions`, `TokenFilePath`, etc.) use `Oidc` instead of `OIDC`. The same applies to struct fields like `requestToken`, `requestUrl`, and others, which should be `RequestToken`, `RequestURL`, etc. for consistency and clarity.

## Impact
Low severity. Not following naming conventions can lead to inconsistency with the Go ecosystem, reduce code readability, and increase cognitive load for contributors familiar with Go idioms.

## Location
Multiple places throughout the file:

## Code Issue
```go
type OidcCredential struct { ... }
type OidcCredentialOptions struct { ... }
// Plus their usage and fields
```

## Fix
Rename types and fields to use `OIDC` (uppercase) for the initialism, and use `URL` (uppercase), etc. For example:

```go
type OIDCCredential struct {
	requestToken  string
	requestURL    string
	token         string
	tokenFilePath string
	cred          *azidentity.ClientAssertionCredential
}

// and for the options:
type OIDCCredentialOptions struct {
	azcore.ClientOptions
	TenantID      string
	ClientID      string
	RequestToken  string
	RequestURL    string
	Token         string
	TokenFilePath string
}
```

Update all function and method names and usages for consistency as well.