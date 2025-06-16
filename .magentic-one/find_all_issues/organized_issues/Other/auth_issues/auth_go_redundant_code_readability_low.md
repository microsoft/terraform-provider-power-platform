# Title
Redundant variable assignment and initialization pattern in GetTokenForScopes impacts readability

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
In the function `GetTokenForScopes`, the `token` variable is initialized as `""`, then possibly re-assigned in all cases. The same goes for `tokenExpiry` and `err`. This is unnecessarily verbose and could be simplified for clarity.

## Impact
Low severity. While not a functional error, this pattern reduces readability and can lead to confusion about variable lifetimes or mistaken belief that the initialized zero-value is required.

## Location
In `GetTokenForScopes` method:

## Code Issue
```go
token := ""
var tokenExpiry time.Time
var err error

switch {
	// ...
}
```

## Fix
Declare the variables inline within switch cases, or restructure to:

```go
var (
	token string
	tokenExpiry time.Time
	err error
)

switch {
case ...:
	token, tokenExpiry, err = ...
// ...
}
```

Or:

```go
switch {
case ...:
	returnVal, expiry, err := ...
	// ...
}
```

This is a style/readability improvement. Consider using the code formatter and static analysis for guidance on declaration scope.