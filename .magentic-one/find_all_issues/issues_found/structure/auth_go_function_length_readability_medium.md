# Title
Excessive function length and responsibility in GetTokenForScopes impairs readability

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The `GetTokenForScopes` method handles all logic for choosing the authentication mechanism, test mode, variable initialization, logging, and calling different credential types. This violates the Single Responsibility Principle (SRP), makes the function unnecessarily long, complicated, and more difficult to test or modify as behaviors change or new auth flows are supported.

## Impact
Medium. Hurts maintainability and increases chance for subtle errors or merge conflicts as logic evolves.

## Location
In `GetTokenForScopes` near the end of the file:

## Code Issue
```go
func (client *Auth) GetTokenForScopes(ctx context.Context, scopes []string) (*string, error) {
    // ... lots of mixed logic ...
}
```

## Fix
Refactor method into smaller, purpose-driven helpers/functions. For example, extract a `getTokenWithCredentials` method, separate logging or mode-specific code, etc.

```go
func (client *Auth) GetTokenForScopes(ctx context.Context, scopes []string) (*string, error) {
    if client.config.TestMode {
        // ...
    }
    token, expiry, err := client.getTokenWithCredentials(ctx, scopes)
    // ...
}
```

Each helper can then focus on a single branching strategy or auth path for improved testability and clarity.
