# Title
Anonymous struct in JSON unmarshal lacks type safety and clarity

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The code unmarshals JSON into an anonymous struct inside the `getAssertion` method. This approach inhibits reusability, makes the code harder to maintain, and defeats type safety gained from concrete types used in multiple locations or tests.

## Impact
Medium severity. It makes the code less testable and harder for static analysis or future refactoring, and reduces documentation/readability.

## Location
In `getAssertion` method:

## Code Issue
```go
var tokenRes struct {
    Count *int    `json:"count"`
    Value *string `json:"value"`
}
if err := json.Unmarshal(body, &tokenRes); err != nil {
    return "", fmt.Errorf("getAssertion: cannot unmarshal response: %w", err)
}
```

## Fix
Define a concrete type at package scope:

```go
type oidcTokenResponse struct {
    Count *int    `json:"count"`
    Value *string `json:"value"`
}
```

Then use:
```go
var tokenRes oidcTokenResponse
if err := json.Unmarshal(body, &tokenRes); err != nil {
    return "", fmt.Errorf("getAssertion: cannot unmarshal response: %w", err)
}
```

This improves type safety, documentation, and maintainability.