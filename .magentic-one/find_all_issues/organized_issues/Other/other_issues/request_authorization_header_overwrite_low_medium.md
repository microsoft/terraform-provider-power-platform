# Authorization header may overwrite existing user header without warning

##

/workspaces/terraform-provider-power-platform/internal/api/request.go

## Problem

The current code sets the "Authorization" header if it is not already set. However, if the authorization header is already present (possibly set by the user or higher-level caller), the function silently overwrites it when a token is passed. This may cause unexpected behavior if callers intend to use a different authorization token.

## Impact

Severity: Low/Medium

This can lead to confusion or bugs if a user or another part of the code sets a custom Authorization header before calling this function, expecting it to be respected. Overwriting it without explicit notice could create authentication issues or security problems in APIs that support more than one authentication mechanism.

## Location

Lines where the "Authorization" header is set in `doRequest`:

```go
	if request.Header.Get("Authorization") == "" {
		request.Header.Set("Authorization", "Bearer "+*token)
	}
```

## Code Issue

```go
	if request.Header.Get("Authorization") == "" {
		request.Header.Set("Authorization", "Bearer "+*token)
	}
```

## Fix

Explicitly define whether the Authorization header should ever be overwritten and, if not, log a warning and do not overwrite. If overwriting is intended, log a warning so the user/dev is aware:

```go
	if authHeader := request.Header.Get("Authorization"); authHeader == "" {
		request.Header.Set("Authorization", "Bearer "+*token)
	} else {
		tflog.Warn(ctx, "Authorization header already set, not overwriting", map[string]any{
			"existingHeader": authHeader,
		})
	}
```
