# Title

Constant Naming Convention Issue - Upper Snake Case Used for Non-Exported Constants

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/constants.go`

## Problem

The file uses **Upper Snake Case** (`NOT_SPECIFIED`, `APP`, etc.) for constants which are not meant to be exported (not available outside the package). According to Go conventions, constants meant for local usage within a package should use **Camel Case** or lowercase naming (`notSpecified`, `app`, etc.).

## Impact

- **Severity**: Medium  
Incorrect naming conventions deviate from Go best practices, making the codebase harder for engineers unfamiliar with the repository to work with.
- Potential confusion arises between exported and non-exported constants, complicating maintainability and readability.

## Location

Lines 4–7: 
```go
const (
	NOT_SPECIFIED = "NotSpecified"
	APP           = "App"
)
```

## Code Issue

The problematic constants are:

```go
const (
	NOT_SPECIFIED = "NotSpecified"
	APP           = "App"
)
```

## Fix

The constants should follow lowerCamelCase for convention adherence:

```go
const (
	notSpecified = "NotSpecified"
	app          = "App"
)
```