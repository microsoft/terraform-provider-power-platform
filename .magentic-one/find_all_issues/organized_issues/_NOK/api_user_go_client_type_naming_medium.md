# Title

Inconsistent Naming: Struct Type `client` Should Be Capitalized

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/api_user.go

## Problem

The `client` struct is defined with a lowercase first letter. According to Go conventions, exported types should use CamelCase and start with an uppercase letter. This may cause confusion and limits the type's usability outside its package, violating Go's standard naming conventions.

## Impact

Severity: Medium

This decreases code readability and maintainability, and restricts use if export is ever required. It runs against Go idioms, which may impact onboarding, documentation quality, and tooling support.

## Location

Defined near the top of the file after `newUserClient`:

## Code Issue

```go
type client struct {
	Api               *api.Client
	environmentClient environment.Client
}
```

## Fix

Capitalize the struct name to `Client` to follow Go naming conventions. You should also update references to this struct throughout the file for consistency.

```go
type Client struct {
	Api               *api.Client
	environmentClient environment.Client
}
```

And all relevant usages should also be updated, for instance:

```go
func (client *Client) EnvironmentHasDataverse(...) { ... }
```

This change improves code quality and future extensibility.
