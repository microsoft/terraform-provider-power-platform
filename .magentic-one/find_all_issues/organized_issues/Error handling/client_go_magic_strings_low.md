# Title

Magic strings used for CAE challenge detection

##

/workspaces/terraform-provider-power-platform/internal/api/client.go

## Problem

The function `IsCaeChallengeResponse` uses hard-coded substrings `"claims="` and `"insufficient_claims"` to determine if a WWW-Authenticate header is a CAE (Continuous Access Evaluation) challenge. Re-using these string literals directly in code can lead to duplication and errors if the string needs to be updated or checked elsewhere.

## Impact

Reduced maintainability and risk of bugs if string requirements change. Severity: **low**

## Location

In function `IsCaeChallengeResponse`:

## Code Issue

```go
if resp.StatusCode == http.StatusUnauthorized {
	wwwAuthenticate := resp.Header.Get("WWW-Authenticate")
	if wwwAuthenticate != "" {
		return strings.Contains(wwwAuthenticate, "claims=") &&
			strings.Contains(wwwAuthenticate, "insufficient_claims")
	}
}
```

## Fix

Move the string literals to appropriately named `const` values at the top of the file, e.g.:

```go
const wwwAuthenticateClaimsKey = "claims="
const wwwAuthenticateInsufficientClaims = "insufficient_claims"

if resp.StatusCode == http.StatusUnauthorized {
	wwwAuthenticate := resp.Header.Get("WWW-Authenticate")
	if wwwAuthenticate != "" {
		return strings.Contains(wwwAuthenticate, wwwAuthenticateClaimsKey) &&
			strings.Contains(wwwAuthenticate, wwwAuthenticateInsufficientClaims)
	}
}
```
