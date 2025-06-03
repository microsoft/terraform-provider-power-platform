# Title
Switch statement in GetTokenForScopes not explicit and could be made clearer

##
/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem
The `switch` logic in `GetTokenForScopes` relies on boolean functions to branch which credential method is called. While functional, it lacks any `default` or clear comment on fall-through, and is potentially brittle if new credential types are added. It makes it harder for maintainers to reason whether credentials could be acquired from multiple sources or if only the first matched credential method is chosen.

## Impact
Low. Reduced clarity for maintainers; not a correctness problem, but it could make future bugs more likely if credentials logic is changed or extended.

## Location
Within `GetTokenForScopes` method:

## Code Issue
```go
switch {
case client.config.IsClientSecretCredentialsProvided():
    // ...
}
```

## Fix
Add explicit comments, or use a clearer if-else chain with error handling comments or docstring. E.g.:

```go
// Only one authentication method will be used, in precedence order below
if client.config.IsClientSecretCredentialsProvided() {
    // ...
} else if client.config.IsCliProvided() {
    // ...
} // etc.
else {
    return nil, errors.New("no credentials provided")
}
```

Or add comments in the switch:

```go
switch {
// only one branch is selected, in GCC order
}
```
