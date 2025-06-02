# Title

Repeated hardcoded string for "Bearer" authorization in `getAssertion`.

##

/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem

The `"Bearer"` string is hardcoded multiple times within the Authentication headers in the `getAssertion` method:

```go
req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.requestToken))
```

This hardcoding reduces maintainability and increases the likelihood of errors if there is a need to update or localize the authorization scheme.

## Impact

Makes the code hard to maintain and update, and can lead to inconsistencies if the authorization scheme needs to change (severity is low).

## Location

The issue exists in the `getAssertion` method.

## Code Issue

```go
req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.requestToken))
```

## Fix

Define the `"Bearer"` string as a constant at the beginning of the file or in a centralized configuration file:

```go
// Define at the top of the file:
const AuthorizationBearer = "Bearer"

// Use in the method:
req.Header.Set("Authorization", fmt.Sprintf("%s %s", AuthorizationBearer, w.requestToken))
```

Explanation:
- Constants improve readability and maintainability.
- Updating the authorization scheme would require changes in only one place.
